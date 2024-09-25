package events

import (
	"encoding/json"
	github2 "github.com/google/go-github/v64/github"
	"github.com/sethvargo/go-githubactions"
	"github.com/slack-go/slack"
	"github.com/wehmoen-dev/ronin-buddy/pkg/github"
	"github.com/wehmoen-dev/ronin-buddy/pkg/utils"
	"github.com/wehmoen-dev/ronin-buddy/pkg/validate"
	"path/filepath"
)

func Closed(client *github.Client, request *github2.PullRequest) {

	commit, err := client.GetCommit(request.GetMergeCommitSHA())
	if err != nil {
		client.Tracking().CaptureException(err)
		client.Fatalf("Failed to get commit: %v", err)
	}

	if len(commit.Parents) == 0 {
		client.Fatalf("Commit has no parent. This should not happen.")
	}

	var dataJsonFilesChanged []string

	for _, file := range commit.Files {
		if file.GetStatus() != "removed" && filepath.Base(file.GetFilename()) == "data.json" {
			dataJsonFilesChanged = append(dataJsonFilesChanged, file.GetFilename())
			githubactions.Debugf("Found change to data.json file: %s", file.GetFilename())
		}
	}

	if len(dataJsonFilesChanged) == 0 {
		githubactions.Infof("No data.json files changed. Nothing to do.")
		return
	}

	var addedWebsites []validate.Website
	var removedWebsites []validate.Website

	for _, file := range dataJsonFilesChanged {
		newContent, err := client.GetFileContentAtSha(file, request.GetMergeCommitSHA())
		if err != nil {
			client.Tracking().CaptureException(err)
			client.Fatalf("Failed to get file content: %v", err)
		}
		oldContent, err := client.GetFileContentAtSha(file, commit.Parents[0].GetSHA())
		isNewFile := false
		if err != nil {
			isNewFile = true
			githubactions.Errorf("Failed to get file content: %v", err)
		}

		var newParsed validate.ProjectMetadata
		var oldParsed validate.ProjectMetadata

		err = json.Unmarshal(newContent, &newParsed)
		if err != nil {
			client.Tracking().CaptureException(err)
			client.Fatalf("Failed to unmarshal new content: %v", err)
		}

		err = json.Unmarshal(oldContent, &oldParsed)
		if err != nil || isNewFile {
			oldParsed = validate.ProjectMetadata{
				Websites: []validate.Website{},
			}
		}

		tdsClient := utils.NewTrustedDomainClient()

		for _, newSite := range newParsed.Websites {
			if !utils.InList(oldParsed.Websites, newSite) {

				if !tdsClient.IsWhitelisted(newSite.Url) {
					addedWebsites = append(addedWebsites, newSite)
				}
			}
		}

		for _, oldSite := range oldParsed.Websites {
			if !utils.InList(newParsed.Websites, oldSite) {
				if tdsClient.IsWhitelisted(oldSite.Url) {
					removedWebsites = append(removedWebsites, oldSite)
				}
			}
		}
	}

	if len(addedWebsites) == 0 && len(removedWebsites) == 0 {
		githubactions.Infof("Nothing to report. Nice.")
		return
	}

	webhookMessage := utils.SlackMessage(request, addedWebsites, removedWebsites)

	err = slack.PostWebhook(client.Config().SlackWebhookUrl, webhookMessage)

	if err != nil {
		client.Tracking().CaptureException(err)
		client.Fatalf("Failed to post message to slack: %v", err)
	}

	githubactions.Infof("Successfully reported changes to Slack.")

}
