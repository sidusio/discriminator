package parsing

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	"sidus.io/discriminator/internal/pkg/labels"
	"sidus.io/discriminator/internal/pkg/templates"
)

var (
	//nolint:lll
	validatorRegex      = regexp.MustCompile(`^\s*(\s*|[^\(\s\)\:]+\((\s*|[^\(\s\)\:]+\:\s*[^\(\s\)\:]+(\,\s*[^\(\s\)\:]+\:\s*[^\(\s\)\:]+)*)\)(\s*\|\s*[^\(\s\)\:]+\((\s*|[^\(\s\)\:]+\:\s*[^\(\s\)\:]+(\,\s*[^\(\s\)\:]+\:\s*[^\(\s\)\:]+)*)\))*)\s*$`)
	invalidFormatErrorF = "input \"%s\" does not follow the expected format"
)

// Parser provides functionality to parse input labels to a set of modifiers,
// provided the necessary templates
type Parser struct {
	templateDirectory TemplateDirectory
}

// NewParser creates a new parser
func NewParser(_ context.Context, templateDirectory TemplateDirectory) (Parser, error) {
	return Parser{templateDirectory: templateDirectory}, nil
}

// Process parses a string, calls the necessary templates and returns a list of modifiers.
//
// Provided input string should be on the form "templateName(parameter: value, p:v) | otherTemplate()"
// with an arbitrary number of template calls and arguments
func (p Parser) Process(ctx context.Context, s string, data templates.ContainerData) (labels.Modifiers, error) {
	// make sure the input is valid before further processing
	if !validatorRegex.Match([]byte(s)) {
		return nil, fmt.Errorf(invalidFormatErrorF, s)
	}
	logrus.WithContext(ctx).Debugf("Processing %s", s)

	// one template call at a time
	parts := strings.Split(s, "|")
	var modifiers labels.Modifiers
	for _, part := range parts {
		logrus.WithContext(ctx).Debugf("Processing %s", part)
		// Parse out arguments and template name
		template, arguments, err := parse(part)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse \"%s\"", part)
		}

		// Parse the template for modifiers
		modifier, err := p.templateDirectory.GetModifiers(ctx, template, templates.Data{
			ContainerData: data,
			Arguments:     arguments,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse template %s", template)
		}
		modifiers = append(modifiers, modifier)
	}
	return modifiers, nil
}

func parse(s string) (string, map[string]string, error) {
	// Make sure the input is valid before further processing
	if !validatorRegex.Match([]byte(s)) {
		return "", nil, fmt.Errorf(invalidFormatErrorF, s)
	}
	// Make sure string isn't empty
	s = strings.TrimSpace(s)
	if s == "" {
		return "", nil, fmt.Errorf(invalidFormatErrorF, s)
	}
	// Separate template name from arguments
	parts := strings.SplitN(s, "(", 2)
	if len(parts) != 2 {
		return "", nil, fmt.Errorf(invalidFormatErrorF, s)
	}
	template := strings.TrimSpace(parts[0])
	// Remove ending parenthesis
	argumentsString := strings.TrimRight(parts[1], ")")
	// Split up arguments
	argumentStrings := strings.Split(argumentsString, ",")
	arguments := make(map[string]string)
	for _, argumentString := range argumentStrings {
		// Only parse if there is actually something to parse
		if strings.TrimSpace(argumentString) != "" {
			// split into key and value
			argumentParts := strings.SplitN(argumentString, ":", 2)
			if len(argumentParts) != 2 {
				return "", nil, fmt.Errorf(invalidFormatErrorF, s)
			}
			key := strings.TrimSpace(argumentParts[0])
			value := strings.TrimSpace(argumentParts[1])
			arguments[key] = value
		}
	}
	return template, arguments, nil
}
