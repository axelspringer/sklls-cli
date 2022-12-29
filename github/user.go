package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type UserInfos struct {
	Type        string `json:"type" bson:"Type,omitempty"`
	UserName    string `json:"login" bson:"UserName,omitempty"`
	Email       string `json:"email" bson:"Email,omitempty"`
	AccessToken string `bson:"AccessToken,omitempty"`

	Name            string `json:"name" bson:"Name,omitempty"`
	Bio             string `json:"bio" bson:"Bio,omitempty"`
	TwitterUsername string `json:"twitter_username" bson:"TwitterUsername,omitempty"`
	Blog            string `json:"blog" bson:"Blog,omitempty"`

	Company              string `json:"email" bson:"Company,omitempty"`
	Location             string `json:"location" bson:"Location,omitempty"`
	Hireable             bool   `json:"hireable" bson:"Hireable,omitempty"`
	PublicRepoCount      int    `json:"public_repos" bson:"PublicRepoCount,omitempty"`
	PublicGistCount      int    `json:"public_gists" bson:"PublicGistCount,omitempty"`
	Followers            int    `json:"followers" bson:"Followers,omitempty"`
	Following            int    `json:"following" bson:"Following,omitempty"`
	CreatedGithubProfile string `json:"created_at" bson:"CreatedGithubProfile,omitempty"`
	UpdatedGithubProfile string `json:"updated_at" bson:"UpdatedGithubProfile,omitempty"`
}

func GetUserInfo(accessToken string) (*UserInfos, error) {
	url := fmt.Sprintf("https://api.github.com/user")

	// Create profile & parse response
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawJsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := UserInfos{}
	errRes := TokenErrorResponse{}
	err = json.Unmarshal([]byte(rawJsonBody), &res)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(rawJsonBody), &errRes)
	if err != nil {
		return nil, err
	}

	if errRes.Error != "" {
		return nil, errors.New(fmt.Sprintf("%s - %s", errRes.Error, errRes.ErrorDescription))
	}

	return &res, nil
}
