package main

import (
	"fmt"

	"github.com/M1z23R/google-cli/gmail"
	"github.com/M1z23R/google-cli/gmail/auth"
)

func main() {
	profile := gmail.GmailProfile{}
	err := auth.AddProfile(&profile)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Println(profile)
}
