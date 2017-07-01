package ide

import (
	"errors"
	"os"
	"path/filepath"
)

// ListHooks returns an array of hooks which are enabled for the ide project
func (project *Project) ListHooks() []string {
	hooks := []string{}

	filepath.Walk(filepath.Join(project.repository.Path(), "hooks"), func(path string, f os.FileInfo, err error) error {
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
		"post-checkout",
		"prepare-commit-msg":
		return true
	}
	return false
}

// EnableHook enables a git repository hook
func (project *Project) EnableHook(hook string) error {
	if !project.isValidHook(hook) {
		return errors.New(hook + " is not a valid hook")
	}

	dest := filepath.Join(project.repository.Path(), "hooks", hook)

	if _, err := os.Stat(dest); err == nil {
		return errors.New("Hook " + hook + " already exists for this project\n")
	}

	source, _ := os.Executable()

	return os.Symlink(source, dest)
}

// DisableHook disables a git repository hook
func (project *Project) DisableHook(hook string) error {
	if !project.isValidHook(hook) {
		return errors.New(hook + " is not a valid hook")
	}

	dest := filepath.Join(project.repository.Path(), "hooks", hook)

	if _, err := os.Stat(dest); err == nil {
		if err := os.Remove(dest); err != nil {
			return err
		}
	}

	return nil
}
