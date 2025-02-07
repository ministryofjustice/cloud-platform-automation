package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"ministryofjustice/cloud-platform-automation/utils"
	"os"

	"github.com/google/go-github/v68/github"
)

var (
	ctx        = context.Background()
	ns         utils.Namespace
	nsFile     = flag.String("nsfile", os.Getenv("NAMESPACE_FILE"), "Namespace file string")
	branch     = flag.String("branch", os.Getenv("BRANCH"), "Branch string")
	githubrepo = flag.String("githubrepo", os.Getenv("GITHUB_REPOSITORY"), "Github Repository string")
	githubref  = flag.String("githubref", os.Getenv("GITHUB_REF"), "Github Respository PR ref string")
	key        = flag.String("key", os.Getenv("GITHUB_PRIVATE_KEY"), "Github App private key string")
	appid      = flag.String("appid", os.Getenv("GITHUB_APP_ID"), "Github App ID string")
	installid  = flag.String("installid", os.Getenv("GITHUB_INSTALLATION_ID"), "Github Installation ID string")
)

func main() {
	flag.Parse()
	owner, repoName, pull := utils.GetPRDetails(*githubrepo, *githubref)

	client, err := utils.AppClient(*key, *appid, *installid)
	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
	}

	file := &github.CommitFile{
		Filename: nsFile,
	}

	ns, err = utils.GetFileContent(client, ctx, file, owner, repoName, *branch)
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
