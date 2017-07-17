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

func (project *Project) GetCtrlpCachFile() string {
	if project.ctrlpCacheFile == "" {
		cacheDir, _ := homedir.Expand("~/.cache/ctrlp") // TODO: make location of the cache configurable
		cacheDir, _ = filepath.Abs(cacheDir)

		cacheFilename := project.Location()
		cacheFilename = strings.TrimSuffix(cacheFilename, "/")
		cacheFilename = strings.Replace(cacheFilename, "/", "%", -1)
		cacheFilename = cacheFilename + ".txt"

		project.ctrlpCacheFile = filepath.Join(cacheDir, cacheFilename)
	}

	return project.ctrlpCacheFile
}

func (project *Project) RefreshCtrlp() error {
	if !project.IsConfigured() {
		return errors.New("Project must be configured before you can RefreshCtrlp")
	}

	cacheFile := project.GetCtrlpCachFile()

	file, err := os.Create(cacheFile)
	if err != nil {
		return err
	}

	defer file.Close()

	w := bufio.NewWriter(file)
	location := project.Location() + "/"

	filepath.Walk(location, func(path string, f os.FileInfo, err error) error {
		// TODO: configure what sort of files to exclude
		if !f.IsDir() && !strings.Contains(path, ".git") {
			fmt.Fprintln(w, strings.Replace(path, location, "", -1))
		}
		return nil
	})

	return w.Flush()
}
