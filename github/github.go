package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GhRepository struct {
	Name   string `json:"name"`
	IsFork bool   `json:"fork"`
}

func GetAllGhOrgRepos(orgName string, pat string) ([]string, error) {
	var ghRepositories []GhRepository
	pageNo := 1

	for {
		repos, err := queryRepos(orgName, pat, pageNo)
		if err != nil {
			return nil, err
		}

		if len(repos) == 0 {
			break
		}

		ghRepositories = append(ghRepositories, repos...)
		fmt.Printf("Page %d - Found %d repos (%d)...\n", pageNo, len(repos), len(ghRepositories))

		pageNo = pageNo + 1
	}

	// Get all repo names
	repoNames := []string{}
	for _, ghRepo := range ghRepositories {
		if ghRepo.IsFork {
			continue
		}
		repoNames = append(repoNames, ghRepo.Name)
	}

	return repoNames, nil
}

func queryRepos(orgName string, pat string, page int) ([]GhRepository, error) {
	// Retrieve all repos from this github org
	client := http.Client{}
	url := fmt.Sprintf("https://api.github.com/orgs/%s/repos?type=all&sort=full_name&per_page=100&page=%d", orgName, page)
	// url := fmt.Sprintf("https://api.github.com/user/repos?type=private&access_token=%s&page=%d", pat, page)
	username := "jpeeck-spring"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Basic "+basicAuth(username, pat))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	// retrieve body from response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// parse body to json
	var ghRepositories []GhRepository
	json.Unmarshal(body, &ghRepositories)

	return ghRepositories, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
