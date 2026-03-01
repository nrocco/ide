package tools

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

func getIgnoredDirs(path string) []string {
	items := []string{}

	output, err := exec.Command("git", "-C", path, "ls-files", "--exclude-standard", "--ignored", "--others", "--directory").Output()
	if err != nil {
		return items
	}
	for _, line := range strings.Split(string(output), "\n") {
		if !strings.HasSuffix(line, "/") {
			continue
		}
		items = append(items, strings.TrimRight(line, "/"))
	}
	return items
}

// WalkGitRepositories walks each root and calls fn for each directory that is a git repository.
// When recursive is false, only the root directories themselves are checked.
func WalkGitRepositories(recursive bool, fn func(path string) error, roots ...string) error {
	ignoredDirs := make(map[string][]string)
	for _, root := range roots {
		info, err := os.Stat(root)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			continue
		}
		err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				return nil
			}

			if d.Name() == ".git" {
				return filepath.SkipDir
			}

			parent := filepath.Dir(path)
			if value, ok := ignoredDirs[parent]; ok {
				if slices.Contains(value, d.Name()) {
					return filepath.SkipDir
				}
			}

			if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
				if err := fn(path); err != nil {
					return err
				}
				if !recursive {
					return filepath.SkipDir
				}
				ignoredDirs[path] = getIgnoredDirs(path)
			}

			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// FindGitRepositories walks each root and returns the paths of all directories that are git repositories.
// When recursive is false, only the root directories themselves are checked.
func FindGitRepositories(recursive bool, roots ...string) ([]string, error) {
	var repos []string

	err := WalkGitRepositories(recursive, func(path string) error {
		repos = append(repos, path)
		return nil
	}, roots...)

	return repos, err
}
