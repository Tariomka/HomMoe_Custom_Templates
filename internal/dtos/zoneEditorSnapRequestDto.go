package dtos

import "image"

// ZoneEditorSnapRequestDto asks for the snapped center of the zone named
// DraggedZone. That zone is excluded from the alignment guides it snaps to, so
// a zone never holds onto itself.
type ZoneEditorSnapRequestDto struct {
	Position    image.Point
	Positions   map[string]image.Point
	ZoneRadius  int
	DraggedZone string
}
