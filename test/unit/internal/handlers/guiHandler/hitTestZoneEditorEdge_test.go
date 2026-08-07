package guiHandler_test

import (
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenAnEdgeHitTestRequested_ReturnsTheCurveUnderThePoint(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	edges := []models.ZoneEditorEdge{{
		ConnectionIndex: 0,
		StartPoint:      data.NewVec2(0.0, 0.0),
		ControlPoint:    data.NewVec2(50.0, 0.0),
		EndPoint:        data.NewVec2(100.0, 0.0),
	}}

	// Act
	index := handler.HitTestZoneEditorEdge(image.Pt(50, 5), edges)

	// Assert
	assert.Equal(t, 0, index)
}
