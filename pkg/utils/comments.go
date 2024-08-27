package utils

//
//import (
//	"github.com/wehmoen-dev/ronin-buddy/pkg/github"
//	"github.com/wehmoen-dev/ronin-buddy/pkg/templates"
//)
//
//func WelcomeComment(client *github.Client) {
//	content, err := templates.Render("welcome", nil)
//
//	if err != nil {
//		client.Fatalf("Failed to render template: %v", err)
//	}
//
//	issueComment, err := client.Comment(content)
//
//	if err != nil {
//		client.Fatalf("Failed to create comment: %v", err)
//	}
//
//	client.SetReviewComment(issueComment)
//}
