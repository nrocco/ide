package ide

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

type Ctrlp struct {
	project   Project
	CacheFile string
}

func LoadCtrlp(project Project) (Ctrlp, error) {
	if !project.IsConfigured() {
		return Ctrlp{}, errors.New("Project must be configured.")
	}

	if project.Location == "" {
		return Ctrlp{}, errors.New("Project directory cannot be empty")
	}

	cacheDir, _ := homedir.Expand("~/.cache/ctrlp") // TODO: make this configurable
	cacheDir, _ = filepath.Abs(cacheDir)

	cacheFilename := project.Location
	cacheFilename = strings.TrimSuffix(cacheFilename, "/")
	cacheFilename = strings.Replace(cacheFilename, "/", "%", -1)
	cacheFilename = cacheFilename + ".txt"

	return Ctrlp{
		project:   project,
		CacheFile: filepath.Join(cacheDir, cacheFilename),
	}, nil
}

func (this *Ctrlp) Refresh() error {
	if this.CacheFile == "" {
		return errors.New("Cache file should not be empty")
	}

	file, err := os.Create(this.CacheFile)
	if err != nil {
		return err
	}

	defer file.Close()

	w := bufio.NewWriter(file)

	filepath.Walk(this.project.Location, func(path string, f os.FileInfo, err error) error {
		fmt.Fprintln(w, strings.Replace(path, this.project.Location, "", -1))
		return nil
	})

	return w.Flush()
}
