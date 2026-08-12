package models

import "image"

// ZoneEditorSnapResult is the outcome of snapping a dragged zone's center.
// GuideX and GuideY carry the alignment line the zone is holding onto and are
// only meaningful when the matching HasGuide flag is set.
type ZoneEditorSnapResult struct {
	Position  image.Point
	GuideX    float64
	GuideY    float64
	HasGuideX bool
	HasGuideY bool
}
