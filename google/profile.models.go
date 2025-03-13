package google

import (
	"time"
)

type GoogleProfile struct {
	ID            int       `json:"id"`
	EmailAddress  string    `json:"emailAddress"`
	Tokens        Tokens    `json:"tokens"`
	LastUpdatedAt time.Time `json:"timestamp"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

