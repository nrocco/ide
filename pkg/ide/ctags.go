package ide

import (
	"bufio"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

// CtagsEntry represents a single ctag tag
type CtagsEntry struct {
	Type      string `json:"_type"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Pattern   string `json:"pattern"`
	Line      int    `json:"line"`
	Kind      string `json:"kind"`
	Scope     string `json:"scope"`
	Signature string `json:"signature"`
	Typeref   string `json:"typeref"`
	ScopeKind string `json:"scopeKind"`
	Language  string `json:"language"`
	Roles     string `json:"roles"`
	Access    string `json:"access"`
	End       int    `json:"end"`
}

// IsPublic returns true if the function is considered a public function
func (entry *CtagsEntry) IsPublic() bool {
	if entry.Access == "public" {
		return true
	} else if entry.Access == "private" || entry.Access == "protected" {
		return false
	}

	if entry.Language == "go" {
		return unicode.IsUpper(rune(entry.Name[0]))
	}

	return false
}

// IsPrivate returns true if the function is considered not a public function
func (entry *CtagsEntry) IsPrivate() bool {
	return !entry.IsPublic()
}

// CtagsFile returns the path to the ctags file of the current ide project
func (project *Project) CtagsFile() string {
	return filepath.Join(project.location, ".git", "tags")
}

// CtagsFileExists returns true if a ctags file exists
func (project *Project) CtagsFileExists() bool {
	if _, err := os.Stat(project.CtagsFile()); err == nil {
		return true
	}
	return false
}

// CtagsFileAge returns the time the ctags file was last modified
func (project *Project) CtagsFileAge() time.Time {
	file, err := os.Open(project.CtagsFile())
	if err != nil {
		return time.Time{}
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return time.Time{}
	}

	return stat.ModTime()
}

// CtagsFileSize returns the size of the ctags file in bytes
func (project *Project) CtagsFileSize() uint64 {
	file, err := os.Open(project.CtagsFile())
	if err != nil {
		return 0
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return 0
	}

	return uint64(stat.Size())
}

// CtagsParseCode parses ctags from the given files
func (project *Project) CtagsParseCode(walker func(CtagsEntry), files ...string) error {
	args := []string{
		"-f-",
		"--excmd=number",
		"--recurse=yes",
		"--sort=no",
		"--totals=no",
		"--with-list-header=no",
		"--machinable=yes",
		"--kinds-php=f",
		"--kinds-go=f",
		"--kinds-python=cfm",
		"--fields=aCeEfFikKlmnNpPrRsStxzZ",
		"--output-format=json",
	}
	args = append(args, files...)

	ctags := exec.Command("ctags", args...)

	stdout, err := ctags.StdoutPipe()
	if err != nil {
		return err
	}

	if err := ctags.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var entry CtagsEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return err
		}
		if entry.Kind == "class" {
			continue
		}
		entry.Language = strings.ToLower(entry.Language)
		entry.Typeref = strings.TrimPrefix(entry.Typeref, "typename:")
		walker(entry)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return ctags.Wait()
}
