package main

import (
	"github.com/wehmoen-dev/ronin-buddy/pkg/events"
	"github.com/wehmoen-dev/ronin-buddy/pkg/github"
	"os"
	"strings"
)

func main() {
	client := github.New()

	// We can only handle pull_request* events
	if !strings.HasPrefix(client.Config().GithubContext.GitHubContext.EventName, "pull_request") {
		client.Fatalf("Event %s not supported", client.Config().GithubContext.GitHubContext.EventName)
	}

	pullRequest, err := client.GetPullRequest()

	if err != nil {
		client.Fatalf("Failed to get pull request: %v", err)
	}

	if !pullRequest.GetMerged() && pullRequest.GetState() != "closed" {
		events.Default(client)
		os.Exit(0)
	}

	if pullRequest.GetMerged() && pullRequest.GetState() == "closed" {
		events.Closed(client, pullRequest)
		os.Exit(0)
	}

	client.Fatalf("Pull request is in an unknown state: %s", pullRequest.GetState())
}
