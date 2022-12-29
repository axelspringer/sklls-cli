package repo

import (
	"sklls/git"
	"strings"
)

type Repo struct {
	GhOrg     string
	Reponame  string
	Repopath  string
	Files     []string
	Blames    map[string][]git.BlameLine
	OriginUrl string
}

func (repo *Repo) LoadFiles(excludeSuffix []string) (filesLoaded int, err error) {
	filesLoaded = 0
	files, err := git.LoadFiles(repo.Repopath)
	if err != nil {
		return filesLoaded, err
	}

	var filteredFiles []string
	for _, file := range files {
		includeFile := true
		for _, suffix := range excludeSuffix {
			if strings.HasSuffix(file, suffix) {
				includeFile = false
			}
		}

		if includeFile {
			filteredFiles = append(filteredFiles, file)
		}
	}

	repo.Files = filteredFiles
	filesLoaded = len(filteredFiles)
	return filesLoaded, nil
}

func LoadBlame(repoPath string, filePath string) ([]git.BlameLine, error) {
	blames, err := git.Blame(repoPath, filePath)
	if err != nil {
		return nil, err
	}

	return blames, nil
}
