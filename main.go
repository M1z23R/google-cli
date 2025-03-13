package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/M1z23R/google-cli/google"
	"github.com/M1z23R/google-cli/persistence"
)

func main() {
	persistence.InitDb()

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "login":
			{
				addGoogleAccount()
				break
			}
		case "rm":
			{
				idStr := os.Args[2]
				if idStr == "all" {
					persistence.ClearDb()
				} else {
					id, err := strconv.Atoi(os.Args[2])
					if err != nil {
						break
					}
					err = persistence.RemoveProfile(id)
					if err != nil {
						fmt.Println("Failed to remove profile.")
					} else {
						fmt.Println("Profile removed.")
					}
				}
				break
			}
		case "profiles":
			{
				printProfiles()
				break
			}
		case "messages":
			{
				var profile google.GoogleProfile
				var err error
				idStr := os.Args[2]
				if idStr == "" {
					err = persistence.GetFirstProfile(&profile)
				} else {
					err = persistence.GetProfile(idStr, &profile)
				}

				if err != nil {
					fmt.Println(err)
					break
				}
				printMessages(&profile)
				break
			}
		case "calendar":
			{
				var profile google.GoogleProfile
				var err error
				if len(os.Args) > 2 {
					err = persistence.GetProfile(os.Args[2], &profile)
				} else {
					err = persistence.GetFirstProfile(&profile)
				}

				if err != nil {
					fmt.Println(err)
					break
				}

				printCalendar(&profile)
				break
			}
		}
		return
	} else {
		printProfiles()
	}
}
