package gmail

import "github.com/M1z23R/google-cli/google"

func ExtractHeader(m *google.GmailMessage, h string) string {
	for _, v := range m.Payload.Headers {
		if v.Name == h {
			return v.Value
		}
	}

	return ""
}
