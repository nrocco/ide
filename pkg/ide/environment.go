package ide

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"os/user"
)

type Environment struct {
	Project Project
	Image   string
	Command string
}

func (env *Environment) Create() error {
	ctx := context.Background()
	client, clientErr := client.NewEnvClient()
	if clientErr != nil {
		return clientErr
	}

	user, userErr := user.Current()
	if userErr != nil {
		return userErr
	}

	// TODO --init option?
	// TODO --network host option?
	config := container.Config{
		Image:     env.Image,
		Cmd:       []string{"cat"},
		User:      fmt.Sprintf("%s:%s", user.Uid, user.Gid),
		OpenStdin: true,
		// Volumes: map[string]struct{}, --volume '%s:%s'
		WorkingDir: env.Project.Location(),
		// Labels: map[string]string, --label ide.project=Name
	}

	_, err := client.ContainerCreate(ctx, &config, nil, nil, "")
	if err != nil {
		return err
	}

	return nil
}
