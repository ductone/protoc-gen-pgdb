package pgdb

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"text/template"
)

var (
	//go:embed templates/*.tmpl templates/layouts/*.tmpl
	templateFiles embed.FS
	templates     map[string]*template.Template
)

func init() {
	err := loadTemplates()
	if err != nil {
		panic(fmt.Errorf("pgdb.loadTemplates failed; %w", err))
	}
}

func loadTemplates() error {
	templates = make(map[string]*template.Template)
	tmplFiles, err := fs.ReadDir(templateFiles, "templates")
	if err != nil {
		return err
	}

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.ParseFS(templateFiles, path.Join("templates", tmpl.Name()), "templates/layouts/*.tmpl")
		if err != nil {
			return err
		}

		templates[tmpl.Name()] = pt
	}
	return nil
}
