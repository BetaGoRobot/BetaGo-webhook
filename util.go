package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func deployNewContainer(containerName, imageName string) {
	ctx := context.Background()
	var timeout time.Duration = 10 * time.Second
	cli, err := client.NewEnvClient()
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
	err = cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		fmt.Println("Remove Container Fails：", err.Error())
	}

	fmt.Println("Removing Previous Image...")
	_, err = cli.ImageRemove(ctx, imageName, types.ImageRemoveOptions{Force: true})
	if err != nil {
		fmt.Println("Removing Image Fails：", err.Error())
	}

	fmt.Println("Pulling Image...")
	reader, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		fmt.Println("Pull Image Fails：", err.Error())
	}
	io.Copy(os.Stdout, reader)
	reader.Close()

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

func PullImg(username, password, imgurl string) error {

	authConfig := types.AuthConfig{
		Username: username, Password: password}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	reader, err := cli.ImagePull(ctx, imgurl, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		return err
	}
	wr, err := io.Copy(os.Stdout, reader)
	fmt.Println(wr)
	if err != nil {
		return err
	}

	return nil
}
