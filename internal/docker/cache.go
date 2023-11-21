package docker

import "github.com/docker/docker/api/types"

var cache = make(map[string]types.ContainerJSON)
