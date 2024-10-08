package events

import (
	"github.com/sethvargo/go-githubactions"
	"github.com/wehmoen-dev/ronin-buddy/pkg/github"
	"github.com/wehmoen-dev/ronin-buddy/pkg/templates"
	"github.com/wehmoen-dev/ronin-buddy/pkg/utils"
	"github.com/wehmoen-dev/ronin-buddy/pkg/validate"
	"os"
	"strings"
)

func Default(client *github.Client) {

	githubactions.Infof("Beginning review of pull request: %s/%s - ID: %d", client.Config().GithubContext.Owner, client.Config().GithubContext.Repository, client.IssueNumber())

	// Set the status to pending - adds label to pr
	client.SetStatus(github.ReviewStatusPending)

	// Get the pull request details
	pullRequest, err := client.GetPullRequest()

	if err != nil {
		client.Tracking().CaptureException(err)
		client.Fatalf("Failed to get pull request: %v", err)
	}

	// GitHub Api won't return more than 3k files per PR. so we couldnt validate all files... So we reject it.
	if *pullRequest.ChangedFiles > 3000 {
		client.SetStatus(github.ReviewStatusFailed)
		client.Fatalf("Oh god. This PR is too big for automated verification. Please consider splitting it into smaller pull requests.")
	}

	content, err := templates.Render("welcome", nil)

	if err != nil {
		client.Tracking().CaptureException(err)
		client.Fatalf("Failed to render template: %v", err)
	}

	issueComment, err := client.Comment(content)

	if err != nil {
		client.Tracking().CaptureException(err)
		client.Fatalf("Failed to create comment: %v", err)
	}

	client.SetReviewComment(issueComment)

	githubactions.Debugf("Comment created: %s", *client.GetReviewComment().HTMLURL)

	// perform actual validation
	validationResults := validate.Run(client.ListFiles())

	var projectResultContent []string

	for project, result := range validationResults.Projects {

		table := utils.NewTable([]string{"File", "Result", "Details"})

		projectStatus := validate.ResultFailed

		if result.IsValid() {
			githubactions.Debugf("Project %s is considered valid", project)
			projectStatus = validate.ResultPassed
		}

		if result.MetaData.IsValid() {
			table.AddRow([]string{"data.json", validate.ResultPassed, "-"})
		} else {
			table.AddRow([]string{"data.json", validate.ResultFailed, "We found some issues with your data.json file. Please check the details."})

			for _, err := range result.MetaData.Errors {
				table.AddRow([]string{"", "", err.Error()})
			}
		}

		if result.Logo.IsValid() {
			table.AddRow([]string{"logo.png", validate.ResultPassed, "-"})
		} else {
			table.AddRow([]string{"logo.png", validate.ResultFailed, "We found some issues with your logo.png file. Please check the details."})

			for _, err := range result.Logo.Errors {
				table.AddRow([]string{"", "", err.Error()})
			}
		}

		unasoociatedFiles := "#### 👻 Unassociated files"

		if len(result.UnassociatedFiles) > 0 {
			unasoociatedFiles += "\n"
			for _, file := range result.UnassociatedFiles {
				unasoociatedFiles += "\n- " + file.GetFilename()
			}

			unasoociatedFiles += "\n\nPlease remove all unassociated files from the PR."
		} else {
			unasoociatedFiles = ""
		}

		results := table.Render()

		details, err := templates.Render("project_details", map[string]string{
			"RESULTS":      results,
			"UNASSOCIATED": unasoociatedFiles,
		})

		if err != nil {
			client.Tracking().CaptureException(err)
			client.Fatalf("Failed to render template: %v", err)
		}

		resultContent, err := templates.Render("project", map[string]string{
			"PROJECT": project,
			"STATUS":  projectStatus,
			"DETAILS": details,
		})

		if err != nil {
			client.Tracking().CaptureException(err)
			client.Fatalf("Failed to render template: %v", err)
		}

		projectResultContent = append(projectResultContent, resultContent)
	}

	globalUnassociatedFiles := "\n\n#### 👻 Unassociated files"

	if len(validationResults.UnassociatedFiles) > 0 {
		globalUnassociatedFiles += "\n"
		for _, file := range validationResults.UnassociatedFiles {
			globalUnassociatedFiles += "\n- " + file.GetFilename()
		}

		globalUnassociatedFiles += "\n\nPlease remove all unassociated files from the PR."
	} else {
		globalUnassociatedFiles = ""
	}

	reviewContent, err := templates.Render("review", map[string]string{
		"RESULTS": strings.Join(projectResultContent, "\n"),
		"NOTES":   globalUnassociatedFiles,
	})

	if err != nil {
		client.Tracking().CaptureException(err)
		client.Fatalf("Failed to render template: %v", err)
	}

	if validationResults.IsValid() {
		_, err := client.Approve()

		if err != nil {
			client.Tracking().CaptureException(err)
			client.Fatalf("Failed to approve pull request: %v", err)
		}

		issueComment, err = client.EditComment(client.GetReviewComment(), reviewContent)

		if err != nil {
			client.Tracking().CaptureException(err)
			client.Fatalf("Failed to edit comment: %v", err)
		}

		client.SetReviewComment(issueComment)

		if validationResults.IsValid() {
			client.SetStatus(github.ReviewStatusPassed)
		} else {
			client.SetStatus(github.ReviewStatusFailed)
		}

		os.Exit(0)
	} else {
		_, err := client.Review(reviewContent)

		if err != nil {
			client.Tracking().CaptureException(err)
			client.Fatalf("Failed to create review: %v", err)
		}

		if validationResults.IsValid() {
			client.SetStatus(github.ReviewStatusPassed)
		} else {
			client.SetStatus(github.ReviewStatusFailed)
		}

		os.Exit(1)
	}

}
