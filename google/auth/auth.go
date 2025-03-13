package auth

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	TokenUrl     = "https://oauth2.googleapis.com/token"
	ClientID     string
	ClientSecret string
	CallbackUrl  = "http://localhost:1337/google/callback"
	BaseURL      = "https://accounts.google.com/o/oauth2/v2/auth/oauthchooseaccount"
)

func init() {
	godotenv.Load()
	ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
}

