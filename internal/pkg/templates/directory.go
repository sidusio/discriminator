package templates

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	"sidus.io/discriminator/internal/pkg/labels"
)

// Directory is a collection of templates
type Directory struct {
	templates *template.Template
	extension string
}

// LoadTemplatesFromPath loads all templates in the given path
func LoadTemplatesFromPath(ctx context.Context, path, extension string) (*template.Template, error) {
	tmpl := template.New("collection")
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), extension) {
				tmpl, err = tmpl.ParseFiles(path)
				if err != nil {
					logrus.WithContext(ctx).WithError(err).Warnf("failed to parse template %s, this template will not be loaded", path)
					return nil
				}
			}
			return nil
		})
	if err != nil {
		return nil, errors.Wrapf(err, "error while processing templates in %s", path)
	}
	return tmpl, nil
}

// NewDirectory creates a directory
func NewDirectory(_ context.Context, tmpl *template.Template, extension string) (*Directory, error) {
	return &Directory{
		templates: tmpl,
		extension: extension,
	}, nil
}

// GetModifiers parses the templates and get modifiers for the specified name and data
//
// The name has to be in the template collection for this method to work
func (d Directory) GetModifiers(ctx context.Context, name string, data Data) (labels.Modifier, error) {
	var text bytes.Buffer
	err := d.templates.ExecuteTemplate(&text, name+d.extension, data)
	if err != nil {
		return labels.Modifier{}, errors.Wrapf(err, "failed to parse template %s with data %+v", name, data)
	}
	return labels.NewModifier(ctx, bytes.NewReader(text.Bytes()))
}

func (d Directory) Count(ctx context.Context) int {
	return len(d.templates.Templates())
}
