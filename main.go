package main

import (
	"fmt"
	"os"
	
	"github.com/opf/openproject-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}