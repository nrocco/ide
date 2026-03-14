package main

import (
	"os"
	"path"
	"strings"
)

func main() {
	if strings.Contains(os.Args[0], ".git/hooks") {
		args := []string{os.Args[0], "hook", "run", path.Base(os.Args[0]), "--"}
		os.Args = append(args, os.Args[1:]...)
	} else if os.Args[0] == "rgit" {
		args := []string{"ide", "tool", "rgit"}
		os.Args = append(args, os.Args[1:]...)
	} else if os.Args[0] == "rrgit" {
		args := []string{"ide", "tool", "rgit", "--rgit-recursive"}
		os.Args = append(args, os.Args[1:]...)
	} else if strings.Contains(os.Args[0], "/go-build") {
		// this exists for development only
	} else if !strings.Contains(os.Args[0], "ide") {
		args := []string{os.Args[0], "shim", "run", "--", path.Base(os.Args[0])}
		os.Args = append(args, os.Args[1:]...)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
