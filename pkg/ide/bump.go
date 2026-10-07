package ide

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ErrNoVersion is returned by a VersionStrategy when no current version could be found
var ErrNoVersion = errors.New("no existing version found")

// semVerRegexp is the official regular expression from https://semver.org
var semVerRegexp = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

// SemVer represents a semantic version
type SemVer struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease string
	Build      string
}

// ParseSemVer parses a semantic version string like 1.2.3, 1.2.3-rc.1 or 1.2.3+build.5
func ParseSemVer(value string) (SemVer, error) {
	matches := semVerRegexp.FindStringSubmatch(value)
	if matches == nil {
		return SemVer{}, fmt.Errorf("invalid semantic version: %q", value)
	}

	var parts [3]int
	for i := range parts {
		number, err := strconv.Atoi(matches[i+1])
		if err != nil {
			return SemVer{}, fmt.Errorf("invalid semantic version: %q: %w", value, err)
		}
		parts[i] = number
	}

	return SemVer{
		Major:      parts[0],
		Minor:      parts[1],
		Patch:      parts[2],
		PreRelease: matches[4],
		Build:      matches[5],
	}, nil
}

// Bump returns the next version for the given part (major, minor or patch).
// Pre-release and build metadata are dropped.
func (version SemVer) Bump(part string) (SemVer, error) {
	switch part {
	case "major":
		return SemVer{Major: version.Major + 1}, nil
	case "minor":
		return SemVer{Major: version.Major, Minor: version.Minor + 1}, nil
	case "patch":
		return SemVer{Major: version.Major, Minor: version.Minor, Patch: version.Patch + 1}, nil
	default:
		return SemVer{}, fmt.Errorf("invalid version part %q, must be one of major, minor or patch", part)
	}
}

// String returns the string representation of the semantic version
func (version SemVer) String() string {
	value := fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)
	if version.PreRelease != "" {
		value += "-" + version.PreRelease
	}
	if version.Build != "" {
		value += "+" + version.Build
	}
	return value
}

// VersionStrategy reads and writes the version of a project
type VersionStrategy interface {
	// Name returns the name used to select the strategy
	Name() string
	// Applies reports whether the strategy can be used for the project
	Applies() bool
	// Current returns the current version, or ErrNoVersion if there is none
	Current() (SemVer, error)
	// Write persists the given version
	Write(version SemVer) error
}

// VersionStrategyAuto selects the first version strategy that applies to the project
const VersionStrategyAuto = "auto"

// VersionStrategies returns all version strategies in order of precedence.
// More specific strategies must come before more generic ones, as auto
// detection picks the first strategy that applies.
func (project *Project) VersionStrategies() []VersionStrategy {
	return []VersionStrategy{
		project.GitVersionStrategy(),
	}
}

// VersionStrategy returns the version strategy with the given name, or the
// first one that applies to the project if name is VersionStrategyAuto
func (project *Project) VersionStrategy(name string) (VersionStrategy, error) {
	strategies := project.VersionStrategies()

	names := []string{VersionStrategyAuto}
	for _, strategy := range strategies {
		names = append(names, strategy.Name())
	}

	for _, strategy := range strategies {
		if name == VersionStrategyAuto && strategy.Applies() {
			return strategy, nil
		}
		if name == strategy.Name() {
			if !strategy.Applies() {
				return nil, fmt.Errorf("strategy %q does not apply to this project", name)
			}
			return strategy, nil
		}
	}

	if name == VersionStrategyAuto {
		return nil, errors.New("no applicable version strategy found")
	}

	return nil, fmt.Errorf("unknown strategy %q, must be one of: %s", name, strings.Join(names, ", "))
}
