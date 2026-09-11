package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenTypeChangeIsRequested_ForwardsItToTheEditor(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	request := dtos.ZoneEditorConnectionTypeRequestDto{
		Connection:     template_model.Connection{From: gofakeit.Word(), To: gofakeit.Word()},
		ConnectionType: "Portal",
		GenerateRoads:  true,
	}
	fixture.zoneEditor.
		On("ChangeConnectionType", mock.Anything, request.ConnectionType, request.GenerateRoads).
		Return()

	// Act
	fixture.handler.ChangeZoneEditorConnectionType(request)

	// Assert
	fixture.zoneEditor.AssertExpectations(t)
}

func TestWhenTypeChangeIsApplied_ReturnsTheEditedConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	request := dtos.ZoneEditorConnectionTypeRequestDto{
		Connection:     template_model.Connection{From: "Spawn-A", To: "Spawn-B", ConnectionType: "Direct"},
		ConnectionType: "Portal",
	}
	fixture.zoneEditor.
		On("ChangeConnectionType", mock.Anything, mock.Anything, mock.Anything).
		Run(func(arguments mock.Arguments) {
			connection, _ := arguments.Get(0).(*template_model.Connection)
			connection.ConnectionType = "Portal"
			connection.Road = new(false)
		}).
		Return()

	// Act
	connection := fixture.handler.ChangeZoneEditorConnectionType(request)

	// Assert
	assert.Equal(
		t,
		template_model.Connection{
			From:           "Spawn-A",
			To:             "Spawn-B",
			ConnectionType: "Portal",
			Road:           new(false),
		},
		connection)
}

func TestWhenTypeChangeIsApplied_DoesNotShareThePlacementRulesWithTheRequest(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	request := dtos.ZoneEditorConnectionTypeRequestDto{
		Connection: template_model.Connection{
			From:                     "Spawn-A",
			To:                       "Spawn-B",
			ConnectionType:           "Portal",
			PortalPlacementRulesFrom: []template_model.PlacementRule{{Type: "MainObject"}},
		},
		ConnectionType: "Direct",
	}
	fixture.zoneEditor.
		On("ChangeConnectionType", mock.Anything, mock.Anything, mock.Anything).
		Run(func(arguments mock.Arguments) {
			connection, _ := arguments.Get(0).(*template_model.Connection)
			connection.PortalPlacementRulesFrom[0].Type = "Border"
		}).
		Return()

	// Act
	fixture.handler.ChangeZoneEditorConnectionType(request)

	// Assert
	assert.Equal(t, "MainObject", request.Connection.PortalPlacementRulesFrom[0].Type)
}
