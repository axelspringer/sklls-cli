package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type TokenErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorUri         string `json:"error_uri"`
}

func GetAccessToken(code string, ghClientId string, ghClientSecret string) (string, error) {
	url := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", ghClientId, ghClientSecret, code)

	// Create profile & parse response
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	rawJsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	res := TokenResponse{}
	errRes := TokenErrorResponse{}
	err = json.Unmarshal([]byte(rawJsonBody), &res)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal([]byte(rawJsonBody), &errRes)
	if err != nil {
		return "", err
	}

	if errRes.Error != "" {
		return "", errors.New(fmt.Sprintf("%s - %s", errRes.Error, errRes.ErrorDescription))
	}

	return res.AccessToken, nil
}
