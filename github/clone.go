package github

import (
	"bytes"
	"fmt"
	"os/exec"
)

func CloneFromGithub(ghorg string, reponame string, outdir string, ghpat string) (string, error) {
	cloneurl := fmt.Sprintf("https://sklls-cli:%s@github.com/%s/%s.git", ghpat, ghorg, reponame)

	// Special case: if reponame is set to empty, we assume ghOrg has the format org/repo (e.g. spring-media/ep-curato)
	if reponame == "" {
		cloneurl = fmt.Sprintf("https://sklls-cli:%s@github.com/%s.git", ghpat, ghorg)
	}

	cmd := exec.Command("git", "clone", cloneurl, outdir)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		errorMsg := stderr.String()
		return "", fmt.Errorf("Cannot clone %s/%s: %s", ghorg, reponame, errorMsg)
	}

	cmdOutput := out.String()
	return cmdOutput, nil
}
