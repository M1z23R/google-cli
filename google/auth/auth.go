package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/M1z23R/google-cli/google"
)

var (
	TokenUrl     = "https://oauth2.googleapis.com/token"
	ClientID     = os.Getenv("GOOGLE_CLIENT_ID")
	ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	CallbackUrl  = "http://localhost:1337/google/callback"
	BaseURL      = "https://accounts.google.com/o/oauth2/v2/auth/oauthchooseaccount"
)

func RefreshToken(profile *google.GoogleProfile) error {
	qps := url.Values{}

	qps.Add("client_id", ClientID)
	qps.Add("client_secret", ClientSecret)
	qps.Add("refresh_token", profile.Tokens.RefreshToken)
	qps.Add("grant_type", "refresh_token")

	url := fmt.Sprintf("%s?%s", TokenUrl, qps.Encode())
	err := google.GmailApiCall("POST", url, nil, &profile.Tokens, profile)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}

	fmt.Println("Refreshed Tokens")
	return nil
}

func TokensFromCode(code string, tokens *google.Tokens) error {
	params := url.Values{}

	params.Add("client_id", ClientID)
	params.Add("client_secret", ClientSecret)
	params.Add("redirect_uri", CallbackUrl)
	params.Add("grant_type", "authorization_code")

	url := fmt.Sprintf("%s?%s", TokenUrl, params.Encode())
	payload := map[string]string{}
	payload["code"] = code

	err := google.GmailApiCall("POST", url, payload, &tokens, nil)

	if err != nil {
		return err
	}

	return nil
}

func GenerateConsentUrl() string {
	params := url.Values{}
	params.Add("access_type", "offline")
	params.Add("prompt", "consent")
	params.Add("client_id", ClientID)
	params.Add("redirect_uri", CallbackUrl)
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

func GetOAuthCode() (string, error) {
	if ClientID == "" || ClientSecret == "" {
		return "", errors.New("GOOGLE_CLIENT_ID and/or GOOGLE_CLIENT_SECRET aren't set.\nUse:\nexport GOOGLE_CLIENT_ID=\"your_client_id\"\nexport GOOGLE_CLIENT_SECRET=\"your_client_secret\" or set it in your .bashrc/.zshrc")
	}

	cmd := exec.Command("xdg-open", GenerateConsentUrl())
	cmd.Start()
	err := cmd.Wait()

	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to run browser automatically, please visit the following manually:\n%s\n", GenerateConsentUrl()))
	}

	codeCh := make(chan string)
	http.HandleFunc("/google/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		codeCh <- code
	})
	go http.ListenAndServe(":1337", nil)

	code := <-codeCh

	return code, nil
}
