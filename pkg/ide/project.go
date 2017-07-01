package ide

import (
	"errors"
	"os"
	"os/exec"

	"gopkg.in/libgit2/git2go.v25"
)

type Project struct {
	Repository *git.Repository
	Config     *git.Config
	Location   string
	Language   string
}

func LoadProject(gitdir string) (Project, error) {
	if gitdir == "" {
		return Project{}, errors.New("Project directory cannot be empty")
	}

	repo, err := git.OpenRepository(gitdir)
	if err != nil {
		return Project{}, err
	}

	config, err := repo.Config()
	if err != nil {
		return Project{}, err
	}

	language, _ := config.LookupString("ide.language")

	return Project{
		Repository: repo,
		Config:     config,
		Language:   language,
		Location:   repo.Workdir(),
	}, nil
}

func (this *Project) IsConfigured() bool {
	return this.Language != ""
}

func (this *Project) AutoDetectLanguage() string {
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

func (this *Project) ListHooks() error {
	cmd := exec.Command("find", ".git/hooks", "-perm", "+0111", "-a", "-not", "-type", "d", "-a", "-not", "-name", "\\*.sample")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	return nil
}

func (this *Project) Save() error {
	err := this.Config.SetString("ide.language", this.Language)
	if err != nil {
		return err
	}

	return nil
}

func (this *Project) Destroy() error {
	err := this.Config.Delete("ide.language")
	if err != nil {
		return err
	}

	return nil
}
