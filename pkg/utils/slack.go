package utils

import (
	"fmt"
	"github.com/google/go-github/v64/github"
	"github.com/slack-go/slack"
	"github.com/wehmoen-dev/ronin-buddy/pkg/validate"
)

func SlackMessage(pullrequest *github.PullRequest, addedSites, removedSites []validate.Website) *slack.WebhookMessage {
	attachment := slack.Attachment{}
	attachment.Title = ":tada: Ronin Ecosystem Update :tada:"
	attachment.TitleLink = pullrequest.GetHTMLURL()

	attachment.AuthorName = pullrequest.GetUser().GetLogin()
	attachment.AuthorLink = pullrequest.GetUser().GetHTMLURL()
	attachment.AuthorIcon = pullrequest.GetUser().GetAvatarURL()

	var fields []slack.AttachmentField

	fields = append(fields, slack.AttachmentField{
		Title: "Added Websites",
	})

	appendWebsites(&fields, addedSites)

	fields = append(fields, slack.AttachmentField{
		Title: "Removed Websites",
	})

	appendWebsites(&fields, removedSites)

	fields = append(fields, slack.AttachmentField{
		Value: "Please review these changes and take steps to update the TDS if necessary.",
	})

	attachment.Fields = fields
	attachment.Color = "#36a64f"

	return &slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
		Username:    "Ronin Buddy via GitHub Actions",
		IconURL:     "https://cdn.axie.bot/default.png",
	}
}

func appendWebsites(fields *[]slack.AttachmentField, websites []validate.Website) {

	if len(websites) == 0 {
		*fields = append(*fields, slack.AttachmentField{
			Value: "- no changes",
		})
		return
	}

	for _, website := range websites {
		*fields = append(*fields, slack.AttachmentField{
			Value: fmt.Sprintf("- <%s> - %s", website.Url, website.Description),
		})
	}
}
