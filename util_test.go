package main

import (
	"testing"

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

}

func TestPullImg(t *testing.T) {
	cli, _ := client.NewClientWithOpts(client.FromEnv)
	PullImg(cli, stableImageName)
}

func Test_splitTest(t *testing.T) {
	splitTest()
}
