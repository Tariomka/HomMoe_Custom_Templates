package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenANodeHitTestRequested_NamesTheZoneUnderThePoint(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	request := dtos.ZoneEditorHitTestRequestDto{
		Position:   data.NewVec2(150.0, 350.0),
		Positions:  map[string]models.Position{"Spawn-A": data.NewVec2(140.0, 350.0)},
		ZoneRadius: 38,
	}

	// Act
	name := handler.HitTestZoneEditorNode(request)

	// Assert
	assert.Equal(t, "Spawn-A", name)
}
