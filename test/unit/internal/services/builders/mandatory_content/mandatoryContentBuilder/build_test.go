package mandatoryContentBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/mandatory_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMultipleOptionsAreChained_ReturnsItemWithAllAccumulatedValues(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedSid := gofakeit.Word()
	expectedName := gofakeit.Word()
	expectedRule := entities.PlacementRule{Type: "Road", Weight: gofakeit.Number(1, 100)}
	builder := mandatory_content.NewContentBuilder(expectedSid)

	// Act
	item := builder.
		WithName(expectedName).
		WithGuarded().
		WithMine().
		WithSoloEncounter().
		WithRules(expectedRule).
		Build()

	// Assert
	assert.Equal(t, entities.MandatoryContentItem{
		SID:           expectedSid,
		Name:          expectedName,
		IsGuarded:     true,
		IsMine:        true,
		SoloEncounter: true,
		Rules:         []entities.PlacementRule{expectedRule},
	}, item)
}
