package main

import (
	"context"
	"fmt"

	"github.com/ErickMaria/envcontainer/internal/runtime/docker"
	"github.com/ErickMaria/envcontainer/internal/runtime/docker/types"
	"github.com/ErickMaria/envcontainer/internal/template"
	"github.com/ErickMaria/envcontainer/pkg/cli"
)

var cmd *cli.Command
var cmds cli.CommandConfig

func init() {

	// # TEMPLATE FILE
	err := template.Initialization()
	if err != nil {
		panic(err)
	}

	configFile, err := template.Unmarshal()
	if err != nil {
		panic(err)
	}

	// # DOCKER API
	ctx := context.Background()
	container := docker.NewDocker()

	// CLI
	cmd, cmds = cli.NewCommand(cli.CommandConfig{
		"build": cli.Command{
			Desc: "build a image using envcontainer configuration in the current directory",
			Exec: func() {

				err := container.Build(ctx, types.BuildOptions{
					ImageName:  configFile.Project.Name,
					Dockerfile: configFile.Container.Build,
				})
				if err != nil {
					panic(err)
				}
			},
		},
		"start": cli.Command{
			Flags: cli.Flag{
				Values: map[string]cli.Values{
					"auto-stop": {
						Defaulvalue: "true",
						Description: "terminal shell that must be used",
					},
				},
			},
			Desc: "run the envcontainer configuration to start the container and link it to the current directory",
			Exec: func() {

				autoStop := *cmd.Flags.Values["auto-stop"].ValueBool
				err := container.Start(ctx, autoStop)
				if err != nil {
					panic(err)
				}
			},
		},
		"stop": cli.Command{
			Desc: "stop all envcontainer configuration running in the current directory",
			Exec: func() {
				err := container.Stop(ctx)
				if err != nil {
					panic(err)
				}
			},
		},
		"run": cli.Command{
			Flags: cli.Flag{
				Values: map[string]cli.Values{
					"name": {
						Description: "container name",
					},
					"image": {
						Description: "envcontainer image",
					},
					"command": {
						Description: "command to run inside container",
					},
					"pull-image-always": {
						Defaulvalue: "false",
						Description: "pull image always before run container (DISABLED)",
					},
				},
			},
			Exec: func() {

				image := *cmd.Flags.Values["image"].ValueString
				pullImageAlways := *cmd.Flags.Values["pull-image-always"].ValueBool
				command := *cmd.Flags.Values["command"].ValueString

				err := container.Run(ctx, types.RunOptions{
					ImageName:       image,
					Command:         command,
					PullImageAlways: pullImageAlways,
				})
				if err != nil {
					panic(err)
				}

			},
			Desc: "execute an .envcontainer on the current directory without saving it locally",
		},
		"version": cli.Command{
			Exec: func() {
				fmt.Println("Version: 0.5.0")
			},
			Desc: "show envcontainer version",
		},
		"help": cli.Command{
			Exec: func() {
				cli.Help(cmds)
			},
			Desc: "Run " + cli.ExecutableName() + " COMMAND' for more information on a command. See: '" + cli.ExecutableName() + " help'",
		},
	})

}

func main() {
	cmd.Listener()
}
