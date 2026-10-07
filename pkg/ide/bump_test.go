package ide

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestParseSemVer(t *testing.T) {
	valid := map[string]SemVer{
		"0.0.0":               {},
		"1.2.3":               {Major: 1, Minor: 2, Patch: 3},
		"10.20.30":            {Major: 10, Minor: 20, Patch: 30},
		"1.2.3-rc.1":          {Major: 1, Minor: 2, Patch: 3, PreRelease: "rc.1"},
		"1.2.3+build.5":       {Major: 1, Minor: 2, Patch: 3, Build: "build.5"},
		"1.2.3-alpha+exp.sha": {Major: 1, Minor: 2, Patch: 3, PreRelease: "alpha", Build: "exp.sha"},
	}
	for input, expected := range valid {
		version, err := ParseSemVer(input)
		if err != nil {
			t.Errorf("ParseSemVer(%q) returned error: %s", input, err)
		}
		if version != expected {
			t.Errorf("ParseSemVer(%q) = %+v, expected %+v", input, version, expected)
		}
		if version.String() != input {
			t.Errorf("ParseSemVer(%q).String() = %q", input, version.String())
		}
	}

	for _, input := range []string{"", "1", "1.2", "v1.2.3", "01.2.3", "1.2.3.4", "1.2.3-", "a.b.c"} {
		if _, err := ParseSemVer(input); err == nil {
			t.Errorf("ParseSemVer(%q) expected an error", input)
		}
	}
}

func TestSemVerBump(t *testing.T) {
	version := SemVer{Major: 1, Minor: 2, Patch: 3, PreRelease: "rc.1"}
	expected := map[string]string{"major": "2.0.0", "minor": "1.3.0", "patch": "1.2.4"}
	for part, want := range expected {
		next, err := version.Bump(part)
		if err != nil {
			t.Fatal(err)
		}
		if next.String() != want {
			t.Errorf("Bump(%q) = %s, expected %s", part, next, want)
		}
	}
	if _, err := version.Bump("foo"); err == nil {
		t.Error("Bump(\"foo\") expected an error")
	}
}

func TestGitVersionStrategy(t *testing.T) {
	dir := t.TempDir()
	repository, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	worktree, _ := repository.Worktree()
	when := time.Now().Add(-time.Hour)
	commit := func() {
		when = when.Add(time.Minute)
		signature := &object.Signature{Name: "test", Email: "fake_example@email.com", When: when}
		if _, err := worktree.Commit("commit", &git.CommitOptions{AllowEmptyCommits: true, Author: signature, Committer: signature}); err != nil {
			t.Fatal(err)
		}
	}

	commit()
	project, err := NewProject(dir)
	if err != nil {
		t.Fatal(err)
	}
	strategy := project.GitVersionStrategy()

	if _, err := strategy.Current(); !errors.Is(err, ErrNoVersion) {
		t.Fatalf("expected ErrNoVersion, got %v", err)
	}

	if err := strategy.Write(SemVer{Minor: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err := strategy.Current(); err == nil {
		t.Fatal("expected an error because HEAD is already tagged")
	}

	commit()
	head, _ := repository.Head()
	// non semver tags are ignored
	repository.CreateTag("v9.9.9", head.Hash(), nil)
	// annotated tags use the tagger date
	tagger := &object.Signature{Name: "test", Email: "fake_example@email.com", When: when.Add(time.Minute)}
	if _, err := repository.CreateTag("0.2.0", head.Hash(), &git.CreateTagOptions{Tagger: tagger, Message: "0.2.0"}); err != nil {
		t.Fatal(err)
	}
	commit()

	current, err := strategy.Current()
	if err != nil {
		t.Fatal(err)
	}
	if current.String() != "0.2.0" {
		t.Errorf("expected 0.2.0, got %s", current)
	}

	if err := strategy.Write(SemVer{Major: 1}); err != nil {
		t.Fatal(err)
	}
	ref, err := repository.Tag("1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repository.TagObject(ref.Hash()); err == nil {
		t.Error("expected 1.0.0 to be a lightweight tag")
	}
}

func TestProjectVersionStrategy(t *testing.T) {
	dir := t.TempDir()
	if _, err := git.PlainInit(dir, false); err != nil {
		t.Fatal(err)
	}
	project, err := NewProject(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{VersionStrategyAuto, "git"} {
		strategy, err := project.VersionStrategy(name)
		if err != nil {
			t.Fatalf("VersionStrategy(%q) returned error: %s", name, err)
		}
		if strategy.Name() != "git" {
			t.Errorf("VersionStrategy(%q) = %s, expected git", name, strategy.Name())
		}
	}

	if _, err := project.VersionStrategy("foo"); err == nil {
		t.Error("VersionStrategy(\"foo\") expected an error")
	}
	if _, err := project.VersionStrategy("chef"); err == nil {
		t.Error("VersionStrategy(\"chef\") expected an error without metadata.rb")
	}

	if err := os.WriteFile(filepath.Join(dir, "metadata.rb"), []byte("version '1.0.0'\n"), 0644); err != nil {
		t.Fatal(err)
	}
	strategy, err := project.VersionStrategy(VersionStrategyAuto)
	if err != nil {
		t.Fatal(err)
	}
	if strategy.Name() != "chef" {
		t.Errorf("VersionStrategy(auto) = %s, expected chef", strategy.Name())
	}
}

func TestChefVersionStrategy(t *testing.T) {
	dir := t.TempDir()
	strategy := (&Project{location: dir}).ChefVersionStrategy()
	file := filepath.Join(dir, "metadata.rb")

	if strategy.Applies() {
		t.Error("expected chef strategy not to apply without metadata.rb")
	}

	os.WriteFile(file, []byte("name 'example'\nchef_version '>= 17.0'\n"), 0644)
	if strategy.Applies() {
		t.Error("expected chef strategy not to apply when only chef_version is present")
	}

	metadata := "name 'example'\nchef_version '>= 17.0'\nversion '2.1.11'\nsupports 'ubuntu'\n"
	os.WriteFile(file, []byte(metadata), 0644)
	if !strategy.Applies() {
		t.Fatal("expected chef strategy to apply")
	}

	current, err := strategy.Current()
	if err != nil {
		t.Fatal(err)
	}
	if current.String() != "2.1.11" {
		t.Errorf("expected 2.1.11, got %s", current)
	}

	if err := strategy.Write(SemVer{Major: 2, Minor: 2}); err != nil {
		t.Fatal(err)
	}
	content, _ := os.ReadFile(file)
	expected := "name 'example'\nchef_version '>= 17.0'\nversion '2.2.0'\nsupports 'ubuntu'\n"
	if string(content) != expected {
		t.Errorf("unexpected metadata.rb content:\n%s", content)
	}

	os.WriteFile(file, []byte("version \"1.2.3\"\n"), 0644)
	current, err = strategy.Current()
	if err != nil {
		t.Fatal(err)
	}
	if current.String() != "1.2.3" {
		t.Errorf("expected 1.2.3, got %s", current)
	}
}
