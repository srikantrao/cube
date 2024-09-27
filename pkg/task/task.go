package task

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"time"
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

type Task struct {
	ID            uuid.UUID
	State         State
	Name          string
	Image         string
	Memory        int
	Disk          int
	CPU           int
	ExposedPorts  nat.PortSet
	PortBindings  nat.PortMap
	RestartPolicy string
	StartTime     time.Time
	FinishTime    time.Time
}

type TaskEvent struct {
	ID        uuid.UUID
	State     State
	Timestamp time.Time
	Task      *Task
}

type Config struct {
	Name          string
	AttachStdin   bool
	AttachStdout  bool
	AttachStderr  bool
	ExposedPorts  nat.PortSet
	Cmd           []string
	Image         string
	CPU           float64
	Memory        int64
	Disk          int64
	Env           []string
	RestartPolicy string
}

type Docker struct {
	Client      *client.Client
	Config      Config
	ContainerID string
}

type DockerResult struct {
	Error       error
	Action      string
	ContainerID string
	Result      string
}

func (d *Docker) Run() *DockerResult {
	ctx := context.Background()

	// Pull the Image
	reader, err := d.Client.ImagePull(
		ctx, d.Config.Image, image.PullOptions{})
	if err != nil {
		log.Printf("Error pulling image %s: %v\n", d.Config.Image, err)
		return &DockerResult{Error: err}
	}
	io.Copy(os.Stdout, reader)

	// Create the Container
	restartPolicy := container.RestartPolicy{
		Name: container.RestartPolicyMode(d.Config.RestartPolicy),
	}
	resources := container.Resources{
		Memory: d.Config.Memory,
	}
	containerConfig := &container.Config{
		Image: d.Config.Image,
		Env:   d.Config.Env,
	}
	hostConfig := &container.HostConfig{
		Resources:       resources,
		RestartPolicy:   restartPolicy,
		PublishAllPorts: true,
	}
	// create the container
	resp, err := d.Client.ContainerCreate(
		ctx, containerConfig, hostConfig, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf("Error creating container using image %s: %v\n", d.Config.Image, err)
		return &DockerResult{Error: err}
	}
	// start the container
	if err := d.Client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Printf("Error starting container %s: %v\n", d.Config.Name, err)
		return &DockerResult{Error: err}
	}

	d.ContainerID = resp.ID
	logs, err := d.Client.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		log.Printf("Error getting logs from container %s: %v\n", d.Config.Name, err)
		return &DockerResult{Error: err}
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, logs)

	return &DockerResult{
		Action:      "start",
		ContainerID: resp.ID,
		Result:      "success",
	}
}

func (d *Docker) Stop() *DockerResult {
	ctx := context.Background()
	// stop the container
	if err := d.Client.ContainerStop(ctx, d.ContainerID, container.StopOptions{}); err != nil {
		log.Printf("Error stopping container %s: %v\n", d.Config.Name, err)
		return &DockerResult{Error: err}
	}

	// remove the container
	removeOptions := container.RemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	}
	err := d.Client.ContainerRemove(ctx, d.ContainerID, removeOptions)
	if err != nil {
		log.Printf("Error removing container %s: %v\n", d.Config.Name, err)
		return &DockerResult{Error: err}
	}
	return &DockerResult{
		Action:      "stop",
		ContainerID: d.ContainerID,
		Result:      "success",
	}
}
