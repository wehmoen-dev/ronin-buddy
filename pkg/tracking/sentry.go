package tracking

import (
	"github.com/getsentry/sentry-go"
	"github.com/sethvargo/go-githubactions"
	"github.com/wehmoen-dev/ronin-buddy/pkg/config"
	"log"
	"time"
)

func New(cfg *config.ActionConfig) *Sentry {

	dsn := githubactions.GetInput("sentry_dsn")

	if dsn == "" {
		githubactions.Warningf("No DSN found for Sentry. Skipping initialization.")
		return &Sentry{
			enabled: false,
			config:  nil,
		}
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:                dsn,
		TracesSampleRate:   1.0,
		ProfilesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	return &Sentry{
		enabled: true,
		config:  cfg,
	}
}

func (s *Sentry) CaptureException(err error) {
	if s.enabled {
		sentry.ConfigureScope(func(scope *sentry.Scope) {

			scope.SetContext("action", map[string]interface{}{
				"action":    s.config.GithubContext.GitHubContext.Action,
				"event":     s.config.GithubContext.GitHubContext.EventName,
				"base_ref":  s.config.GithubContext.GitHubContext.BaseRef,
				"head_ref":  s.config.GithubContext.GitHubContext.HeadRef,
				"head_refs": s.config.GithubContext.GitHubContext.HeadRef,
			})

			scope.SetContext("repo", map[string]interface{}{
				"owner": s.config.GithubContext.Owner,
				"repo":  s.config.GithubContext.Repository,
			})

			var number float64
			if s.config.GithubContext.GitHubContext.EventName == "pull_request" || s.config.GithubContext.GitHubContext.EventName == "pull_request_target" {
				if num, ok := s.config.GithubContext.GitHubContext.Event["number"].(float64); ok {
					number = num
				}
			}

			scope.SetContext("pull_request", map[string]interface{}{
				"number": number,
				"sha":    s.config.GithubContext.GitHubContext.SHA,
			})

			scope.SetTags(map[string]string{
				"owner": s.config.GithubContext.Owner,
				"repo":  s.config.GithubContext.Repository,
				"sha":   s.config.GithubContext.GitHubContext.SHA,
			})

			scope.SetUser(sentry.User{
				Username: s.config.GithubContext.GitHubContext.Actor,
				ID:       s.config.GithubContext.GitHubContext.ActorID,
			})

			eventId := sentry.CaptureException(err)
			ok := sentry.Flush(2 * time.Second)
			if ok && eventId != nil {
				githubactions.Debugf("Captured exception. Sentry ID: %s", *eventId)
			} else {
				githubactions.Warningf("Failed to capture exception.")
			}
		})

	} else {
		githubactions.Debugf("Sentry is disabled. Skipping exception capture.")
	}
}
