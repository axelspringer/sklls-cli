package helper

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkListAllGitFoldersRecursively(b *testing.B) {
	cwd, err := os.Getwd()
	if err != nil {
		b.Fatalf(`Could not get cwd:\n%s`, err)
		return
	}

	bmPath := path.Join(cwd, "../../")
	for i := 0; i < b.N; i++ {
		gitFolders, err := ListAllGitFoldersRecursively(bmPath)
		if err != nil {
			b.Fatalf(`Error while tryting to list git folders:\n%s`, err)
			return
		}

		if len(gitFolders) == 0 {
			b.Fatalf(`Expected git folders to be non-empty`)
		}
	}
}

func TestConvertCommitTimestamp(t *testing.T) {
	const dateStrParsedFromCommit = "Thu, 30 Dec 2021 15:48:17 +0100"
	convertedToDate := FromDateString(dateStrParsedFromCommit)
	convertedBackToStr := ToDateString(convertedToDate)

	assert.Equal(t, convertedBackToStr, dateStrParsedFromCommit)
}

// func TestContains(t *testing.T) {
// 	strArr := []string{"aaa", "bB"}
// 	assert.Equal(t, true, Contains(strArr, "aaa"), "Returns true for containing exact elements")
// 	assert.Equal(t, true, Contains(strArr, "BB"), "Returns true for containing exact elements (case-insensitive)")
// }

// func TestMergeStrings(t *testing.T) {
// 	assert.Equal(
// 		t,
// 		[]string{"a", "b", "c"},
// 		MergeStrings([]string{"a", "b"}, []string{"c"}),
// 		"Merge correctly assembles combines the elements from both inputs",
// 	)

// 	assert.Equal(
// 		t,
// 		[]string{"a", "b", "c"},
// 		MergeStrings([]string{"a", "b"}, []string{"c", "b"}),
// 		"Merge avoids merging duplicates",
// 	)

// 	assert.Equal(
// 		t,
// 		[]string{"a", "b"},
// 		MergeStrings([]string{}, []string{"a", "b"}),
// 		"Can deal with empty source array",
// 	)
// }

// // Cannot work!
// // --> See here: https://stackoverflow.com/a/60989330
// // (bad idea to use time as map key)
// func TestMergeLogs(t *testing.T) {
// 	assert.Equal(
// 		t,
// 		map[string]map[time.Time]int{
// 			".js": {
// 				test.TimeFromIsoStr("1999-01-02T04:05:06+00:01"): 3,
// 			},
// 			".ts": {
// 				test.TimeFromIsoStr("1999-01-02T04:05:06+00:01"): 1,
// 			},
// 		},
// 		MergeLogs(
// 			map[string]map[time.Time]int{
// 				".js": {
// 					test.TimeFromIsoStr("1999-01-02T04:05:06+00:01"): 1,
// 				},
// 			},
// 			map[string]map[time.Time]int{
// 				".js": {
// 					test.TimeFromIsoStr("1999-01-02T04:05:06+00:01"): 2,
// 				},
// 				".ts": {
// 					test.TimeFromIsoStr("1999-01-02T04:05:06+00:01"): 1,
// 				},
// 			},
// 		),
// 		"MergeLogs correctly combines & adds up skills",
// 	)
// }

// func TestStringsSimilar(t *testing.T) {
// 	source := "Jonas Peeck"
// 	testCases := []struct {
// 		Target      string
// 		IsSimilar   bool
// 		Description string
// 	}{
// 		{"jonas peeck", true, "Exact strings match (case-insensitive)"},
// 		{"jonas.peeck", true, "Strings that only differ in minor characters match"},
// 		{"jonas.peeck@mypersonalemail.com", true, "When prefix matches, whole string matches"},
// 		{"Jonas Müller", false, "Names that differ in last name don't match"},
// 	}

// 	/*
// 		Threshold can be changed to find out when these test cases pass!
// 		If matching with certain cases fails (when it shouldn't), add them
// 		here and play around with the threshold until test passes again
// 	*/
// 	threshhold := 0.8
// 	for _, testCase := range testCases {
// 		isSimilar, similarity := StringsSimilar(source, testCase.Target, threshhold)
// 		fmt.Printf("Strings match factor for [%s <> %s]:\t%f (isSimilar? %t)\n", source, testCase.Target, similarity, isSimilar)
// 		assert.Equal(t, testCase.IsSimilar, isSimilar, testCase.Description)
// 	}
// }

// func TestEmailPrefixSimilar(t *testing.T) {
// 	source := "jonas.peeck@bild.de"
// 	testCases := []struct {
// 		Target      string
// 		IsSimilar   bool
// 		Description string
// 	}{
// 		{"jonas.peeck@spring-media.de", true, "Exact email prefixes match"},
// 		{"jonaspeeck@axelspringer.com", true, "Email prefixes that only differ in minor characters match"},
// 		{"jonas@mypersonalemail.com", true, "When only part of the email prefix matches, there is no match"},
// 		{"jonas.mueller@axelspringer.de", false, "Names that differ in last name don't match"},
// 	}

