package docker

import (
	"context"
	"fmt"
	"goisolator/internal/labels"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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

	var l []string

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
		id, err := NetworkAtoB(container.Name, link)
		if err != nil {
			logrus.Error(err)
			continue
		}

		l = append(l, id)

		err = LinkBtoNetwork(c2, id)
		if err != nil {
			logrus.Error(err)
			continue
		}
		logrus.Infof("Linked %s to %s", c2.Name, container.Name)
	}

	if labels.Ignore {
		logrus.Debugf("Container %s ignored", container.ID)
		return
	}

	err := LinkAToL(container, l)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("Linked %s to others", container.Name)

}

func DoIsolationBtoA(container types.ContainerJSON) {
	// Check if any container are linked to this container
	// If so, link them
	lbl := labels.MapToLabels(container.Config.Labels)
	logrus.Debugf("Labels: %+v", lbl)

	for _, c := range cache {
		labels := labels.MapToLabels(c.Config.Labels)
		logrus.Debugf("Labels: %+v", labels)

		if len(labels.LinkTo) == 0 {
			logrus.Debugf("No links found, skipping...")
			continue
		}

		var l []string

		for _, link := range labels.LinkTo {
			if link == ContainerName(container.Name) {
				id, err := NetworkAtoB(c.Name, container.Name)
				if err != nil {
					logrus.Error(err)
					continue
				}
				l = append(l, id)
				err = LinkBtoNetwork(c, id)
				if err != nil {
					logrus.Error(err)
					continue
				}

				logrus.Infof("Linked %s to %s", container.Name, c.Name)
			}
		}

		if labels.Ignore {
			logrus.Debugf("Container %s ignored", container.ID)
			continue
		}

		err := LinkAToL(c, l)
		if err != nil {
			logrus.Error(err)
			continue
		}
		logrus.Infof("Linked %s to others", c.Name)
	}
}

func NetworkAtoB(a, b string) (string, error) {
	// Make sure the network exists
	// If not, create it
	networkName := fmt.Sprintf("goisolator+%s+%s", ContainerName(a), ContainerName(b))

	args := filters.NewArgs(filters.Arg("label", "goisolator"))

	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: args,
	})
	if err != nil {
		return "", err
	}

	for _, network := range networks {
		logrus.Debugf("Found network: %s", network.Name)
		if network.Name == networkName {
			logrus.Debugf("Network %s found", networkName)
			return network.ID, nil
		}
	}

	logrus.Debugf("Network %s not found, creating...", networkName)
	// Make sure the two containers will be able to communicate
	resp, err := cli.NetworkCreate(context.Background(), networkName, types.NetworkCreate{
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
		return "", err
	}

	return resp.ID, nil
}

func LinkBtoNetwork(b types.ContainerJSON, nw string) error {
	//Add network to B
	err := cli.NetworkConnect(context.Background(), nw, b.ID, nil)
	if err != nil {
		return err
	}

	return nil
}

func LinkAToL(a types.ContainerJSON, nws []string) error {
	var err error
	for _, nw := range nws {
		cli.NetworkConnect(context.Background(), nw, a.ID, nil)
	}

	a, err = cli.ContainerInspect(context.Background(), a.ID)
	if err != nil {
		return err
	}

	err = cli.ContainerStop(context.Background(), a.ID, container.StopOptions{
		Signal: "SIGKILL",
	})
	if err != nil {
		return err
	}
	err = cli.ContainerRemove(context.Background(), a.ID, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		return err
	}

	// Recreate A with same config
	nwc := a.NetworkSettings.Networks

	a.Config.Labels["goisolator.ignore"] = "true"

	lbl := labels.MapToLabels(a.Config.Labels)
	if lbl.Traefik {
		a.Config.Labels["traefik.enable"] = "true"
	}

	resp, err1 := cli.ContainerCreate(context.Background(), a.Config, a.HostConfig, nil, nil, a.Name)
	// Connect networks
	for _, v := range nwc {
		cli.NetworkConnect(context.Background(), v.NetworkID, resp.ID, nil)
	}

	// Start container
	err = cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return err1
}
