package main

import (
	"fmt"
	"time"

	clicolors "github.com/M1z23R/google-cli/cli-colors"
	"github.com/M1z23R/google-cli/google"
	"github.com/M1z23R/google-cli/google/auth"
	"github.com/M1z23R/google-cli/google/calendar"
	"github.com/M1z23R/google-cli/google/gmail"
	"github.com/M1z23R/google-cli/persistence"
)

func printProfiles() {
	profiles := []google.GoogleProfile{}
	err := persistence.GetProfiles(&profiles)
	if err != nil {
		fmt.Println("Couldn't load profiles")
	}

	if len(profiles) == 0 {
		fmt.Println("No profiles are added.")
	} else {
		for i, profile := range profiles {
			fmt.Printf("%d) %s\n", i+1, profile.EmailAddress)
		}
	}

}

func printMessages(profile *google.GoogleProfile) {
	var r google.GmailMessagesResponse

	auth.RefreshToken(profile)
	persistence.UpsertProfile(profile)

	gmail.MessagesList(profile, &r, 1)

	for _, v := range r.Messages {
		var m google.GmailMessage
		var c int
		gmail.MessagesGetMetadata(profile, &m, v.Id)
		gmail.MessagesUnreadCount(profile, &c)
		lastMessageSubject := fmt.Sprintf("%s%s%s", clicolors.Blue, gmail.ExtractHeader(&m, "Subject"), clicolors.Reset)
		lastMessageFrom := fmt.Sprintf("%s%s%s", clicolors.Cyan, gmail.ExtractHeader(&m, "From"), clicolors.Reset)
		lastMessage := fmt.Sprintf("%s | %s\n", lastMessageFrom, lastMessageSubject)
		unreadCount := fmt.Sprintf("%sUnread: %d%s", clicolors.GetUnreadColor(c), c, clicolors.Reset)
		fmt.Printf("%s | %s\n", unreadCount, lastMessage)
	}
}

func printCalendar(profile *google.GoogleProfile) {
	var r google.CalendarEventResponse
	start := time.Now()

	auth.RefreshToken(profile)
	persistence.UpsertProfile(profile)
	calendar.GetNextEvent(profile, &r, start)
	for _, v := range r.Items {
		fmt.Printf("%s%s | %s%s\n", clicolors.GetEventColor(v.Start.DateTime), v.Summary, joinTimes(v.Start.DateTime, v.End.DateTime), clicolors.Reset)
	}
}

func joinTimes(start time.Time, end time.Time) string {
	return fmt.Sprintf("%s %s-%s",
		start.Format("2006-01-02"),
		start.Format("15:04"),
		end.Format("15:04"),
	)
}

func addGoogleAccount() {
	profile := google.GoogleProfile{}
	err := auth.AddProfile(&profile)
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
