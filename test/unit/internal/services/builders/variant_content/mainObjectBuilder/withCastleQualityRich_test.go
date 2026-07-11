package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenRichCastleQualityIsChosen_SetsRichConstructionSidOnBuiltObject(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithCastleQualityRich().Build()

	// Assert
	assert.Equal(t, entities.MainObject{BuildingsConstructionSid: "rich_buildings_construction"}, mainObject)
}
