package gmail

import (
	"time"
)

type GmailProfile struct {
	ID            int       `json:"id"`
	EmailAddress  string    `json:"emailAddress"`
	MessagesTotal int       `json:"messagesTotal"`
	ThreadsTotal  int       `json:"threadsTotal"`
	HistoryId     string    `json:"historyId"`
	Tokens        Tokens    `json:"tokens"`
	LastUpdatedAt time.Time `json:"timestamp"`
	UnreadCount   int       `json:"unreadCount"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

