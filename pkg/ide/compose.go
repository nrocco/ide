package ide

import (
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/docker/docker/client"
)

// DockerClient returns an instance of a docker client to interface with docker
func (project *Project) DockerClient() (*client.Client, error) {
	if project.dockerClient == nil {
		dockerClient, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return nil, err
		}
		project.dockerClient = dockerClient
	}
	return project.dockerClient, nil
}

// DockerComposeClient returns an instance of the compose service to interface with docker compose
func (project *Project) DockerComposeClient() (api.Service, error) {
	if project.composeClient == nil {
		client, err := project.DockerClient()
		if err != nil {
			return nil, err
		}
		cli, err := command.NewDockerCli(command.WithAPIClient(client))
		if err != nil {
			return nil, err
		}
		if err := cli.Initialize(flags.NewClientOptions()); err != nil {
			return nil, err
		}
		backend := compose.NewComposeService(cli)
		if err != nil {
			return nil, err
		}
		project.composeClient = backend
	}
	return project.composeClient, nil
}
