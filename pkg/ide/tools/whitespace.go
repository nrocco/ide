package tools

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
)

var (
	trailingSpace      = regexp.MustCompile(`\s+$`)
	windowsLineEndings = regexp.MustCompile(`\r\n$`)
	literalTabs        = regexp.MustCompile(`^\t+`)
)

// LintWhitespace checks if a file has trailing whitespace, literal tabs or windows line endings
func LintWhitespace(path string, spaces bool, clrf bool, tabs bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Bytes()
		if spaces && trailingSpace.Match(line) {
			PrintViolation(LinterViolation{Linter: "whitespace", File: path, Line: strconv.Itoa(lineNumber), Col: "0", Severity: "error", Message: "trailing spaces detected"})
		}
		if clrf && windowsLineEndings.Match(line) {
			PrintViolation(LinterViolation{Linter: "whitespace", File: path, Line: strconv.Itoa(lineNumber), Col: "0", Severity: "error", Message: "windows line endings detected"})
		}
		if tabs && literalTabs.Match(line) {
			PrintViolation(LinterViolation{Linter: "whitespace", File: path, Line: strconv.Itoa(lineNumber), Col: "0", Severity: "error", Message: "literal tab characters detected"})
		}
		lineNumber++
	}

	return nil
}
