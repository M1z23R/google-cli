package main

import (
	"fmt"
	"html"
	"os"
	"strconv"
	"strings"
	"time"

	clicolors "github.com/M1z23R/google-cli/cli-colors"
	"github.com/M1z23R/google-cli/google"
	"github.com/M1z23R/google-cli/google/auth"
	"github.com/M1z23R/google-cli/google/calendar"
	"github.com/M1z23R/google-cli/google/gmail"
	"github.com/M1z23R/google-cli/persistence"
)

func printProfiles(colorMode clicolors.ColorMode) {
	profiles := []google.GoogleProfile{}
	err := persistence.GetProfiles(&profiles)
	if err != nil {
		fmt.Println("Couldn't load profiles")
	}

	if len(profiles) == 0 {
		fmt.Println("No profiles are added.")
	} else {
		for i, profile := range profiles {
			fmt.Printf("%s%d) %s%s\n", clicolors.GetColor(colorMode, ""), i+1, profile.EmailAddress, clicolors.GetColor(colorMode, clicolors.Reset))
		}
	}

}

func printLatestMessages(colorMode clicolors.ColorMode, count int) {
	var profile google.GoogleProfile
	profileId := -1

	var err error
	if len(os.Args) > 3 {
		profileId, err = strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println("Couldn't parse profileId")
			return
		}

		err = persistence.GetProfile(profileId, &profile)
	} else {
		err = persistence.GetProfile(-1, &profile)
	}

	if err != nil {
		fmt.Println(err)
		return
	}
	var r google.GmailMessagesResponse

	if profile.Tokens.Expires.Before(time.Now()) {
		err = auth.RefreshToken(&profile)
		if err != nil {
			fmt.Println(err)
			return
		}

		persistence.UpsertProfile(&profile)
	}

	gmail.ListMessages(&profile, &r, count)

	var c int
	gmail.GetUnreadCount(&profile, &c)
	unreadCount := fmt.Sprintf("%sUnread: %d%s", clicolors.GetUnreadColor(colorMode, c), c, clicolors.GetColor(colorMode, clicolors.Reset))
	if count > 1 {
		fmt.Printf("%s\n", unreadCount)
	}

	for i, v := range r.Messages {
		var m google.GmailMessage
		gmail.GetMessageMetadata(&profile, &m, v.Id)

		from := extractEmail(gmail.ExtractHeader(&m, "From"))
		subject := gmail.ExtractHeader(&m, "Subject")
		if len(subject) > 33 {
			subject = subject[:30] + "..."
		}

		lastMessageFrom := fmt.Sprintf("%s%s%s", clicolors.GetColor(colorMode, clicolors.LightGreen), from, clicolors.GetColor(colorMode, clicolors.Reset))
		lastMessageSubject := fmt.Sprintf("%s%s%s", clicolors.GetColor(colorMode, clicolors.LightBlue), html.EscapeString(subject), clicolors.GetColor(colorMode, clicolors.Reset))
		lastMessage := fmt.Sprintf("%s | %s", lastMessageFrom, lastMessageSubject)

		if count > 1 {
			fmt.Printf("%d) %s\n", i+1, lastMessage)
		} else {
			fmt.Printf("%s | %s\n", unreadCount, lastMessage)
		}
	}
}

func printEvents(colorMode clicolors.ColorMode, next bool) {
	var profile google.GoogleProfile
	profileId := -1

	var err error
	if len(os.Args) > 3 {
		profileId, err = strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println("Couldn't parse profileId")
			return
		}

		err = persistence.GetProfile(profileId, &profile)
	} else {
		err = persistence.GetProfile(-1, &profile)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	if profile.Tokens.Expires.Before(time.Now()) {
		err = auth.RefreshToken(&profile)
		if err != nil {
			fmt.Println(err)
			return
		}

		persistence.UpsertProfile(&profile)
	}

	var r google.CalendarEventResponse

	start := time.Now()
	calendar.GetEvents(&profile, &r, start, 10)

	for _, v := range r.Items {
		if next {
			if v.Start.DateTime.After(time.Now()) {
				fmt.Printf("%s%s%s | %s%s%s\n", clicolors.GetColor(colorMode, clicolors.LightGreen), v.Summary, clicolors.GetColor(colorMode, clicolors.Reset), clicolors.GetEventColor(colorMode, v.Start.DateTime), joinTimes(v.Start.DateTime, v.End.DateTime), clicolors.GetColor(colorMode, clicolors.Reset))
				break
			}
		} else {
			fmt.Printf("%s%s%s | %s%s%s\n", clicolors.GetColor(colorMode, clicolors.LightGreen), v.Summary, clicolors.GetColor(colorMode, clicolors.Reset), clicolors.GetEventColor(colorMode, v.Start.DateTime), joinTimes(v.Start.DateTime, v.End.DateTime), clicolors.GetColor(colorMode, clicolors.Reset))
		}
	}
}

func addSecrets() {
	if len(os.Args) < 6 {
		fmt.Println("You must provide ClientId, ClientSecret and RedirectUri.")
		return
	}

	secret := google.GoogleSecret{ClientId: os.Args[3], ClientSecret: os.Args[4], RedirectUri: os.Args[5]}
	err := persistence.UpsertGoogleSecrets(&secret)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Something went wrong.")
	} else {
		fmt.Println("Secrets stored successfully.")
	}
}

func addGoogleAccount() {
	if len(os.Args) < 4 {
		fmt.Println("Must provide google secrets id.")
		return
	}
	id, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("Failed to convert argument `%s` to int.\n", os.Args[3])
		return
	}
	profile := google.GoogleProfile{}

	err = persistence.GetSecret(id, &profile.Secrets)
	if err != nil {
		fmt.Printf("Failed to get secret with id '%s'.\n", os.Args[3])
		return
	}

	err = auth.AddProfile(&profile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = persistence.UpsertProfile(&profile)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Profile %s was upserted successfully.\n", profile.EmailAddress)
}

// utility
func extractEmail(input string) string {
	startIndex := strings.Index(input, "<")
	endIndex := strings.Index(input, ">")

	if startIndex != -1 && endIndex != -1 && startIndex < endIndex {
		return input[startIndex+1 : endIndex]
	}

	return strings.TrimSpace(input)
}

func joinTimes(start time.Time, end time.Time) string {
	return fmt.Sprintf("%s %s-%s",
		start.Format("02-01-2006"),
		start.Format("15:04"),
		end.Format("15:04"),
	)
}
