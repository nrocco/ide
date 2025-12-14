package tools

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

// Linter TODO
type Linter struct {
	Name    string
	Command string
	Args    []string
	Matcher *regexp.Regexp
}

// LinterResult wraps stdout/stderr of a linter to extract violations
type LinterResult struct {
	*bufio.Scanner
	Name string
	File string
	Matcher *regexp.Regexp
}

// LinterViolation is a structure with all information related to a single violation a linter detected
type LinterViolation struct {
	Linter   string
	File     string
	Line     string
	Col      string
	Severity string
	Message  string
}

// ForEachViolation parses the linters stdout/stderr and loops through every violation found
func (r *LinterResult) ForEachViolation(walker func(LinterViolation)) error {
	for r.Scan() {
		match := r.Matcher.FindStringSubmatch(r.Text())
		if nil == match {
			continue
		}
		violation := LinterViolation{
			Linter: r.Name,
			File:   r.File,
		}
		for i, name := range r.Matcher.SubexpNames() {
			if name == "File" {
				violation.File = match[i]
			} else if name == "Line" {
				violation.Line = match[i]
			} else if name == "Col" {
				violation.Col = match[i]
			} else if name == "Severity" {
				violation.Severity = match[i]
			} else if name == "Message" {
				violation.Message = match[i]
			}
		}
		walker(violation)
	}
	return nil
}

// PrintViolation is a walker function to print the violation to stdout
func PrintViolation(violation LinterViolation) {
	fmt.Printf("%s||%s||%s||%s||%s||%s\n", violation.File, violation.Line, violation.Col, violation.Linter, violation.Severity, violation.Message)
}

// Exec runs the linter on the specified file path
func (l *Linter) Exec(path string) *LinterResult {
	var args []string

	for _, arg := range l.Args {
		if arg == "<file>" {
			args = append(args, path)
		} else {
			args = append(args, arg)
		}
	}

	output, _ := exec.Command(l.Command, args...).CombinedOutput()

	return &LinterResult{
		Scanner: bufio.NewScanner(bytes.NewReader(output)),
		Name:    l.Name,
		File:    path,
		Matcher: l.Matcher,
	}
}
