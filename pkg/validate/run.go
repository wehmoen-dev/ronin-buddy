package validate

import (
	"errors"
	"fmt"
	"github.com/google/go-github/v64/github"
	"github.com/sethvargo/go-githubactions"
	"os"
)

func Run(files []*github.CommitFile) *PullRequestValidationResult {

	prResult := PullRequestValidationResult{
		Projects:          map[string]*ProjectValidationResult{},
		UnassociatedFiles: []*github.CommitFile{},
	}

	prRequest := validateFiles(files)

	prResult.UnassociatedFiles = prRequest.UnassociatedFiles

	for project, projectRequest := range prRequest.ValidProjects {

		githubactions.Debugf("Begin validation of project %s", project)

		result := ProjectValidationResult{
			UnassociatedFiles: projectRequest.InvalidFiles,
			MetaData: &ProjectValidationResultInner{
				Validated: false,
				Valid:     true,
				Errors:    []error{},
			},
			Logo: &ProjectValidationResultInner{
				Validated: false,
				Valid:     true,
				Errors:    []error{},
			},
		}

		if projectRequest.MetaData != nil {
			result.MetaData.Validated = true
			content, err := os.ReadFile(projectRequest.MetaData.GetFilename())
			if err != nil {
				result.MetaData.Errors = append(result.MetaData.Errors, err)
				githubactions.Debugf("Failed to read data.json for project %s: %v", project, err)
				continue
			}
			metaDataResult := ValidateProjectMetadata(content)
			if len(metaDataResult) > 0 {
				githubactions.Debugf("Project %s has invalid data.json: %v", project, metaDataResult)
				result.MetaData.Errors = metaDataResult
			} else {
				result.MetaData.Valid = true
				githubactions.Debugf("Project %s has valid data.json", project)
			}
		}

		if projectRequest.Logo != nil {
			result.Logo.Validated = true
			logoResult := validateLogo(projectRequest.Logo.GetFilename())
			if len(logoResult) > 0 {
				githubactions.Debugf("Project %s has invalid logo.png: %v", project, logoResult)
				result.Logo.Errors = logoResult
			} else {
				githubactions.Debugf("Project %s has valid logo.png", project)
			}
		} else {

			_, err := os.Stat(fmt.Sprintf("projects/%s/logo.png", project))

			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					githubactions.Debugf("Project %s is missing logo.png", project)
					result.Logo.Valid = false
					result.Logo.Errors = append(result.Logo.Errors, errors.New("logo.png is missing but required"))
				}
			} else {
				result.Logo.Valid = false
				result.Logo.Errors = append(result.Logo.Errors, errors.New("logo.png is missing but required"))
				githubactions.Debugf("Project %s is missing logo.png", project)
			}

		}

		githubactions.Debugf("Project %s validated as %t", project, result.IsValid())

		prResult.Projects[project] = &result

	}

	return &prResult

}
