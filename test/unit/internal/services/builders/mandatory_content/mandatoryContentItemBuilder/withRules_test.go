package mandatoryContentItemBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/mandatory_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenRulesAreProvidedTwice_AppendsAllRulesOnBuiltItem(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedSid := gofakeit.Word()
	firstRule := template_model.PlacementRule{Type: "Road", Weight: gofakeit.Number(1, 100)}
	secondRule := template_model.PlacementRule{Type: "Crossroads", Weight: gofakeit.Number(1, 100)}
	thirdRule := template_model.PlacementRule{Type: "MainObject", Weight: gofakeit.Number(1, 100)}
	builder := mandatory_content.NewContentItemBuilder(expectedSid)

	// Act
	item := builder.WithRules(firstRule, secondRule).WithRules(thirdRule).Build()

	// Assert
	assert.Equal(t, template_model.MandatoryContentItem{
		SID:   expectedSid,
		Rules: []template_model.PlacementRule{firstRule, secondRule, thirdRule},
	}, item)
}
