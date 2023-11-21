package main

import (
	"goisolator/internal/config"
	"goisolator/internal/docker"
)

func main() {
	config.GetConfig()

	for {
		docker.StartListener()
	}
}
