package tools

import (
	"regexp"
)

var (
	reLintEslin = regexp.MustCompile(`^[^:]+: line (\d+), col \d+, (.*)$`)
	reLintFlake8 = regexp.MustCompile(`^(.*):(\d+):\d+:(.*)$`)
	reLintFoodcritic = regexp.MustCompile(`^(.*): ([^:]+):(\d+)$`)
	reLintGobuild = regexp.MustCompile(`^([^:]+):(\d+):\d+: (.*)$`)
	reLintGolint = regexp.MustCompile(`^([^:]+):(\d+):\d+: (.*)$`)
	reLintGovet = regexp.MustCompile(`^([^:]+):(\d+):\d+: (.*)$`)
	reLintJq = regexp.MustCompile(`^parse error: (.*) at line (\d+), .*$`)
	reLintPhp = regexp.MustCompile(`^Parse error: (.*) in (.*) on line (\d+)$`)
	reLintRubocop = regexp.MustCompile(`^([^:]*):(\d+):\d+: (.*)$`)
	reLintShellcheck = regexp.MustCompile(`^(.*):(\d+):\d+: (.*)$`)
	reLintYaml = regexp.MustCompile(`^([^:]+):(\d+):\d+: (.*)$`)
)

// LintYaml uses `yamllint` to lint yaml files
func LintYaml(path string) error {
	return execLinter("yamllint", "--strict", "--format=parsable", path).ForEach(reLintYaml, func(err []string) {
		reportViolation(err[1], err[2], err[3])
	})
}

// LintShellcheck uses `shellcheck` to format bash/sh files
func LintShellcheck(path string) error {
	return execLinter("shellcheck", "--format", "gcc", path).ForEach(reLintShellcheck, func(err []string) {
		reportViolation(err[1], err[2], err[3])
	})
}

// LintRubocop uses `rubocop` to lint ruby files
func LintRubocop(path string) error {
	return execLinter("jq", ".", path).ForEach(reLintRubocop, func(err []string) {
		reportViolation(err[1], err[2], err[3])
	})
}

// LintPhpstan TODO
func LintPhpstan(path string) error {
	return nil
}

// LintPhp uses `php` to lint php files
func LintPhp(path string) error {
	return execLinter("php", "-l", path).ForEach(reLintPhp, func(err []string) {
		reportViolation(err[2], err[3], err[1])
	})
}

// LintJq uses the `jq` tool to lint json files
func LintJq(path string) error {
	return execLinter("jq", ".", path).ForEach(reLintJq, func(err []string) {
		reportViolation(path, err[2], err[1])
	})
}

// LintFoodcritic uses `foodcritic` to lint chef ruby files
func LintFoodcritic(path string) error {
	return execLinter("jq", ".", path).ForEach(reLintFoodcritic, func(err []string) {
		reportViolation(err[2], err[3], err[1])
	})
}

// LintGobuild uses `go build` to lint go files
func LintGobuild(path string) error {
	return execLinter("go", "build", "./...").ForEach(reLintGobuild, func(err []string) {
		reportViolation(err[1], err[2], err[3])
	})
}

// LintFlake8 uses `flake8` to lint python files
func LintFlake8(path string) error {
	return execLinter("flake8", "--extend-ignore=E501", path).ForEach(reLintFlake8, func(err []string) {
		reportViolation(err[1], err[2], err[3])
	})
}

// LintEslint uses `eslint` to lint typescript files
func LintEslint(path string) error {
	return execLinter("eslint", "--format=compact", path).ForEach(reLintEslin, func(err []string) {
		reportViolation(path, err[1], err[2])
	})
}

// LintGovet uses the `go vet` tool to lint go files
func LintGovet(path string) error {
	return execLinter("go", "vet", "./...").ForEach(reLintJq, func(err []string) {
		reportViolation(err[1], err[2], err[3])
	})
}

// LintGolint uses `golint` to lint go files
func LintGolint(path string) error {
	return execLinter("golint", "./...").ForEach(reLintGolint, func(err []string) {
		reportViolation(err[1], err[2], err[3])
	})
}
