package ide

import (
	"errors"
	"os"
	"path/filepath"
)

// ListBinaries returns an array of binaries which are added to this ide project
func (project *Project) ListBinaries() map[string]string {
	binaries := map[string]string{}

	for _, option := range project.config.Raw.Section("ide").Subsection("binaries").Options {
		binaries[option.Key] = option.Value
	}

	return binaries
}

// GetBinary returns the command for a binary
func (project *Project) GetBinary(binary string) string {
	return project.config.Raw.Section("ide").Subsection("binaries").Option(binary)
}

// RefreshBinaries syncs the binaries from .git/config with .git/bin
func (project *Project) RefreshBinaries() error {
	binaries := project.ListBinaries()

	if len(binaries) == 0 {
		return nil
	}

	if _, err := os.Stat(".git/bin"); os.IsNotExist(err) {
		if err := os.Mkdir(".git/bin", 0755); err != nil {
			return err
		}
	}

	source, _ := os.Executable()

	for binary := range binaries {
		dest := filepath.Join(project.location, ".git", "bin", binary)
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

// AddBinary adds a binary to this project
func (project *Project) AddBinary(binary string, command string) error {
	dest := filepath.Join(project.location, ".git", "bin", binary)

	if _, err := os.Stat(dest); err == nil {
		return errors.New("Binary " + binary + " already exists for this project")
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

	project.config.Raw.SetOption("ide", "binaries", binary, command)

	return project.repository.Storer.SetConfig(project.config)
}

// RemoveBinary removes a binary from this project
func (project *Project) RemoveBinary(binary string) error {
	dest := filepath.Join(project.location, ".git", "bin", binary)

	if _, err := os.Lstat(dest); err == nil {
		if err := os.Remove(dest); err != nil {
			return err
		}
	}

	project.config.Raw.Section("ide").Subsection("binaries").RemoveOption(binary)

	return project.repository.Storer.SetConfig(project.config)
}
