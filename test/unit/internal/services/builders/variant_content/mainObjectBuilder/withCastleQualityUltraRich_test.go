package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenUltraRichCastleQualityIsChosen_SetsUltraRichConstructionSidOnBuiltObject(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithCastleQualityUltraRich().Build()

	// Assert
	assert.Equal(t, entities.MainObject{BuildingsConstructionSid: "ultra_rich_buildings_construction"}, mainObject)
}
