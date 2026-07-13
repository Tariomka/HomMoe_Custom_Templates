package models

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
)

// VariantMapping stores the selectable item variants for a single content SID.
type VariantMapping struct {
	// Content is the SID this mapping applies to.
	Content SidMapping
	// Variants maps each variant id to its human-readable description.
	Variants []data.Tuple[int, string]
}

// NewVariantMapping builds a VariantMapping for the given content and variants.
func NewVariantMapping(content SidMapping, variants []data.Tuple[int, string]) VariantMapping {
	return VariantMapping{Content: content, Variants: variants}
}

func (this VariantMapping) GetVariantByID(id int) (string, bool) {
	for _, tuple := range this.Variants {
		if tuple.Key == id {
			return tuple.Value, true
		}
	}
	return "", false
}

func (this VariantMapping) GetVariantIDsInOrder() []int {
	keys := make([]int, 0, len(this.Variants))
	for _, tuple := range this.Variants {
		keys = append(keys, tuple.Key)
	}
	slices.Sort(keys)
	return keys
}

// DisplayText returns the description of the lowest variant id, or the content
// name when no variants are defined. Using the lowest id keeps the result
// deterministic regardless of Go map iteration order.
func (this VariantMapping) DisplayText() string {
	if len(this.Variants) == 0 {
		return this.Content.Name
	}

	keys := this.GetVariantIDsInOrder()
	description, _ := this.GetVariantByID(keys[0])
	return description
}

// String implements [fmt.Stringer] so a VariantMapping renders as its display text.
func (this VariantMapping) String() string {
	return this.DisplayText()
}
