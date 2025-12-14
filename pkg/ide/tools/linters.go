package tools

import (
	"regexp"
)

var (
	// CookstyleLinter uses `cookstyle` to lint chef ruby files
	CookstyleLinter = Linter{
		Name:    "cookstyle",
		Command: "cookstyle",
		Args:    []string{"--display-cop-names", "--no-color", "--format", "emacs"},
		Matcher: regexp.MustCompile(`^(.*):(?P<Line>[\d]+):[\d]+:\s+(?<Message>.+)$`),
	}

	// EsLintLinter uses `eslint` to lint javascript, typescript and vue files
	EsLintLinter = Linter{
		Name:    "eslint",
		Command: "eslint",
		Args:    []string{"--format=compact"},
		Matcher: regexp.MustCompile(`^[^:]+: line (\d+), col \d+, (.*)$`),
	}

	// Flake8Linter uses `flake8` to lint python files
	Flake8Linter = Linter{
		Name:    "flake8",
		Command: "flake8",
		Args:    []string{"--extend-ignore=E501"},
		Matcher: regexp.MustCompile(`^(?P<File>.*):(?P<Line>\d+):(?P<Col>\d+):(?<Message>.*)$`),
	}

	// GobuildLinter uses `go build` to lint go files
	GobuildLinter = Linter{
		Name:    "go",
		Command: "go",
		Args:    []string{"build", "./..."},
		Matcher: regexp.MustCompile(`^(?P<File>[^:]+):(?P<Line>\d+):(?P<Col>\d+): (?P<Message>.*)$`),
	}

	// GolintLinter uses `golint` to lint go files
	GolintLinter = Linter{
		Name:    "golint",
		Command: "golint",
		Args:    []string{"./..."},
		Matcher: regexp.MustCompile(`^(?P<File>[^:]+):(?P<Line>\d+):(?P<Col>\d+): (?P<Message>.*)$`),
	}

	// GovetLinter uses the `go vet` tool to lint go files
	GovetLinter = Linter{
		Name:    "go",
		Command: "go",
		Args:    []string{"vet", "./..."},
		Matcher: regexp.MustCompile(`^(?P<File>[^:]+):(?P<Line>\d+):(?P<Col>\d+): (?P<Message>.*)$`),
	}

	// JqLinter uses the `jq` tool to lint json files
	JqLinter = Linter{
		Name:    "jq",
		Command: "jq",
		Args:    []string{"."},
		Matcher: regexp.MustCompile(`^parse (?<Severity>error): (?<Message>.*) at line (?<Line>\d+), .*$`),
	}

	// PhpLinter uses `php` to lint php files
	PhpLinter = Linter{
		Name:    "php",
		Command: "php",
		Args:    []string{"-l"},
		Matcher: regexp.MustCompile(`^PHP Parse (?<Severity>error): (?<Message>.*) in (?<File>.*) on line (?<Line>\d+)$`),
	}

	// ShellcheckLinter uses `shellcheck` to format bash/sh files
	ShellcheckLinter = Linter{
		Name:    "shellcheck",
		Command: "shellcheck",
		Args:    []string{"--format", "gcc"},
		Matcher: regexp.MustCompile(`^(?<File>.*):(?<Line>\d+):(?<Col>\d+): (?<Message>.*)$`),
	}

	// YamlLinter uses the python tool yamllint to lint yaml files
	YamlLinter = Linter{
		Name:    "yamllint",
		Command: "yamllint",
		Args:    []string{"--strict", "--format=parsable", "<file>"},
		Matcher: regexp.MustCompile(`^(?P<File>[^:]+):(?P<Line>\d+):(?P<Col>\d+): \[(?P<Severity>[a-z]+)\] (?P<Message>.*)$`),
	}
)
