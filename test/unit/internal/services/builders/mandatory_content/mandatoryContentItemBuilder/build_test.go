package mandatoryContentItemBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/mandatory_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMultipleOptionsAreChained_ReturnsItemWithAllAccumulatedValues(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedSid := gofakeit.Word()
	expectedName := gofakeit.Word()
	expectedRule := template_model.PlacementRule{Type: "Road", Weight: gofakeit.Number(1, 100)}
	builder := mandatory_content.NewContentItemBuilder(expectedSid)

	// Act
	item := builder.
		WithName(expectedName).
		WithGuarded().
		WithMine().
		WithSoloEncounter().
		WithRules(expectedRule).
		Build()

	// Assert
	assert.Equal(t, template_model.MandatoryContentItem{
		SID:           expectedSid,
		Name:          expectedName,
		IsGuarded:     true,
		IsMine:        true,
		SoloEncounter: true,
		Rules:         []template_model.PlacementRule{expectedRule},
	}, item)
}
