package main

import (
	"fmt"
	"time"

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

	gmail.MessagesList(profile, &r)

	for i, v := range r.Messages {
		var m google.GmailMessage
		gmail.MessagesGet(profile, &m, v.Id)
		fmt.Printf("%d) %s\n", i, gmail.ExtractSubject(&m))
	}
}

func printCalendar(profile *google.GoogleProfile) {
	var r google.CalendarEventResponse
	now := time.Now()
	tomorrow := time.Now().Add(24 * time.Hour)
	start := now
	end := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())

	auth.RefreshToken(profile)
	persistence.UpsertProfile(profile)
	calendar.GetEvents(profile, &r, start, end)
	for _, v := range r.Items {
		fmt.Printf("%s | %s\n", joinTimes(v.Start.DateTime, v.End.DateTime), v.Summary)
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
