package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/google/go-github/v45/github"
	"github.com/warent/calzone/service/structures"
	"gopkg.in/yaml.v3"
)

var yamlCache map[string]string = map[string]string{}

type Calzone int

func (t *Calzone) BeginInstall(BeginInstallArgs *structures.BeginInstallArgs, BeginInstallResponse *structures.BeginInstallResponse) error {

	// list public repositories for org "github"
	client := github.NewClient(nil)
	_, directory, _, err := client.Repositories.GetContents(context.Background(), "warent", "calzone-repository", BeginInstallArgs.Calzone, nil)
	if err != nil {
		return err
	}

	files := map[string]string{}

	for _, file := range directory {
		files[path.Base(*file.DownloadURL)] = *file.DownloadURL
	}

	if _, ok := files["parameters.yaml"]; ok {
		resp, err := http.Get(files["parameters.yaml"])
		if err != nil {
			return fmt.Errorf("GET error: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Status error: %v", resp.StatusCode)
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Read body: %v", err)
		}
		resp.Body.Close()

		err = yaml.Unmarshal(data, BeginInstallResponse)
		if err != nil {
			return err
		}
	}

	resp, err := http.Get(files["calzone.yaml"])
	if err != nil {
		return fmt.Errorf("GET error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Read body: %v", err)
	}
	resp.Body.Close()
	yamlCache[BeginInstallArgs.Calzone] = string(data)

	return nil
}

func (t *Calzone) CompleteInstall(CompleteInstallArgs *structures.CompleteInstallArgs, CompleteInstallResponse *structures.CompleteInstallResponse) error {
	configData := yamlCache[CompleteInstallArgs.Calzone]
	config := structures.CalzoneConfig{}
	tmpl, err := template.New("config").Parse(configData)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, *CompleteInstallArgs)
	if err != nil {
		return err
	}

	configParsed := buf.Bytes()

	err = yaml.Unmarshal(configParsed, &config)
	if err != nil {
		return err
	}

	reader, writer := io.Pipe()
	cmdCtx, cmdDone := context.WithCancel(context.Background())

	scannerStopped := make(chan struct{})
	go func() {
		defer close(scannerStopped)

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			fmt.Printf("SCANNED LINE %s\n", scanner.Text())
		}
	}()

	safeId := strings.ReplaceAll(CompleteInstallArgs.Calzone, "/", "--")

	os.Mkdir(fmt.Sprintf("/mnt/data/%s", safeId), 0755)

	cmd := exec.Command("bash", "-c", fmt.Sprintf("minikube start --driver=docker -p %s --mount-string=\"/mnt/data/%s:/mnt/data\" --mount --memory=%v --cpus=%v -o json", safeId, safeId, config.System.Memory, config.System.Cpus))
	cmd.Stdout = writer
	_ = cmd.Start()
	go func() {
		_ = cmd.Wait()
		cmdDone()
		writer.Close()
	}()
	<-cmdCtx.Done()

	<-scannerStopped

	CompleteInstallResponse.Port = 30000
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
