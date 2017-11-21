package main

import (
	"os"

	"github.com/nrocco/ide/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
