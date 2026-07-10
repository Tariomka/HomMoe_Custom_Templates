package models

import "sort"

// VariantMapping stores the selectable item variants for a single content SID.
// It mirrors the C# `VariantMapping` class: a content reference plus a map of
// variant id → human-readable description.
type VariantMapping struct {
	// Content is the SID this mapping applies to.
	Content SidMapping
	// Variants maps each variant id to its human-readable description.
	Variants map[int]string // TODO: Does this need to be a map?
}

// NewVariantMapping builds a VariantMapping for the given content and variants.
func NewVariantMapping(content SidMapping, variants map[int]string) VariantMapping {
	return VariantMapping{Content: content, Variants: variants}
}

// DisplayText returns the description of the lowest variant id, or the content
// name when no variants are defined. Using the lowest id keeps the result
// deterministic regardless of Go map iteration order.
func (this VariantMapping) DisplayText() string {
	if len(this.Variants) == 0 {
		return this.Content.Name
	}
	keys := make([]int, 0, len(this.Variants))
	for key := range this.Variants {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	return this.Variants[keys[0]]
}

// String implements [fmt.Stringer] so a VariantMapping renders as its display text.
func (this VariantMapping) String() string {
	return this.DisplayText()
}
