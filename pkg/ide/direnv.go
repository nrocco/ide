package ide

import (
	"bufio"
	"os"
)

// DirEnvAddLayoutIde updates your local .envrc and adds .git/bin to the $PATH
func (project *Project) DirEnvAddLayoutIde() error {
	file, err := os.OpenFile(".envrc", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == "layout_ide" {
			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if _, err := file.Write([]byte("layout_ide\n")); err != nil {
		return err
	}

	return nil
}

// DirEnvHasLayoutIde checks if .envrc contains layout_ide
func (project *Project) DirEnvHasLayoutIde() bool {
	file, err := os.Open(".envrc")
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == "layout_ide" {
			return true
		}
	}

	return false
}

// DirEnvExists checks if the current project has a .envrc file
func (project *Project) DirEnvExists() bool {
	_, err := os.Stat(".envrc")
	return !os.IsNotExist(err)
}
