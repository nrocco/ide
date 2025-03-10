package cmd

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/anmitsu/go-shlex"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

type shimContext struct {
	User     string
	Group    string
	UID      int
	GID      int
	Project  string
	Location string
}

// RelDir calculates the relative directory in the ide/git repo
func (b shimContext) RelDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	rel, err := filepath.Rel(b.Location, dir)
	if err != nil {
		return ""
	}
	return rel
}

// IsTTY detects if the current file descriptors are attached to a TTY
func (b shimContext) IsTTY() bool {
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

func runPlainShim(command string, args []string) error {
	tmpl, err := template.New("shim").Parse(command)
	if err != nil {
		return err
	}

	var b bytes.Buffer

	err = tmpl.Execute(&b, shimContext{
		UID:      os.Getuid(),
		GID:      os.Getgid(),
		Project:  project.Name(),
		Location: project.Location(),
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

	runner := exec.Command(parts[0], parts[1:]...)
	runner.Stdin = os.Stdin
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	if err := runner.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			os.Exit(exiterr.ExitCode())
		}
		return err
	}

	return nil
}

func runComposeShim(command string, args []string) error {
	re := regexp.MustCompile(`\[(.+)\]:(.+)`)
	matches := re.FindStringSubmatch(command)
	if len(matches) != 3 {
		return errors.New("no bueno")
	}

	service := matches[1]
	command = os.ExpandEnv(matches[2])
	parts, err := shlex.Split(command, true)
	if err != nil {
		return err
	}

	parts = append(parts, args[1:]...)

	composeClient, err := project.DockerComposeClient()
	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(project.Location(), dir)
	if err != nil {
		return err
	}

	containers, err := composeClient.Ps(context.Background(), project.Name(), api.PsOptions{
		Services: []string{service},
	})
	if err != nil {
		return err
	}
	if len(containers) != 1 {
		return errors.New("blaaaat") // TODO
	}

	dockerClient, err := project.DockerClient()
	if err != nil {
		return err
	}

	containerSpec, err := dockerClient.ContainerInspect(context.Background(), containers[0].Name)
	if err != nil {
		return err
	}

	// TODO: add support for project name other then the directory name
	_, err = composeClient.Exec(context.Background(), project.Name(), api.RunOptions{
		Service:    service,
		Command:    parts,
		WorkingDir: filepath.Join(containerSpec.Config.WorkingDir, rel),
	})

	return err
}

var runShimCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a shim in the context of an ide project",
	Long:  "Run a shim in the context of an ide project",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		shims := []string{}
		if len(args) == 0 {
			if err := loadProject(cmd, args); err == nil {
				for shim := range project.ListShims() {
					shims = append(shims, shim)
				}
			}
		}
		return shims, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		command := project.GetShim(args[0])
		if command == "" {
			return errors.New("shim does not exist")
		} else if strings.HasPrefix(command, "compose[") {
			return runComposeShim(command, args)
		}
		return runPlainShim(command, args)
	},
}

func init() {
	shimCmd.AddCommand(runShimCmd)
}
