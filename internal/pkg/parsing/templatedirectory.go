package parsing

import (
	"context"

	"sidus.io/discriminator/internal/pkg/labels"
	"sidus.io/discriminator/internal/pkg/templates"
)

// TemplateDirectory provides abstraction for the requirements om the template directory
type TemplateDirectory interface {
	GetModifiers(ctx context.Context, name string, data templates.Data) (labels.Modifier, error)
}
