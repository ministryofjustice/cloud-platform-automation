package main

import (
	"context"
	"fmt"
	"log"
	"ministryofjustice/cloud-platform-automation/utils"
	"os"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

var (
	ctx = context.Background()
	d   = Data{}
	ns  utils.Namespace
)

type Data struct {
	PublicRepo bool
	ValidTeam  bool
}

func main() {
	owner, repoName, pull := utils.GetPRDetails(os.Getenv("GITHUB_REF"), os.Getenv("GITHUB_REPOSITORY"))

	client, err := utils.AppClient()
	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
	}

	prFiles, branch, err := utils.GetPullRequestDetails(client, owner, repoName, pull)
	if err != nil {
		log.Fatalf("Error getting pull request files: %v\n", err)
	}

	for _, file := range prFiles {
		if strings.Contains(file.GetFilename(), "namespace") {
			ns, err = utils.GetFileContent(client, ctx, file, owner, repoName, branch)
			if err != nil {
				log.Fatalf("Error getting file content: %v\n", err)
			}
		}
	}

	b, err := utils.CheckRepoPublic(client, ns.MetaData.Annotations.SourceCodeURL)
	if err != nil {
		log.Fatalf("Error checking repo is public: %v\n", err)
	}
	d.PublicRepo = b

	team, err := utils.CheckTeamName(client, owner)
	if err != nil {
		log.Fatalf("Error checking team name: %v\n", err)
	}
	d.ValidTeam = team

	message := fmt.Sprintf("Team name: %s\n - Valid: %s\n\nRepository: %s\n - Public: %s\n", ns.MetaData.Annotations.TeamName, fmt.Sprintf("%v", d.ValidTeam), ns.MetaData.Annotations.SourceCodeURL, fmt.Sprintf("%v", d.PublicRepo))

	if d.PublicRepo && d.ValidTeam {
		githubactions.SetOutput("valid", "true")
		err := utils.CreateComment(client, owner, repoName, message, pull)
		if err != nil {
			log.Fatalf("Error creating comment: %v\n", err)
		}
	} else if !d.PublicRepo || !d.ValidTeam {
		githubactions.SetOutput("valid", "false")
		err := utils.CreateComment(client, owner, repoName, message, pull)
		if err != nil {
			log.Fatalf("Error creating comment: %v\n", err)
		}
	}
}
