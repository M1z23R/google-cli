package main

import (
	"fmt"
	"os"
	"strconv"

	clicolors "github.com/M1z23R/google-cli/cli-colors"
	"github.com/M1z23R/google-cli/persistence"
)

func Interpret(colorMode clicolors.ColorMode) {
	if len(os.Args) >= 3 {
		switch os.Args[2] {
		case "secrets":
			{
				addSecrets()
				break
			}
		case "login":
			{
				addGoogleAccount()
				break
			}
		case "rm":
			{
				rm()
				break
			}

		case "ls":
			{
				ls(colorMode)
				break
			}
		default:
			{
				single(colorMode)
				break
			}
		}
	} else {
		printProfiles(colorMode)
	}
}

func single(colorMode clicolors.ColorMode) {
	switch os.Args[2] {
	case "messages":
		{
			printLatestMessages(colorMode, 1)
			break
		}
	case "calendar":
		{
			printEvents(colorMode, true)
			break
		}
	}
}

func rm() {
	idStr := os.Args[3]
	if idStr == "all" {
		persistence.ClearDb()
	} else {
		id, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Printf("Failed to convert argument `%s` to int.\n", os.Args[3])
			return
		}
		err = persistence.RemoveProfile(id)
		if err != nil {
			fmt.Println("Failed to remove profile.")
		} else {
			fmt.Println("Profile removed.")
		}
	}
}

func ls(colorMode clicolors.ColorMode) {
	switch os.Args[3] {
	case "profiles":
		{
			printProfiles(colorMode)
			break
		}
	case "messages":
		{
			printLatestMessages(colorMode, 10)
			break
		}
	case "calendar":
		{
			printEvents(colorMode, false)
			break
		}
	}
}
