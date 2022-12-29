package analyzer

import (
	"fmt"
	"regexp"
	"sklls/helper"
)

type Origin string
type ContributorId string
type Ext string
type ParserName string
type Dependency string
type DateStr string
type LineCount int

type SkllsAnalysis map[Origin]map[ContributorId]Sklls

func (skllsAnalysis *SkllsAnalysis) AddExtCount(origin Origin, contributorId ContributorId, ext Ext, dateStr DateStr, lineCount LineCount) {
	if _, ok := (*skllsAnalysis)[origin]; !ok {
		(*skllsAnalysis)[origin] = map[ContributorId]Sklls{}
	}

	if _, ok := (*skllsAnalysis)[origin][contributorId]; !ok {
		(*skllsAnalysis)[origin][contributorId] = Sklls{
			Ext: map[Ext]map[DateStr]LineCount{},
			Dep: map[ParserName]map[Dependency]map[DateStr]LineCount{},
		}
	}

	(*skllsAnalysis)[origin][contributorId] = (*skllsAnalysis)[origin][contributorId].AddExt(ext, dateStr, lineCount)
}

func (skllsAnalysis *SkllsAnalysis) AddDepCount(origin Origin, contributorId ContributorId, parserName ParserName, dependency Dependency, dateStr DateStr, lineCount LineCount) {
	if _, ok := (*skllsAnalysis)[origin]; !ok {
		(*skllsAnalysis)[origin] = map[ContributorId]Sklls{}
	}

	if _, ok := (*skllsAnalysis)[origin][contributorId]; !ok {
		(*skllsAnalysis)[origin][contributorId] = Sklls{
			Ext: map[Ext]map[DateStr]LineCount{},
			Dep: map[ParserName]map[Dependency]map[DateStr]LineCount{},
		}
	}

	(*skllsAnalysis)[origin][contributorId] = (*skllsAnalysis)[origin][contributorId].AddDep(parserName, dependency, dateStr, lineCount)
}

func (skllsAnalysis *SkllsAnalysis) AddUsername(origin Origin, contributorId ContributorId, username string) {
	if _, ok := (*skllsAnalysis)[origin]; !ok {
		(*skllsAnalysis)[origin] = map[ContributorId]Sklls{}
	}

	if _, ok := (*skllsAnalysis)[origin][contributorId]; !ok {
		(*skllsAnalysis)[origin][contributorId] = Sklls{
			Ext: map[Ext]map[DateStr]LineCount{},
			Dep: map[ParserName]map[Dependency]map[DateStr]LineCount{},
		}
	}

	(*skllsAnalysis)[origin][contributorId] = (*skllsAnalysis)[origin][contributorId].AddUsernameToList(username)
}

func (skllsAnalysis *SkllsAnalysis) Merge(src *SkllsAnalysis) {
	for origin, _ := range *src {
		for contributorEmail, _ := range (*src)[origin] {
			// Go through exts
			for ext, _ := range (*src)[origin][contributorEmail].Ext {
				for dateStr, lineCount := range (*src)[origin][contributorEmail].Ext[ext] {
					skllsAnalysis.AddExtCount(origin, contributorEmail, ext, dateStr, lineCount)
				}
			}

			// Go through deps
			for parserName, _ := range (*src)[origin][contributorEmail].Dep {
				for dependency, _ := range (*src)[origin][contributorEmail].Dep[parserName] {
					for dateStr, lineCount := range (*src)[origin][contributorEmail].Dep[parserName][dependency] {
						skllsAnalysis.AddDepCount(origin, contributorEmail, parserName, dependency, dateStr, lineCount)
					}
				}
			}

			// Go through emails
			for _, email := range (*src)[origin][contributorEmail].Usernames {
				skllsAnalysis.AddUsername(origin, contributorEmail, email)
			}
		}
	}
}

type Sklls struct {
	Ext       ExtMap
	Dep       DepMap
	Usernames []string
}

type ExtMap map[Ext]map[DateStr]LineCount
type DepMap map[ParserName]map[Dependency]map[DateStr]LineCount

func (sklls Sklls) AddExt(ext Ext, dateStr DateStr, lineCount LineCount) Sklls {
	if sklls.Ext == nil {
		sklls.Ext = map[Ext]map[DateStr]LineCount{}
	}

	if _, ok := sklls.Ext[ext]; !ok {
		sklls.Ext[ext] = map[DateStr]LineCount{}
	}

	if _, ok := sklls.Ext[ext]; !ok {
		sklls.Ext[ext] = map[DateStr]LineCount{}
	}

	if _, ok := sklls.Ext[ext][dateStr]; !ok {
		sklls.Ext[ext][dateStr] = lineCount
	} else {
		sklls.Ext[ext][dateStr] += lineCount
	}

	return sklls
}

func (sklls Sklls) AddDep(parserName ParserName, dep Dependency, dateStr DateStr, lineCount LineCount) Sklls {
	if sklls.Dep == nil {
		sklls.Dep = map[ParserName]map[Dependency]map[DateStr]LineCount{}
	}

	if _, ok := sklls.Dep[parserName]; !ok {
		sklls.Dep[parserName] = map[Dependency]map[DateStr]LineCount{}
	}

	if _, ok := sklls.Dep[parserName][dep]; !ok {
		sklls.Dep[parserName][dep] = map[DateStr]LineCount{}
	}

	if _, ok := sklls.Dep[parserName][dep][dateStr]; !ok {
		sklls.Dep[parserName][dep][dateStr] = lineCount
	} else {
		sklls.Dep[parserName][dep][dateStr] += lineCount
	}

	return sklls
}

func (sklls Sklls) AddUsernameToList(username string) Sklls {
	if sklls.Usernames == nil {
		sklls.Usernames = []string{}
	}

	if helper.Contains(sklls.Usernames, username) {
		return sklls
	}

	sklls.Usernames = append(sklls.Usernames, username)
	return sklls
}

type OriginDeconstructed struct {
	Scm  string
	Org  string
	Repo string
}

func ParseOrigin(origin Origin) (*OriginDeconstructed, error) {
	gitOriginMatcher := regexp.MustCompile(`.*@(.*):(.*)\/(.*).git`)
	httpsOriginMatcher := regexp.MustCompile(`https://.*@(.*)\/(.*)\/(.*).git`)

	matches := gitOriginMatcher.FindStringSubmatch(string(origin))

	// If no matches were found, try the httpsOriginMatcher
	if len(matches) == 0 {
		matches = httpsOriginMatcher.FindStringSubmatch(string(origin))
	}

	// Still no matches? Seems like an invalid origin string
	if len(matches) == 0 {
		return nil, fmt.Errorf("Cannot parse origin '%s' (regex found no match)\n", origin)
	}

	scm := matches[1]
	org := matches[2]
	repo := matches[3]

	return &OriginDeconstructed{scm, org, repo}, nil
}
