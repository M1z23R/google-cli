package auth

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	TokenUrl     = "https://oauth2.googleapis.com/token"
	ClientID     string
	ClientSecret string
	CallbackUrl  = "http://localhost:1337/microsoft/callback"
	BaseURL      = "https://accounts.google.com/o/oauth2/v2/auth/oauthchooseaccount"
)

func init() {
	godotenv.Load()
	ClientID = os.Getenv("MICROSOFT_CLIENT_ID")
	ClientSecret = os.Getenv("MICROSOFT_CLIENT_SECRET")
}
