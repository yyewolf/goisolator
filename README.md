<p align="center">
  <p align="center">
    <a href="https://github.com/yyewolf/goisolator/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/yyewolf/goisolator.svg?style=flat-square"></a>
    <a href="/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
    <a href="https://codeclimate.com/github/yyewolf/goisolator/test_coverage"><img src="https://api.codeclimate.com/v1/badges/d9fcf617937d6026221f/test_coverage" /></a>
    <a href="https://codeclimate.com/github/yyewolf/goisolator/maintainability"><img src="https://api.codeclimate.com/v1/badges/d9fcf617937d6026221f/maintainability" /></a>
    <a href="https://goreportcard.com/report/github.com/yyewolf/goisolator/backend"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/yyewolf/goisolator/backend"></a>
    <a href="https://godoc.org/github.com/yyewolf/goisolator/backend"><img src="https://godoc.org/github.com/yyewolf/goisolator/backend?status.svg" alt="GoDoc"></a>
  </p>
</p>

# Go Isolator

**Go Isolator** is a utility for infrastructures using `docker compose` as a method of deployments. It allows you to isolate containers and avoid one network behind your reverse proxy.

# Installation

You can use the docker pull command to fetch the latest version of Go Isolator from the GitHub Container Registry (ghcr.io):

```bash
docker pull ghcr.io/yyewolf/goisolator:v1.0.9
```

# Usage

Go Isolator operates by leveraging container labels and names to enable or restrict network access between containers. To link containers A and B to container C without interconnecting them, use the following label:

```bash
goisolator.linkto=<container_name>
```

Replace <container_name> with the actual name of the container you want to link to. This ensures that containers A and B can access container C without being interconnected.

# Example

Suppose you have containers named container_a, container_b, and container_c. To allow container_a and container_b to access container_c, add the following labels to the respective container definitions in your Docker Compose file:

```yaml
version: '3'
services:
  container_a:
    image: your_image_a
    labels:
      - "goisolator.linkto=container_c"
    # other configurations for container_a

  container_b:
    image: your_image_b
    labels:
      - "goisolator.linkto=container_c"
    # other configurations for container_b

  container_c:
    image: your_image_c
    # other configurations for container_c
```

Now, `container_a` and `container_b` can access `container_c` without being interconnected.

# Contributing

If you'd like to contribute to Go Isolator, please follow the contribution guidelines.
License

This project is licensed under the MIT License - see the LICENSE file for details.

# Acknowledgments

Special thanks to [@Butanal](https://github.com/Butanal) for the help !

# Support

For any issues or questions, please open an issue.
