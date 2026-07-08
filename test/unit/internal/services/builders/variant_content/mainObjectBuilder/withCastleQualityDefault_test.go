package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenDefaultCastleQualityIsChosen_SetsDefaultConstructionSidOnBuiltObject(t *testing.T) {
	// Arrange
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithCastleQualityDefault().Build()

	// Assert
	assert.Equal(t, entities.MainObject{BuildingsConstructionSid: "default_buildings_construction"}, mainObject)
}
