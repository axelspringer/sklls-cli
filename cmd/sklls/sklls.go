package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sklls/analyzer"
	"sklls/github"
	"sklls/helper"
	"sklls/repo"
	"strings"
	"time"
)

func main() {
	// Setup logging (if verbose flag is set)
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	concurrency := flag.Int("concurrency", 10, "Dictates how many files are analyzed in parallel")
	perfLogThreshold := flag.Int("perfLogThreshold", 500, "Performance data is logged for any file analysis that takes longer than perfLogThreshold (in ms)")
	dir := flag.String("dir", "", "Directory where to look for repos for. Default is CWD.")
	out := flag.String("out", "./", "Output folder")
	exclude := flag.String("exclude", "package-lock.json,yarn.lock", "Comma separated list of file suffixes to exclude from scanning (can be both extensions and filenames)")
	ghRepos := flag.String("ghrepos", "", "Comma separated list of Github repos (org/repo - e.g. spring-media/ep-curato) to clone from Github & analyze (requires ghpat to be set as well, if repos are non-public)")
	ghPat := flag.String("ghpat", "", "Github Personal Access Token for cloning non-public repos")
	cloneDir := flag.String("cloneDir", "", "Directory to use for cloning repos into (default: Temp dir is being used")
	timeout := flag.Int("timeout", 30, "Timeout in seconds for analyzing repositieries (default: 30 [seconds])")
	flag.Parse()
	fmt.Println("")

	excludeSuffix := strings.Split(*exclude, ",")
	githubRepos := strings.Split(*ghRepos, ",")
	useGithubRepos := *ghRepos != ""

	// Create temp folder for github repos to be cloned into
	tempCloneDir := ""
	if useGithubRepos {
		var err error

		tempCloneDir, err = ioutil.TempDir("", "sklls")
		if *cloneDir != "" {
			tempCloneDir = filepath.Join(*cloneDir, "sklls")
			err = os.MkdirAll(tempCloneDir, os.ModePerm)
		}

		if err != nil {
			fmt.Printf("Cannot create clone dir:\n%s\n", err)
			os.Exit(1)
		}
		defer os.RemoveAll(tempCloneDir)

		// Clone repos from Github...
		for _, ghOrgAndRepo := range githubRepos {
			fmt.Printf("Cloning %s...", ghOrgAndRepo)
			_, err := github.CloneFromGithub(ghOrgAndRepo, "", filepath.Join(tempCloneDir, ghOrgAndRepo), *ghPat)
			if err != nil {
				fmt.Printf("Skipping %s (Cannot clone repo: %s)\n", ghOrgAndRepo, err)
				continue
			}
			fmt.Printf("Done\n")
		}
	}

	// Set search dir
	cwd, _ := os.Getwd()
	searchDir := *dir
	if searchDir == "" {
		searchDir = cwd
	}
	if useGithubRepos {
		searchDir = tempCloneDir
	}

	// TODO: Use the cloned repos as the repo-input for analyzeRepos()
	fmt.Printf("Finding all repos in %s... ", searchDir)
	gitFolders, err := helper.ListAllGitFoldersRecursively(searchDir)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d repo(s)\n", len(gitFolders))
	for _, gitFolder := range gitFolders {
		fmt.Printf("    %s\n", gitFolder)
	}

	fmt.Printf("Analyzing repos...\n")
	if len(excludeSuffix) > 0 {
		fmt.Printf("Excluding file suffixes in analysis: %v\n", excludeSuffix)
	}

	skllsProfiles := analyzeReposWithCb(
		gitFolders,
		excludeSuffix,
		*concurrency,
		*perfLogThreshold,
		*verbose,
		time.Duration(*timeout)*time.Second,
		func(repoFinished string, reposFinished int, total int) {
			if *verbose {
				fmt.Printf("  Finished analyzing %s [%d/%d]\n", repoFinished, reposFinished, total)
			}
		},
	)
	fmt.Printf("Done\n")
	outputPath := path.Join(cwd, *out)
	err = os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		fmt.Printf("Cannot create output directory %s:\n%s\n", outputPath, err)
		os.Exit(1)
	}

	fmt.Printf("Writing analysis to %s... ", outputPath)

	for _, skllsProfile := range *skllsProfiles {
		for originRaw, originData := range skllsProfile {
			for contributor, contributorData := range originData {
				origin, err := analyzer.ParseOrigin(originRaw)

				if err != nil {
					fmt.Printf("Cannot parse origin info %s:\n%s\n", originRaw, err)
					continue
				}

				dstPath := path.Join(outputPath, origin.Org, origin.Repo)
				err = os.MkdirAll(dstPath, os.ModePerm)
				if err != nil {
					fmt.Printf("Error: Cannot create output dir %s: %s\n", dstPath, err)
				}

				filePath := path.Join(dstPath, fmt.Sprintf("%s.json", contributor))
				contributorData, err := json.Marshal(contributorData)
				if err != nil {
					fmt.Printf("Error while trying to write %s:\n%s\n", filePath, err)
					break
				}

				err = ioutil.WriteFile(filePath, contributorData, 0644)
				if err != nil {
					fmt.Printf("Error while trying to write %s:\n%s\n", filePath, err)
					break
				}
			}
		}
	}

	fmt.Printf("Done\n")
}

func analyzeReposWithCb(gitFolders []string, excludeSuffix []string, concurrency int, perfLogThreshold int, verbose bool, timeout time.Duration, progress func(repoFinished string, reposFinished int, total int)) *[]analyzer.SkllsAnalysis {
	skllsProfiles := []analyzer.SkllsAnalysis{}
	for i, gitFolder := range gitFolders {
		fmt.Printf("Analyzing [%d/%d] %s... ", i+1, len(gitFolders), gitFolder)
		repo := repo.NewFromLocal(gitFolder)
		resultCh := make(chan RepoAnalysis)

		go loadRepoAnalysis(resultCh, repo, concurrency, perfLogThreshold, verbose, excludeSuffix)
		var res RepoAnalysis

		select {
		case res = <-resultCh:
		case <-time.After(timeout):
			fmt.Printf("\nSkipping %s (Timed out after %s)\n", gitFolder, timeout)
			continue
		}

		// If an error occured, skip that one
		if res.Error != nil {
			fmt.Printf("\nSkipping %s (%s)\n", gitFolder, res.Error)
			continue
		}

		skllsProfiles = append(skllsProfiles, res.Analysis)

		fmt.Print("Done\n")
		progress(gitFolder, i+1, len(gitFolders))
	}
	return &skllsProfiles
}

type RepoAnalysis struct {
	Analysis analyzer.SkllsAnalysis
	Error    error
}

func loadRepoAnalysis(result chan RepoAnalysis, repo *repo.Repo, concurrency int, perfLogThreshold int, verbose bool, excludeSuffix []string) {
	repoAnalysis, err := analyzer.LoadAnalysisForRepo(repo, analyzer.AnalysisOptions{
		MaxConcurrentGitBlame: concurrency,
		PerfLogThreshold:      perfLogThreshold,
		EnableLogging:         verbose,
		ExcludeSuffix:         excludeSuffix,
	})
	result <- RepoAnalysis{repoAnalysis, err}
}
