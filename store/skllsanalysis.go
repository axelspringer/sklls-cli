package store

import (
	"encoding/json"
	"sklls/analyzer"
)

type ProfileType string

const UserProfileReadOnlyType = ProfileType("UserProfileReadOnly")
const UserProfileEditType = ProfileType("UserProfileEdit")

type UserProfile struct {
	Type           ProfileType
	Id             string
	Email          string
	EditToken      string
	Profile        string
	HiddenOrigins  []string
	HiddenExts     []string
	HiddenDeps     []string
	SelectedEmails []string
	Ext            analyzer.ExtMap        `json:",omitempty"`
	Dep            analyzer.DepMap        `json:",omitempty"`
	SkllsAnalysis  analyzer.SkllsAnalysis `json:",omitempty"`
}

func (profile *UserProfile) FromJson(rawJson []byte) error {
	err := json.Unmarshal(rawJson, profile)
	return err
}

func (profile *UserProfile) ToJson() ([]byte, error) {
	return json.Marshal(profile)
}
