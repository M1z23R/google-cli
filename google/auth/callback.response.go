package auth

import (
	"net/http"
)

func CallbackResponse(codeCh chan string) {
	http.HandleFunc("/google/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		htmlResponse := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Authentication Successful</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					text-align: center;
					margin-top: 50px;
				}
				.container {
					max-width: 400px;
					margin: auto;
					padding: 20px;
					box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
					border-radius: 8px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h2>Logged in successfully</h2>
				<p>You may close this tab.</p>
			</div>
		</body>
		</html>`

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlResponse))

		codeCh <- code
	})
}
