package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneIsRemoved_ReturnsTheRemainingZones(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	request := removeRequest()
	remainingZones := []entities.Zone{{Name: gofakeit.Word()}}
	fixture.zoneEditor.
		On("RemoveZone", request.Zones, request.Connections, request.ZoneName).
		Return(remainingZones, []entities.Connection{})

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
	remainingConnections := []entities.Connection{{From: gofakeit.Word()}}
	fixture.zoneEditor.
		On("RemoveZone", request.Zones, request.Connections, request.ZoneName).
		Return([]entities.Zone{}, remainingConnections)

	// Act
	mutation := fixture.handler.RemoveZoneEditorZone(request)

	// Assert
	assert.Equal(t, remainingConnections, mutation.Connections)
}

func removeRequest() dtos.ZoneEditorRemoveRequestDto {
	return dtos.ZoneEditorRemoveRequestDto{
		Zones:       []entities.Zone{{Name: gofakeit.Word()}},
		Connections: []entities.Connection{{From: gofakeit.Word()}},
		ZoneName:    gofakeit.Word(),
	}
}
