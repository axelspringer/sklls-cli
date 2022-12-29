package analyzer

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"sklls/repo"
	"strings"
	"sync"
	"time"
)

type AnalysisOptions struct {
	MaxConcurrentGitBlame int
	PerfLogThreshold      int
	EnableLogging         bool
	ExcludeSuffix         []string
}

func LoadAnalysisForRepo(r *repo.Repo, opts AnalysisOptions) (SkllsAnalysis, error) {
	// Make sure we have files to analyze!
	if len(r.Files) == 0 {
		r.LoadFiles(opts.ExcludeSuffix)
	}

	if !opts.EnableLogging {
		log.SetOutput(ioutil.Discard)
	}

	// Return if no origin URL was found
	if r.OriginUrl == "" {
		return nil, fmt.Errorf("No origin URL found")
	}

	log.Printf("[ANALYZE] Max concurrent git blames: %d\n", opts.MaxConcurrentGitBlame)
	log.Printf("[ANALYZE] %s - %d files found\n", r.OriginUrl, len(r.Files))

	// This is still a bit messy! But prepare the packages for the node.js parser
	ecmaScriptPackages := NewEcmaScriptPackages()
	for _, repoFile := range r.Files {
		if ecmaScriptPackages.IsPackageFile(FilePath(repoFile)) {
			filePath := path.Join(r.Repopath, repoFile)
			content, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			ecmaScriptPackages.AddPackageFile(string(content))
		}
	}

	// Prepare dependency parsers
	parsers := []DependencyParser{
		NewEcmaScriptParser(ecmaScriptPackages),
	}

	// Setup results slots
	skllsAnalysisResults := make([]SkllsAnalysis, len(r.Files))

	// Perform analysis
	var wg sync.WaitGroup
	rate := make(chan struct{}, opts.MaxConcurrentGitBlame)
	filesCompleted := 0

	for index, file := range r.Files {
		wg.Add(1)
		go func(filePath string, results *SkllsAnalysis) {
			defer wg.Done()
			rate <- struct{}{}
			loadAnalysisForFile(results, r.OriginUrl, r.Repopath, filePath, parsers, time.Duration(opts.PerfLogThreshold)*time.Millisecond)
			filesCompleted += 1
			<-rate
		}(file, &skllsAnalysisResults[index])
	}

	// Give updates at an interval of a few seconds
	ticker := time.NewTicker(5 * time.Second)
	firstTick := true
	go func() {
		for range ticker.C {
			if opts.EnableLogging {
				if firstTick {
					fmt.Print("\n")
					firstTick = false
				}
				fmt.Printf("\t%d/%d files analyzed\n", filesCompleted, len(r.Files))
			}
		}
	}()

	wg.Wait()
	ticker.Stop()

	// Merge results
	results := SkllsAnalysis{}
	for _, fileAnalysisResult := range skllsAnalysisResults {
		results.Merge(&fileAnalysisResult)
	}

	return results, nil
}

func loadAnalysisForFile(results *SkllsAnalysis, gitOriginUrl string, repoPath string, filePath string, dependencyParsers []DependencyParser, perfLogThreshold time.Duration) {
	// Load blame for file
	repoName := filepath.Base(repoPath)
	blameStart := time.Now()
	blameLines, _ := repo.LoadBlame(repoPath, filePath)
	blameEnd := time.Since(blameStart)

	// Extract infos from the blame
	ext := Ext(strings.ToLower(filepath.Ext(filePath)))
	skllsAnalysis := SkllsAnalysis{}

	// Extract all the dependencies
	depsInFile := map[ParserName][]Dependency{}
	commitHashes := map[string]bool{}
	for _, blameLine := range blameLines {
		commitHashes[blameLine.CommitHash] = true

		for _, dp := range dependencyParsers {
			if dp.ExtSupported(string(ext)) {
				depString, version := dp.ParseLine(blameLine.LineContent)
				dep := Dependency(depString + "@" + version)
				if dep != "@" {
					parserName := ParserName(dp.GetDisplayName())

					// Add to deps
					if _, ok := depsInFile[parserName]; !ok {
						depsInFile[parserName] = []Dependency{}
					}

					if !containsDependency(depsInFile[parserName], dep) {
						depsInFile[parserName] = append(depsInFile[parserName], dep)
					}
				}
			}
		}
	}

	// To not make the log explode, only log files that took longer than a certain amount of time
	if blameEnd > perfLogThreshold {
		log.Printf(
			`[ANALYZE] {"File":"%s/%s/%s","BlameLineCount":"%d","Duration":"%s","DurationPerLine":"%s","CommitCount":"%d"}`,
			gitOriginUrl,
			repoName,
			filePath,
			len(blameLines),
			blameEnd,
			time.Duration(int(blameEnd)/int(math.Max(float64(len(blameLines)), 1))),
			len(commitHashes),
		)
	}

	// Add all of lines to the sklls analysis
	for _, blameLine := range blameLines {
		// Contributors are identified thorugh their emails
		contributor := ContributorId(strings.ToLower(blameLine.AuthorEmail))

		// Add ext
		origin := Origin(gitOriginUrl)
		dateStr := DateStr(blameLine.IsoDateTime)
		skllsAnalysis.AddExtCount(origin, contributor, ext, dateStr, 1)

		// Add dep
		for parserName, dependencies := range depsInFile {
			for _, dependency := range dependencies {
				skllsAnalysis.AddDepCount(origin, contributor, parserName, dependency, dateStr, 1)
			}
		}

		// Add email
		skllsAnalysis.AddUsername(origin, contributor, blameLine.AuthorName)
	}

	*results = skllsAnalysis
}

func containsDependency(s []Dependency, b Dependency) bool {
	for _, a := range s {
		if strings.EqualFold(string(a), string(b)) {
			return true
		}
	}
	return false
}
