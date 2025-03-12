package ide

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ShimList returns an array of shims which are added to this ide project
func (project *Project) ShimList() map[string]string {
	shims := map[string]string{}

	for _, option := range project.config.Raw.Section("ide").Subsection("shims").Options {
		option.Key = strings.ReplaceAll(option.Key, "-----", ".")
		shims[option.Key] = option.Value
	}

	return shims
}

// ShimGet returns the command for a shim
func (project *Project) ShimGet(shim string) string {
	shim = strings.ReplaceAll(shim, ".", "-----")
	return project.config.Raw.Section("ide").Subsection("shims").Option(shim)
}

// ShimRefresh syncs the shims from .git/config with .git/bin
func (project *Project) ShimRefresh() error {
	shims := project.ShimList()

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

// ShimAdd adds a shim to this project
func (project *Project) ShimAdd(shim string, command string) error {
	validShimRegexp := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9._-]*$")

	if !validShimRegexp.MatchString(shim) {
		return errors.New("not a valid shim name: " + shim)
	}

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

	project.config.Raw.SetOption("ide", "shims", strings.ReplaceAll(shim, ".", "-----"), command)

	return project.repository.Storer.SetConfig(project.config)
}

// ShimRemove removes a shim from this project
func (project *Project) ShimRemove(shim string) error {
	dest := filepath.Join(project.location, ".git", "bin", shim)

	if _, err := os.Lstat(dest); err == nil {
		if err := os.Remove(dest); err != nil {
			return err
		}
	}

	project.config.Raw.Section("ide").Subsection("shims").RemoveOption(strings.ReplaceAll(shim, ".", "-----"))

	return project.repository.Storer.SetConfig(project.config)
}
