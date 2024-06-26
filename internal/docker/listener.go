package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/sirupsen/logrus"
)

func StartListener() {
	// Listen for new containers
	// When a new container is created, call the function below
	events, _ := cli.Events(context.Background(), types.EventsOptions{})

	for event := range events {
		if event.Action == "start" {
			container, err := cli.ContainerInspect(context.Background(), event.ID)
			if err != nil {
				logrus.Error(err)
				continue
			}

			logrus.Infof("Container found: %s", container.Name)

			DoIsolationAtoB(container)
			DoIsolationBtoA(container)

			Reconciliate()
		}
		if event.Action == "destroy" {
			for _, container := range cache {
				if container.ID == event.ID {
					logrus.Infof("Container destroyed: %s", container.Name)
					delete(cache, container.Name)

					// Prune networks
					cli.NetworksPrune(context.Background(), filters.NewArgs(filters.Arg("label", "goisolator")))
					break
				}
			}
		}
	}
}
