package gmail

import (
	"fmt"
	"net/url"

	"github.com/M1z23R/google-cli/google"
)

const gmailApiUrl = "https://gmail.googleapis.com/gmail/v1/users/me/messages"

func ListMessages(profile *google.GoogleProfile, messagesList *google.GmailMessagesResponse, count int) error {
	values := url.Values{}
	values.Add("labelIds", "INBOX")
	values.Add("maxResults", fmt.Sprintf("%d", count))
	url := fmt.Sprintf("%s?%s", gmailApiUrl, values.Encode())

	err := google.GmailApiCall("GET", url, nil, &messagesList, profile)
	if err != nil {
		return err
	}

	return nil
}

func GetMessage(profile *google.GoogleProfile, message *google.GmailMessage, id string) error {
	url := fmt.Sprintf("%s/%s", gmailApiUrl, id)

	err := google.GmailApiCall("GET", url, nil, &message, profile)
	if err != nil {
		return err
	}

	return nil
}

func GetMessageMetadata(profile *google.GoogleProfile, message *google.GmailMessage, id string) error {
	values := url.Values{}
	values.Add("metadataHeaders", "From")
	values.Add("metadataHeaders", "Subject")
	url := fmt.Sprintf("%s/%s?%s", gmailApiUrl, id, values.Encode())

	err := google.GmailApiCall("GET", url, nil, &message, profile)
	if err != nil {
		return err
	}

	return nil
}

func GetUnreadCount(profile *google.GoogleProfile, count *int) error {
	values := url.Values{}
	values.Add("labelIds", "INBOX")
	values.Add("q", "is:unread")
	values.Add("maxResults", "20")
	url := fmt.Sprintf("%s?%s", gmailApiUrl, values.Encode())

	var messagesList google.GmailMessagesResponse
	err := google.GmailApiCall("GET", url, nil, &messagesList, profile)
	if err != nil {
		return err
	}
	*count = len(messagesList.Messages)
	return nil
}

func SendMessage(profile *google.GoogleProfile, from, to, subject, body, references, inReplyTo string) error {
	rawMessage, err := createMimeMessage(from, to, subject, body, references, inReplyTo)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/send", gmailApiUrl)

	payload := map[string]string{}
	payload["Raw"] = rawMessage

	var r google.GmailSendMessageResponse

	err = google.GmailApiCall("POST", url, payload, &r, profile)
	if err != nil {
		return err
	}
	return nil
}
