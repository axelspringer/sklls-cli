package analyzer

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const skllsAnalysisJson = `{
	"git@github.com:spring-media/mock-project.git": {
		"jonas.peeck@axelspringer.com": {
			"Ext": {
				".groovy": {
					"Mon, 01 Nov 2021 14:48:58 +0100": 13
				}
			},
			"Dep": {
				"EcmaScript": {
					"@apollo/client": {
						"Thu, 26 Aug 2021 10:55:55 +0200": 8
					}
				}
			}
		}
	}
}`

const mergerJson = `{
	"git@github.com:spring-media/mock-project.git": {
		"jonas.peeck@axelspringer.com": {
			"Ext": {
				".groovy": {
					"Mon, 01 Nov 2021 14:48:58 +0100": 2
				}
			},
			"Dep": {
				"EcmaScript": {
					"@apollo/client": {
						"Thu, 26 Aug 2021 10:55:55 +0200": 7
					}
				}
			}
		}
	}
}`

func TestAddExtCount(t *testing.T) {
	sa := SkllsAnalysis{}
	err := json.Unmarshal([]byte(skllsAnalysisJson), &sa)
	if err != nil {
		panic(err)
	}

	const origin = Origin("git@github.com:spring-media/mock-project.git")
	const contributorEmail = ContributorId("jonas.peeck@axelspringer.com")
	const ext = Ext(".groovy")
	const dateStr = DateStr("Mon, 01 Nov 2021 14:48:58 +0100")

	sa.AddExtCount(origin, contributorEmail, ext, dateStr, 2)
	assert.Equal(t, LineCount(15), sa[origin][contributorEmail].Ext[ext][dateStr], "Line count was increased by the provided value")
}

func TestAddDepCount(t *testing.T) {
	sa := SkllsAnalysis{}
	err := json.Unmarshal([]byte(skllsAnalysisJson), &sa)
	if err != nil {
		panic(err)
	}

	const origin = Origin("git@github.com:spring-media/mock-project.git")
	const contributorEmail = ContributorId("jonas.peeck@axelspringer.com")
	const parserName = ParserName("EcmaScript")
	const dep = Dependency("@apollo/client")
	const dateStr = DateStr("Thu, 26 Aug 2021 10:55:55 +0200")

	sa.AddDepCount(origin, contributorEmail, parserName, dep, dateStr, 7)
	assert.Equal(t, LineCount(15), sa[origin][contributorEmail].Dep[parserName][dep][dateStr], "Line count was increased by the provided value")
}

func TestAddEmail(t *testing.T) {
	sa := SkllsAnalysis{}
	err := json.Unmarshal([]byte(skllsAnalysisJson), &sa)
	if err != nil {
		panic(err)
	}

	const origin = Origin("git@github.com:spring-media/mock-project.git")
	const contributor = ContributorId("jonas.peeck@axelspringer.com")

	sa.AddUsername(origin, contributor, "funny@sunny.com")
	sa.AddUsername(origin, contributor, "funny@sunny.com")
	sa.AddUsername(origin, contributor, "hello@hi.com")
	assert.Equal(t, []string{"funny@sunny.com", "hello@hi.com"}, sa[origin][contributor].Usernames, "Emails were correctly added")
}

func TestAddToEmailList(t *testing.T) {
	s := Sklls{}
	s = s.AddUsernameToList("funny@sunny.com")
	s = s.AddUsernameToList("funny@sunny.com")
	s = s.AddUsernameToList("hello@hi.com")

	assert.Equal(t, []string{"funny@sunny.com", "hello@hi.com"}, s.Usernames, "Adds emails in a deduplicated way")
}

func TestMerge(t *testing.T) {
	sa := SkllsAnalysis{}
	err := json.Unmarshal([]byte(skllsAnalysisJson), &sa)
	if err != nil {
		panic(err)
	}

	merger := SkllsAnalysis{}
	err = json.Unmarshal([]byte(mergerJson), &merger)
	if err != nil {
		panic(err)
	}

	sa.Merge(&merger)

	// Ensure sklls count is correct
	const origin = Origin("git@github.com:spring-media/mock-project.git")
	const contributorEmail = ContributorId("jonas.peeck@axelspringer.com")
	const ext = Ext(".groovy")
	const dateStrA = DateStr("Mon, 01 Nov 2021 14:48:58 +0100")
	assert.Equal(t, LineCount(15), sa[origin][contributorEmail].Ext[ext][dateStrA], "Line count (sklls) was increased by the provided value")

	// Ensure dependency count is correct
	const parserName = ParserName("EcmaScript")
	const dep = Dependency("@apollo/client")
	const dateStrB = DateStr("Thu, 26 Aug 2021 10:55:55 +0200")
	assert.Equal(t, LineCount(15), sa[origin][contributorEmail].Dep[parserName][dep][dateStrB], "Line count (dep) was increased by the provided value")

}

func TestOriginToFilename(t *testing.T) {
	assert.Equal(t, "github.com_spring-media_ps-zoom-enhance", OriginToFilename("git@github.com:spring-media/ps-zoom-enhance.git"), "Converts origin into a useable filename")
}
