package main

import (
	"os"

	"github.com/lambdal/guest-agent/tools/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
