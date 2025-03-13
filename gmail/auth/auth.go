package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/M1z23R/google-cli/gmail"
)

var (
	TokenUrl     = "https://oauth2.googleapis.com/token"
	ClientID     = os.Getenv("GOOGLE_CLIENT_ID")
	ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	CallbackUrl  = "http://localhost:1337/callback"
	BaseURL      = "https://accounts.google.com/o/oauth2/v2/auth/oauthchooseaccount"
)

func RefreshToken(profile *gmail.GmailProfile) error {
	qps := url.Values{}

	qps.Add("client_id", ClientID)
	qps.Add("client_secret", ClientSecret)
	qps.Add("refresh_token", profile.Tokens.RefreshToken)
	qps.Add("grant_type", "refresh_token")

	url := fmt.Sprintf("%s?%s", TokenUrl, qps.Encode())
	err := gmail.GmailApiCall("POST", url, nil, &profile.Tokens, profile)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}

	return nil
}

func TokensFromCode(code string, tokens *gmail.Tokens) error {
	params := url.Values{}

	params.Add("client_id", ClientID)
	params.Add("client_secret", ClientSecret)
	params.Add("redirect_uri", CallbackUrl)
	params.Add("grant_type", "authorization_code")

	url := fmt.Sprintf("%s?%s", TokenUrl, params.Encode())
	payload := map[string]string{}
	payload["code"] = code

	err := gmail.GmailApiCall("POST", url, payload, &tokens, nil)

	if err != nil {
		return err
	}

	return nil
}

func GenerateConsentUrl() string {
	params := url.Values{}
	params.Add("access_type", "offline")
	params.Add("client_id", ClientID)
	params.Add("redirect_uri", CallbackUrl)
	params.Add("response_type", "code")
	params.Add("service", "lso")
	params.Add("o2v", "2")
	params.Add("flowName", "GeneralOAuthFlow")
	scopes := []string{
		"https://mail.google.com",
	}
	params.Add("scope", strings.Join(scopes, " "))

	consentURL := fmt.Sprintf("%s?%s", BaseURL, params.Encode())
	return consentURL
}

func GetOAuthCode() (string, error) {
	if ClientID == "" || ClientSecret == "" {
		return "", errors.New("GOOGLE_CLIENT_ID and/or GOOGLE_CLIENT_SECRET aren't set.\nUse:\nexport GOOGLE_CLIENT_ID=\"your_client_id\"\nexport GOOGLE_CLIENT_SECRET=\"your_client_secret\" or set it in your .bashrc/.zshrc")
	}

	exec.Command("xdg-open", GenerateConsentUrl()).Start()
	GenerateConsentUrl()
	codeCh := make(chan string)
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		codeCh <- code
	})
	go http.ListenAndServe(":1337", nil)

	code := <-codeCh

	return code, nil
}

