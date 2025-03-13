package auth

import (
	"errors"

	"github.com/M1z23R/google-cli/gmail"
	"github.com/M1z23R/google-cli/gmail/profiles"
)

const gmailApiUrl = "https://gmail.googleapis.com/gmail/v1"

func AddProfile(profile *gmail.GmailProfile) error {
	code, err := GetOAuthCode()
	if err != nil {
		return err 
	}

	err = TokensFromCode(code, &profile.Tokens)
	if err != nil {
		return errors.New("error occured while getting tokens from code.")
	}

	err = profiles.GetProfile(profile)

	if err != nil {
		return errors.New("error while getting profile")
	}
	return nil
}
