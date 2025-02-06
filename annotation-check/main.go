package main

import (
	"context"
	"fmt"
	"log"
	"ministryofjustice/cloud-platform-automation/utils"
	"os"

	"github.com/google/go-github/v68/github"
)

var (
	ctx    = context.Background()
	ns     utils.Namespace
	nsFile = os.Getenv("NAMESPACE_FILE")
	branch = os.Getenv("BRANCH")
)

func main() {
	owner, repoName, pull := utils.GetPRDetails(os.Getenv("GITHUB_REF"), os.Getenv("GITHUB_REPOSITORY"))

	client, err := utils.AppClient()
	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
	}

	file := &github.CommitFile{
		Filename: &nsFile,
	}

	ns, err = utils.GetFileContent(client, ctx, file, owner, repoName, branch)
	if err != nil {
		log.Fatalf("Error getting file content: %v\n", err)
	}

	publicRepo, err := utils.CheckRepoPublic(client, ns.MetaData.Annotations.SourceCodeURL)
	if err != nil {
		log.Fatalf("Error checking repo is public: %v\n", err)
	}

	teamValidation, err := utils.CheckTeamName(client, owner)
	if err != nil {
		log.Fatalf("Error checking team name: %v\n", err)
	}

	resultsErr := utils.Results(client, teamValidation, publicRepo)
	if resultsErr != nil {
		log.Printf("Results: %v", resultsErr)
	}

	message := fmt.Sprintf("Team name: %s\n - Valid: %s\n\nRepository: %s\n - Public: %s\n", ns.MetaData.Annotations.TeamName, fmt.Sprintf("%v", teamValidation), ns.MetaData.Annotations.SourceCodeURL, fmt.Sprintf("%v", publicRepo))
	err = utils.CreateComment(client, owner, repoName, message, pull)
	if err != nil {
		log.Fatalf("Error creating comment: %v\n", err)
	}
}
