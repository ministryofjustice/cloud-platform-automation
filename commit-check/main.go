package main

import (
	"fmt"
	"ministryofjustice/cloud-platform-automation/utils"
	"os"
	"strconv"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

var (
	ref  = os.Getenv("GITHUB_REF")
	repo = os.Getenv("GITHUB_REPOSITORY")
	file = os.Getenv("FILE_NAME")
)

// main function to check if the pull request contains the file with the name specified and no other commited files.
func main() {

	githubrefS := strings.Split(ref, "/")
	prnum := githubrefS[2]
	pull, _ := strconv.Atoi(prnum)

	repoS := strings.Split(repo, "/")
	owner := repoS[0]
	repoName := repoS[1]

	f, _, err := utils.GetPullRequestDetails(owner, repoName, pull)
	if err != nil {
		panic(err)
	}

	b := utils.Files(file, f)
	if b {
		fmt.Printf("pull request commit contains the file with the name %v and no other commited files. This is valid for auto-approval.", file)
		githubactions.New().SetOutput("approval", "approval_not_needed")
	}
	if !b {
		fmt.Printf("pull request commit does not contain the file with the name %v or contains other commited files. This is not valid for auto-approval.", file)
		githubactions.New().SetOutput("approval", "approval_needed")
	}
}
