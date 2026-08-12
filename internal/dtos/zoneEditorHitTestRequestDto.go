package dtos

import "image"

// ZoneEditorHitTestRequestDto asks which zone node of a laid-out canvas covers
// Position, given the node centers and their shared radius.
type ZoneEditorHitTestRequestDto struct {
	Position   image.Point
	Positions  map[string]image.Point
	ZoneRadius int
}
