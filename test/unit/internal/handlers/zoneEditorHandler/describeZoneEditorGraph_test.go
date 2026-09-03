package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenGraphIsDescribed_ReturnsTheEditorsErrorFlag(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	zones := []template_model.Zone{{Name: gofakeit.Word()}}
	connections := []entities.Connection{{}}
	fixture.connectionEditor.On("ComputeHasErrors", zones, connections).Return(true)
	fixture.connectionEditor.On("FindIsolatedZones", zones, connections).Return([]string{})

	// Act
	graph := fixture.handler.DescribeZoneEditorGraph(zones, connections)

	// Assert
	assert.True(t, graph.HasErrors)
}

func TestWhenGraphIsDescribed_ReturnsTheIsolatedZoneCount(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	isolated := []string{gofakeit.Word(), gofakeit.Word(), gofakeit.Word()}
	fixture.connectionEditor.On("ComputeHasErrors", mock.Anything, mock.Anything).Return(false)
	fixture.connectionEditor.On("FindIsolatedZones", mock.Anything, mock.Anything).Return(isolated)

	// Act
	graph := fixture.handler.DescribeZoneEditorGraph(nil, nil)

	// Assert
	assert.Equal(t, len(isolated), graph.IsolatedZoneCount)
}
