package templates

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/sethvargo/go-githubactions"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed views/*
var templateFiles embed.FS

var Templates *template.Template

func init() {
	var err error
	Templates, err = template.ParseFS(templateFiles, "views/*.tmpl")
	if err != nil {
		githubactions.Fatalf("Failed to parse templates: %v", err)
	}
	for _, t := range Templates.Templates() {
		githubactions.Debugf("Loaded template: %s", t.Name())
	}
}

func Render(name string, data map[string]string) (string, error) {

	if filepath.IsAbs(name) || strings.Contains(name, "..") || strings.ContainsAny(name, "\\/:") {
		return "", fmt.Errorf("invalid template name: %s", name)
	}

	if strings.Contains(name, "/") {
		name = strings.ReplaceAll(name, "/", "_")
	}

	if !strings.HasSuffix(name, ".tmpl") {
		name = fmt.Sprintf("%s.tmpl", name)
	}

	var buf bytes.Buffer
	err := Templates.ExecuteTemplate(&buf, name, data)
	return buf.String(), err
}
