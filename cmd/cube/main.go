package main

import (
	"context"
	"cube/pkg/task"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"time"
)

func main() {

	// print message that the container is starting
	fmt.Println("Starting container")
	// Create a container
	d, result := createContainer()
	if result.Error != nil {
		return
	}

	// Stop the container
	stopContainer(context.Background(), d)
	// Remove the container
	err := d.Client.ContainerRemove(context.Background(), d.ContainerID, container.RemoveOptions{})
	if err != nil {
		fmt.Printf("Error removing container %s: %v\n", d.Config.Name, err)
		return

	}
}

func createContainer() (*task.Docker, *task.DockerResult) {
	// Create a container with the specified configuration
	config := task.Config{
		Name:  "test-container-1",
		Image: "postgres:latest",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
		},
	}

	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d := task.Docker{
		Client: dc,
		Config: config,
	}

	result := d.Run()
	if result.Error != nil {
		fmt.Printf("Error running container: %v\n", result.Error)
		return nil, result
	}
	fmt.Printf("Container %s created successfully\n", result.ContainerID)
	return &d, result
}

func stopContainer(ctx context.Context, d *task.Docker) {
	// Stop the container
	timeout := time.Second * 10
	err := d.Client.ContainerStop(ctx, d.ContainerID, &timeout)
	if err != nil {
		fmt.Printf("Error stopping container %s: %v\n", d.Config.Name, err)
		return
	}
	fmt.Printf("Container %s stopped successfully\n", d.Config.Name)
}
