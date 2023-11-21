package main

import (
	"goisolator/internal/docker"
)

func main() {
	for {
		docker.StartListener()
	}
}
