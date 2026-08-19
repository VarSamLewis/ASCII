package main

import (
	"os"

	"ascii/ascii/src/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
