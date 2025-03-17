package calendar

import (
	"fmt"
	"net/url"
	"time"

	"github.com/M1z23R/google-cli/google"
)

const getEventsUrl = "https://www.googleapis.com/calendar/v3/calendars/primary/events"

func GetEvents(profile *google.GoogleProfile, eventsResponse *google.CalendarEventResponse, start time.Time, count int) error {
	timeMin := start.Format(time.RFC3339)
	values := url.Values{}
	values.Add("maxResults", fmt.Sprintf("%d", count))
	values.Add("calendarId", profile.EmailAddress)
	values.Add("timeMin", timeMin)
	values.Add("singleEvents", "true")
	values.Add("orderBy", "startTime")
	url := fmt.Sprintf("%s?%s", getEventsUrl, values.Encode())

	err := google.GmailApiCall("GET", url, nil, &eventsResponse, profile)
	if err != nil {
		return err
	}

	return nil
}
