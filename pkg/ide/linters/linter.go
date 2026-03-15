package linters

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"regexp"
)

// Matcher parses raw linter output into violations
type Matcher interface {
	Parse(output []byte, name, file string) []LinterViolation
}

// RegexMatcher parses linter output line-by-line using a regular expression
type RegexMatcher struct {
	re *regexp.Regexp
}

// NewRegexMatcher creates a RegexMatcher from a regexp pattern string
func NewRegexMatcher(pattern string) *RegexMatcher {
	return &RegexMatcher{re: regexp.MustCompile(pattern)}
}

// Parse scans output line by line and extracts violations using the regexp
func (m *RegexMatcher) Parse(output []byte, name, file string) []LinterViolation {
	var violations []LinterViolation
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		match := m.re.FindStringSubmatch(scanner.Text())
		if match == nil {
			continue
		}
		v := LinterViolation{
			Linter: name,
			File:   file,
		}
		for i, subname := range m.re.SubexpNames() {
			switch subname {
			case "File":
				v.File = match[i]
			case "Line":
				v.Line = match[i]
			case "Col":
				v.Col = match[i]
			case "Severity":
				v.Severity = match[i]
			case "Message":
				v.Message = match[i]
			}
		}
		violations = append(violations, v)
	}
	return violations
}

// JSONMatcher parses linter output as JSON using a provided parse function
type JSONMatcher struct {
	parse func(output []byte, name, file string) []LinterViolation
}

// NewJSONMatcher creates a JSONMatcher with the given parse function
func NewJSONMatcher(parse func(output []byte, name, file string) []LinterViolation) *JSONMatcher {
	return &JSONMatcher{parse: parse}
}

// Parse delegates to the provided parse function
func (m *JSONMatcher) Parse(output []byte, name, file string) []LinterViolation {
	return m.parse(output, name, file)
}

// NewLinterResult creates a LinterResult by reading all bytes from r
func NewLinterResult(r io.Reader, name, file string, matcher Matcher) (*LinterResult, error) {
	output, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &LinterResult{
		output:  output,
		Name:    name,
		File:    file,
		Matcher: matcher,
	}, nil
}

// Linter TODO
type Linter struct {
	Name    string
	Command string
	Args    []string
	Matcher Matcher
}

// LinterResult wraps stdout/stderr of a linter to extract violations
type LinterResult struct {
	output []byte
	Name   string
	File   string
	Matcher Matcher
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
	for _, v := range r.Matcher.Parse(r.output, r.Name, r.File) {
		walker(v)
	}
	return nil
}

// PrintViolation is a walker function to print the violation to stdout
func PrintViolation(violation LinterViolation) {
	fmt.Printf("%s||%s||%s||%s||%s||%s\n", violation.File, violation.Line, violation.Col, violation.Linter, violation.Severity, violation.Message)
}

// Exec runs the linter on the specified file path
func (l *Linter) Exec(path string, debug bool) *LinterResult {
	var args []string

	for _, arg := range l.Args {
		if arg == "<file>" {
			args = append(args, path)
		} else {
			args = append(args, arg)
		}
	}

	output, err := exec.Command(l.Command, args...).CombinedOutput() // TODO handle err here

	if debug {
		fmt.Printf("command: %s\nargs: %v\nerror: %s\noutput:\n%s\n", l.Command, args, err, output)
	}

	return &LinterResult{
		output:  output,
		Name:    l.Name,
		File:    path,
		Matcher: l.Matcher,
	}
}
