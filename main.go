package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {

	dockerCLI := NewDockerClient()
	if dockerCLI == nil {
		fmt.Println("Check if docker is installed in your system")
		return
	}

	cmd := &cli.Command{
		Name:    "klir",
		Version: "v1.0.0",
		Usage:   "Clean up Docker container & images with ease",
		Commands: []*cli.Command{
			{
				Name:  "stop",
				Usage: "stop running docker containers",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "verbose",
						Aliases: []string{"v"},
						Usage:   "enable verbose output",
					},
				},
				Action: dockerCLI.StopRunningContainers,
			},
			{
				Name:  "rm",
				Usage: "delete docker containers",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "verbose",
						Aliases: []string{"v"},
						Usage:   "enable verbose output",
					},
				},
				Action: dockerCLI.RemoveDockerContainers,
			},
			{
				Name:  "rmi",
				Usage: "delete docker images",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "verbose",
						Aliases: []string{"v"},
						Usage:   "enable verbose output",
					},
				},
				Action: dockerCLI.RemoveDockerImages,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
