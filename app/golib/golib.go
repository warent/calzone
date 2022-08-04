package main

import "C"
import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/mitchellh/go-homedir"
)

//export GetKey
func GetKey() {
	defer func() {
		if x := recover(); x != nil {
			// open output file
			_, err := os.Stat("err.golib.txt")
			if err == os.ErrNotExist {
				os.Create("err.golib.txt")
			} else {
				panic(err)
			}
			fo, err := os.Open("err.golib.txt")
			if err != nil {
				panic(err)
			}
			// close fo on exit and check for its returned error
			defer func() {
				if err := fo.Close(); err != nil {
					panic(err)
				}
			}()
			fo.WriteString(fmt.Sprintf("%+v", x))
			panic(x)
		}
	}()

	const ADDR = "192.168.52.1"
	const PORT = 61895
	var err error
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, _ := cli.ContainerList(ctx, types.ContainerListOptions{})
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
		// cmd.Execute()
		return
	}

	fmt.Println("It looks like Calzone is not running on your machine. Setting it up now...")

	home, err := homedir.Dir()
	os.Mkdir(home+"/.calzone", 0755)
	os.Mkdir(home+"/.calzone/data", 0755)

	_, err = cli.NetworkCreate(ctx, "calzone-net", types.NetworkCreate{
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

}

func main() {}
