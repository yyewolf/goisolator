package docker

import (
	"context"
	"fmt"
	"goisolator/internal/labels"

	"github.com/docker/docker/api/types/container"
)

func (svc *DockerService) RecreateContainer(ctx context.Context, c *Container) error {
	a, err := svc.client.ContainerInspect(context.Background(), c.ID)
	if err != nil {
		return err
	}

	c.ContainerJSON = a

	err = svc.client.ContainerStop(context.Background(), c.ID, container.StopOptions{Signal: "SIGKILL"})
	if err != nil {
		return fmt.Errorf("error stopping the container: %w", err)
	}

	c.Recreation = true
	err = svc.client.ContainerRemove(context.Background(), c.ID, container.RemoveOptions{Force: true})
	if err != nil {
		return fmt.Errorf("error removing the container: %w", err)
	}

	// Recreate A with same config
	nwc := c.NetworkSettings.Networks

	c.Config.Labels["goisolator.ignore"] = "true"

	lbl := labels.MapToLabels(c.Config.Labels)
	if lbl.Traefik {
		c.Config.Labels["traefik.enable"] = "true"
	}

	resp, err := svc.client.ContainerCreate(context.Background(), c.Config, c.HostConfig, nil, nil, c.ContainerJSON.Name)
	if err != nil {
		return fmt.Errorf("error creating the container: %w", err)
	}

	// Connect networks
	for _, v := range nwc {
		err := svc.client.NetworkConnect(context.Background(), v.NetworkID, resp.ID, nil)
		if err != nil {
			return fmt.Errorf("error connecting the container to one of its network: %w", err)
		}
	}

	// Start container
	err = svc.client.ContainerStart(context.Background(), resp.ID, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("error starting the container: %w", err)
	}

	return nil
}
