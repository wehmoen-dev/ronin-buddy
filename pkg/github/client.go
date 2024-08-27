package github

import (
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v64/github"
	"github.com/sethvargo/go-githubactions"
	"github.com/wehmoen-dev/ronin-buddy/pkg/config"
	"net/http"
	"strconv"
)

type Client struct {
	client *github.Client
	cfg    *config.ActionConfig

	reviewComment *github.IssueComment
}

// New /**
//   - Create a new github client
//   - @param cfg
//   - @returns {Client}
func New() *Client {

	cfg := config.FromContext()
	privateKey := []byte(cfg.AppPrivateKey)

	appIdInt, err := strconv.Atoi(cfg.AppId)
	if err != nil {
		githubactions.Fatalf("Failed to convert app id to int: %v", err)
	}
	appInstallationIdInt, err := strconv.Atoi(cfg.AppInstallationId)
	if err != nil {
		githubactions.Fatalf("Failed to convert app installation id to int: %v", err)
	}

	itr, err := ghinstallation.New(http.DefaultTransport, int64(appIdInt), int64(appInstallationIdInt), privateKey)
	if err != nil {
		githubactions.Fatalf("Failed to create ghinstallation client: %v", err)
	}

	innerClient := github.NewClient(&http.Client{Transport: itr})

	client := &Client{
		client: innerClient,
		cfg:    cfg,
	}

	client.createMissingLabels()

	return client
}

func (c *Client) R() *github.Client {
	return c.client
}

func (c *Client) Config() *config.ActionConfig {
	return c.cfg
}

func (c *Client) IssueNumber() int {
	var number float64
	if c.cfg.GithubContext.GitHubContext.EventName == "pull_request" || c.cfg.GithubContext.GitHubContext.EventName == "pull_request_target" {
		if num, ok := c.cfg.GithubContext.GitHubContext.Event["number"].(float64); ok {
			number = num
		}
	}
	return int(number)
}

func (c *Client) Labels() *LabelList {
	return ActionLabels()
}

func (c *Client) SetStatus(status ReviewStatus) {
	githubactions.Debugf("Setting pr status to: %s", status)
	switch status {
	case ReviewStatusPending:
		c.SetLabel(c.Labels().GetByName(LabelValidationPending), true)
		c.SetLabel(c.Labels().GetByName(LabelValidationFailed), false)
		c.SetLabel(c.Labels().GetByName(LabelValidationPassed), false)
		break
	case ReviewStatusFailed:
		c.SetLabel(c.Labels().GetByName(LabelValidationPending), false)
		c.SetLabel(c.Labels().GetByName(LabelValidationFailed), true)
		c.SetLabel(c.Labels().GetByName(LabelValidationPassed), false)
		break
	case ReviewStatusPassed:
		c.SetLabel(c.Labels().GetByName(LabelValidationPending), false)
		c.SetLabel(c.Labels().GetByName(LabelValidationFailed), false)
		c.SetLabel(c.Labels().GetByName(LabelValidationPassed), true)
		break
	}
}

func (c *Client) GetReviewComment() *github.IssueComment {
	return c.reviewComment
}

func (c *Client) SetReviewComment(comment *github.IssueComment) {
	c.reviewComment = comment
}
