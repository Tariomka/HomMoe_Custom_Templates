package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenConnectionIsCreated_ReturnsTheEditorsConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	request := dtos.ZoneEditorConnectionRequestDto{
		From:            gofakeit.Word(),
		To:              gofakeit.Word(),
		Zones:           []template_model.Zone{{Name: gofakeit.Word()}},
		PlayerZoneNames: map[string]bool{gofakeit.Word(): true},
	}
	expected := entities.Connection{From: request.From, To: request.To}
	fixture.connectionEditor.
		On("NewDefaultConnection", request.From, request.To, request.Zones, request.PlayerZoneNames).
		Return(expected)

	// Act
	connection := fixture.handler.CreateZoneEditorConnection(request)

	// Assert
	assert.Equal(t, expected, connection)
}
