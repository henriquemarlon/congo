package main

import (
	"os"

	"github.com/henriquemarlon/congo/cmd/congo/root"
)

func main() {
	err := root.Cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
