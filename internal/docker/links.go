package docker

import (
	"context"
	"fmt"
	"goisolator/internal/labels"

	"github.com/sirupsen/logrus"
)

func (svc *DockerService) regenerateGraph() {
	// Empty the Links to regenerate them
	for _, container := range svc.containers {
		container.Links = make(map[string]*Container)
	}

	for _, container := range svc.containers {
		lbl := labels.MapToLabels(container.Config.Labels)
		if !lbl.Enabled {
			continue
		}

		for _, linksTo := range lbl.LinkTo {
			linkedTo, found := svc.containers[linksTo]
			if !found {
				logrus.Warnf("Container %v tried to link to %v but this container has not been found", container.Name(), linksTo)
				continue
			}

			container.Links[linkedTo.Name()] = linkedTo
		}
	}
}

func (svc *DockerService) HandleContainer(ctx context.Context, container *Container) error {
	lbls := labels.MapToLabels(container.Config.Labels)
	if lbls.Ignore || !lbls.Enabled {
		return nil
	}

	// We create the networks for each containers
	err := svc.CreateContainerNetworks(ctx, container)
	if err != nil {
		return fmt.Errorf("failed creating containers networks: %w", err)
	}

	// We connect every link to their networks
	err = svc.ConnectNetworks(ctx, container)
	if err != nil {
		return fmt.Errorf("failed connecting containers links to network: %w", err)
	}

	err = svc.RecreateContainer(ctx, container)
	if err != nil {
		return fmt.Errorf("failed connecting containers links to network: %w", err)
	}

	return nil
}
