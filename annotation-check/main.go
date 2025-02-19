package main

import (
	"context"
	"fmt"
	"log"
	"ministryofjustice/cloud-platform-automation/utils"
	"os"
	"strings"

	"github.com/google/go-github/v68/github"
)

var (
	ctx       = context.Background()
	ns        utils.Namespace
	nsFile    = os.Getenv("NAMESPACE_FILE")
	branch    = os.Getenv("BRANCH")
	key       = os.Getenv("GITHUB_PRIVATE_KEY")
	appid     = os.Getenv("GITHUB_APP_ID")
	installid = os.Getenv("GITHUB_INSTALLATION_ID")
	result    Results
)

type Results struct {
	Repo   []string
	Public []bool
}

func main() {
	owner, repoName, pull := utils.GetPRDetails(os.Getenv("GITHUB_REF"), os.Getenv("GITHUB_REPOSITORY"))

	client, err := utils.AppClient(key, appid, installid)
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

	teamValidation, err := utils.CheckTeamName(client, owner)
	if err != nil {
		log.Fatalf("Error checking team name: %v\n", err)
	}

	if strings.Contains(ns.MetaData.Annotations.SourceCodeURL, ",") {
		urls := strings.Split(ns.MetaData.Annotations.SourceCodeURL, ",")
		for _, url := range urls {
			publicRepo, err := utils.CheckRepoPublic(client, url)
			if err != nil {
				log.Fatalf("Error checking repo is public: %v\n", err)
			}

			result = Results{
				Repo:   append(result.Repo, url),
				Public: append(result.Public, publicRepo),
			}
		}
	} else {
		publicRepo, err := utils.CheckRepoPublic(client, ns.MetaData.Annotations.SourceCodeURL)
		if err != nil {
			log.Printf("Error checking repo is public: %v\n", err)
		}

		result = Results{
			Repo:   append(result.Repo, ns.MetaData.Annotations.SourceCodeURL),
			Public: append(result.Public, publicRepo),
		}
	}

	resultsErr := utils.Results(client, teamValidation, result.Public)
	if resultsErr != nil {
		log.Printf("Results: %v", resultsErr)
	}

	message := fmt.Sprintf("Team name: %s\n - Valid: %s\n\nRepository: %s\n - Public: %s\n", ns.MetaData.Annotations.TeamName, fmt.Sprintf("%v", teamValidation), result.Repo, fmt.Sprintf("%v", result.Public))
	err = utils.CreateComment(client, owner, repoName, message, pull)
	if err != nil {
		log.Fatalf("Error creating comment: %v\n", err)
	}
}
