package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
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
	assert.Equal(t,
		template_model.MainObject{BuildingsConstructionSid: "ultra_rich_buildings_construction"},
		mainObject)
}
