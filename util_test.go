package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func Test_deployNewContainer(t *testing.T) {
	type args struct {
		containerName string
		imageName     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				containerName: nightlyContainerName,
				imageName:     nightlyImageName,
			},
		},
		{
			name: "test2",
			args: args{
				containerName: stableContainerName,
				imageName:     stableImageName,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deployNewContainer(tt.args.containerName, tt.args.imageName)
		})
	}
}

func Test_test(t *testing.T) {

	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		return
	}
	closer, err := cli.ImagePull(context.Background(), "kevinmatt/betago:latest", types.ImagePullOptions{})
	if err != nil {
		fmt.Println(err.Error())
	}
	err = closer.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
}
