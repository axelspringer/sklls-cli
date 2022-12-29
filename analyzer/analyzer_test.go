package analyzer

import (
	"os"
	"path"
	"sklls/repo"
	"testing"
)

// Benchmark #1: No limit on parallel git blames: 	29,01s
// Benchmark #2: Max 1 parallel git blame: 			30.67s
// Benchmark #3: Max 10 parallel git blame: 		10.81s
// Benchmark #4: Max 100 parallel git blame:		14.96s
// Benchmark #5: Max 1000 parallel git blame:		18.68s
// Benchmark #6: Max 10000 parallel git blame:		18.68s
// Benchmark #6: Max 100000 parallel git blame:		17.41s

// Benchmark #7: Max 10 parallel git blames, 10.000 lines max	11.511s
// Benchmark #*: Max 10 parallel git blames, line max disabled	11.495s <-- Didn't seem to make ANY difference!
func BenchmarkLoadBlamesForRepo(b *testing.B) {
	cwd, err := os.Getwd()
	if err != nil {
		b.Fatalf(`Could not get cwd:\n%s`, err)
		return
	}

	repoPath := path.Join(cwd, "../../")
	r := repo.NewFromLocal(repoPath)
	_, err = r.LoadFiles()
	if err != nil {
		b.Fatalf(`Could not load files:\n%s`, err)
		return
	}
	b.ResetTimer()

	if err != nil {
		b.Fatalf(`Error while trying to load files:\n%s`, err)
		return
	}

	for i := 0; i < b.N; i++ {
		LoadAnalysisForRepo(r, AnalysisOptions{
			MaxConcurrentGitBlame: 10,
			EnableLogging:         false,
		})
	}

}
