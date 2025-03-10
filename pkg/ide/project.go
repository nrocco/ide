package ide

import (
	"path/filepath"

	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/docker/client"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

// Project represents an ide project
type Project struct {
	repository    *git.Repository
	config        *config.Config
	dockerClient  *client.Client
	composeClient api.Service
	location      string
}

// NewProject instantiates a new instance of Project for a given directory
func NewProject(location string) (*Project, error) {
	repository, openErr := git.PlainOpenWithOptions(location, &git.PlainOpenOptions{DetectDotGit: true})
	if openErr != nil {
		return &Project{}, openErr
	}

	config, configErr := repository.Config()
	if configErr != nil {
		return &Project{}, configErr
	}

	workTree, workTreeErr := repository.Worktree()
	if configErr != nil {
		return &Project{}, workTreeErr
	}

	return &Project{
		repository: repository,
		config:     config,
		location:   workTree.Filesystem.Root(),
	}, nil
}

// Name returns the name of the ide project, extracted from the parent directory name
func (project *Project) Name() string {
	return filepath.Base(project.Location())
}

// Branch returns the currently checked out branch of the ide project
func (project *Project) Branch() string {
	head, headErr := project.repository.Head()
	if headErr != nil {
		return ""
	}

	if head == nil {
		return ""
	}

	return head.Name().Short()
}

// Email returns the email of the user of the ide project
func (project *Project) Email() string {
	return project.config.Raw.Section("user").Option("email")
}

// Location returns the absolute file path of the ide project
func (project *Project) Location() string {
	return project.location
}

// Destroy removes any trace of ide configuration from .git/config file
func (project *Project) Destroy() error {
	for _, hook := range project.ListHooks() {
		if err := project.RemoveHook(hook); err != nil {
			return err
		}
	}

	for shim := range project.ListShims() {
		if err := project.RemoveShim(shim); err != nil {
			return err
		}
	}

	project.config.Raw.RemoveSection("ide")

	return project.repository.Storer.SetConfig(project.config)
}
