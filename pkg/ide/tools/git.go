package tools

import (
	"io/fs"
	"os"
	"path/filepath"
)

// WalkGitRepositories walks each root and calls fn for each directory that is a git repository.
// When recursive is false, only the root directories themselves are checked.
func WalkGitRepositories(recursive bool, fn func(path string) error, roots ...string) error {
	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				return nil
			}

			if d.Name() == ".git" {
				return filepath.SkipDir
			}

			if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
				if err := fn(path); err != nil {
					return err
				}
				if !recursive {
					return filepath.SkipDir
				}
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
