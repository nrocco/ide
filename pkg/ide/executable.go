package ide

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Executable struct {
	name      string
	container string
	program   string
}

func NewExecutable(definition string, name string) (*Executable, error) {
	var container string
	var program string

	parts := strings.SplitN(definition, ":", 2)

	if len(parts) == 2 {
		container = parts[0]
		program = parts[1]
	} else {
		if filepath.IsAbs(parts[0]) {
			program = parts[0]
		} else {
			fullPath, err := exec.Command("which", parts[0]).Output()
			if err != nil {
				return &Executable{}, err
			}

			program = strings.TrimSpace(string(fullPath))
		}
	}

	if name == "" {
		name = filepath.Base(program)
	}

	return &Executable{name, container, program}, nil
}

func (exe *Executable) Name() string {
	return exe.name
}

func (exe *Executable) Program() string {
	return exe.program
}

func (exe *Executable) Container() string {
	return exe.container
}

func (exe *Executable) Definition() string {
	if exe.Containerized() {
		return exe.Container() + ":" + exe.Program()
	}

	return exe.Program()
}

func (exe *Executable) Source() string {
	if exe.Containerized() {
		source, _ := os.Executable()
		return source
	}

	return exe.Program()
}

func (exe *Executable) Target() string {
	return filepath.Join(".git", "bin", exe.Name())
}

func (exe *Executable) Containerized() bool {
	return exe.Container() != ""
}

func (exe *Executable) Create() error {
	os.MkdirAll(filepath.Join(".git", "bin"), os.ModePerm)

	return os.Symlink(exe.Source(), exe.Target())
}

func (exe *Executable) Destroy() error {
	return os.Remove(exe.Target())
}

func (project *Project) ListExecutables() []string {
	var executables []string

	for _, option := range project.config.Raw.Section("ide").Subsection("executables").Options {
		executables = append(executables, option.Key)
	}

	return executables
}

func (project *Project) NewExecutable(definition string, name string) (*Executable, error) {
	if definition == "" {
		return &Executable{}, errors.New("You must specify a executable definition")
	}

	// TODO check if executable already exists

	executable, err := NewExecutable(definition, name)
	if err != nil {
		return &Executable{}, err
	}

	err = executable.Create()
	if err != nil {
		return &Executable{}, err
	}

	project.config.Raw.Section("ide").Subsection("executables").AddOption(executable.Name(), executable.Definition())
	project.repository.Storer.SetConfig(project.config)

	return executable, nil
}

func (project *Project) GetExecutable(name string) (*Executable, error) {
	definition := project.config.Raw.Section("ide").Subsection("executables").Option(name)

	if definition == "" {
		return &Executable{}, errors.New("Executable " + name + " does not exist")
	}

	return NewExecutable(definition, name)
}

func (project *Project) RemoveExecutable(name string) error {
	executable, err := project.GetExecutable(name)
	if err != nil {
		return err
	}

	project.config.Raw.Section("ide").Subsection("executables").RemoveOption(name)
	project.repository.Storer.SetConfig(project.config)

	return executable.Destroy()
}
