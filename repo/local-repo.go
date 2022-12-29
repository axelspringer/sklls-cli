package repo

import (
	"path"
	"sklls/git"
)

func NewFromLocal(localpath string) *Repo {
	ghorg := ""
	reponame := path.Base(localpath)
	originUrl, _ := git.GetOrigin(localpath)

	repo := Repo{
		GhOrg:     ghorg,
		Reponame:  reponame,
		Repopath:  localpath,
		Files:     []string{},
		Blames:    map[string][]git.BlameLine{},
		OriginUrl: originUrl,
	}
	return &repo
}
