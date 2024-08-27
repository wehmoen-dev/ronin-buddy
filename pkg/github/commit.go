package github

import (
	"context"
	"encoding/base64"
	"github.com/google/go-github/v64/github"
)

func (c *Client) GetCommit(sha string) (*github.RepositoryCommit, error) {
	commit, _, err := c.client.Repositories.GetCommit(context.Background(), c.cfg.GithubContext.Owner, c.cfg.GithubContext.Repository, sha, nil)

	return commit, err
}

func (c *Client) GetFileContentAtSha(filename string, sha string) ([]byte, error) {
	file, _, _, err := c.client.Repositories.GetContents(context.Background(), c.cfg.GithubContext.Owner, c.cfg.GithubContext.Repository, filename, &github.RepositoryContentGetOptions{
		Ref: sha,
	})

	if err != nil {
		return nil, err
	}

	content, err := base64.StdEncoding.DecodeString(*file.Content)

	return content, nil
}
