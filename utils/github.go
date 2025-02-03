package utils

import (
	"context"
	"fmt"

	"github.com/google/go-github/v68/github"
)

var (
	client = github.NewClient(nil)
	ctx    = context.Background()
)

func GetPullRequestFiles(o string, r string, n int) ([]*github.CommitFile, *github.Response, error) {
	files, resp, err := client.PullRequests.ListFiles(ctx, o, r, n, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching files: %w", err)
	}

	return files, resp, err
}
