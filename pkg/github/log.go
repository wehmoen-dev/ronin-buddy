package github

import (
	"fmt"
	"github.com/sethvargo/go-githubactions"
	"github.com/wehmoen-dev/ronin-buddy/pkg/templates"
	"strings"
)

// Wraps FatalF to potentially edit any already made review comment.
func (c *Client) Fatalf(format string, args ...interface{}) {

	if c.reviewComment != nil && strings.HasPrefix(c.cfg.GithubContext.GitHubContext.EventName, "pull_request") {

		content, err := templates.Render("error", map[string]string{
			"ERROR": fmt.Sprintf(format, args...),
		})

		if err != nil {
			githubactions.Debugf("Failed to render error template: %v", err)
			githubactions.Fatalf(format, args...)
		}

		_, err = c.EditComment(c.reviewComment, content)

		if err != nil {
			githubactions.Debugf("Failed to edit review comment: %v", err)
		}
	}

	githubactions.Fatalf(format, args...)
}
