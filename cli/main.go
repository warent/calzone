package main

import (
	"context"
	"fmt"
	"log"
	"net/rpc"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/warent/calzone/cli/v2/cmd"
)

const ADDR = "127.0.0.1"
const PORT = 61895

func main() {
	rpcConn, err := rpc.Dial("tcp", fmt.Sprintf("%v:%v", ADDR, PORT))
	if err == nil {
		if err != nil {
			log.Fatal("arith error:", err)
		}
		rpcConn.Close()
		cmd.Execute()
		return
	}

	fmt.Println("It looks like Calzone is not running on your machine. Setting it up now...")

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	// if err != nil {
	// 	panic(err)
	// }
	// io.Copy(os.Stdout, reader)

	config := &container.Config{
		Image: "calzone-service:latest",
		ExposedPorts: nat.PortSet{
			"61895/tcp": struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"61895/tcp": []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "61895",
				},
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)

}
