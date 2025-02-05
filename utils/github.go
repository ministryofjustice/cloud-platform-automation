package utils

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v68/github"
	"github.com/jferrl/go-githubauth"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

var (
	client = github.NewClient(nil)
	ctx    = context.Background()
	ns     = Namespace{}
)

type Namespace struct {
	MetaData struct {
		Annotations struct {
			SourceCodeURL string `yaml:"cloud-platform.justice.gov.uk/source-code"`
			TeamName      string `yaml:"cloud-platform.justice.gov.uk/team-name"`
		} `yaml:"annotations"`
	} `yaml:"metadata"`
}

func AppClient() (*github.Client, error) {
	key := os.Getenv("GITHUB_PRIVATE_KEY")
	appid := os.Getenv("GITHUB_APP_ID")
	installid := os.Getenv("GITHUB_INSTALLATION_ID")
	privateKey := []byte(key)

	appIDInt, err := strconv.ParseInt(appid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting app ID to int64: %w", err)
	}

	installIDInt, err := strconv.ParseInt(installid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting installation ID to int64: %w", err)
	}

	appTokenSource, err := githubauth.NewApplicationTokenSource(appIDInt, privateKey)
	if err != nil {
		return nil, fmt.Errorf("error creating application token source: %w", err)
	}

	installationTokenSource := githubauth.NewInstallationTokenSource(installIDInt, appTokenSource)

	// oauth2.NewClient uses oauth2.ReuseTokenSource to reuse the token until it expires.
	// The token will be automatically refreshed when it expires.
	// InstallationTokenSource has the mechanism to refresh the token when it expires.
	httpClient := oauth2.NewClient(context.Background(), installationTokenSource)

	client := github.NewClient(httpClient)
	return client, nil
}

func GetPullRequestDetails(o string, r string, n int) ([]*github.CommitFile, string, error) {
	files, _, err := client.PullRequests.ListFiles(ctx, o, r, n, nil)
	if err != nil {
		return nil, "", fmt.Errorf("error fetching files: %w", err)
	}

	branch, _, err := client.PullRequests.Get(ctx, o, r, n)
	if err != nil {
		return nil, "", fmt.Errorf("error fetching branch: %w", err)
	}

	return files, branch.GetHead().GetRef(), err
}

func GetFileContent(client *github.Client, ctx context.Context, file *github.CommitFile, owner, repo, ref string) (Namespace, error) {
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	content, _, _, err := client.Repositories.GetContents(ctx, owner, repo, *file.Filename, opts)
	if err != nil {
		return Namespace{}, err
	}

	getcon, _ := content.GetContent()

	yaml.Unmarshal([]byte(getcon), &ns)

	return ns, nil
}

func CheckRepoPublic(client *github.Client, url string) (bool, error) {
	surl := strings.Split(url, "/")
	owner := surl[3]
	prrepo := surl[4]

	repo, _, err := client.Repositories.Get(ctx, owner, prrepo)
	if err != nil {
		return false, fmt.Errorf("error fetching repo: %v", err)
	}

	if repo.GetPrivate() {
		return false, nil
	} else {
		return true, nil
	}
}

func CheckTeamName(client *github.Client, owner string) (bool, error) {
	teams, _, err := client.Teams.GetTeamBySlug(ctx, owner, ns.MetaData.Annotations.TeamName)
	if err != nil {
		return false, fmt.Errorf("error fetching teams: %v", err)
	}

	if teams.GetSlug() == ns.MetaData.Annotations.TeamName {
		return true, nil
	}
	return false, err
}

func CreateComment(client *github.Client, owner, repoName, message string, pull int) error {
	comment := &github.IssueComment{
		Body: github.Ptr(message),
	}

	_, _, err := client.Issues.CreateComment(ctx, owner, repoName, pull, comment)
	if err != nil {
		return err
	}

	return nil
}
