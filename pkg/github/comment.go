package github

import (
	"context"
	"github.com/google/go-github/v64/github"
	"github.com/sethvargo/go-githubactions"
)

func (c *Client) Comment(body string) (*github.IssueComment, error) {
	githubactions.Debugf("Creating comment: %s/%s issue_number: %d", c.cfg.GithubContext.Owner, c.cfg.GithubContext.Repository, c.IssueNumber())
	cmt, _, err := c.client.Issues.CreateComment(context.Background(), c.cfg.GithubContext.Owner, c.cfg.GithubContext.Repository, c.IssueNumber(), &github.IssueComment{
		Body: &body,
	})

	return cmt, err
}

func (c *Client) EditComment(comment *github.IssueComment, body string) (*github.IssueComment, error) {
	githubactions.Debugf("Editing comment: %s/%s issue_number: %d comment_id: %d", c.cfg.GithubContext.Owner, c.cfg.GithubContext.Repository, c.IssueNumber(), *comment.ID)

	cmt, _, err := c.client.Issues.EditComment(context.Background(), c.cfg.GithubContext.Owner, c.cfg.GithubContext.Repository, *comment.ID, &github.IssueComment{
		Body: &body,
	})

	return cmt, err
}

func (c *Client) RemoveComment(comment *github.IssueComment) error {
	githubactions.Debugf("Removing comment: %s/%s issue_number: %d comment_id: %d", c.cfg.GithubContext.Owner, c.cfg.GithubContext.Repository, c.IssueNumber(), *comment.ID)
	_, err := c.client.Issues.DeleteComment(context.Background(), c.cfg.GithubContext.Owner, c.cfg.GithubContext.Repository, *comment.ID)
	return err
}
