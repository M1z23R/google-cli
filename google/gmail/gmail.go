package gmail

import (
	"fmt"
	"net/url"

	"github.com/M1z23R/google-cli/google"
)

const gmailApiUrl = "https://gmail.googleapis.com/gmail/v1/users/me/messages"

func MessagesList(profile *google.GoogleProfile, messagesList *google.GmailMessagesResponse) error {
	values := url.Values{}
	values.Add("labelIds", "INBOX")
	values.Add("maxResults", fmt.Sprintf("%d", 10))
	url := fmt.Sprintf("%s?%s", gmailApiUrl, values.Encode())

	err := google.GmailApiCall("GET", url, nil, &messagesList, profile)
	if err != nil {
		return err
	}

	return nil
}

func MessagesGet(profile *google.GoogleProfile, message *google.GmailMessage, id string) error {
	url := fmt.Sprintf("%s/%s", gmailApiUrl, id)

	err := google.GmailApiCall("GET", url, nil, &message, profile)
	if err != nil {
		return err
	}

	return nil
}

func MessagesGetMetadata(profile *google.GoogleProfile, message *google.GmailMessage, id string) error {
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
