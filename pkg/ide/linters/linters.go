package linters

import (
	"encoding/json"
	"strconv"
)

var (
	// CookstyleLinter uses `cookstyle` to lint chef ruby files
	CookstyleLinter = Linter{
		Name:    "cookstyle",
		Command: "cookstyle",
		Args:    []string{"--display-cop-names", "--no-color", "--format", "emacs", "<file>"},
		Matcher: NewRegexMatcher(`^(.*):(?P<Line>[\d]+):[\d]+:\s+(?<Message>.+)$`),
	}

	// EsLintLinter uses `eslint` to lint javascript, typescript and vue files
	EsLintLinter = Linter{
		Name:    "eslint",
		Command: "eslint",
		Args:    []string{"--format=json", "<file>"},
		Matcher: NewJSONMatcher(func(output []byte, name, file string) []LinterViolation {
			var results []struct {
				FilePath string `json:"filePath"`
				Messages []struct {
					Line     int    `json:"line"`
					Column   int    `json:"column"`
					Severity int    `json:"severity"` // 1 = warn, 2 = error
					Message  string `json:"message"`
				} `json:"messages"`
			}
			if err := json.Unmarshal(output, &results); err != nil {
				return nil
			}
			var violations []LinterViolation
			for _, r := range results {
				for _, m := range r.Messages {
					severity := "warning"
					if m.Severity == 2 {
						severity = "error"
					}
					violations = append(violations, LinterViolation{
						Linter:   name,
						File:     r.FilePath,
						Line:     strconv.Itoa(m.Line),
						Col:      strconv.Itoa(m.Column),
						Severity: severity,
						Message:  m.Message,
					})
				}
			}
			return violations
		}),
	}

	// Flake8Linter uses `flake8` to lint python files
	Flake8Linter = Linter{
		Name:    "flake8",
		Command: "flake8",
		Args:    []string{"--extend-ignore=E501", "<file>"},
		Matcher: NewRegexMatcher(`^(?P<File>.*):(?P<Line>\d+):(?P<Col>\d+):(?<Message>.*)$`),
	}

	// GobuildLinter uses `go build` to lint go files
	GobuildLinter = Linter{
		Name:    "go",
		Command: "go",
		Args:    []string{"build", "./..."},
		Matcher: NewRegexMatcher(`^(?P<File>[^:]+):(?P<Line>\d+):(?P<Col>\d+): (?P<Message>.*)$`),
	}

	// GolintLinter uses `golint` to lint go files
	GolintLinter = Linter{
		Name:    "golint",
		Command: "golint",
		Args:    []string{"./..."},
		Matcher: NewRegexMatcher(`^(?P<File>[^:]+):(?P<Line>\d+):(?P<Col>\d+): (?P<Message>.*)$`),
	}

	// GovetLinter uses the `go vet` tool to lint go files
	GovetLinter = Linter{
		Name:    "go",
		Command: "go",
		Args:    []string{"vet", "./..."},
		Matcher: NewRegexMatcher(`^(?P<File>[^:]+):(?P<Line>\d+):(?P<Col>\d+): (?P<Message>.*)$`),
	}

	// JqLinter uses the `jq` tool to lint json files
	JqLinter = Linter{
		Name:    "jq",
		Command: "jq",
		Args:    []string{".", "<file>"},
		Matcher: NewRegexMatcher(`^jq: parse (?<Severity>error): (?<Message>.*) at line (?<Line>\d+), column (?<Col>\d+)$`),
	}

	// PhpLinter uses `php` to lint php files
	PhpLinter = Linter{
		Name:    "php",
		Command: "php",
		Args:    []string{"-l", "<file>"},
		Matcher: NewRegexMatcher(`^PHP Parse (?<Severity>error): (?<Message>.*) in (?<File>.*) on line (?<Line>\d+)$`),
	}

	// ShellcheckLinter uses `shellcheck` to format bash/sh files
	ShellcheckLinter = Linter{
		Name:    "shellcheck",
		Command: "shellcheck",
		Args:    []string{"--format", "gcc", "<file>"},
		Matcher: NewRegexMatcher(`^(?<File>.*):(?<Line>\d+):(?<Col>\d+): (?<Message>.*)$`),
	}

	// YamlLinter uses the python tool yamllint to lint yaml files
	YamlLinter = Linter{
		Name:    "yamllint",
		Command: "yamllint",
		Args:    []string{"--strict", "--format=parsable", "<file>"},
		Matcher: NewRegexMatcher(`^(?P<File>[^:]+):(?P<Line>\d+):(?P<Col>\d+): \[(?P<Severity>[a-z]+)\] (?P<Message>.*)$`),
	}
)
