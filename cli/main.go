package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/mitchellh/go-homedir"

	"github.com/warent/calzone/cli/v2/cmd"
)

const ADDR = "192.168.52.1"
const PORT = 61895

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	var running bool
outer:
	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/calzone-service" {
				running = true
				break outer
			}
		}
	}
	if running {
		cli.Close()
		cmd.Execute()
		return
	}

	fmt.Println("It looks like Calzone is not running on your machine. Setting it up now...")

	home, err := homedir.Dir()
	os.Mkdir(home+"/.calzone", 0755)
	os.Mkdir(home+"/.calzone/data", 0755)

	netwk, err := cli.NetworkCreate(ctx, "calzone-net", types.NetworkCreate{
		CheckDuplicate: true,
		IPAM: &network.IPAM{
			Config: []network.IPAMConfig{{
				Subnet:  "192.168.52.0/24",
				Gateway: "192.168.52.254",
			}},
		},
	})
	if err != nil {
		panic(err)
	}

	svcPort, _ := nat.NewPort("tcp", "61895")
	portSet := nat.PortSet{
		svcPort: struct{}{},
	}
	for i := 30000; i <= 40000; i++ {
		appPort, _ := nat.NewPort("tcp", fmt.Sprintf("%v", i))
		portSet[appPort] = struct{}{}
	}

	config := &container.Config{
		Image:        "calzone-service:latest",
		ExposedPorts: portSet,
	}

	hostConfig := &container.HostConfig{
		Runtime: "sysbox-runc",
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: home + "/.calzone/data",
				Target: "/mnt/data",
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "calzone-service")
	if err != nil {
		panic(err)
	}

	err = cli.NetworkConnect(ctx, netwk.ID, resp.ID, nil)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)

}
