package zoneContentEditorService_test

import (
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zone_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoRuleOfThatTypeExists_TheRuleIsAppended(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	existing := models.ContentRuleRow{Name: gofakeit.Word()}
	added := models.ContentRuleRow{Name: gofakeit.Sentence(2)}

	// Act
	rules := service.UpsertContentRule([]models.ContentRuleRow{existing}, added)

	// Assert
	assert.Equal(t, []models.ContentRuleRow{existing, added}, rules)
}

func TestWhenARuleOfThatTypeExists_ItIsReplacedInPlace(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	name := gofakeit.Word()
	replacement := models.ContentRuleRow{Name: name, DistanceName: gofakeit.Word()}

	// Act
	rules := service.UpsertContentRule([]models.ContentRuleRow{{Name: name}}, replacement)

	// Assert
	assert.Equal(t, []models.ContentRuleRow{replacement}, rules)
}

func TestWhenTheExistingRuleNameDiffersOnlyByCase_ItIsStillReplaced(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	name := gofakeit.Word()
	replacement := models.ContentRuleRow{Name: strings.ToUpper(name)}

	// Act
	rules := service.UpsertContentRule([]models.ContentRuleRow{{Name: name}}, replacement)

	// Assert
	assert.Equal(t, []models.ContentRuleRow{replacement}, rules)
}
