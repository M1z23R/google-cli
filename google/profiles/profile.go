package profiles

import (
	"fmt"

	"github.com/M1z23R/google-cli/google"
)

const getProfileUrl = "https://gmail.googleapis.com/gmail/v1/users/me/profile"

func GetProfile(profile *google.GoogleProfile) error {
	err := google.GmailApiCall("GET", getProfileUrl, nil, &profile, profile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
