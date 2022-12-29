package git

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFiles(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	basePath := path.Join(cwd, "fixtures/test-repo")
	files, err := LoadFiles(basePath)

	assert.Nil(t, err, "LoadFiles does not throw any errors")
	// It's weird - for some reason I'm getting more than one file even though the original command also only shows one file...
	assert.Equal(t, "git-blame.test.txt", files[0], "Returns expected file")
}

func TestBlame(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	basePath := path.Join(cwd, "fixtures/test-repo")
	filePath := "git-blame.test.txt"
	blameLines, err := Blame(basePath, filePath)

	assert.Nil(t, err, "Blame does not throw any errors")
	blame := blameLines[0]

	assert.Equal(t, "0bc4c685bb1a9f19f5eee95d2597b4cfd8d153b5", blame.CommitHash, "BlameLine.CommitHash is correctly parsed")
	assert.Equal(t, 1, blame.LineNumber, "BlameLine.LineNumber is correctly parsed")
	assert.Equal(t, "jpeeck-spring", blame.AuthorName, "BlameLine.AuthorName is correctly parsed")
	assert.Equal(t, "jonas.peeck@axelspringer.com", blame.AuthorEmail, "BlameLine.AuthorEmail is correctly parsed")
	assert.Equal(t, "Sun, 20 Feb 2022 14:05:39 +0100", blame.IsoDateTime, "BlameLine.IsoDateTime is correctly parsed")
	assert.Equal(t, "Testcommit", blame.Summary, "BlameLine.Summary is correctly parsed")
	assert.Equal(t, "git-blame.test.txt", blame.Filename, "BlameLine.Filename is correctly parsed")
	assert.Equal(t, "1234567890", blame.LineContent, "BlameLine.LineContent is correctly parsed")
}

func TestGetOrigin(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	basePath := path.Join(cwd, ".")
	origin, err := GetOrigin(basePath)

	assert.Nil(t, err, "Getorigin does not throw any errors")
	assert.Equal(t, "git@github.com:spring-media/sklls-cli.git", origin, "Origin returns the expected value")
}
