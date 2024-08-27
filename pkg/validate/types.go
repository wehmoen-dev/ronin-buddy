package validate

import (
	"github.com/google/go-github/v64/github"
	"github.com/sethvargo/go-githubactions"
)

type ProjectMetadata struct {
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description" validate:"required"`
	Websites    []Website   `json:"websites" validate:"required"`
	Contracts   *[]Contract `json:"contracts" validate:"required,omitempty"`
	Categories  []Category  `json:"categories" validate:"required"`
	Email       *string     `json:"email,omitempty" validate:"email,omitempty"`
	Social      *Social     `json:"social,omitempty"`
}

type Website struct {
	Url         string `json:"url"`
	Description string `json:"description"`
}

type Contract struct {
	Address     string `json:"address"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type Platform string

type Category string

type Social map[Platform]string

type ProjectValidationRequest struct {
	MetaData     *github.CommitFile
	Logo         *github.CommitFile
	InvalidFiles []*github.CommitFile
}

type PullRequestValidationRequest struct {
	ValidProjects     map[string]*ProjectValidationRequest
	UnassociatedFiles []*github.CommitFile
}

type ProjectValidationResultInner struct {
	Validated bool
	Valid     bool
	Errors    []error
}

func (c *ProjectValidationResultInner) IsValid() bool {
	if !c.Validated {
		return true
	}
	return c.Valid && len(c.Errors) == 0
}

type ProjectValidationResult struct {
	MetaData          *ProjectValidationResultInner
	Logo              *ProjectValidationResultInner
	UnassociatedFiles []*github.CommitFile
}

type PullRequestValidationResult struct {
	Projects          map[string]*ProjectValidationResult
	UnassociatedFiles []*github.CommitFile
}

func (r *ProjectValidationResult) IsValid() bool {
	// Invalid if there are unassociated files
	if len(r.UnassociatedFiles) > 0 {
		return false
	}

	// Check MetaData validity if it was validated
	if r.MetaData != nil && !r.MetaData.IsValid() {
		return false
	}

	// Check Logo validity if it was validated
	if r.Logo != nil && !r.Logo.IsValid() {
		githubactions.Debugf("Logo Valid: %v", r.Logo.Valid)
		githubactions.Debugf("Logo Errors: %v", r.Logo.Errors)
		return false
	}

	// If none of the invalid conditions are met, return true
	return true
}

func (r *PullRequestValidationResult) IsValid() bool {

	if len(r.UnassociatedFiles) > 0 {
		return false
	}

	for _, project := range r.Projects {
		if !project.IsValid() {
			return false
		}
	}
	return true
}
