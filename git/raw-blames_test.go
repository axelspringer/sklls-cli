package git

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRawBlameLine(t *testing.T) {
	hashHeaderLine := "96c10005fd3f360ef10173514f2291f27b04273e 2 2"
	newRawBlameLine := NewRawBlameLine(hashHeaderLine)

	assert.Equal(t, hashHeaderLine, newRawBlameLine.HashHeader, "Hash header line is set to input argument")
	assert.Equal(t, 0, len(newRawBlameLine.RawDataLines), "Data lines are set to empty slice")
	assert.Equal(t, "", newRawBlameLine.LineContents, "LineContents are set to empty string")
}

func TestAddLine(t *testing.T) {
	hashHeaderLine := "96c10005fd3f360ef10173514f2291f27b04273e 2 2"
	authorTimeLine := "author-time 1627993497"
	contentLine := "	import (" // Note: There's a tab character before "import": `\timport (`
	newRawBlameLine := NewRawBlameLine(hashHeaderLine)

	newRawBlameLine.AddLine(authorTimeLine)
	assert.Equal(t, 1, len(newRawBlameLine.RawDataLines), "author-time line is correctly added to data lines")
	assert.Equal(t, authorTimeLine, newRawBlameLine.RawDataLines[0], "author-time line is added as is")

	newRawBlameLine.AddLine(contentLine)
	assert.Equal(t, 1, len(newRawBlameLine.RawDataLines), "content line is not added to data lines (size stays the same)")
	assert.Equal(t, contentLine, newRawBlameLine.LineContents, "content line is added to LineContents as is (tab character not removed - that's accomplished by ParseLineContent())")
}

func TestIsHashHeaderLine(t *testing.T) {
	assert.Equal(t, true, IsHashHeaderLine("96c10005fd3f360ef10173514f2291f27b04273e 2 2"), "Hash header line is recognized correctly")

	assert.Equal(t, false, IsHashHeaderLine("96c10005fd3f360ef10173514f2291f27b04273e"), "Hash header line missing line counts is not recognized")
	assert.Equal(t, false, IsHashHeaderLine("author Jonas Peeck"), "Author line is not recognized")
}

func TestIsDataLine(t *testing.T) {
	assert.Equal(t, true, IsDataLine("committer GitHub"), "Data line starting with known prefix is recognized")
	assert.Equal(t, true, IsDataLine("author-time 1627993497"), "Data line starting with known prefix is recognized")

	assert.Equal(t, false, IsDataLine("96c10005fd3f360ef10173514f2291f27b04273e 1 1 4"), "Hash header line recognized as not a data line")
	assert.Equal(t, false, IsDataLine("       package main"), "Line contents line not recognized as data line")
}

func TestIsLineContents(t *testing.T) {
	contentLine := "	import ("
	assert.Equal(t, true, IsLineContents(contentLine))
}

func TestFindAndParseDataLine(t *testing.T) {
	hashHeaderLine := "96c10005fd3f360ef10173514f2291f27b04273e 2 2"
	authorTimeLine := "author-time 1627993497"
	authorTimeValue := "1627993497"
	authorEmailLine := "author-mail <jane@doe.com>"
	authorEmailValue := "jane@doe.com"

	newRawBlameLine := NewRawBlameLine(hashHeaderLine)
	newRawBlameLine.AddLine(authorTimeLine)
	newRawBlameLine.AddLine(authorEmailLine)

	assert.Equal(t, newRawBlameLine.FindAndParseDataLine("author-time"), authorTimeValue, "Extracts correct value for author-time")
	assert.Equal(t, newRawBlameLine.FindAndParseDataLine("author-mail"), authorEmailValue, "author-email is returned without angle brackets")
	assert.Equal(t, newRawBlameLine.FindAndParseDataLine("does-not-exist"), "", "Returns empty string when dataPrefix is not found")
}

func TestParseCommitHashLineNumber(t *testing.T) {
	hashHeaderLine := "96c10005fd3f360ef10173514f2291f27b04273e 123 2"
	newRawBlameLine := NewRawBlameLine(hashHeaderLine)
	commitHash, lineNumber := newRawBlameLine.ParseCommitHashLineNumber()

	assert.Equal(t, "96c10005fd3f360ef10173514f2291f27b04273e", commitHash, "Correctly parsed commit hash")
	assert.Equal(t, 123, lineNumber, "Correctly parsed line number")
}

func TestParseLineContent(t *testing.T) {
	hashHeaderLine := "96c10005fd3f360ef10173514f2291f27b04273e 123 2"
	contentLine := "	import ("
	contentLineNoTab := "import ("

	newRawBlameLine := NewRawBlameLine(hashHeaderLine)
	newRawBlameLine.AddLine(contentLine)

	assert.Equal(t, contentLineNoTab, newRawBlameLine.ParseLineContent(), "Returns line contents without leading tab character")
}

func TestParseTimeDat(t *testing.T) {
	hashHeaderLine := "96c10005fd3f360ef10173514f2291f27b04273e 123 2"
	/*
		For ease of testing, this unix timestamp represents the following date / time:
		2012-11-10 at 09:08:07
	*/
	authorTimeLine := "author-time 1352534887"
	authorTzLine := "author-tz +0200"

	newRawBlameLine := NewRawBlameLine(hashHeaderLine)
	newRawBlameLine.AddLine(authorTimeLine)
	newRawBlameLine.AddLine(authorTzLine)

	authorDateTime := newRawBlameLine.ParseTimeDate()
	dateTime, err := time.Parse(time.RFC1123Z, authorDateTime)
	assert.Nil(t, err, "dateStr can be parsed without problems")
	assert.Equal(t, 2012, dateTime.Year(), "Year is parsed correctly")
	assert.Equal(t, time.Month(11), dateTime.Month(), "Month is parsed correctly")
	assert.Equal(t, 10, dateTime.Day(), "Day is parsed correctly")
	assert.Equal(t, 9, dateTime.Hour(), "Hour is parsed correctly")
	assert.Equal(t, 8, dateTime.Minute(), "Minute is parsed correctly")
	assert.Equal(t, 7, dateTime.Second(), "Second is parsed correctly")
}
