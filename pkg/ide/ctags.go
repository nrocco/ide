package ide

import (
	"os"
	"os/exec"
	"path/filepath"
)

type Ctags struct {
	project        Project
	TagsFile       string
	VendorTagsFile string
}

func LoadCtags(project Project) (Ctags, error) {
	gitDir := project.Repository.Path()

	return Ctags{
		project:        project,
		TagsFile:       filepath.Join(gitDir, "tags"),
		VendorTagsFile: filepath.Join(gitDir, "tags_vendors"),
	}, nil
}

func (this *Ctags) Refresh() error {
	options := []string{
		"--recurse", "--tag-relative=yes", "--sort=yes",
		"--exclude=.git", "--exclude=.hg", "--exclude=.svn",
		"-f", this.TagsFile,
		"--kinds-php=cif", "--kinds-python=-i",
		"--languages=" + this.project.Language,
	}

	cmd := exec.Command("ctags", options...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	return nil
}
