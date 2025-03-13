package profiles

import "github.com/M1z23R/google-cli/gmail"

const getProfileUrl = "https://gmail.googleapis.com/gmail/v1/users/me/profile"

func GetProfile(profile *gmail.GmailProfile) error {
	err := gmail.GmailApiCall("GET", getProfileUrl, nil, &profile, profile)
	if err != nil {
		return err
	}

	return nil
}
