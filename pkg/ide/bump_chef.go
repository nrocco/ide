package ide

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// chefVersionRegexp matches the version line in a chef cookbook metadata.rb file, e.g. version '1.2.3'
var chefVersionRegexp = regexp.MustCompile(`(?m)^[ \t]*version[ \t]+(?:'([^'\n]*)'|"([^"\n]*)")`)

// ChefVersionStrategy stores the version of a chef cookbook in its metadata.rb file
type ChefVersionStrategy struct {
	file string
}

// ChefVersionStrategy returns a VersionStrategy backed by the metadata.rb file in the project root
func (project *Project) ChefVersionStrategy() *ChefVersionStrategy {
	return &ChefVersionStrategy{
		file: filepath.Join(project.Location(), "metadata.rb"),
	}
}

// Name returns the name of the strategy
func (strategy *ChefVersionStrategy) Name() string {
	return "chef"
}

// Applies reports whether there is a metadata.rb file containing a version
func (strategy *ChefVersionStrategy) Applies() bool {
	_, _, err := strategy.read()
	return err == nil
}

// Current returns the version from the metadata.rb file
func (strategy *ChefVersionStrategy) Current() (SemVer, error) {
	content, location, err := strategy.read()
	if err != nil {
		return SemVer{}, err
	}

	return ParseSemVer(string(content[location[0]:location[1]]))
}

// Write replaces the version in the metadata.rb file
func (strategy *ChefVersionStrategy) Write(version SemVer) error {
	content, location, err := strategy.read()
	if err != nil {
		return err
	}

	info, err := os.Stat(strategy.file)
	if err != nil {
		return err
	}

	updated := append([]byte{}, content[:location[0]]...)
	updated = append(updated, version.String()...)
	updated = append(updated, content[location[1]:]...)

	return os.WriteFile(strategy.file, updated, info.Mode().Perm())
}

// read returns the content of the metadata.rb file and the start and end
// offsets of the version value within that content
func (strategy *ChefVersionStrategy) read() ([]byte, [2]int, error) {
	content, err := os.ReadFile(strategy.file)
	if err != nil {
		return nil, [2]int{}, err
	}

	match := chefVersionRegexp.FindSubmatchIndex(content)
	if match == nil {
		return nil, [2]int{}, fmt.Errorf("no version found in %s", strategy.file)
	}

	// the version is either single (group 1) or double (group 2) quoted
	if match[2] >= 0 {
		return content, [2]int{match[2], match[3]}, nil
	}
	return content, [2]int{match[4], match[5]}, nil
}
