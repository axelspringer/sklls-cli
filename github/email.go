package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GithubUserEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

// The email provided by getGithubUserInfo() is the public email, that's why we have to make a separate call to retrieve the actually usable email
func GetPrimaryAccountEmail(accessToken string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/user/emails")

	// Create profile & parse response
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

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

	res := []GithubUserEmail{}
	err = json.Unmarshal([]byte(rawJsonBody), &res)
	if err != nil {
		return "", err
	}

	email := ""
	for _, ghUserEmailEntry := range res {
		if ghUserEmailEntry.Primary {
			email = ghUserEmailEntry.Email
		}
	}

	if email == "" {
		return "", errors.New("Cannot find primary email address for user (probably authentication issue with Github)")
	}

	return email, nil
}
