package repo

import (
	"fmt"
	"os"
	"path"
	"testing"
)

func BenchmarkNewFromLocal(b *testing.B) {
	cwd, err := os.Getwd()
	if err != nil {
		b.Fatalf(`Could not get cwd:\n%s`, err)
		return
	}

	repoPath := path.Join(cwd, "../../")
	for i := 0; i < b.N; i++ {
		NewFromLocal(repoPath)
	}
}

func BenchmarkLoadFiles(b *testing.B) {
	cwd, err := os.Getwd()
	if err != nil {
		b.Fatalf(`Could not get cwd:\n%s`, err)
		return
	}

	repoPath := path.Join(cwd, "../../")
	filesLoaded := 0
	for i := 0; i < b.N; i++ {
		repo := NewFromLocal(repoPath)
		filesLoaded, err = repo.LoadFiles()
		if err != nil {
			b.Fatalf(`Error while trying to load files:\n%s`, err)
			return
		}
	}

	fmt.Printf("\nLoaded %d files for repo.LoadFiles() Benchmark\n", filesLoaded)
}

// TODO: Refactor these tests!
// func BenchmarkLoadBlames(b *testing.B) {
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		b.Fatalf(`Could not get cwd:\n%s`, err)
// 		return
// 	}

// 	repoPath := path.Join(cwd, "../../")
// 	repo := NewFromLocal(repoPath)
// 	filesLoaded, err := repo.LoadFiles()
// 	linesLoaded := 0
// 	b.ResetTimer()

// 	if err != nil {
// 		b.Fatalf(`Error while trying to load files:\n%s`, err)
// 		return
// 	}

// 	for i := 0; i < b.N; i++ {
// 		linesLoaded, _ = repo.LoadBlame(func(fileCount int, filesLoaded int, linesLoaded int, currentFileName string) {})
// 	}

// 	fmt.Printf("\nLoaded %d files, %d lines for repo.LoadBlames() Benchmark\n", filesLoaded, linesLoaded)
// }

// func BenchmarkCreatePeopleProfiles(b *testing.B) {
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		b.Fatalf(`Could not get cwd:\n%s`, err)
// 		return
// 	}

// 	repoPath := path.Join(cwd, "../../")
// 	repo := NewFromLocal(repoPath)
// 	filesLoaded, err := repo.LoadFiles()
// 	linesLoaded, _ := repo.LoadBlame(func(fileCount int, filesLoaded int, linesLoaded int, currentFileName string) {})
// 	b.ResetTimer()

// 	if err != nil {
// 		b.Fatalf(`Error while trying to load files:\n%s`, err)
// 		return
// 	}

// 	for i := 0; i < b.N; i++ {
// 		repo.PeopleProfilesFromBlames()
// 	}

// 	fmt.Printf("\nLoaded %d files, %d lines for repo.PeopleProfilesFromBlames() Benchmark\n", filesLoaded, linesLoaded)
// }

// func TestRepo_GeneratePeopleProfilesFromBlameResults(t *testing.T) {

// 	timestamp := time.Now()

// 	repo := Repo{GhOrg: "NMT", Reponame: "foobar"}
// 	repo.Blames = map[string][]git.BlameLine{
// 		"/ignore/me.gif":     {git.BlameLine{AuthorTime: timestamp, AuthorName: "Foo Bar", AuthorEmail: "foo@bar.com", LineContent: "GIF"}},
// 		"/foo/bar/README.md": {git.BlameLine{AuthorTime: timestamp, AuthorName: "Foo Bar", AuthorEmail: "foo@bar.com", LineContent: "# Hello World"}},
// 		"/bar/foo/index.js": {
// 			git.BlameLine{AuthorTime: timestamp, AuthorName: "Foo Bar", AuthorEmail: "foo@bar.com", LineContent: "import React from 'react'"},
// 			git.BlameLine{AuthorTime: timestamp.AddDate(-1, 0, 0), AuthorName: "Foo Bar", AuthorEmail: "foo@bar.com", LineContent: "import React from 'react'"},
// 		},
// 		"/bar/foo/local.js": {git.BlameLine{AuthorTime: timestamp, AuthorName: "Foo Bar", AuthorEmail: "foo@bar.com", LineContent: "import Stulle from './index'"}},
// 	}

// 	repo.GeneratePeopleProfilesFromBlameResults()
// 	assert.NotNil(t, repo.People)

// 	person := repo.People.FindPerson("foo@bar.com", "Foo Bar", 4711)
// 	assert.NotNil(t, person, "Expected author to exist")

// 	fileSkill := person.Skills[".md"]
// 	assert.NotNil(t, fileSkill)
// 	assert.EqualValues(t, 1, fileSkill[timestamp])

// 	moduleSkill := person.Libraries["react"]
// 	assert.NotNil(t, moduleSkill)
// 	assert.EqualValues(t, 2, moduleSkill[timestamp])

// 	localModuleSkill := person.Libraries["index"]
// 	assert.Nil(t, localModuleSkill)

// }
