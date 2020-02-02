package labels

import (
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

// Modifiers is a collection of Modifier
type Modifiers []Modifier

// Modifier representing labels to add and remove from a set of labels
//
// deletions always trumps additions
type Modifier struct {
	additions map[string]string
	deletions []string
}

// NewModifier parses a text for modifiers
//
// Every line that starts with a "+" will be parsed as an addition,
// additions have a key and a value separated by a "="
// ex. "+my.key=value"
// Every line that starts with a "-" will be parsed as a deletion,
// deletions only has a key
// ex. "-my.key"
func NewModifier(ctx context.Context, text io.Reader) (Modifier, error) {
	m := Modifier{
		additions: make(map[string]string),
	}
	scanner := bufio.NewScanner(text)
	for scanner.Scan() {
		row := strings.TrimSpace(scanner.Text())
		if len(row) > 0 {
			logrus.WithContext(ctx).Debugf("Parsing row: \"%s\"", row)
			if strings.HasPrefix(row, "-") {
				row = strings.TrimLeft(row, "-")
				logrus.WithContext(ctx).Debugf("Found deletion: \"%s\"", row)
				m.deletions = append(m.deletions, row)
			} else if strings.HasPrefix(row, "+") {
				row = strings.TrimLeft(row, "+")

				parts := strings.SplitN(row, "=", 2)
				if len(parts) == 2 {
					logrus.WithContext(ctx).Debugf("Found addition: \"%s\"", row)
					m.additions[parts[0]] = parts[1]
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return Modifier{}, errors.Wrapf(err, "failed to parse input for modifier")
	}
	return m, nil
}

// Apply applies a modifier to a set of labels
//
// deletions always trumps additions
func (m Modifier) Apply(labels map[string]string) {
	for key, value := range m.additions {
		labels[key] = value
	}
	for _, key := range m.deletions {
		delete(labels, key)
	}
}

// Apply applies a set of modifers to a set of labels in sequential order
func (ms Modifiers) Apply(labels map[string]string) {
	for _, m := range ms {
		m.Apply(labels)
	}
}
