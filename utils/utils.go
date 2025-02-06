package utils

import (
	"strconv"
	"strings"
)

// GetPRDetails uses pull request details and splits them into repository owner, repository name and pull request number and return the values
func GetPRDetails(ref, repo string) (string, string, int) {
	githubrefS := strings.Split(ref, "/")
	prnum := githubrefS[2]
	pull, _ := strconv.Atoi(prnum)

	repoS := strings.Split(repo, "/")
	owner := repoS[0]
	repoName := repoS[1]

	return owner, repoName, pull
}
