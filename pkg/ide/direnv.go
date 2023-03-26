package ide

import (
	"bufio"
	"path/filepath"
	"os"
	"strings"
)

// AddGitBinToPath updates your local .envrc and adds .git/bin to the $PATH
func (project *Project) AddGitBinToPath() error {
	file, err := os.OpenFile(".envrc", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "PATH_add .git/bin") {
			return nil
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if _, err := file.Write([]byte("PATH_add .git/bin\n")); err != nil {
		return err
	}

	return nil
}

// HasGitBinInPath checks if $PATH contains the current .git/bin directory
func (project *Project) HasGitBinInPath() bool {
	return strings.Contains(os.Getenv("PATH"), filepath.Join(project.location, ".git", "bin"))
}

// HasDirEnv checks if the current project has a .envrc file
func (project *Project) HasDirEnv() bool {
	_, err := os.Stat(".envrc")
	return !os.IsNotExist(err)
}
