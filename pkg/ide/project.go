package ide

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/libgit2/git2go.v25"
)

// Project represents an ide project
type Project struct {
	repository     *git.Repository
	config         *git.Config
	ctrlpCacheFile string
	ctagsFile      string
}

// LoadProject instantiates a new instance of Project for a given directory
func LoadProject(directory string) (Project, error) {
	if directory == "" {
		return Project{}, errors.New("Project directory cannot be empty")
	}

	repo, err := git.OpenRepository(directory)
	if err != nil {
		return Project{}, err
	}

	config, err := repo.Config()
	if err != nil {
		return Project{}, err
	}

	return Project{
		repository: repo,
		config:     config,
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

	//find the branch name
	branch := ""
	branchElements := strings.Split(head.Name(), "/")
	if len(branchElements) == 3 {
		branch = branchElements[2]
	}

	return branch
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
	language, _ := project.config.LookupString("ide.language")

	return language
}

// Location returns the absolute file path of the ide project
func (project *Project) Location() string {
	return project.repository.Workdir()
}

// SetLanguage stores the given language in the .git/config file of the ide project
func (project *Project) SetLanguage(language string) error {
	err := project.config.SetString("ide.language", language)
	if err != nil {
		return err
	}

	return nil
}

// Destroy removes any trace of ide configuration from .git/config file
func (project *Project) Destroy() error {
	err := project.config.Delete("ide.language")
	if err != nil {
		return err
	}

	return nil
}
