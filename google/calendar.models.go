package google

import "time"

type CalendarEventResponse struct {
	Kind             string              `json:"kind"`
	Etag             string              `json:"etag"`
	Summary          string              `json:"summary"`
	Description      string              `json:"description"`
	Updated          time.Time           `json:"updated"`
	TimeZone         string              `json:"timeZone"`
	AccessRole       string              `json:"accessRole"`
	DefaultReminders []DefaultReminder   `json:"defaultReminders"`
	Items            []CalendarEventItem `json:"items"`
}

type DefaultReminder struct {
	Method  string `json:"method"`
	Minutes int    `json:"minutes"`
}

type EventItemCreator struct {
	Email string `json:"email"`
	Self  bool   `json:"self"`
}
type EventItemOrganizer struct {
	Email string `json:"email"`
	Self  bool   `json:"self"`
}
type EventItemTime struct {
	DateTime time.Time `json:"dateTime"`
	TimeZone string    `json:"timeZone"`
}
type EventItemReminder struct {
	UseDefault bool `json:"useDefault"`
}
type EventItemAttendee struct {
	Email          string `json:"email"`
	Self           bool   `json:"self,omitempty"`
	ResponseStatus string `json:"responseStatus"`
	DisplayName    string `json:"displayName,omitempty"`
	Optional       bool   `json:"optional,omitempty"`
	Organizer      bool   `json:"organizer,omitempty"`
}

type EventItemConferenceData struct {
	EntryPoints        []EventItemEntryPoint       `json:"entryPoints"`
	ConferenceSolution EventItemConferenceSolution `json:"conferenceSolution"`
	ConferenceID       string                      `json:"conferenceId"`
}
type EventItemEntryPoint struct {
	EntryPointType string `json:"entryPointType"`
	URI            string `json:"uri"`
	Label          string `json:"label,omitempty"`
	Pin            string `json:"pin,omitempty"`
	RegionCode     string `json:"regionCode,omitempty"`
}
type EventItemConferenceSolution struct {
	Key     EventItemKey `json:"key"`
	Name    string       `json:"name"`
	IconURI string       `json:"iconUri"`
}
type EventItemKey struct {
	Type string `json:"type"`
}

type CalendarEventItem struct {
	Kind                  string                  `json:"kind"`
	Etag                  string                  `json:"etag"`
	ID                    string                  `json:"id"`
	Status                string                  `json:"status"`
	HTMLLink              string                  `json:"htmlLink"`
	Created               time.Time               `json:"created"`
	Updated               time.Time               `json:"updated"`
	Summary               string                  `json:"summary"`
	Creator               EventItemCreator        `json:"creator"`
	Organizer             EventItemOrganizer      `json:"organizer"`
	Start                 EventItemTime           `json:"start"`
	End                   EventItemTime           `json:"end"`
	RecurringEventID      string                  `json:"recurringEventId,omitempty"`
	OriginalStartTime     EventItemTime           `json:"originalStartTime,omitempty"`
	ICalUID               string                  `json:"iCalUID"`
	Sequence              int                     `json:"sequence"`
	Reminders             EventItemReminder       `json:"reminders"`
	EventType             string                  `json:"eventType"`
	Attendees             []EventItemAttendee     `json:"attendees,omitempty"`
	HangoutLink           string                  `json:"hangoutLink,omitempty"`
	ConferenceData        EventItemConferenceData `json:"conferenceData,omitempty"`
	Description           string                  `json:"description,omitempty"`
	Location              string                  `json:"location,omitempty"`
	GuestsCanInviteOthers bool                    `json:"guestsCanInviteOthers,omitempty"`
	PrivateCopy           bool                    `json:"privateCopy,omitempty"`
}
