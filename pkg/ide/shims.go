package ide

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/anmitsu/go-shlex"
	"github.com/mattn/go-isatty"
)

// ShimList returns an array of shims which are added to this ide project
func (project *Project) ShimList() map[string]string {
	shims := map[string]string{}

	for _, option := range project.config.Raw.Section("ide").Subsection("shims").Options {
		option.Key = strings.ReplaceAll(option.Key, "-----", ".")
		shims[option.Key] = option.Value
	}

	return shims
}

// ShimGet returns the command for a shim
func (project *Project) ShimGet(shim string) string {
	shim = strings.ReplaceAll(shim, ".", "-----")
	return project.config.Raw.Section("ide").Subsection("shims").Option(shim)
}

// ShimRefresh syncs the shims from .git/config with .git/bin
func (project *Project) ShimRefresh() error {
	shims := project.ShimList()

	if len(shims) == 0 {
		return nil
	}

	if _, err := os.Stat(".git/bin"); os.IsNotExist(err) {
		if err := os.Mkdir(".git/bin", 0755); err != nil {
			return err
		}
	}

	source, _ := os.Executable()

	for shim := range shims {
		dest := filepath.Join(project.location, ".git", "bin", shim)
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

// ShimAdd adds a shim to this project
func (project *Project) ShimAdd(shim string, command string) error {
	validShimRegexp := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9._-]*$")

	if !validShimRegexp.MatchString(shim) {
		return fmt.Errorf("not a valid shim name: %s", shim)
	}

	dest := filepath.Join(project.location, ".git", "bin", shim)

	if _, err := os.Stat(dest); err == nil {
		return fmt.Errorf("shim %s already exists for this project", shim)
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

	project.config.Raw.SetOption("ide", "shims", strings.ReplaceAll(shim, ".", "-----"), command)

	return project.repository.Storer.SetConfig(project.config)
}

// ShimRemove removes a shim from this project
func (project *Project) ShimRemove(shim string) error {
	dest := filepath.Join(project.location, ".git", "bin", shim)

	if _, err := os.Lstat(dest); err == nil {
		if err := os.Remove(dest); err != nil {
			return err
		}
	}

	project.config.Raw.Section("ide").Subsection("shims").RemoveOption(strings.ReplaceAll(shim, ".", "-----"))

	return project.repository.Storer.SetConfig(project.config)
}

// ShimRun runs a shim by name, dispatching to the appropriate runner
func (project *Project) ShimRun(shim string, args []string) error {
	if strings.HasPrefix(shim, "compose[") {
		return project.runComposeShim(shim, args)
	}
	return project.runPlainShim(shim, args)
}

// IsTTY detects if the current file descriptors are attached to a TTY
func IsTTY() bool {
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

func (project *Project) runPlainShim(command string, args []string) error {
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

func (project *Project) runComposeShim(command string, args []string) error {
	re := regexp.MustCompile(`\[(.+)\]:(.+)`)
	matches := re.FindStringSubmatch(command)
	if len(matches) != 3 {
		return fmt.Errorf("invalid compose[service]:command string: %s", command)
	}

	service := matches[1]

	parts, err := shlex.Split(os.ExpandEnv(matches[2]), true)
	if err != nil {
		return err
	}

	runningContainers, _ := exec.Command("docker", "compose", "ps", "--quiet", "--filter", "status=running", service).Output()

	runArgs := []string{"docker", "compose"}
	if len(runningContainers) == 0 {
		runArgs = append(runArgs, "run", "--rm")
	} else {
		runArgs = append(runArgs, "exec")
	}
	if !IsTTY() {
		runArgs = append(runArgs, "-T")
	}
	if workDir := project.containerWorkDir(service); workDir != "" {
		runArgs = append(runArgs, "-w", workDir)
	}
	runArgs = append(runArgs, service)
	runArgs = append(runArgs, parts...)
	runArgs = append(runArgs, args[1:]...)

	runner := exec.Command(runArgs[0], runArgs[1:]...)
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

func (project *Project) containerWorkDir(service string) string {
	out, err := exec.Command("docker", "compose", "config", "--format", "json").Output()
	if err != nil {
		return ""
	}

	var config struct {
		Services map[string]struct {
			Volumes []struct {
				Type   string `json:"type"`
				Source string `json:"source"`
				Target string `json:"target"`
			} `json:"volumes"`
		} `json:"services"`
	}

	if err := json.Unmarshal(out, &config); err != nil {
		return ""
	}

	svc, ok := config.Services[service]
	if !ok {
		return ""
	}

	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	repoPath, _ := filepath.EvalSymlinks(project.Location())
	cwd, _ = filepath.EvalSymlinks(cwd)

	for _, v := range svc.Volumes {
		if v.Type != "bind" {
			continue
		}
		source, _ := filepath.EvalSymlinks(v.Source)

		// Check that the repo path falls within this mount source
		if _, err := filepath.Rel(source, repoPath); err != nil || !strings.HasPrefix(repoPath, source) {
			continue
		}

		// Map the host cwd to the corresponding container path
		cwdRel, err := filepath.Rel(source, cwd)
		if err != nil || strings.HasPrefix(cwdRel, "..") {
			continue
		}

		return filepath.Join(v.Target, cwdRel)
	}

	return ""
}

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
	return IsTTY()
}
