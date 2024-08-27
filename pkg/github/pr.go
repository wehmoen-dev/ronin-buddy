package github

import (
	"context"
	"github.com/google/go-github/v64/github"
)

func (c *Client) GetPullRequest() (*github.PullRequest, error) {
	pr, _, err := c.client.PullRequests.Get(context.Background(), c.cfg.GithubContext.Owner, c.cfg.GithubContext.Repository, c.IssueNumber())
	return pr, err
}
