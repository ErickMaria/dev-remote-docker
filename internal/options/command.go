package options

import (
	"io/ioutil"
	"os"

	"github.com/ErickMaria/envcontainer/internal/config"
)

const (
	PATH_DEFAULT   = ".envcontainer/compose/env"
	DOCKERFILE     = ".envcontainer/Dockerfile"
	DOCKER_COMPOSE = ".envcontainer/compose/docker-compose.yaml"
	ENV            = ".envcontainer/compose/env/.variables"
)

const (
	INIT string = "init"
	RUN  string = "run"
	STOP string = "stop"
	HELP string = "help"
)

type Command struct {
	Flags Flag
	Exec  func()
}

func Init(flags Flag) {

	values := flags.Values

	var ports = []string{}

	if *values["listener"].value != "" {
		ports = append(ports, *values["listener"].value)
	}

	dc := config.DockerCompose{
		Version: "3.6",
		Services: config.Services{
			Environment: config.Environment{
				ContainerName: *values["project"].value,
				Volumes: []config.Volumes{
					config.Volumes{
						Type:   "bind",
						Source: "../../",
						Target: "/home/envcontainer/" + *values["project"].value,
					},
				},
				Build: config.Build{
					Dockerfile: "Dockerfile",
					Context:    "../",
				},
				Ports:      ports,
				WorkingDir: "/home/envcontainer/" + *values["project"].value,
				EnvFile: []string{
					*values["envfile"].value,
				},
				Privileged: true,
				StdinOpen:  true,
				Tty:        true,
			},
		},
	}

	data := dc.Marshal()

	err := os.MkdirAll(PATH_DEFAULT, 0755)
	check(err)

	createFile := func(name string, data []byte) {
		check(ioutil.WriteFile(name, data, 0644))
	}

	createFile(DOCKERFILE, []byte("FROM "+*values["image"].value))
	createFile(DOCKER_COMPOSE, data)
	createFile(ENV, []byte(""))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
