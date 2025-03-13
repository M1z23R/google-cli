package gmail

import "github.com/M1z23R/google-cli/google"

func ExtractSubject(m *google.GmailMessage) string {
	for _, v := range m.Payload.Headers {
		if v.Name == "Subject" {
			return v.Value
		}
	}

	return ""
}
