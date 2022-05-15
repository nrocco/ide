package ide

import (
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

// Project represents an ide project
type Project struct {
	repository *git.Repository
	config     *config.Config
	location   string
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

// IsConfigured returns true if the current git repository is setup as an ide project
func (project *Project) IsConfigured() bool {
	return project.config.Raw.HasSection("ide")
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
		project.DisableHook(hook)
	}

	for binary := range project.ListBinaries() {
		project.DisableBinary(binary)
	}

	project.config.Raw.RemoveSection("ide")

	return project.repository.Storer.SetConfig(project.config)
}
