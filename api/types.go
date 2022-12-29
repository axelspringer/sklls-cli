package api

import (
	"encoding/json"
	"sklls/analyzer"
)

type ProfileCreationRequest struct {
	ApiKey  string
	Profile SkllsProfile
}

type ProfileUpdateRequest struct {
	Id           string
	EditToken    string
	HiddenSkills []string
	HiddenLibs   []string
}

type SkllsProfile struct {
	Id            string
	EditToken     string
	Email         string
	Profile       string
	Sklls         Developer
	SkllsAnalysis analyzer.SkllsAnalysis
}

type Developer struct {
	Id              string
	Emails          []string
	UserNames       []string
	OrgRepoNames    []string
	SkillsLineCount map[string]int
	SkillLog        map[string]map[string]int
	LibsLineCount   map[string]int
	LibsLog         map[string]map[string]int
	LastCommitDate  string
	Skills          []Skill
	Libraries       []Skill

	HiddenSkills []string
	HiddenLibs   []string
}

type Skill struct {
	Skill  string  `json:"skill"`
	Factor float64 `json:"factor"`
}

func NewSkllsProfileFromJson(jsonString string) (SkllsProfile, error) {
	var profile SkllsProfile
	err := json.Unmarshal([]byte(jsonString), &profile)
	return profile, err
}

type ProfileCreationResponse struct {
	StatusCode int
	Success    bool
	Message    string
	Body       string
}

type JoinWaitingListRequest struct {
	Name      string
	WorkEmail string
	Company   string
	City      string
	Country   string
	Reason    string
}

type TrackingEventCreationRequest struct {
	SessionId string
	ApiKey    string
	Type      string
	UserId    string
	MetaJson  string
}
