package mandatoryContentBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/mandatory_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenRulesAreProvidedTwice_AppendsAllRulesOnBuiltItem(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedSid := gofakeit.Word()
	firstRule := entities.PlacementRule{Type: "Road", Weight: gofakeit.Number(1, 100)}
	secondRule := entities.PlacementRule{Type: "Crossroads", Weight: gofakeit.Number(1, 100)}
	thirdRule := entities.PlacementRule{Type: "MainObject", Weight: gofakeit.Number(1, 100)}
	builder := mandatory_content.NewContentBuilder(expectedSid)

	// Act
	item := builder.WithRules(firstRule, secondRule).WithRules(thirdRule).Build()

	// Assert
	assert.Equal(t, entities.MandatoryContentItem{
		SID:   expectedSid,
		Rules: []entities.PlacementRule{firstRule, secondRule, thirdRule},
	}, item)
}
