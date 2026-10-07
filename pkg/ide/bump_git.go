package ide

import (
	"errors"
	"fmt"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// GitVersionStrategy stores the version of a project as a lightweight git tag
type GitVersionStrategy struct {
	repository *git.Repository
}

// GitVersionStrategy returns a VersionStrategy backed by the git tags of the project
func (project *Project) GitVersionStrategy() *GitVersionStrategy {
	return &GitVersionStrategy{
		repository: project.repository,
	}
}

// Name returns the name of the strategy
func (strategy *GitVersionStrategy) Name() string {
	return "git"
}

// Applies reports whether the project is a git repository
func (strategy *GitVersionStrategy) Applies() bool {
	return strategy.repository != nil
}

// Current returns the most recently created semantic version tag.
// It returns an error if HEAD is already tagged with that version.
func (strategy *GitVersionStrategy) Current() (SemVer, error) {
	tags, err := strategy.repository.Tags()
	if err != nil {
		return SemVer{}, err
	}

	var latest SemVer
	var latestCommit plumbing.Hash
	var latestDate time.Time
	found := false

	err = tags.ForEach(func(ref *plumbing.Reference) error {
		version, err := ParseSemVer(ref.Name().Short())
		if err != nil {
			return nil
		}

		commit, date, err := strategy.resolveTag(ref)
		if err != nil {
			return nil
		}

		if !found || date.After(latestDate) {
			latest, latestCommit, latestDate, found = version, commit, date, true
		}
		return nil
	})
	if err != nil {
		return SemVer{}, err
	}

	if !found {
		return SemVer{}, ErrNoVersion
	}

	head, err := strategy.repository.Head()
	if err != nil {
		return SemVer{}, err
	}
	if head.Hash() == latestCommit {
		return SemVer{}, fmt.Errorf("HEAD is already tagged as %s", latest)
	}

	return latest, nil
}

// resolveTag returns the commit a tag points to and its creation date,
// which is the tagger date for annotated tags and the committer date for lightweight tags
func (strategy *GitVersionStrategy) resolveTag(ref *plumbing.Reference) (plumbing.Hash, time.Time, error) {
	tag, err := strategy.repository.TagObject(ref.Hash())
	switch {
	case err == nil:
		commit, err := tag.Commit()
		if err != nil {
			return plumbing.ZeroHash, time.Time{}, err
		}
		return commit.Hash, tag.Tagger.When, nil
	case errors.Is(err, plumbing.ErrObjectNotFound):
		commit, err := strategy.repository.CommitObject(ref.Hash())
		if err != nil {
			return plumbing.ZeroHash, time.Time{}, err
		}
		return commit.Hash, commit.Committer.When, nil
	default:
		return plumbing.ZeroHash, time.Time{}, err
	}
}

// Write creates a lightweight tag for the given version on HEAD
func (strategy *GitVersionStrategy) Write(version SemVer) error {
	head, err := strategy.repository.Head()
	if err != nil {
		return err
	}

	// Passing no options creates a lightweight (unsigned) tag
	_, err = strategy.repository.CreateTag(version.String(), head.Hash(), nil)
	return err
}
