package ide

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// CtagsFile returns the path to the ctags file of the current ide project
func (project *Project) CtagsFile() string {
	return filepath.Join(project.location, ".git", "tags") // TODO: make location of the tags file configurable
}

// CtagsFileAge returns the time the ctags file was last modified
func (project *Project) CtagsFileAge() time.Time {
	file, err := os.Open(project.CtagsFile())
	if err != nil {
		return time.Time{}
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return time.Time{}
	}

	return stat.ModTime()
}

// CtagsFileSize returns the size of the ctags file in bytes
func (project *Project) CtagsFileSize() uint64 {
	file, err := os.Open(project.CtagsFile())
	if err != nil {
		return 0
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return 0
	}

	return uint64(stat.Size())
}

// RefreshCtags generates a new ctags file for the current project
func (project *Project) RefreshCtags() error {
	if !project.IsConfigured() {
		return errors.New("Project must be configured before you can RefreshCtags")
	}

	// TODO: make default options for ctags configurable
	options := []string{
		"--tag-relative=yes", "--sort=yes", "--totals=yes",
		"--exclude=.git", "--exclude=.hg", "--exclude=.svn",
		"--recurse", "-f", project.CtagsFile(),
		"--kinds-php=cif", "--kinds-python=-i",
		"--languages=" + project.Language(),
	}

	cmd := exec.Command("ctags", options...)
	cmd.Dir = project.Location()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	return nil
}
