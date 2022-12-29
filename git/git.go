package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// TODO: Add git hashes to this!
func LoadFiles(repopath string) ([]string, error) {
	cmd := exec.Command("git", "ls-tree", "--full-tree", "-r", "--name-only", "HEAD")
	cmd.Dir = repopath
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		errorMsg := stderr.String()
		return []string{}, fmt.Errorf("Cannot list files for in %s: %s", repopath, errorMsg)
	}

	cmdOutput := out.String()
	files := strings.Split(cmdOutput, "\n")

	return files, nil
}

type BlameLine struct {
	CommitHash  string
	LineNumber  int
	AuthorName  string
	AuthorEmail string
	IsoDateTime string
	Summary     string
	Filename    string
	LineContent string
}

func Blame(basePath string, filePath string) ([]BlameLine, error) {
	// Execute git blame command (producing porcelain aka machine-readable output)
	cmd := exec.Command("git", "blame", "-e", "--line-porcelain", filePath)
	cmd.Dir = basePath
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		errorMsg := stderr.String()
		return nil, fmt.Errorf("Cannot load blame for file %s:\n%s", filepath.Join(basePath, filePath), errorMsg)
	}

	cmdOutput := out.String()
	rawOutputLines := strings.Split(cmdOutput, "\n")

	// Bundle together individual blame-lines that belong to the same code line
	rawBlameLines := []*RawBlameLine{}
	rawBlameLine := &RawBlameLine{}
	for _, rawOutputLine := range rawOutputLines {
		isBeginningOfBlame := IsHashHeaderLine(rawOutputLine)
		isEndOfBlame := IsLineContents(rawOutputLine)

		if isBeginningOfBlame {
			rawBlameLine = NewRawBlameLine(rawOutputLine)
			continue
		}

		rawBlameLine.AddLine(rawOutputLine)

		if isEndOfBlame {
			rawBlameLines = append(rawBlameLines, rawBlameLine)
		}
	}

	// Extract the infos from the rawBlameLines
	blameLines := []BlameLine{}
	for _, rawBlameLine := range rawBlameLines {
		commitHash, lineNumber := rawBlameLine.ParseCommitHashLineNumber()

		blameLine := BlameLine{
			CommitHash:  commitHash,
			LineNumber:  lineNumber,
			AuthorName:  rawBlameLine.FindAndParseDataLine("author"),
			AuthorEmail: rawBlameLine.FindAndParseDataLine("author-mail"),
			IsoDateTime: rawBlameLine.ParseTimeDate(),
			Summary:     rawBlameLine.FindAndParseDataLine("summary"),
			Filename:    rawBlameLine.FindAndParseDataLine("filename"),
			LineContent: rawBlameLine.ParseLineContent(),
		}

		blameLines = append(blameLines, blameLine)
	}
	return blameLines, nil
}

func GetOrigin(basePath string) (string, error) {
	// Execute git blame command (producing porcelain aka machine-readable output)
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = basePath
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		errorMsg := stderr.String()
		return "", fmt.Errorf("Cannot load origin for repo %s:\n%s", basePath, errorMsg)
	}

	outStr := out.String()
	remoteUrls := strings.Split(outStr, "\n")
	return remoteUrls[0], nil
}
