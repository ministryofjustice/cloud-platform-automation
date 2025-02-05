package main

import (
	"fmt"
	"ministryofjustice/cloud-platform-automation/utils"
	"os"

	"github.com/sethvargo/go-githubactions"
)

var (
	file = os.Getenv("FILE_NAME")
)

// main function to check if the pull request contains the file with the name specified and no other commited files.
func main() {
	owner, repoName, pull := utils.GetPRDetails(os.Getenv("GITHUB_REF"), os.Getenv("GITHUB_REPOSITORY"))

	client, err := utils.AppClient()
	if err != nil {
		panic(err)
	}

	f, _, err := utils.GetPullRequestDetails(client, owner, repoName, pull)
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
