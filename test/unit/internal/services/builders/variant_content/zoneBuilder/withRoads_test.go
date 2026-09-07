package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenRoadsAreProvided_SetsRoadsOnBuiltZone(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedRoads := []template_model.Road{
		{Type: "Stone", From: template_model.TypedRef{Type: gofakeit.Word()}},
		{Type: "Dirt", To: template_model.TypedRef{Type: gofakeit.Word()}},
	}
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithRoads(expectedRoads).Build()

	// Assert
	assert.Equal(t, template_model.Zone{Roads: expectedRoads}, zone)
}
