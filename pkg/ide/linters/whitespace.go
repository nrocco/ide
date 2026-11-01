package linters

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
	leadingSpaces      = regexp.MustCompile(`^ +`)
)

// LintWhitespace checks if a file has trailing whitespace, windows line endings or the wrong indentation.
// If tabs is true lines indented with spaces are reported, otherwise lines indented with tabs are reported.
func LintWhitespace(path string, spaces bool, clrf bool, tabs bool) error {
	file, err := os.Open(path) // #nosec G304 -- path is supplied by the user on the command line
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
		if tabs && leadingSpaces.Match(line) {
			PrintViolation(LinterViolation{Linter: "whitespace", File: path, Line: strconv.Itoa(lineNumber), Col: "0", Severity: "error", Message: "indentation with spaces detected"})
		}
		if !tabs && literalTabs.Match(line) {
			PrintViolation(LinterViolation{Linter: "whitespace", File: path, Line: strconv.Itoa(lineNumber), Col: "0", Severity: "error", Message: "literal tab characters detected"})
		}
		lineNumber++
	}

	return nil
}