// 	/*
// 		Threshold can be changed to find out when these test cases pass!
// 		If matching with certain cases fails (when it shouldn't), add them
// 		here and play around with the threshold until test passes again
// 	*/
// 	threshhold := 0.9
// 	for _, testCase := range testCases {
// 		isSimilar, similarity := EmailPrefixSimilar(source, testCase.Target, threshhold)
// 		fmt.Printf("Strings match factor for [%s <> %s]:\t%f (isSimilar? %t)\n", source, testCase.Target, similarity, isSimilar)
// 		assert.Equal(t, testCase.IsSimilar, isSimilar, testCase.Description)
// 	}
// }

// func TestContainsSimilarThreshhold(t *testing.T) {
// 	source := "Jonas Peeck"
// 	matchAgainst := []string{"jonas peeck", "jonas-peeck", "jonas.peeck"}
// 	threshhold := 0.8

// 	assert.Equal(t, true, ContainsSimilarThreshhold(matchAgainst, source, threshhold))
// }

// func TestEmpty(t *testing.T) {
// 	assert.Equal(t, true, Empty(""), "Empty string is considered empty")
// 	assert.Equal(t, true, Empty("   "), "Strings are trimmed before checking for emptiness")
// 	assert.Equal(t, false, Empty("abc"), "Non-empty is recognized as such")
// }

// func TestContainMatchingElement(t *testing.T) {
// 	a := []string{"github.com/userA/repo", "github.com/userB/repo", "github.com/userC/repo", "   "}
// 	b := []string{"github.com/userC/repo", "github.com/userD/repo", "github.com/userE/repo"}
// 	c := []string{"github.com/userD/repo", "github.com/userE/repo", "   ", ""}

// 	assert.Equal(t, true, ContainMatchingElement(a, b), "a & b contain matching element")
// 	assert.Equal(t, true, ContainMatchingElement(b, c), "b & c contain matching element")
// 	assert.Equal(t, false, ContainMatchingElement(a, c), "a & c do not contain matching element")
// }

// func TestContainSimilarElement(t *testing.T) {
// 	testCases := []struct {
// 		Source              []string
// 		Target              []string
// 		TreatSourceAsEmails bool
// 		TreatTargetAsEmails bool
// 		ExpectedToBeSimilar bool
// 		Description         string
// 	}{
// 		{[]string{"Jonas Peeck"}, []string{"jonas peeck"}, false, false, true, "Exact strings match"},
// 		{[]string{"Jonas Peeck"}, []string{"jonas.peeck"}, false, false, true, "Slightly different strings match"},

// 		{[]string{"JonasPeeck@email.com"}, []string{"jonaspeeck"}, true, false, true, "Exact strings match (a treated as emails)"},
// 		{[]string{"JonasPeeck@email.com"}, []string{"jonas.peeck"}, true, false, true, "Slightly different strings match (a treated as emails)"},

// 		{[]string{"jonaspeeck"}, []string{"JonasPeeck@email.com"}, false, true, true, "Exact strings match (b treated as emails)"},
// 		{[]string{"jonas.peeck"}, []string{"JonasPeeck@email.com"}, false, true, true, "Slightly different strings match (b treated as emails)"},
// 	}

// 	/*
// 		Threshold can be changed to find out when these test cases pass!
// 		If matching with certain cases fails (when it shouldn't), add them
// 		here and play around with the threshold until test passes again
// 	*/
// 	threshhold := 0.8
// 	for _, testCase := range testCases {
// 		containSimilarElement := ContainSimilarElement(
// 			testCase.Source,
// 			testCase.Target,
// 			testCase.TreatSourceAsEmails,
// 			testCase.TreatTargetAsEmails,
// 			threshhold,
// 		)
// 		fmt.Printf("Contain similar element? [%v <> %v]:\t (isSimilar? %t)\n", testCase.Source, testCase.Target, containSimilarElement)
// 		assert.Equal(t, testCase.ExpectedToBeSimilar, containSimilarElement, testCase.Description)
// 	}
// }

// func TestGetEmailPrefix(t *testing.T) {
// 	assert.Equal(t, "jonas", GetEmailPrefix("jonas@jones.com"), "Gets correct Email prefix")
// 	assert.Equal(t, "jonas", GetEmailPrefix("JOnAS@jones.com"), "Lowercases email prefix")
// }

// func TestGetGhOrgRepoName(t *testing.T) {
// 	ghOrgRepoName, ok := GetGhOrgRepoName("org", "repo")
// 	assert.Equal(t, "github.com/org/repo", ghOrgRepoName, "Github URL is correctly assembled")
// 	assert.Equal(t, true, ok, "Second return value is true when both parameters were non empty")

// 	orgMissingValue, orgMissingOk := GetGhOrgRepoName("", "repo")
// 	assert.Equal(t, "", orgMissingValue, "When org or repo are empty, return empty value")
// 	assert.Equal(t, false, orgMissingOk, "When org or repo are empty, return OK flag as false")

// 	repoMissingValue, repoMissingOk := GetGhOrgRepoName("org", "")
// 	assert.Equal(t, "", repoMissingValue, "When org or repo are empty, return empty value")
// 	assert.Equal(t, false, repoMissingOk, "When org or repo are empty, return OK flag as false")
// }

// func TestMaskLeft(t *testing.T) {
// 	assert.Equal(t, "1234", MaskLeft("1234"), "Does not touch the last four characters")
// 	assert.Equal(t, "xxxxxxxxxxxx7890", MaskLeft("ghpat_1234567890"), "Replaces evertyhing except for the last four characters with a lowercase 'x'")
// }
