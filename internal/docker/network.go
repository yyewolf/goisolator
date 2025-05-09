package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/sirupsen/logrus"
)

type Network struct {
	network.CreateOptions

	NetworkID string
	For       *Container
}

func (svc *DockerService) listDockerNetworks() (map[string]*network.Summary, error) {
	dockerNetworks, err := svc.client.NetworkList(context.Background(), network.ListOptions{
		Filters: filters.NewArgs(filters.Arg("label", "goisolator")),
	})
	if err != nil {
		return nil, err
	}

	mappedDockerNetworks := make(map[string]*network.Summary)
	for _, network := range dockerNetworks {
		mappedDockerNetworks[network.Name] = &network
	}

	return mappedDockerNetworks, nil
}

func (svc *DockerService) CreateContainerNetworks(ctx context.Context, container *Container) error {
	dockerNetworks, err := svc.listDockerNetworks()
	if err != nil {
		return fmt.Errorf("failed listing docker networks: %w", err)
	}

	networks := container.GetNetworks()
	for networkName, network := range networks {
		dockerNetwork, found := dockerNetworks[networkName]
		if found {
			network.NetworkID = dockerNetwork.ID
			logrus.Debugf("Skipped network creation because it already exists: %v", networkName)
			continue
		}

		newDockerNetwork, err := svc.client.NetworkCreate(ctx, networkName, network.CreateOptions)
		if err != nil {
			logrus.Errorf("Network %v creation failed: %v", networkName, err)
			continue
		}
		logrus.Debugf("Network created successfully: %v", networkName)

		network.NetworkID = newDockerNetwork.ID
	}

	return nil
}

func (svc *DockerService) FetchContainerNetworks(ctx context.Context, container *Container) error {
	dockerNetworks, err := svc.listDockerNetworks()
	if err != nil {
		return fmt.Errorf("failed listing docker networks: %w", err)
	}

	networks := container.GetNetworks()
	for networkName, network := range networks {
		dockerNetwork, found := dockerNetworks[networkName]
		if found {
			network.NetworkID = dockerNetwork.ID
			continue
		}
	}

	return nil
}

func (svc *DockerService) ConnectNetworks(ctx context.Context, container *Container) error {
	networks := container.GetNetworks()
	for networkName, network := range networks {
		err := svc.client.NetworkConnect(context.Background(), network.NetworkID, network.For.ID, nil)
		if err != nil {
			logrus.Errorf("Could not connect '%v' to network '%v': %v", network.For.Name(), networkName, err)
		} else {
			logrus.Debugf("Connected '%v' to network '%v' successfully", network.For.Name(), networkName)
		}

		err = svc.client.NetworkConnect(context.Background(), network.NetworkID, container.ID, nil)
		if err != nil {
			logrus.Errorf("Could not connect '%v' to network '%v': %v", container.Name(), networkName, err)
		} else {
			logrus.Debugf("Connected '%v' to network '%v' successfully", container.Name(), networkName)
		}
	}

	return nil
}

func (svc *DockerService) DisconnectLinksFromNetwork(ctx context.Context, container *Container) error {
	networks := container.GetNetworks()
	for networkName, network := range networks {
		err := svc.client.NetworkDisconnect(context.Background(), network.NetworkID, network.For.ID, true)
		if err != nil {
			logrus.Errorf("Could not disconnect '%v' from network '%v': %v", network.For.Name(), networkName, err)
			continue
		}

		logrus.Debugf("Disconnected '%v' from network '%v' successfully", network.For.Name(), networkName)
	}

	return nil
}
