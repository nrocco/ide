package ide

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// CtagsFile returns the path to the ctags file of the current ide project
func (project *Project) CtagsFile() string {
	return filepath.Join(project.location, ".git", "tags")
}

// CtagsOptions returns project specific ctags flags from .git/config file
func (project *Project) CtagsOptions() []string {
	return strings.Fields(project.config.Raw.Section("ide").Option("ctags"))
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
	tmpCtagsFile := project.CtagsFile() + ".new"

	os.Remove(tmpCtagsFile)

	options := append([]string{
		"--tag-relative=yes", "--sort=yes", "--totals=yes",
		"--exclude=.git", "--exclude=.hg", "--exclude=.svn",
		"--recurse", "-f", tmpCtagsFile,
	}, project.CtagsOptions()...)

	cmd := exec.Command("ctags", options...)
	cmd.Dir = project.Location()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	if err := os.Rename(tmpCtagsFile, project.CtagsFile()); err != nil {
		return err
	}

	return nil
}
