package docker

import (
	"context"
	"fmt"
	"goisolator/internal/labels"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/sirupsen/logrus"
)

func ContainerName(cn string) string {
	if cn[0] == '/' {
		cn = cn[1:]
	}

	if cn[len(cn)-2] == '-' {
		cn = cn[:len(cn)-2]
	}
	return cn
}

func DoIsolationAtoB(container types.ContainerJSON) {
	labels := labels.MapToLabels(container.Config.Labels)
	logrus.Debugf("Labels: %+v", labels)

	if len(labels.LinkTo) == 0 {
		logrus.Debugf("No links found, skipping...")
		return
	}

	for _, link := range labels.LinkTo {
		logrus.Debugf("Checking link from %s to %s", container.ID, link)

		var c2 types.ContainerJSON

		for _, c := range cache {
			// Chech link == c.name while ignoring / and -1
			if link == ContainerName(c.Name) {
				c2 = c
				break
			}
		}

		if c2.ContainerJSONBase == nil {
			logrus.Debugf("Container %s not found", link)
			continue
		}

		// Create network (doesn't do anything if already exists)
		NetworkAtoB(container.Name, link)
		err := LinkAandB(container, c2)
		if err != nil {
			logrus.Error(err)
			continue
		}
		logrus.Infof("Linked %s to %s", container.Name, c2.Name)
	}

}

func DoIsolationBtoA(container types.ContainerJSON) {
	// Check if any container are linked to this container
	// If so, link them

	for _, c := range cache {
		labels := labels.MapToLabels(c.Config.Labels)
		logrus.Debugf("Labels: %+v", labels)

		if len(labels.LinkTo) == 0 {
			logrus.Debugf("No links found, skipping...")
			continue
		}

		for _, link := range labels.LinkTo {
			if link == ContainerName(container.Name) {
				NetworkAtoB(c.Name, container.Name)
				err := LinkAandB(c, container)
				if err != nil {
					logrus.Error(err)
					continue
				}

				logrus.Infof("Linked %s to %s", c.Name, container.Name)
			}
		}
	}
}

func NetworkAtoB(a, b string) error {
	// Make sure the network exists
	// If not, create it
	networkName := fmt.Sprintf("goisolator+%s+%s", ContainerName(a), ContainerName(b))

	args := filters.NewArgs(filters.Arg("label", "goisolator"))

	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: args,
	})
	if err != nil {
		return err
	}

	for _, network := range networks {
		logrus.Debugf("Found network: %s", network.Name)
		if network.Name == networkName {
			logrus.Debugf("Network %s found", networkName)
			return nil
		}
	}

	logrus.Debugf("Network %s not found, creating...", networkName)
	// Make sure the two containers will be able to communicate
	_, err = cli.NetworkCreate(context.Background(), networkName, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "bridge",
		Attachable:     true,
		Labels: map[string]string{
			"goisolator":   "true",
			"goisolator.a": ContainerName(a),
			"goisolator.b": ContainerName(b),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func LinkAandB(a, b types.ContainerJSON) error {
	// Make sure the network exists
	// If not, create it
	networkName := fmt.Sprintf("goisolator+%s+%s", ContainerName(a.Name), ContainerName(b.Name))

	args := filters.NewArgs(filters.Arg("label", "goisolator"))

	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: args,
	})
	if err != nil {
		logrus.Fatal(err)
	}

	nw := types.NetworkResource{}

	for _, network := range networks {
		if network.Name == networkName {
			logrus.Debugf("Network %s found", networkName)
			break
		}
		nw = network
	}

	// Add network to A and B
	err1 := cli.NetworkConnect(context.Background(), nw.ID, a.ID, nil)
	err2 := cli.NetworkConnect(context.Background(), nw.ID, b.ID, nil)

	// Make sure containers can communicate using

	if err1 != nil {
		return err1
	}

	if err2 != nil {
		return err2
	}

	return nil
}
