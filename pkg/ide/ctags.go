package ide

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

func (project *Project) GetCtagsFile() string {
	if project.ctagsFile == "" {
		project.ctagsFile = filepath.Join(project.repository.Path(), "tags")
	}

	return project.ctagsFile
}

func (project *Project) RefreshCtags() error {
	if !project.IsConfigured() {
		return errors.New("Project must be configured before you can RefreshCtags")
	}

	options := []string{
		"--recurse", "--tag-relative=yes", "--sort=yes",
		"--exclude=.git", "--exclude=.hg", "--exclude=.svn",
		"-f", project.GetCtagsFile(),
		"--kinds-php=cif", "--kinds-python=-i",
		"--languages=" + project.Language(),
	}

	cmd := exec.Command("ctags", options...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	return nil
}
