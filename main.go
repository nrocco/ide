package main

import (
	"fmt"
	"os"

	"github.com/nrocco/ide/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
