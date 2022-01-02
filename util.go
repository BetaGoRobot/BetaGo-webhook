package main

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func deployNewContainer() {
	var timeout time.Duration = 5 * time.Second
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Make Docker Client Fails", err.Error())
	}
	err = cli.ContainerStop(context.Background(), "betago", &timeout)
	if err != nil {
		fmt.Println("Stop Container Fails", err.Error())
	}
	err = cli.ContainerRemove(context.Background(), "betago", types.ContainerRemoveOptions{RemoveVolumes: true})
	if err != nil {
		fmt.Println("Remove Container Fails", err.Error())
	}
	_, err = cli.ImagePull(context.Background(), "kevinmatt/betago", types.ImagePullOptions{})
	if err != nil {
		fmt.Println("Pull Container Fails", err.Error())
	}
	_, err = cli.ContainerCreate(context.Background(), &container.Config{Image: "kevinmatt/betago"}, &container.HostConfig{}, &network.NetworkingConfig{}, &v1.Platform{OS: "amd64"}, "betago")
	if err != nil {
		fmt.Println("Create Container Fails", err.Error())
	}
	err = cli.ContainerStart(context.Background(), "betago", types.ContainerStartOptions{})
	if err != nil {
		fmt.Println("Start Container Fails", err.Error())
	}
}
