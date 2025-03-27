package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	"github.com/M1z23R/google-cli/google"
	"github.com/M1z23R/google-cli/google/profiles"
)

const gmailApiUrl = "https://gmail.googleapis.com/gmail/v1"

func RefreshToken(profile *google.GoogleProfile) error {
	qps := url.Values{}

	qps.Add("client_id", profile.Secrets.ClientId)
	qps.Add("client_secret", profile.Secrets.ClientSecret)
	qps.Add("refresh_token", profile.Tokens.RefreshToken)
	qps.Add("grant_type", "refresh_token")

	url := fmt.Sprintf("%s?%s", TokenUrl, qps.Encode())
	err := google.GmailApiCall("POST", url, nil, &profile.Tokens, profile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func TokensFromCode(code string, profile *google.GoogleProfile) error {
	params := url.Values{}

	params.Add("client_id", profile.Secrets.ClientId)
	params.Add("client_secret", profile.Secrets.ClientSecret)
	params.Add("redirect_uri", profile.Secrets.RedirectUri)
	params.Add("grant_type", "authorization_code")

	url := fmt.Sprintf("%s?%s", TokenUrl, params.Encode())
	payload := map[string]string{}
	payload["code"] = code

	err := google.GmailApiCall("POST", url, payload, &profile.Tokens, nil)
	return err
}

func GenerateConsentUrl(secrets *google.GoogleSecret) string {
	params := url.Values{}
	params.Add("access_type", "offline")
	params.Add("prompt", "consent")
	params.Add("client_id", secrets.ClientId)
	params.Add("redirect_uri", secrets.RedirectUri)
	params.Add("response_type", "code")
	params.Add("service", "lso")
	params.Add("o2v", "2")
	params.Add("flowName", "GeneralOAuthFlow")
	scopes := []string{
		"openid",
		"email",
		"profile",
		"https://mail.google.com",
		"https://www.googleapis.com/auth/calendar",
	}
	params.Add("scope", strings.Join(scopes, " "))

	consentURL := fmt.Sprintf("%s?%s", BaseURL, params.Encode())
	return consentURL
}

func AddProfile(profile *google.GoogleProfile) error {
	code, err := GetOAuthCode(&profile.Secrets)
	if err != nil {
		return err
	}

	err = TokensFromCode(code, profile)
	if err != nil {
		return err
	}

	err = profiles.GetProfile(profile)

	return err
}

func GetOAuthCode(secrets *google.GoogleSecret) (string, error) {
	cmd := exec.Command("xdg-open", GenerateConsentUrl(secrets))
	cmd.Start()
	err := cmd.Wait()

	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to run browser automatically, please visit the following manually:\n%s\n", GenerateConsentUrl(secrets)))
	}

	codeCh := make(chan string)
	CallbackResponse(codeCh)
	go http.ListenAndServe(":1337", nil)

	code := <-codeCh

	return code, nil
}
