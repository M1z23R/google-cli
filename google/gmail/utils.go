package gmail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/M1z23R/google-cli/google"
)

func ExtractHeader(m *google.GmailMessage, h string) string {
	for _, v := range m.Payload.Headers {
		if v.Name == h {
			return v.Value
		}
	}

	return ""
}

func createMimeMessage(from, to, subject, body, references, inReplyTo string) (string, error) {
	var mimeMessage bytes.Buffer

	correctedSubject := subject
	mimeMessage.WriteString("MIME-Version: 1.0\r\n")
	mimeMessage.WriteString("From: " + from + "\r\n")
	mimeMessage.WriteString("To: " + to + "\r\n")
	if references != "" && inReplyTo != "" {
		mimeMessage.WriteString("References: " + references + "\r\n")
		mimeMessage.WriteString("In-Reply-To: " + inReplyTo + "\r\n")
		correctedSubject = fmt.Sprintf("RE:%s", subject)
	}
	mimeMessage.WriteString("Subject: " + correctedSubject + "\r\n")
	mimeMessage.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	mimeMessage.WriteString("\r\n")

	mimeMessage.WriteString(body)

	encoded := base64.StdEncoding.EncodeToString(mimeMessage.Bytes())

	encoded = strings.Replace(encoded, "+", "-", -1)
	encoded = strings.Replace(encoded, "/", "_", -1)
	encoded = strings.TrimRight(encoded, "=")

	return encoded, nil
}
