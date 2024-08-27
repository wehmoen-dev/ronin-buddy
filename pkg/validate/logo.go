package validate

import (
	"bytes"
	"fmt"
	"github.com/sethvargo/go-githubactions"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func validateLogo(filePath string) []error {

	var errors []error

	file, err := os.ReadFile(filePath)
	if err != nil {
		githubactions.Debugf("error opening logo: %v", err)
		return append(errors, fmt.Errorf("error opening logo: %v", err))
	}

	img, filetype, err := image.DecodeConfig(bytes.NewReader(file))

	if err != nil {
		githubactions.Debugf("error reading logo stats: %v", err)
		return append(errors, fmt.Errorf("error reading logo: %v", err))
	}

	fileStats, err := os.Stat(filePath)

	if err != nil {
		githubactions.Debugf("error reading logo stats: %v", err)
		return append(errors, fmt.Errorf("error reading logo stats: %v", err))
	}

	if fileStats.Size() > ImageMaxFileSize {
		errors = append(errors, fmt.Errorf("logo must be at most %d bytes - is %d", ImageMaxFileSize, fileStats.Size()))
	}

	if filetype != "png" {
		errors = append(errors, fmt.Errorf("logo must be a PNG file. Got %s", filetype))
	}

	if img.Width < ImageMinSideLength || img.Height < ImageMinSideLength {
		errors = append(errors, fmt.Errorf("logo must be at least %dx%d", ImageMinSideLength, ImageMinSideLength))
	}

	if img.Width > ImageMaxSideLength || img.Height > ImageMaxSideLength {
		errors = append(errors, fmt.Errorf("logo must be at most %dx%d", ImageMaxSideLength, ImageMaxSideLength))
	}

	if img.Width != img.Height {
		errors = append(errors, fmt.Errorf("logo must be a square"))
	}

	return errors
}
