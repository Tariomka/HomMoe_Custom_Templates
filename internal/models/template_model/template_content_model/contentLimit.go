package template_content_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type ContentLimit struct {
	SID          string
	IncludeLists []string
	Content      []WeightedContent
	Variant      *int
	MaxCount     int
}

func (this ContentLimit) Clone() ContentLimit {
	clone := this
	clone.IncludeLists = slices.Clone(this.IncludeLists)
	clone.Content = slices.Clone(this.Content)
	clone.Variant = helpers.ClonePointer(this.Variant)
	return clone
}
