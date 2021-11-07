package tools

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

type execLinterResult struct {
	*bufio.Scanner
}

func (r *execLinterResult) ForEach(re *regexp.Regexp, walker func([]string)) error {
	for r.Scan() {
		matches := re.FindStringSubmatch(r.Text())
		if nil == matches {
			continue
		}
		walker(matches)
	}
	return nil
}

func reportViolation(file string, line string, message string) {
	fmt.Printf("%s||%s||%s\n", file, line, message)
}

func execLinter(args ...string) *execLinterResult {
	output, _ := exec.Command(args[0], args[1:]...).CombinedOutput()
	return &execLinterResult{bufio.NewScanner(bytes.NewReader(output))}
}
