package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/google/go-github/v45/github"
	"github.com/warent/calzone/service/structures/args"
)

type Calzone int

func (t *Calzone) Install(InstallArgs *args.Install, response *int) error {

	// list public repositories for org "github"
	client := github.NewClient(nil)
	_, directory, _, err := client.Repositories.GetContents(context.Background(), "warent", "calzone-repository", InstallArgs.Calzone, nil)
	if err != nil {
		return err
	}

	for _, file := range directory {
		fmt.Println(*file.DownloadURL)
	}
	return nil
}

func main() {
	calzone := new(Calzone)
	server := rpc.NewServer()
	err := server.RegisterName("Calzone", calzone)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	l, err := net.Listen("tcp", ":61895")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	server.Accept(l)
}
