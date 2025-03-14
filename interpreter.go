package main

import (
	"fmt"
	"os"
	"strconv"

	clicolors "github.com/M1z23R/google-cli/cli-colors"
	"github.com/M1z23R/google-cli/google"
	"github.com/M1z23R/google-cli/persistence"
)

func Interpret(colorMode clicolors.ColorMode) {
	if len(os.Args) >= 3 {
		switch os.Args[2] {
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

			// single
		case "messages":
			{
				var profile google.GoogleProfile
				var err error
				if len(os.Args) > 3 {
					err = persistence.GetProfile(os.Args[3], &profile)
				} else {
					err = persistence.GetFirstProfile(&profile)
				}

				if err != nil {
					fmt.Println(err)
					break
				}
				printLatestMessages(colorMode, &profile, 1)
				break
			}
		case "calendar":
			{
				var profile google.GoogleProfile
				var err error
				if len(os.Args) > 3 {
					err = persistence.GetProfile(os.Args[3], &profile)
				} else {
					err = persistence.GetFirstProfile(&profile)
				}

				if err != nil {
					fmt.Println(err)
					break
				}

				printEvents(colorMode, &profile, 1)
				break
			}
		}
		return
	} else {
		printProfiles(colorMode)
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
			var profile google.GoogleProfile
			var err error
			if len(os.Args) > 3 {
				err = persistence.GetProfile(os.Args[4], &profile)
			} else {
				err = persistence.GetFirstProfile(&profile)
			}

			if err != nil {
				fmt.Println(err)
				break
			}
			printLatestMessages(colorMode, &profile, 10)
			break
		}
	case "calendar":
		{
			var profile google.GoogleProfile
			var err error
			if len(os.Args) > 3 {
				err = persistence.GetProfile(os.Args[4], &profile)
			} else {
				err = persistence.GetFirstProfile(&profile)
			}

			if err != nil {
				fmt.Println(err)
				break
			}

			printEvents(colorMode, &profile, 10)
			break
		}
	}
}
