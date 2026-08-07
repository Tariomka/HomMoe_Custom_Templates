package zoneEditorGeometryService_test

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/mock"
)

// The fixtures below place nodes on a 700px canvas at the unscaled zone radius,
// so every expected coordinate in these tests is an exact integer.
const (
	fixtureCanvasSide = 700
	fixtureZoneRadius = 38
)

// newGeometryFixture builds the service over a preview layout stub that always
// reports the given node positions, so the tests pin the curve math rather than
// the node placement (which the preview layout service owns and tests).
func newGeometryFixture(
	positions map[string]image.Point,
) (connection_editor.IZoneEditorGeometryService, *test_helpers.PreviewLayoutServiceMock) {
	previewLayout := &test_helpers.PreviewLayoutServiceMock{}
	previewLayout.
		On("BuildPreviewLayout", mock.Anything, mock.Anything, mock.Anything).
		Return(preview.Layout{Positions: positions, ZoneRadius: fixtureZoneRadius})

	return connection_editor.NewZoneEditorGeometryService(previewLayout), previewLayout
}

// chordPositions places two zones on a horizontal chord with nothing between
// them, so a lone connection has no reason to bend.
func chordPositions() map[string]image.Point {
	return map[string]image.Point{
		"A": image.Pt(140, 350),
		"B": image.Pt(560, 350),
	}
}

// trianglePositions adds a third zone well clear of the A-B chord.
func trianglePositions() map[string]image.Point {
	positions := chordPositions()
	positions["C"] = image.Pt(350, 140)

	return positions
}

// obstaclePositions parks a zone 14px off the middle of the A-B chord, inside
// the clearance a curve must keep from a node.
func obstaclePositions() map[string]image.Point {
	positions := chordPositions()
	positions["D"] = image.Pt(350, 364)

	return positions
}

func newConnection(name, from, to string) entities.Connection {
	return entities.Connection{Name: name, From: from, To: to}
}

func connectionIndices(edges []models.ZoneEditorEdge) []int {
	indices := make([]int, 0, len(edges))
	for _, edge := range edges {
		indices = append(indices, edge.ConnectionIndex)
	}

	return indices
}
