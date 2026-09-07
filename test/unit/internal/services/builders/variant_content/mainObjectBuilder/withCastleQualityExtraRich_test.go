package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenExtraRichCastleQualityIsChosen_SetsExtraRichConstructionSidOnBuiltObject(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithCastleQualityExtraRich().Build()

	// Assert
	assert.Equal(t,
		template_model.MainObject{BuildingsConstructionSid: "extra_rich_buildings_construction"},
		mainObject)
}
