package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestWhenGraphIsDescribed_ReturnsServiceEquivalentStatus(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	zones := []entities.Zone{{Name: "A"}, {Name: "B"}, {Name: "C"}}
	connections := []entities.Connection{{From: "A", To: "Missing"}}
	expected := dtos.ZoneEditorGraphDto{HasErrors: true, IsolatedZoneCount: 2}

	// Act
	result := handler.DescribeZoneEditorGraph(zones, connections)

	// Assert
	assert.Equal(t, expected, result)
}
