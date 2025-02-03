package main

import (
	"context"
	"fmt"
	"log"
	"ministryofjustice/cloud-platform-automation/utils"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-github/v68/github"
	"github.com/sethvargo/go-githubactions"
	"gopkg.in/yaml.v2"
)

var (
	ctx   = context.Background()
	owner = "jackstockley89"
	repo  = "cloud-platform-environments"
	ns    = Namespace{}
	d     = Data{}
	token = os.Getenv("GITHUB_OAUTH_TOKEN")
)

type Namespace struct {
	MetaData struct {
		Annotations struct {
			SourceCodeURL string `yaml:"cloud-platform.justice.gov.uk/source-code"`
			TeamName      string `yaml:"cloud-platform.justice.gov.uk/team-name"`
		} `yaml:"annotations"`
	} `yaml:"metadata"`
}

type Data struct {
	ValidURL   bool
	PublicRepo bool
	ValidTeam  bool
}

func Client(token string) *github.Client {
	client := github.NewTokenClient(ctx, token)
	return client
}

func GetFileContent(client *github.Client, ctx context.Context, file *github.CommitFile, owner, repo, ref string) (*github.RepositoryContent, error) {
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	content, _, _, err := client.Repositories.GetContents(ctx, owner, repo, *file.Filename, opts)
	if err != nil {
		fmt.Printf("Error fetching file content: %v\n", err)
		return nil, err
	}

	return content, nil
}

func CheckRepoPublic(url string) (bool, error) {
	client := Client(token)

	surl := strings.Split(url, "/")
	owner := surl[3]
	prrepo := surl[4]

	fmt.Printf("Owner: %s\n", owner)
	fmt.Printf("Repo: %s\n", prrepo)

	repo, _, err := client.Repositories.Get(ctx, owner, prrepo)
	if err != nil {
		fmt.Printf("Error fetching repo: %v\n", err)
		return false, err
	}

	if repo.GetPrivate() {
		return false, nil
	} else {
		return true, nil
	}
}

func CheckTeamName() (bool, error) {
	client := Client(token)

	teams, _, err := client.Teams.ListTeams(ctx, "ministryofjustice", nil)
	if err != nil {
		fmt.Printf("Error fetching teams: %v\n", err)
		return false, err
	}

	for _, team := range teams {
		if team.GetName() == ns.MetaData.Annotations.TeamName {
			return true, nil
		} else {
			return false, nil
		}
	}

	return false, err
}

func main() {
	client := github.NewClient(nil)

	prFiles, _, err := utils.GetPullRequestFiles(owner, repo, 66)
	if err != nil {
		log.Fatalf("Error getting pull request files: %v\n", err)
	}

	for _, file := range prFiles {
		if strings.Contains(file.GetFilename(), "00-namespace.yaml") {
			content, err := GetFileContent(client, ctx, file, owner, repo, "source-code-test")
			if err != nil {
				log.Fatalf("Error getting file content: %v\n", err)
			}

			getcon, _ := content.GetContent()

			yaml.Unmarshal([]byte(getcon), &ns)
		}
	}

	url, err := url.Parse(ns.MetaData.Annotations.SourceCodeURL)
	if err != nil {
		log.Fatalf("Error parsing url: %v\n", err)
	}

	if url.Host == "github.com" {
		d.ValidURL = true
		b, err := CheckRepoPublic(ns.MetaData.Annotations.SourceCodeURL)
		if err != nil {
			log.Fatalf("Error checking repo is public: %v\n", err)
		}
		if b {
			d.PublicRepo = true
		} else {
			d.PublicRepo = false
		}
	} else {
		d.ValidURL = false
	}

	team, err := CheckTeamName()
	if err != nil {
		log.Fatalf("Error checking team name: %v\n", err)
	}
	if team {
		d.ValidTeam = true
	} else {
		d.ValidTeam = false
	}

	switch {
	case d.ValidURL && d.PublicRepo && d.ValidTeam:
		githubactions.SetOutput("valid", "true")
		githubactions.SetOutput("source-code-url", ns.MetaData.Annotations.SourceCodeURL)
		githubactions.SetOutput("team-name", ns.MetaData.Annotations.TeamName)
	case !d.ValidURL || !d.PublicRepo || !d.ValidTeam:
		githubactions.SetOutput("valid", "false")
		githubactions.SetOutput("source-code-url", ns.MetaData.Annotations.SourceCodeURL)
		githubactions.SetOutput("team-name", ns.MetaData.Annotations.TeamName)
	}
}
