package config

import (
	"github.com/sethvargo/go-githubactions"
)

type ActionConfig struct {
	GithubToken     string
	SlackWebhookUrl string

	AppInstallationId string
	AppId             string
	AppPrivateKey     string

	GithubContext struct {
		*githubactions.GitHubContext
		Owner      string
		Repository string
	}
}

func FromContext() *ActionConfig {
	cfg := &ActionConfig{
		GithubToken:       githubactions.GetInput("github_token"),
		SlackWebhookUrl:   githubactions.GetInput("slack_webhook_url"),
		AppInstallationId: githubactions.GetInput("app_installation_id"),
		AppId:             githubactions.GetInput("app_id"),
		AppPrivateKey:     githubactions.GetInput("app_private_key"),
	}

	if cfg.SlackWebhookUrl != "" {
		githubactions.AddMask(cfg.SlackWebhookUrl)
	}

	ghContext, err := githubactions.Context()

	if err != nil {
		githubactions.Fatalf("Failed to get github context: %v", err)
	}

	owner, repo := ghContext.Repo()

	cfg.GithubContext = struct {
		*githubactions.GitHubContext
		Owner      string
		Repository string
	}{
		GitHubContext: ghContext,
		Owner:         owner,
		Repository:    repo,
	}

	if cfg.GithubToken == "" {
		githubactions.Fatalf(`Missing required input: %s`, "github_token")
	}

	githubactions.AddMask(cfg.GithubToken)

	if cfg.AppInstallationId == "" {
		githubactions.Fatalf(`Missing required input: %s`, "app_installation_id")
	}

	githubactions.AddMask(cfg.AppInstallationId)

	if cfg.AppId == "" {
		githubactions.Fatalf(`Missing required input: %s`, "app_id")
	}

	githubactions.AddMask(cfg.AppId)

	if cfg.AppPrivateKey == "" {
		githubactions.Fatalf(`Missing required input: %s`, "app_private_key")
	}

	githubactions.AddMask(cfg.AppPrivateKey)

	return cfg
}
