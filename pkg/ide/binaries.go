package ide

import (
	"errors"
	"os"
	"path/filepath"
)

// ListBinaries returns an array of binaries which are enabled for the ide project
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

// EnableBinary adds a binary to this project
func (project *Project) EnableBinary(binary string, command string) error {
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

// DisableBinary removes a binary from this project
func (project *Project) DisableBinary(binary string) error {
	dest := filepath.Join(project.location, ".git", "bin", binary)

	if _, err := os.Stat(dest); err == nil {
		if err := os.Remove(dest); err != nil {
			return err
		}
	}

	project.config.Raw.Section("ide").Subsection("binaries").RemoveOption(binary)

	return project.repository.Storer.SetConfig(project.config)
}
