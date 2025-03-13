package google

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GmailApiCall(method string, url string, requestBody any, responseBody any, profile *GoogleProfile) error {
	payload, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	var req *http.Request

	if requestBody == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(payload))
	}

	if err != nil {
		return err
	}

	req.Header.Set("content-type", "application/json")
	if profile != nil && profile.Tokens.AccessToken != "" {
		req.Header.Set("authorization", fmt.Sprintf("Bearer %s", profile.Tokens.AccessToken))
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, responseBody)

	if err != nil {
		return err
	}

	return nil
}
