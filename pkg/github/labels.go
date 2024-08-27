package github

import (
	"context"
	"github.com/google/go-github/v64/github"
	"github.com/sethvargo/go-githubactions"
	"github.com/wehmoen-dev/ronin-buddy/pkg/utils"
)

func (l LabelName) String() string {
	return string(l)
}

func (l LabelList) GetByName(name LabelName) *Label {
	for _, label := range l {
		if label.Name == name {
			return &label
		}
	}
	return nil
}

func (c *Client) createMissingLabels() {
	labels := ActionLabels()

	repoLabels, _, err := c.client.Issues.ListLabels(context.Background(), c.cfg.GithubContext.Owner, c.cfg.GithubContext.Repository, nil)

	if err != nil {
		c.Fatalf("Failed to list labels: %v", err)
	}

	for _, label := range *labels {
		if !repoLabelsInclude(repoLabels, label.Name) {
			result, _, err := c.client.Issues.CreateLabel(context.Background(), c.cfg.GithubContext.Owner, c.cfg.GithubContext.Repository, &github.Label{
				Name:  utils.Ptr(label.Name.String()),
				Color: &label.Color,
			})

			if err != nil {
				c.Fatalf("Failed to create label: %v", err)
			}

			githubactions.Debugf("Created label: %s", *result.Name)
		}
	}
}

func repoLabelsInclude(labels []*github.Label, labelName LabelName) bool {
	for _, label := range labels {
		if label.GetName() == labelName.String() {
			return true
		}
	}
	return false
}

func (c *Client) SetLabel(label *Label, active bool) {

	if label == nil {
		c.Fatalf("Label is nil")
	}

	if active {
		if !c.issueHasLabel(label) {
			c.addLabelToIssue(label)
			githubactions.Debugf("Added label %s", label.Name)
		}
	} else {
		if c.issueHasLabel(label) {
			c.removeLabelFromIssue(label)
			githubactions.Debugf("Removed label %s", label.Name)
		}
	}
}

func (c *Client) issueHasLabel(label *Label) bool {

	issue := c.IssueNumber()

	if issue == 0 {
		c.Fatalf("Issue number is 0")
	}

	labels, _, err := c.client.Issues.ListLabelsByIssue(context.Background(), c.cfg.GithubContext.Owner, c.cfg.GithubContext.Repository, issue, nil)

	if err != nil {
		c.Fatalf("Failed to list labels: %v", err)
	}

	return repoLabelsInclude(labels, label.Name)
}
func (c *Client) addLabelToIssue(label *Label) {

	issue := c.IssueNumber()

	if issue == 0 {
		c.Fatalf("Issue number is 0")
	}

	_, _, err := c.client.Issues.AddLabelsToIssue(context.Background(), c.cfg.GithubContext.Owner, c.cfg.GithubContext.Repository, issue, []string{label.Name.String()})

	if err != nil {
		c.Fatalf("Failed to add label: %v", err)
	}
}

func (c *Client) removeLabelFromIssue(label *Label) {

	issue := c.IssueNumber()

	if issue == 0 {
		c.Fatalf("Issue number is 0")
	}

	_, err := c.client.Issues.RemoveLabelForIssue(context.Background(), c.cfg.GithubContext.Owner, c.cfg.GithubContext.Repository, issue, label.Name.String())

	if err != nil {
		c.Fatalf("Failed to remove label: %v", err)
	}
}
