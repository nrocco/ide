package ide

import (
	"os"
	"path/filepath"
	"time"
)

// CtagsFile returns the path to the ctags file of the current ide project
func (project *Project) CtagsFile() string {
	return filepath.Join(project.location, ".git", "tags")
}

// HasCtagsFile returns true if a ctags file exists
func (project *Project) HasCtagsFile() bool {
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
