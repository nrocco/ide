package cmd

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/anmitsu/go-shlex"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

type binaryContext struct {
	User    string
	Group   string
	UID     int
	GID     int
	Project string
	Dir     string
}

func (b binaryContext) RelDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	rel, err := filepath.Rel(b.Dir, dir)
	if err != nil {
		return ""
	}
	return rel
}

func (b binaryContext) DockerRunOrExec(container string) string {
	fafa, _ := exec.Command("docker", "container", "ps", "--quiet", "--filter", "status=running", container).Output()
	if len(fafa) == 0 {
		return "docker container run --rm"
	}
	return "docker container exec"
}

func (b binaryContext) DockerComposeRunOrExec(service string) string {
	fafa, _ := exec.Command("docker", "compose", "ps", "--quiet", "--filter", "status=running", service).Output()
	if len(fafa) == 0 {
		return "docker compose run --rm"
	}
	return "docker compose exec"
}

func (b binaryContext) IsTTY() bool {
	if !isatty.IsTerminal(os.Stdin.Fd()) {
		return false
	}
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return false
	}
	if !isatty.IsTerminal(os.Stderr.Fd()) {
		return false
	}
	return true
}

var runBinaryCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a binary in the context of an ide project",
	Long:  "Run a binary in the context of an ide project",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		binaries := []string{}
		if len(args) == 0 {
			if err := loadProject(cmd, args); err == nil {
				for binary := range project.ListBinaries() {
					binaries = append(binaries, binary)
				}
			}
		}
		return binaries, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		command := project.GetBinary(args[0])
		if command == "" {
			return errors.New("Binary does not exist. Did you forget to add it?")
		}

		tmpl, err := template.New("command").Parse(command)
		if err != nil {
			return err
		}

		var b bytes.Buffer

		err = tmpl.Execute(&b, binaryContext{
			UID:     os.Getuid(),
			GID:     os.Getgid(),
			Project: project.Name(),
			Dir:     project.Location(),
		})
		if err != nil {
			return err
		}

		command = os.ExpandEnv(b.String())

		parts, err := shlex.Split(command, true)
		if err != nil {
			return err
		}

		parts = append(parts, args[1:]...)

		fafa := exec.Command(parts[0], parts[1:]...)
		fafa.Stdin = os.Stdin
		fafa.Stdout = os.Stdout
		fafa.Stderr = os.Stderr
		if err := fafa.Run(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				os.Exit(exiterr.ExitCode())
			}
			return err
		}

		return nil
	},
}

func init() {
	binaryCmd.AddCommand(runBinaryCmd)
}
