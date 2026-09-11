package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenConnectionIsCreated_ReturnsTheEditorsConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	request := newConnectionRequest()
	expected := template_model.Connection{From: request.From, To: request.To}
	fixture.connectionEditor.
		On("NewDefaultConnection", request.From, request.To, request.Zones, request.PlayerZoneNames).
		Return(expected)
	fixture.zoneEditor.On("ApplyConnectionRoadPolicy", mock.Anything, mock.Anything).Return()

	// Act
	connection := fixture.handler.CreateZoneEditorConnection(request)

	// Assert
	assert.Equal(t, expected, connection)
}

func TestWhenConnectionIsCreated_AppliesTheRoadSettingToIt(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	request := newConnectionRequest()
	request.GenerateRoads = true
	fixture.connectionEditor.
		On("NewDefaultConnection", request.From, request.To, request.Zones, request.PlayerZoneNames).
		Return(template_model.Connection{From: request.From, To: request.To})
	fixture.zoneEditor.On("ApplyConnectionRoadPolicy", mock.Anything, true).Return()

	// Act
	fixture.handler.CreateZoneEditorConnection(request)

	// Assert
	fixture.zoneEditor.AssertExpectations(t)
}

func TestWhenTheCreatedConnectionIsStamped_ReturnsTheStampedRoadFlag(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	request := newConnectionRequest()
	fixture.connectionEditor.
		On("NewDefaultConnection", request.From, request.To, request.Zones, request.PlayerZoneNames).
		Return(template_model.Connection{From: request.From, To: request.To})
	fixture.zoneEditor.
		On("ApplyConnectionRoadPolicy", mock.Anything, mock.Anything).
		Run(func(arguments mock.Arguments) {
			connection, _ := arguments.Get(0).(*template_model.Connection)
			connection.Road = new(false)
		}).
		Return()

	// Act
	connection := fixture.handler.CreateZoneEditorConnection(request)

	// Assert
	assert.Equal(t, new(false), connection.Road)
}

func newConnectionRequest() dtos.ZoneEditorConnectionRequestDto {
	return dtos.ZoneEditorConnectionRequestDto{
		From:            gofakeit.Word(),
		To:              gofakeit.Word(),
		Zones:           []template_model.Zone{{Name: gofakeit.Word()}},
		PlayerZoneNames: map[string]bool{gofakeit.Word(): true},
	}
}
