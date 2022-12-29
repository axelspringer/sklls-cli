package repo

import "sklls/github"

func NewFromGithub(ghorg string, reponame string, outdir string, ghpat string) (*Repo, error) {
	_, err := github.CloneFromGithub(ghorg, reponame, outdir, ghpat)
	if err != nil {
		return nil, err
	}

	repo := NewFromLocal(outdir)
	return repo, nil
}
