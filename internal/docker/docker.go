package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

var cli *client.Client

func init() {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Docker client initialized")

	// Fill up containers
	c, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		logrus.Fatal(err)
	}

	for _, container := range c {
		logrus.Infof("Container found: %s", container.Names[0])

		inspect, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			logrus.Error(err)
		}

		cache[container.Names[0]] = inspect
	}

	for _, container := range cache {
		DoIsolationAtoB(container)
		DoIsolationBtoA(container)
	}
}
