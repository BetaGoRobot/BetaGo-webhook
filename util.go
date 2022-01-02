package main

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func deployNewContainer(containerName, imageName string) {
	ctx := context.Background()
	var timeout time.Duration = 10 * time.Second
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	defer cli.Close()
	if err != nil {
		fmt.Println("Make Docker Client Fails", err.Error())
	}

	fmt.Println("Stopping Container...")
	err = cli.ContainerStop(ctx, containerName, &timeout)
	if err != nil {
		fmt.Println("Stop Container Fails：", err.Error())
	}

	fmt.Println("Removing Container...")
	err = cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{RemoveVolumes: true, RemoveLinks: true, Force: true})
	if err != nil {
		fmt.Println("Remove Container Fails：", err.Error())
	}

	fmt.Println("Pulling Image...")
	reader, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{All: true})
	if err != nil {
		fmt.Println("Pull Container Fails：", err.Error())
	}
	defer reader.Close()

	fmt.Println("Creating Container...")
	_, err = cli.ContainerCreate(ctx, &container.Config{Image: imageName}, nil, nil, nil, containerName)
	if err != nil {
		fmt.Println("Create Container Fails：", err.Error())
	}

	fmt.Println("Starting Container...")
	if err = cli.ContainerStart(ctx, containerName, types.ContainerStartOptions{}); err != nil {
		fmt.Println("Start Container Fails：", err.Error())
	}
	return
}
