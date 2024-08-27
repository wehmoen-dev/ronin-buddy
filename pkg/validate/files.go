package validate

import (
	"github.com/google/go-github/v64/github"
	"github.com/sethvargo/go-githubactions"
	"strings"
)

func validateFiles(files []*github.CommitFile) PullRequestValidationRequest {
	result := PullRequestValidationRequest{
		ValidProjects:     make(map[string]*ProjectValidationRequest),
		UnassociatedFiles: []*github.CommitFile{},
	}

	for _, file := range files {
		filename := file.GetFilename()

		// Check if the file is within /projects/
		if !strings.HasPrefix(filename, "projects/") {
			githubactions.Debugf("Unassociated file: %s", filename)
			result.UnassociatedFiles = append(result.UnassociatedFiles, file)
			continue
		}

		// Ensure the file is directly under /projects/<project>/ and extract project name
		segments := strings.Split(filename, "/")
		if len(segments) != 3 {
			githubactions.Debugf("Invalid file structure: %s", filename)
			result.UnassociatedFiles = append(result.UnassociatedFiles, file)
			continue
		}

		project := segments[1]
		baseFilename := segments[2]

		// Initialize the request struct for this project if it doesn't exist
		if _, exists := result.ValidProjects[project]; !exists {
			result.ValidProjects[project] = &ProjectValidationRequest{}
		}

		// Validate the base filename
		switch baseFilename {
		case "data.json":
			githubactions.Debugf("Found data.json for project %s", project)
			result.ValidProjects[project].MetaData = file

		case "logo.png":
			githubactions.Debugf("Found logo.png for project %s", project)
			result.ValidProjects[project].Logo = file

		default:
			githubactions.Debugf("Invalid file in project %s: %s", project, filename)
			result.ValidProjects[project].InvalidFiles = append(result.ValidProjects[project].InvalidFiles, file)
		}
	}

	return result
}
