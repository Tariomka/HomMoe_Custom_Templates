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
	existing := models.ContentRuleRowSave{Name: gofakeit.Word()}
	added := models.ContentRuleRowSave{Name: gofakeit.Sentence(2)}

	// Act
	rules := service.UpsertContentRule([]models.ContentRuleRowSave{existing}, added)

	// Assert
	assert.Equal(t, []models.ContentRuleRowSave{existing, added}, rules)
}

func TestWhenARuleOfThatTypeExists_ItIsReplacedInPlace(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	name := gofakeit.Word()
	replacement := models.ContentRuleRowSave{Name: name, DistanceName: gofakeit.Word()}

	// Act
	rules := service.UpsertContentRule([]models.ContentRuleRowSave{{Name: name}}, replacement)

	// Assert
	assert.Equal(t, []models.ContentRuleRowSave{replacement}, rules)
}

func TestWhenTheExistingRuleNameDiffersOnlyByCase_ItIsStillReplaced(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	name := gofakeit.Word()
	replacement := models.ContentRuleRowSave{Name: strings.ToUpper(name)}

	// Act
	rules := service.UpsertContentRule([]models.ContentRuleRowSave{{Name: name}}, replacement)

	// Assert
	assert.Equal(t, []models.ContentRuleRowSave{replacement}, rules)
}
