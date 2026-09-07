package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneIsRemoved_ReturnsTheRemainingZones(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	request := removeRequest()
	remainingZones := []template_model.Zone{{Name: gofakeit.Word()}}
	fixture.zoneEditor.
		On("RemoveZone", request.Zones, request.Connections, request.ZoneName).
		Return(remainingZones, []template_model.Connection{})

	// Act
	mutation := fixture.handler.RemoveZoneEditorZone(request)

	// Assert
	assert.Equal(t, remainingZones, mutation.Zones)
}

func TestWhenZoneIsRemoved_ReturnsTheRemainingConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	request := removeRequest()
	remainingConnections := []template_model.Connection{{From: gofakeit.Word()}}
	fixture.zoneEditor.
		On("RemoveZone", request.Zones, request.Connections, request.ZoneName).
		Return([]template_model.Zone{}, remainingConnections)

	// Act
	mutation := fixture.handler.RemoveZoneEditorZone(request)

	// Assert
	assert.Equal(t, remainingConnections, mutation.Connections)
}

func removeRequest() dtos.ZoneEditorRemoveRequestDto {
	return dtos.ZoneEditorRemoveRequestDto{
		Zones:       []template_model.Zone{{Name: gofakeit.Word()}},
		Connections: []template_model.Connection{{From: gofakeit.Word()}},
		ZoneName:    gofakeit.Word(),
	}
}
