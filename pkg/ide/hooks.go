package ide

import (
	"errors"
	"os"
	"path/filepath"
)

// ListHooks returns an array of hooks which are added to this ide project
func (project *Project) ListHooks() []string {
	hooks := []string{}

	filepath.Walk(filepath.Join(project.location, ".git", "hooks"), func(path string, f os.FileInfo, err error) error {
		if f.Mode()&os.ModeSymlink != 0 {
			hooks = append(hooks, f.Name())
		}
		return nil
	})

	return hooks
}

func (project *Project) isValidHook(hook string) bool {
	switch hook {
	case
		"commit-msg",
		"prepare-commit-msg":
		return true
	}
	return false
}

// AddHook adds a git repository hook
func (project *Project) AddHook(hook string) error {
	if !project.isValidHook(hook) {
		return errors.New(hook + " is not a valid hook")
	}

	dest := filepath.Join(project.location, ".git", "hooks", hook)

	if _, err := os.Stat(dest); err == nil {
		return errors.New("hook " + hook + " already exists for this project")
	}

	source, _ := os.Executable()

	return os.Symlink(source, dest)
}

// RemoveHook removes a git repository hook
func (project *Project) RemoveHook(hook string) error {
	if !project.isValidHook(hook) {
		return errors.New(hook + " is not a valid hook")
	}

	dest := filepath.Join(project.location, ".git", "hooks", hook)

	if _, err := os.Stat(dest); err == nil {
		if err := os.Remove(dest); err != nil {
			return err
		}
	}

	return nil
}
