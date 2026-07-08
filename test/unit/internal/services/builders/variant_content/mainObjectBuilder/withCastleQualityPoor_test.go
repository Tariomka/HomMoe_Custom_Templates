package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenPoorCastleQualityIsChosen_SetsPoorConstructionSidOnBuiltObject(t *testing.T) {
	// Arrange
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithCastleQualityPoor().Build()

	// Assert
	assert.Equal(t, entities.MainObject{BuildingsConstructionSid: "poor_buildings_construction"}, mainObject)
}
