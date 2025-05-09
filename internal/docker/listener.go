package docker

import (
	"context"
	"errors"
	"goisolator/internal/labels"
	"io"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/sirupsen/logrus"
)

func (svc *DockerService) StartListener(ctx context.Context) {
	svc.wg.Add(1)
	defer svc.wg.Done()

restartListen:
	// Listen for new containers
	// When a new container is created, call the function below
	eventChan, errorChan := svc.client.Events(ctx, events.ListOptions{})

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-errorChan:
			logrus.Errorf("Got error from event handler: %v", err)
			if errors.Is(err, io.EOF) {
				goto restartListen
			}
			continue
		case event := <-eventChan:
			if event.Type != events.ContainerEventType {
				continue
			}

			switch event.Action {
			case "start":
				err := svc.AddContainer(ctx, event.ID)
				if err != nil {
					logrus.Errorf("Got error when adding container: %v", err)
					continue
				}
			case "destroy":
				for _, container := range svc.containers {
					if event.ID != container.ID {
						continue
					}

					if container.Recreation {
						continue
					}

					logrus.Infof("Container destroyed: %v", container.Name())

					lbls := labels.MapToLabels(container.Config.Labels)
					if !lbls.Enabled {
						continue
					}

					delete(svc.containers, container.Name())

					// Clean useless goisolator networks
					svc.FetchContainerNetworks(ctx, container)
					svc.DisconnectLinksFromNetwork(ctx, container)
					svc.client.NetworksPrune(context.Background(), filters.NewArgs(filters.Arg("label", "goisolator")))
				}
			}
		}
	}

}
