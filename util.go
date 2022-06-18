package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func deployNewContainer(containerName, imageName string) {
	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = stopContainer(cli, containerName, -1)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = removeContainer(cli, containerName)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = removeImage(cli, imageName)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = PullImg(cli, imageName)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = createContainer(cli, containerName, imageName)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = startContainer(cli, containerName)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Finished")
	return
}

// PullImg  pulls an image from the docker registry
//  @param cli
//  @param imageName
//  @return err
func PullImg(cli *client.Client, imageName string) (err error) {
	fmt.Println("Pulling Image...")
	reader, err := cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		return
	}
	io.Copy(os.Stdout, reader)
	defer reader.Close()
	return
}

func stopContainer(cli *client.Client, containerName string, timeout int64) (err error) {
	var expired = time.Duration(timeout)
	fmt.Println("Stopping Container...")
	err = cli.ContainerStop(context.Background(), containerName, &expired)
	if err != nil {
		return
	}
	return
}

func removeContainer(cli *client.Client, containerName string) (err error) {
	fmt.Println("Removing Container...")
	err = cli.ContainerRemove(context.Background(), containerName, types.ContainerRemoveOptions{})
	if err != nil {
		return
	}
	return
}

func removeImage(cli *client.Client, imageName string) (err error) {
	fmt.Println("Removing Previous Image...")
	_, err = cli.ImageRemove(context.Background(), imageName, types.ImageRemoveOptions{Force: true})
	if err != nil {
		return
	}
	return
}

func createContainer(cli *client.Client, containerName, imageName string) (err error) {
	fmt.Println("Creating Container...")
	_, err = cli.ContainerCreate(context.Background(), &container.Config{Image: imageName, Env: []string{"COM_MES=" + GitRes.Commit.Message, "HTML_URL=" + GitRes.HTMLURL, "COM_URL=" + GitRes.CommentsURL}}, &container.HostConfig{AutoRemove: false, NetworkMode: "betago"}, nil, nil, containerName)
	if err != nil {
		return
	}
	return
}

func startContainer(cli *client.Client, containerName string) (err error) {
	fmt.Println("Starting Container...")
	if err = cli.ContainerStart(context.Background(), containerName, types.ContainerStartOptions{}); err != nil {
		return
	}
	return
}
