package tools

import (
	"os"
	"os/exec"
)

// FixClrf TODO
func FixClrf(path string) error {
	return exec.Command("dos2unix", "--verbose", "--safe", path).Run()
}

// FixWhitespace TODO
func FixWhitespace(path string) error {
	return exec.Command("sed", "-i", "", "-e", "s/[ \t]*$//", path).Run()
}

// FixPhpcsfixer TODO
func FixPhpcsfixer(path string) error {
	return exec.Command("bin/php-cs-fixer", "fix", "--no-ansi", "--no-interaction", path).Run()
}

// FixJq TODO
func FixJq(path string) error {
	file, err := os.Create(path + ".tmp")
	if err != nil {
		return err
	}
	defer file.Close()
	cmd := exec.Command("jq", "--sort-keys", "--monochrome-output", ".", path)
	cmd.Stdout = file
	if err := cmd.Run(); err != nil {
		return err
	}
	return os.Rename(path+".tmp", path)
}

// FixGoimports TODO
func FixGoimports(path string) error {
	return exec.Command("goimports", "-w", path).Run()
}

// FixGofmt TODO
func FixGofmt(path string) error {
	return exec.Command("gofmt", "-w", "-s", path).Run()
}
