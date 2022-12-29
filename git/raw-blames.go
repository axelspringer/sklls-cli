package git

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

/*
	git blame --line-porcelain have the following structure:

	HashHeader
	RawDataLine
	RawDataLine
	RawDataLine
	RawDataLine
	...
	LineContents

	--> Example output of git blame --line-porcelain:

	96c10005fd 3f360ef101 73514f2291 f27b04273e 1 1 4
	author Jonas Peeck
	author-mail <36203457+jpeeck-spring@users.noreply.github.com>
	author-time 1627993497
	author-tz +0200
	committer GitHub
	committer-mail <noreply@github.com>
	committer-time 1627993497
	committer-tz +0200
	summary Add hackathon prototype (#1)
	filename skillshub-analyzer/cmd/lambda/lambda.go
					package main


	See here for more details:
	https://git-scm.com/docs/git-blame#_the_porcelain_format
*/

// Helper functions to detect various parts of a porcelain git blame output

// Data lines start with one of the following prefixes:
var DataLinePrefixes = [...]string{
	"author",
	"author-mail",
	"author-time",
	"author-tz",
	"committer",
	"committer-mail",
	"committer-time",
	"committer-tz",
	"summary",
	"filename",
}

type RawBlameLine struct {
	HashHeader   string
	RawDataLines []string
	LineContents string
}

/*
	Functions to assemble raw blame lines from git blame --line-porcelain -e
*/
func NewRawBlameLine(hashHeaderLine string) *RawBlameLine {
	return &RawBlameLine{
		HashHeader:   hashHeaderLine,
		RawDataLines: []string{},
		LineContents: "",
	}
}

func (rawBlameLine *RawBlameLine) AddLine(line string) {
	if IsLineContents(line) {
		rawBlameLine.LineContents = line
		return
	}

	rawBlameLine.RawDataLines = append(rawBlameLine.RawDataLines, line)
}

// See detailled git-blame output + explanations at the top ☝🏻
func IsHashHeaderLine(line string) bool {
	var commitHash string
	var lineNumber int
	argsParsed, _ := fmt.Sscanf(line, "%s %d", &commitHash, &lineNumber)

	return argsParsed == 2 && len(commitHash) == 40
}

func IsDataLine(line string) bool {
	for _, dataLinePrefix := range DataLinePrefixes {
		if strings.HasPrefix(line, dataLinePrefix) {
			return true
		}
	}

	return false
}

func IsLineContents(line string) bool {
	return strings.HasPrefix(line, "\t")
}

/*
	Functions to consume raw blame lines
*/
func (rawBlameLine *RawBlameLine) FindAndParseDataLine(dataPrefix string) string {
	for _, dataLine := range rawBlameLine.RawDataLines {
		if strings.HasPrefix(dataLine, dataPrefix) {
			dataLineValue := strings.Join(strings.Split(dataLine, dataPrefix)[1:], dataPrefix)
			dataLineValue = strings.ReplaceAll(dataLineValue, "<", "")
			dataLineValue = strings.ReplaceAll(dataLineValue, ">", "")
			dataLineValue = strings.Trim(dataLineValue, " ")

			return dataLineValue
		}
	}

	return ""
}

func (rawBlameLine *RawBlameLine) ParseCommitHashLineNumber() (string, int) {
	var commitHash string
	var lineNumber int
	fmt.Sscanf(rawBlameLine.HashHeader, "%s %d", &commitHash, &lineNumber)

	return commitHash, lineNumber
}

func (rawBlameLine *RawBlameLine) ParseLineContent() string {
	return strings.Trim(rawBlameLine.LineContents, "\t")
}

// Returns a RFC1123Z timestamp string
func (rawBlameLine *RawBlameLine) ParseTimeDate() string {
	authorTime := rawBlameLine.FindAndParseDataLine("author-time")
	authorTz := rawBlameLine.FindAndParseDataLine("author-tz")

	i, _ := strconv.ParseInt(authorTime, 10, 64)
	timestamp := time.Unix(i, 0)

	rawLayout := "Mon, 02 Jan 2006 15:04:05"
	rawDateStr := timestamp.Format(rawLayout)

	// This little appendage makes this a RFC1123Z timestamp string
	// See here: https://pkg.go.dev/time#pkg-constants
	rawIsoDateStr := fmt.Sprintf("%s %s", rawDateStr, authorTz)
	return rawIsoDateStr
}
