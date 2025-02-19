package ide

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// ListShims returns an array of shims which are added to this ide project
func (project *Project) ListShims() map[string]string {
	shims := map[string]string{}

	for _, option := range project.config.Raw.Section("ide").Subsection("shims").Options {
		option.Key = strings.Replace(option.Key, "-----", ".", -1)
		shims[option.Key] = option.Value
	}

	return shims
}

// GetShim returns the command for a shim
func (project *Project) GetShim(shim string) string {
	shim = strings.Replace(shim, ".", "-----", -1)
	return project.config.Raw.Section("ide").Subsection("shims").Option(shim)
}

// RefreshShims syncs the shims from .git/config with .git/bin
func (project *Project) RefreshShims() error {
	shims := project.ListShims()

	if len(shims) == 0 {
		return nil
	}

	if _, err := os.Stat(".git/bin"); os.IsNotExist(err) {
		if err := os.Mkdir(".git/bin", 0755); err != nil {
			return err
		}
	}

	source, _ := os.Executable()

	for shim := range shims {
		dest := filepath.Join(project.location, ".git", "bin", shim)
		if _, err := os.Lstat(dest); err == nil {
			if err := os.Remove(dest); err != nil {
				return err
			}
		}
		if err := os.Symlink(source, dest); err != nil {
			return err
		}
	}

	return nil
}

// AddShim adds a shim to this project
func (project *Project) AddShim(shim string, command string) error {
	dest := filepath.Join(project.location, ".git", "bin", shim)

	if _, err := os.Stat(dest); err == nil {
		return errors.New("shim " + shim + " already exists for this project")
	}

	if _, err := os.Stat(".git/bin"); os.IsNotExist(err) {
		if err := os.Mkdir(".git/bin", 0755); err != nil {
			return err
		}
	}

	source, _ := os.Executable()

	if err := os.Symlink(source, dest); err != nil {
		return err
	}

	project.config.Raw.SetOption("ide", "shims", strings.Replace(shim, ".", "-----", -1), command)

	return project.repository.Storer.SetConfig(project.config)
}

// RemoveShim removes a shim from this project
func (project *Project) RemoveShim(shim string) error {
	dest := filepath.Join(project.location, ".git", "bin", shim)

	if _, err := os.Lstat(dest); err == nil {
		if err := os.Remove(dest); err != nil {
			return err
		}
	}

	project.config.Raw.Section("ide").Subsection("shims").RemoveOption(strings.Replace(shim, ".", "-----", -1))

	return project.repository.Storer.SetConfig(project.config)
}
