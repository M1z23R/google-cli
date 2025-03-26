package main

import (
	"os"

	clicolors "github.com/M1z23R/google-cli/cli-colors"
	"github.com/M1z23R/google-cli/persistence"
)

func main() {
	persistence.InitDb()

	var colorMode clicolors.ColorMode
	if len(os.Args) > 1 && os.Args[1] == "html" {
		colorMode = clicolors.HTML
	} else {
		colorMode = clicolors.Terminal
	}

	Interpret(colorMode)
}
