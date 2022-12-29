package helper

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
)

func ListAllGitFoldersRecursively(rootPath string) ([]string, error) {
	gitFolders := []string{}
	err := filepath.WalkDir(rootPath,
		func(currentPath string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if path.Base(currentPath) == ".git" {
				gitFolders = append(gitFolders, path.Dir(currentPath))
			}

			return nil
		})

	return gitFolders, err
}

func ListAllGitFoldersRecursivelyCb(rootPath string, cb func(currentDir string, gitFoldersFound int)) ([]string, error) {
	gitFolders := []string{}
	err := filepath.WalkDir(rootPath,
		func(currentPath string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if path.Base(currentPath) == ".git" {
				gitFolders = append(gitFolders, path.Dir(currentPath))
			}

			cb(currentPath, len(gitFolders))

			return nil
		})

	return gitFolders, err
}

func ListFolders(path string) (error, []string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err, nil
	}

	folders := []string{}
	for _, file := range files {
		if file.IsDir() {
			folders = append(folders, file.Name())
		}
	}

	return nil, folders
}

func DoesFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// Stolen from: https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func Openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		panic(err)
	}
}

func Contains(s []string, b string) bool {
	for _, a := range s {
		if strings.EqualFold(a, b) {
			return true
		}
	}
	return false
}

func MergeStrings(existing []string, new []string) []string {
	merged := existing
	for _, newElem := range new {
		if !Contains(merged, newElem) {
			merged = append(merged, newElem)
		}
	}

	return merged
}

func MergeLogs(existing map[string]map[time.Time]int, new map[string]map[time.Time]int) map[string]map[time.Time]int {
	merged := existing
	for groupingKey, logEntries := range new {
		if _, ok := merged[groupingKey]; !ok {
			merged[groupingKey] = map[time.Time]int{}
		}

		for logEntryTimeDate, count := range logEntries {
			if _, ok := merged[groupingKey][logEntryTimeDate]; !ok {
				merged[groupingKey][logEntryTimeDate] = count
				continue
			}

			merged[groupingKey][logEntryTimeDate] += count
		}
	}

	return merged
}

func StringsSimilar(a string, b string, threshold float64) (bool, float64) {
	swg := metrics.NewSmithWatermanGotoh()
	similarity := strutil.Similarity(strings.ToLower(a), strings.ToLower(b), swg)
	return similarity >= threshold, similarity
}

func ToDateString(dateTime time.Time) string {
	return dateTime.Format(time.RFC1123Z)
}

func FromDateString(dateTimeStr string) time.Time {
	dateTime, _ := time.Parse(time.RFC1123Z, dateTimeStr)
	return dateTime
}

func EmailPrefixSimilar(emailA string, emailB string, threshold float64) (bool, float64) {
	isSimilar, similarity := StringsSimilar(GetEmailPrefix(emailA), GetEmailPrefix(emailB), threshold)
	return isSimilar, similarity
}

func ContainsSimilarThreshhold(s []string, e string, threshold float64) bool {
	for _, a := range s {
		similar, _ := StringsSimilar(a, e, threshold)
		if similar {
			return true
		}
	}
	return false
}

func ContainsSimilarEmailPrefixThreshold(s []string, e string, threshold float64) bool {
	swg := metrics.NewSmithWatermanGotoh()
	for _, a := range s {
		if strutil.Similarity(strings.ToLower(strings.Split(a, "@")[0]), strings.ToLower(strings.Split(e, "@")[0]), swg) >= threshold {
			return true
		}
	}
	return false
}

func ContainMatchingElement(a []string, b []string) bool {
	for _, elemA := range a {
		for _, elemB := range b {
			if strings.EqualFold(elemA, elemB) && !Empty(elemA) && !Empty(elemB) {
				return true
			}
		}
	}

	return false
}

func ContainSimilarElement(a []string, b []string, aTreatAsEmails bool, bTreatAsEmails bool, threshhold float64) bool {
	for _, elemA := range a {
		for _, elemB := range b {
			strA := elemA
			if aTreatAsEmails {
				strA = GetEmailPrefix(elemA)
			}

			strB := elemB
			if bTreatAsEmails {
				strB = GetEmailPrefix(elemB)
			}

			if Empty(strA) || Empty(strB) {
				continue
			}

			similar, _ := StringsSimilar(strA, strB, 0.8)
			if similar {
				return true
			}
		}
	}

	return false
}

// Stolen from:
// https://stackoverflow.com/a/48227081
func ContainsEmptyString(ss ...string) bool {
	for _, s := range ss {
		if s == "" {
			return true
		}
	}
	return false
}

func EnvEmpty(envVar string) bool {
	if os.Getenv(envVar) == "" {
		return true
	} else {
		return false
	}
}

func Empty(value string) bool {
	trim := strings.Trim(value, " ")
	return trim == ""
}

func GetEmailPrefix(fullEmail string) string {
	return strings.ToLower(strings.Split(fullEmail, "@")[0])
}

func GetGhOrgRepoName(ghOrg string, repoName string) (string, bool) {
	if Empty(ghOrg) || Empty(repoName) {
		return "", false
	}
	return fmt.Sprintf("github.com/%s/%s", ghOrg, repoName), true
}

func DurationInMs(duration time.Duration) int64 {
	return int64(duration / time.Millisecond)
}

func UserNameEmailToAuthor(userName string, email string) string {
	return fmt.Sprintf("%s %s", userName, email)
}

func MaskLeft(s string) string {
	rs := []rune(s)
	for i := 0; i < len(rs)-4; i++ {
		rs[i] = 'x'
	}
	return string(rs)
}
