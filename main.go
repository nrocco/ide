package main

import (
	"log"
	"os"
	"strings"

	"github.com/nrocco/ide/cmd"
)

func main() {
	if strings.Contains(os.Args[0], ".git/hooks") {
		elems := strings.Split(os.Args[0], "/")

		args := []string{os.Args[0], "hook", "exec", elems[2]}
		os.Args = append(args, os.Args[1:]...)
	}

	if err := cmd.RootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
