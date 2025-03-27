package google

import (
	"time"
)

type GoogleSecret struct {
	ID           int    `json:"id"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RedirectUri  string `json:"redirectUri"`
}

type GoogleProfile struct {
	ID            int          `json:"id"`
	EmailAddress  string       `json:"emailAddress"`
	Tokens        Tokens       `json:"tokens"`
	Secrets       GoogleSecret `json:"secrets"`
	LastUpdatedAt time.Time    `json:"timestamp"`
}

type Tokens struct {
	AccessToken  string    `json:"access_token"`
	ExpiresIn    int       `json:"expires_in"`
	IdToken      string    `json:"id_token"`
	RefreshToken string    `json:"refresh_token"`
	Scope        string    `json:"scope"`
	TokenType    string    `json:"token_type"`
	Expires      time.Time `json:"expires"`
	LastUpdated  time.Time `json:"last_updated"`
}
