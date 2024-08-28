package github

import (
	"context"
	"github.com/google/go-github/v64/github"
)

func (c *Client) ListFiles() []*github.CommitFile {

	files, _, err := c.client.PullRequests.ListFiles(context.Background(), c.cfg.GithubContext.Owner, c.cfg.GithubContext.Repository, c.IssueNumber(), nil)
	if err != nil {
		c.Tracking().CaptureException(err)
		c.Fatalf("Failed retrieve files: %v", err)
	}

	return files

}
