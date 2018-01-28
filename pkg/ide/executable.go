package ide

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Executable represents an executable program in the context of the current ide project
type Executable struct {
	name      string
	container string
	program   string
}

// NewExecutable Creates a new Executable object
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

// Name returns the name of the Executable
func (exe *Executable) Name() string {
	return exe.name
}

// Program returns the program of the Executable
func (exe *Executable) Program() string {
	return exe.program
}

// Container returns the container (if applicable) of the Executable
func (exe *Executable) Container() string {
	return exe.container
}

// Definition returns the full definition of the Executable
func (exe *Executable) Definition() string {
	if exe.Containerized() {
		return exe.Container() + ":" + exe.Program()
	}

	return exe.Program()
}

// Source returns the name of the executable that needs to be invoked to start
// it
func (exe *Executable) Source() string {
	if exe.Containerized() {
		source, _ := os.Executable()
		return source
	}

	return exe.Program()
}

// Target returns the file path where the executable should be created
func (exe *Executable) Target() string {
	return filepath.Join(".git", "bin", exe.Name())
}

// Containerized determines if the Executable is containerized
func (exe *Executable) Containerized() bool {
	return exe.Container() != ""
}

// Create creates a symlink to the executable in the target location (see Executable.Target())
func (exe *Executable) Create() error {
	os.MkdirAll(filepath.Join(".git", "bin"), os.ModePerm)

	return os.Symlink(exe.Source(), exe.Target())
}

// Destroy removes the symlink for this Executable
func (exe *Executable) Destroy() error {
	return os.Remove(exe.Target())
}

// ListExecutables lists all the Executables for the current IDE project
func (project *Project) ListExecutables() []string {
	var executables []string

	for _, option := range project.config.Raw.Section("ide").Subsection("executables").Options {
		executables = append(executables, option.Key)
	}

	return executables
}

// NewExecutable adds a new executable to the current project
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

// GetExecutable returns an existing Executable for the current project
func (project *Project) GetExecutable(name string) (*Executable, error) {
	definition := project.config.Raw.Section("ide").Subsection("executables").Option(name)

	if definition == "" {
		return &Executable{}, errors.New("Executable " + name + " does not exist")
	}

	return NewExecutable(definition, name)
}

// RemoveExecutable removes an executable from the current project
func (project *Project) RemoveExecutable(name string) error {
	executable, err := project.GetExecutable(name)
	if err != nil {
		return err
	}

	project.config.Raw.Section("ide").Subsection("executables").RemoveOption(name)
	project.repository.Storer.SetConfig(project.config)

	return executable.Destroy()
}
