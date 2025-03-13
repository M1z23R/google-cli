package calendar

import (
	"fmt"
	"net/url"
	"time"

	"github.com/M1z23R/google-cli/google"
)

const getEventsUrl = "https://www.googleapis.com/calendar/v3/calendars/primary/events"

func GetNextEvent(profile *google.GoogleProfile, eventsResponse *google.CalendarEventResponse, start time.Time) error {
	timeMin := start.Format(time.RFC3339)

	values := url.Values{}
	values.Add("maxResults", "1")
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

func GetEvents(profile *google.GoogleProfile, eventsResponse *google.CalendarEventResponse, start time.Time, end time.Time) error {
	timeMin := start.Format(time.RFC3339)
	timeMax := end.Format(time.RFC3339)

	values := url.Values{}
	values.Add("calendarId", profile.EmailAddress)
	values.Add("timeMin", timeMin)
	values.Add("timeMax", timeMax)
	values.Add("singleEvents", "true")
	values.Add("orderBy", "startTime")
	url := fmt.Sprintf("%s?%s", getEventsUrl, values.Encode())

	err := google.GmailApiCall("GET", url, nil, &eventsResponse, profile)
	if err != nil {
		return err
	}

	return nil
}
