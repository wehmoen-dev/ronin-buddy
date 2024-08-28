package github

import (
	"context"
	"github.com/google/go-github/v64/github"
	"github.com/sethvargo/go-githubactions"
	"github.com/wehmoen-dev/ronin-buddy/pkg/utils"
)

func (c *Client) Review(body string) (*github.PullRequestReview, error) {

	if c.reviewComment != nil {
		githubactions.Debugf("Removing review comment: %s", *c.reviewComment.HTMLURL)
		err := c.RemoveComment(c.reviewComment)
		if err != nil {
			c.Tracking().CaptureException(err)
			githubactions.Debugf("Failed to remove review comment: %v", err)
		}
		c.SetReviewComment(nil) // Reset review comment
	}

	cmt, _, err := c.client.PullRequests.CreateReview(context.Background(), c.cfg.GithubContext.Owner, c.cfg.GithubContext.Repository, c.IssueNumber(), &github.PullRequestReviewRequest{
		Body:  &body,
		Event: utils.Ptr("REQUEST_CHANGES"),
	})

	if cmt != nil {
		githubactions.Errorf("Oh no. Something seems to be wrong with your pull request. Please check %s for more information.", *cmt.HTMLURL)
	}

	return cmt, err
}

func (c *Client) Approve() (*github.PullRequestReview, error) {
	githubactions.Debugf("Approving pull request 🎉")
	cmt, _, err := c.client.PullRequests.CreateReview(context.Background(), c.cfg.GithubContext.Owner, c.cfg.GithubContext.Repository, c.IssueNumber(), &github.PullRequestReviewRequest{
		Event: utils.Ptr("APPROVE"),
	})

	return cmt, err
}
