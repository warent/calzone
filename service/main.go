package main

import (
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
	"path"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/go-github/v45/github"
	"github.com/warent/calzone/service/structures"
	"gopkg.in/yaml.v3"
)

var yamlCache map[string]string = map[string]string{}
var messageQueue map[string][]string = map[string][]string{}

type Calzone int

func (t *Calzone) GetMessages(calzone string, Messages *[]string) error {
	msgs := make([]string, len(messageQueue[calzone]))
	copy(msgs, messageQueue[calzone])
	messageQueue[calzone] = []string{}
	*Messages = msgs
	return nil
}

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

	// reader, writer := io.Pipe()
	// cmdCtx, cmdDone := context.WithCancel(context.Background())

	// scannerStopped := make(chan struct{})
	// go func() {
	// 	defer close(scannerStopped)
	// 	messageQueue[CompleteInstallArgs.Calzone] = []string{}

	// 	scanner := bufio.NewScanner(reader)
	// 	for scanner.Scan() {
	// 		blob := scanner.Text()
	// 		parsed := map[string]interface{}{}
	// 		err := json.Unmarshal([]byte(blob), &parsed)
	// 		if err != nil {
	// 			fmt.Println(err, blob)
	// 			continue
	// 		}
	// 		data := parsed["data"].(map[string]interface{})
	// 		// Use data["message"] for more verbose output
	// 		if data != nil && data["name"] != nil {
	// 			messageQueue[CompleteInstallArgs.Calzone] = append(messageQueue[CompleteInstallArgs.Calzone], data["name"].(string))
	// 		}
	// 	}
	// }()

	safeId := strings.ReplaceAll(CompleteInstallArgs.Calzone, "/", "--")

	os.Mkdir(fmt.Sprintf("/mnt/data/%s", safeId), 0755)

	// cmd := exec.Command("bash", "-c", fmt.Sprintf("minikube start --driver=docker -p %s --mount-string=\"/mnt/data/%s:/mnt/data\" --mount --memory=%v --cpus=%v -o json", safeId, safeId, config.System.Memory, config.System.Cpus))
	// cmd.Stdout = writer
	// _ = cmd.Start()
	// go func() {
	// 	_ = cmd.Wait()
	// 	cmdDone()
	// 	writer.Close()
	// }()
	// <-cmdCtx.Done()

	// <-scannerStopped

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	for depname, depconf := range config.Deployments {
		mounts := []mount.Mount{}
		for _, vol := range depconf.Volumes {
			parts := strings.Split(vol, ":")
			volname := parts[0]
			path := parts[1]
						{
			// 	Type:   mount.TypeBind,
			// 	Source: home + "/.calzone/data",
			// 	Target: "/mnt/data",
			// },
			mounts = append(mounts, mount.Mount{
				Type:   mount.TypeBind,
				Source: "/mnt/data/%s/%s",
				Target: path,
			})
		}

		stream, err := cli.ImagePull(context.Background(), depconf.Image, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}
		io.ReadAll(stream)
		stream.Close()

		img, _, err := cli.ImageInspectWithRaw(context.Background(), depconf.Image)
		if err != nil {
			panic(err)
		}

		containerConfig := &container.Config{
			Image:        depconf.Image,
			ExposedPorts: img.Config.ExposedPorts,
		}

		bindings := nat.PortMap{}
		for port := range containerConfig.ExposedPorts {
			bindings[port] = []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port.Port(),
				},
			}
		}

		hostConfig := &container.HostConfig{
			Mounts:       mounts,
			PortBindings: bindings,
		}

		resp, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, nil, depname)
		if err != nil {
			panic(err)
		}

		if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}

	}

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
