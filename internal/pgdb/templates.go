package pgdb

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"strconv"
	"strings"
	"text/template"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	//go:embed templates/*.tmpl templates/layouts/*.tmpl
	templateFiles embed.FS
	templates     map[string]*template.Template
)

// templateFuncs contains custom functions for templates.
var templateFuncs = template.FuncMap{
	"formatFieldPath":  formatFieldPath,
	"formatSourceKind": formatSourceKind,
	"formatProtoKind":  formatProtoKind,
}

// formatFieldPath converts a []int32 to Go source code representation.
func formatFieldPath(path []int32) string {
	if len(path) == 0 {
		return "nil"
	}
	parts := make([]string, len(path))
	for i, n := range path {
		parts[i] = strconv.FormatInt(int64(n), 10)
	}
	return "[]int32{" + strings.Join(parts, ", ") + "}"
}

// formatSourceKind converts a ColumnSourceKind to Go source code representation.
func formatSourceKind(kind pgdb_v1.ColumnSourceKind) string {
	return fmt.Sprintf("pgdb_v1.ColumnSourceKind(%d)", kind)
}

// formatProtoKind converts a protoreflect.Kind to Go source code representation.
func formatProtoKind(kind protoreflect.Kind) string {
	return fmt.Sprintf("protoreflect.Kind(%d)", kind)
}

//nolint:gochecknoinits // compling templates from embed
func init() {
	err := loadTemplates()
	if err != nil {
		panic(fmt.Errorf("pgdb.loadTemplates failed; %w", err))
	}
}

func templateExecToString(name string, c interface{}) (string, error) {
	buf := bytes.Buffer{}
	err := templates[name].Execute(&buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
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

		pt, err := template.New(tmpl.Name()).Funcs(templateFuncs).ParseFS(templateFiles, path.Join("templates", tmpl.Name()), "templates/layouts/*.tmpl")
		if err != nil {
			return err
		}

		templates[tmpl.Name()] = pt
	}
	return nil
}
