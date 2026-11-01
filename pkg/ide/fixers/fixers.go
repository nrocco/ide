package fixers

import (
	"os"
	"os/exec"
)

// FixClrf TODO
func FixClrf(path string) error {
	return exec.Command("dos2unix", "--verbose", "--safe", path).Run() // #nosec G204 -- fixed binary, arguments are passed without a shell
}

// FixWhitespace TODO
func FixWhitespace(path string) error {
	return exec.Command("sed", "-i", "", "-e", "s/[ \t]*$//", path).Run() // #nosec G204 -- fixed binary, arguments are passed without a shell
}

// FixPhpcsfixer TODO
func FixPhpcsfixer(path string) error {
	return exec.Command("bin/php-cs-fixer", "fix", "--no-ansi", "--no-interaction", path).Run() // #nosec G204 -- fixed binary, arguments are passed without a shell
}

// FixJq TODO
func FixJq(path string) error {
	file, err := os.Create(path + ".tmp") // #nosec G304 -- path is supplied by the user on the command line
	if err != nil {
		return err
	}
	defer file.Close()
	cmd := exec.Command("jq", "--sort-keys", "--monochrome-output", ".", path) // #nosec G204 -- fixed binary, arguments are passed without a shell
	cmd.Stdout = file
	if err := cmd.Run(); err != nil {
		_ = os.Remove(path + ".tmp")
		return err
	}
	return os.Rename(path+".tmp", path)
}

// FixGoimports TODO
func FixGoimports(path string) error {
	return exec.Command("goimports", "-w", path).Run() // #nosec G204 -- fixed binary, arguments are passed without a shell
}

// FixGofmt TODO
func FixGofmt(path string) error {
	return exec.Command("gofmt", "-w", "-s", path).Run() // #nosec G204 -- fixed binary, arguments are passed without a shell
}

// FixEslint TODO
func FixEslint(path string) error {
	return exec.Command("eslint", "--fix", path).Run() // #nosec G204 -- fixed binary, arguments are passed without a shell
}
