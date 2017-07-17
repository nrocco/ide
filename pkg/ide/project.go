package ide

import (
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
)

// Project represents an ide project
type Project struct {
	repository     *git.Repository
	config         *config.Config
	location       string
	ctrlpCacheFile string
	ctagsFile      string
}

// LoadProject instantiates a new instance of Project for a given directory
func LoadProject(directory string) (Project, error) {
	directory, _ = homedir.Expand(directory)
	directory, _ = filepath.Abs(directory)

	repo, openErr := git.PlainOpen(directory)
	if openErr != nil {
		return Project{}, openErr
	}

	config, configErr := repo.Config()
	if configErr != nil {
		return Project{}, configErr
	}

	return Project{
		repository: repo,
		config:     config,
		location:   directory,
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
	} else if _, err := os.Stat("manage.py"); err == nil {
		return "python"
	} else if _, err := os.Stat("main.go"); err == nil {
		return "go"
	}

	return "plain"
}

// Language returns the language of the ide project as stored in .git/config file
func (project *Project) Language() string {
	return project.config.Raw.Section("ide").Option("language")
}

// Location returns the absolute file path of the ide project
func (project *Project) Location() string {
	return project.location
}

// SetLanguage stores the given language in the .git/config file of the ide project
func (project *Project) SetLanguage(language string) error {
	project.config.Raw.Section("ide").SetOption("language", language)
	project.repository.Storer.SetConfig(project.config)

	return nil
}

func (project *Project) GetExecutable(executable string) string {
	return project.config.Raw.Section("ide").Subsection("executables").Option(executable)
}

func (project *Project) AddExecutable(executable string, container string) error {
	project.config.Raw.Section("ide").Subsection("executables").AddOption(executable, container)
	project.repository.Storer.SetConfig(project.config)

	return nil
}

func (project *Project) RemoveExecutable(executable string) error {
	project.config.Raw.Section("ide").Subsection("executables").RemoveOption(executable)
	project.repository.Storer.SetConfig(project.config)

	return nil
}

// Destroy removes any trace of ide configuration from .git/config file
func (project *Project) Destroy() error {
	for _, hook := range project.ListHooks() {
		project.DisableHook(hook)
	}

	project.config.Raw.RemoveSection("ide")
	project.repository.Storer.SetConfig(project.config)

	return nil
}
