package guiHandler_test

import (
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func TestWhenANodeHitTestRequested_NamesTheZoneUnderThePoint(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	request := dtos.ZoneEditorHitTestRequestDto{
		Position:   image.Pt(150, 350),
		Positions:  map[string]image.Point{"Spawn-A": image.Pt(140, 350)},
		ZoneRadius: 38,
	}

	// Act
	name := handler.HitTestZoneEditorNode(request)

	// Assert
	assert.Equal(t, "Spawn-A", name)
}
