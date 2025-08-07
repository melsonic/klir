package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/charmbracelet/huh"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/urfave/cli/v3"
)

type ContainerItem struct {
	ID    string
	Name  string
	Image string
}

type ImageItem struct {
	ID   string
	Name string
	Size int64
}

type DockerClient struct {
	client *client.Client
}

func NewDockerClient() *DockerClient {
	dockerCli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		slog.Debug("error initializing new docker client", "error", err.Error())
		return nil
	}

	_, err = dockerCli.Ping(context.Background())
	if err != nil {
		slog.Debug("error pinging docker server", "error", err.Error())
		return nil
	}

	defer dockerCli.Close()

	return &DockerClient{
		client: dockerCli,
	}
}

func (dc *DockerClient) StopRunningContainers(ctx context.Context, cmd *cli.Command) error {
	if cmd.Bool("verbose") {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	containers, err := dc.client.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		slog.Debug("error fetching container list", "error", err.Error())
		return err
	}

	if len(containers) == 0 {
		fmt.Println("No running containers found.")
		return nil
	}

	var runningContainers []*ContainerItem
	var maxNameLen int = 0

	for i := range containers {
		runningContainers = append(runningContainers, &ContainerItem{
			ID:    containers[i].ID,
			Name:  containers[i].Names[0],
			Image: containers[i].Image,
		})
		maxNameLen = max(maxNameLen, len(containers[i].Names[0]))
	}

	var selectedContainers []*ContainerItem

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[*ContainerItem]().
				Title("Select the containers to STOP").
				Options(ParseContainerItemList(runningContainers, maxNameLen)...).
				Value(&selectedContainers),
		),
	)

	err = form.Run()
	if err != nil {
		slog.Debug("error creating a form", "error", err.Error())
		return err
	}

	for i := range selectedContainers {
		err = dc.client.ContainerStop(context.Background(), selectedContainers[i].ID, container.StopOptions{})
		if err != nil {
			fmt.Printf("\x1b[31mx\x1b[0m Error Stopping Container %s\n", selectedContainers[i].Name)
			slog.Debug("Error Stopping Container", "Container ID", selectedContainers[i].ID, "Error", err.Error())
		} else {
			fmt.Printf("\x1b[32m✓\x1b[0m Container %s Stopped\n", selectedContainers[i].Name)
			slog.Debug("Stopped Container", "Container ID", selectedContainers[i].ID)
		}
	}

	return nil
}

func (dc *DockerClient) RemoveDockerContainers(ctx context.Context, cmd *cli.Command) error {

	if cmd.Bool("verbose") {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	forceRemoval := cmd.Bool("force")

	containers, err := dc.client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		slog.Debug("Error fetching container list", "error", err.Error())
		return err
	}

	if len(containers) == 0 {
		fmt.Println("No containers found.")
		return nil
	}

	var stoppedContainers []*ContainerItem
	var maxNameLen int = 0

	for i := range containers {
		if forceRemoval || containers[i].State == container.StatePaused || containers[i].State == container.StateExited || containers[i].State == container.StateDead {
			stoppedContainers = append(stoppedContainers, &ContainerItem{
				ID:    containers[i].ID,
				Name:  containers[i].Names[0],
				Image: containers[i].Image,
			})
			maxNameLen = max(maxNameLen, len(containers[i].Names[0]))
		}
	}

	if len(stoppedContainers) == 0 {
		fmt.Println("No inactive containers found")
		return nil
	}

	var selectedContainers []*ContainerItem

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[*ContainerItem]().
				Title("Select the containers to REMOVE").
				Options(ParseContainerItemList(stoppedContainers, maxNameLen)...).
				Value(&selectedContainers),
		),
	)

	err = form.Run()
	if err != nil {
		slog.Debug("error creating a form", "error", err.Error())
		return err
	}

	for i := range selectedContainers {
		err = dc.client.ContainerRemove(context.Background(), selectedContainers[i].ID, container.RemoveOptions{
			Force: forceRemoval,
		})
		if err != nil {
			fmt.Printf("\x1b[31mx\x1b[0m Error Removing Container %s\n", selectedContainers[i].Name)
			slog.Debug("Error Removing Container", "Container ID", selectedContainers[i].ID, "Error", err.Error())
		} else {
			fmt.Printf("\x1b[32m✓\x1b[0m Container %s Removed\n", selectedContainers[i].Name)
			slog.Debug("Removed Container", "Container ID", selectedContainers[i].ID)
		}
	}

	return nil
}

func (dc *DockerClient) RemoveDockerImages(ctx context.Context, cmd *cli.Command) error {

	if cmd.Bool("verbose") {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	forceRemoval := cmd.Bool("force")

	images, err := dc.client.ImageList(context.Background(), image.ListOptions{All: true})
	if err != nil {
		slog.Debug("Error fetching image list", "error", err.Error())
		return err
	}

	if len(images) == 0 {
		fmt.Println("No Docker images found.")
		return nil
	}

	var imagesList []*ImageItem

	maxNameLen := 0

	for i := range images {
		if forceRemoval || images[i].Containers <= 0 {
			imagesList = append(imagesList, &ImageItem{
				ID:   images[i].ID,
				Name: images[i].RepoTags[0],
				Size: images[i].Size,
			})
			maxNameLen = max(maxNameLen, len(images[i].RepoTags[0]))
		}
	}

	if len(imagesList) == 0 {
		fmt.Println("No orphaned images found")
		return nil
	}

	var selectedImages []*ImageItem

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[*ImageItem]().
				Title("Select the images to REMOVE").
				Options(ParseImageItemList(imagesList, maxNameLen)...).
				Value(&selectedImages),
		),
	)

	err = form.Run()
	if err != nil {
		slog.Debug("error creating a form", "error", err.Error())
		return err
	}

	for i := range selectedImages {
		_, err = dc.client.ImageRemove(context.Background(), selectedImages[i].ID, image.RemoveOptions{
			Force: forceRemoval,
		})
		if err != nil {
			fmt.Printf("\x1b[31mx\x1b[0m Error Removing Image %s\n", selectedImages[i].Name)
			slog.Debug("Error Removing Image", "Image ID", selectedImages[i].ID, "Error", err.Error())
		} else {
			fmt.Printf("\x1b[32m✓\x1b[0m Image %s Removed\n", selectedImages[i].Name)
			slog.Debug("Removed Image", "Image ID", selectedImages[i].ID)
		}
	}

	return nil
}

func ParseContainerItemList(inputList []*ContainerItem, maxNameLen int) []huh.Option[*ContainerItem] {
	var result []huh.Option[*ContainerItem]
	for i := range inputList {
		key := fmt.Sprintf("%-*.*s | %s", maxNameLen, maxNameLen, inputList[i].Name, inputList[i].Image)
		result = append(result, huh.Option[*ContainerItem]{
			Key:   key,
			Value: inputList[i],
		})
	}
	return result
}

func ParseImageItemList(inputList []*ImageItem, maxNameLen int) []huh.Option[*ImageItem] {
	var result []huh.Option[*ImageItem]
	for i := range inputList {
		key := fmt.Sprintf("%-*.*s | %.2f MB", maxNameLen, maxNameLen, inputList[i].Name, float64(inputList[i].Size)/(1024*1024))
		result = append(result, huh.Option[*ImageItem]{
			Key:   key,
			Value: inputList[i],
		})
	}
	return result
}
