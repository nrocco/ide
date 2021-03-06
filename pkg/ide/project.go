package ide

import (
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	homedir "github.com/mitchellh/go-homedir"
)

// Project represents an ide project
type Project struct {
	repository *git.Repository
	config     *config.Config
	location   string
}

// LoadProject instantiates a new instance of Project for a given directory
func LoadProject(location string) (*Project, error) {
	location, _ = homedir.Expand(location)
	location, _ = filepath.Abs(location)

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
	return project.Language() != ""
}

// AutoDetectLanguage guesses the language of a git repository based on some
// standard files and defaults to plain
func (project *Project) AutoDetectLanguage() string {
	if _, err := os.Stat("setup.py"); err == nil {
		return "python"
	} else if _, err := os.Stat("composer.json"); err == nil {
		return "php"
	} else if _, err := os.Stat("package.json"); err == nil {
		return "javascript"
	} else if _, err := os.Stat("main.go"); err == nil {
		return "go"
	} else if _, err := os.Stat("go.mod"); err == nil {
		return "go"
	}

	return "plain"
}

// Language returns the language of the ide project as stored in .git/config file
func (project *Project) Language() string {
	return project.config.Raw.Section("ide").Option("language")
}

// Email returns the email of the user of the ide project
func (project *Project) Email() string {
	return project.config.Raw.Section("user").Option("email")
}

// Location returns the absolute file path of the ide project
func (project *Project) Location() string {
	return project.location
}

// SetLanguage stores the given language in the .git/config file of the ide project
func (project *Project) SetLanguage(language string) error {
	project.config.Raw.Section("ide").SetOption("language", language)

	return project.repository.Storer.SetConfig(project.config)
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
