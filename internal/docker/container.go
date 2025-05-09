package docker

import (
	"fmt"

	"github.com/docker/docker/api/types"
)

type Container struct {
	types.ContainerJSON

	Links      map[string]*Container
	Recreation bool

	Networks map[string]*Network
}

func (container *Container) Name() string {
	cn := container.ContainerJSON.Name
	if cn[0] == '/' {
		cn = cn[1:]
	}

	if cn[len(cn)-2] == '-' {
		cn = cn[:len(cn)-2]
	}
	return cn
}

func (container *Container) GetNetworks() map[string]*Network {
	if container.Networks != nil {
		return container.Networks
	}

	output := make(map[string]*Network)

	for _, link := range container.Links {
		networkName := fmt.Sprintf("goisolator+%s+%s", container.Name(), link.Name())

		output[networkName] = &Network{
			For: link,
			NetworkCreate: types.NetworkCreate{
				CheckDuplicate: true,
				Driver:         "bridge",
				Attachable:     true,
				Labels: map[string]string{
					"goisolator":           "true",
					"goisolator.container": container.Name(),
					"goisolator.link":      link.Name(),
				},
			},
		}
	}

	container.Networks = output

	return output
}
