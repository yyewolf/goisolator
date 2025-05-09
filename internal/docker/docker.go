package docker

import (
	"context"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

type DockerService struct {
	client *client.Client

	containers map[string]*Container

	wg sync.WaitGroup
}

func NewDockerService(cli *client.Client) *DockerService {
	logrus.Info("Docker client initialized")

	return &DockerService{
		client:     cli,
		containers: make(map[string]*Container),
	}
}

func (svc *DockerService) Start(ctx context.Context) {
	svc.ReconciliateLoop(ctx)
	svc.StartListener(ctx)

	svc.wg.Wait()
}

func (svc *DockerService) ReconciliateLoop(ctx context.Context) {
	svc.wg.Add(1)
	defer svc.wg.Done()

	t := time.NewTicker(10 * time.Minute)
	go func() {
		svc.Reconciliate(ctx)

		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				svc.Reconciliate(ctx)
			}
		}
	}()
}

func (svc *DockerService) Reconciliate(ctx context.Context) {
	containerList, err := svc.client.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		logrus.Fatal(err)
	}

	for _, container := range containerList {
		err := svc.AddContainer(ctx, container.ID)
		if err != nil {
			logrus.Errorf("Reconciliation got error when adding container: %v", err)
			continue
		}
	}
}

func (svc *DockerService) AddContainer(ctx context.Context, id string) error {
	dockerContainer, err := svc.client.ContainerInspect(ctx, id)
	if err != nil {
		return err
	}

	container := &Container{
		ContainerJSON: dockerContainer,
	}

	logrus.Infof("Container found: %s", container.Name())

	svc.containers[container.Name()] = container

	svc.regenerateGraph()
	return svc.HandleContainer(ctx, container)
}
