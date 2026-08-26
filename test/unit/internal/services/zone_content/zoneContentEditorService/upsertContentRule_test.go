package zoneContentEditorService_test

import (
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zone_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoRuleOfThatTypeExists_TheRuleIsAppended(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	existing := editor_state_model.ContentRuleRow{Name: gofakeit.Word()}
	added := editor_state_model.ContentRuleRow{Name: gofakeit.Sentence(2)}

	// Act
	rules := service.UpsertContentRule([]editor_state_model.ContentRuleRow{existing}, added)

	// Assert
	assert.Equal(t, []editor_state_model.ContentRuleRow{existing, added}, rules)
}

func TestWhenARuleOfThatTypeExists_ItIsReplacedInPlace(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	name := gofakeit.Word()
	replacement := editor_state_model.ContentRuleRow{Name: name, DistanceName: gofakeit.Word()}

	// Act
	rules := service.UpsertContentRule([]editor_state_model.ContentRuleRow{{Name: name}}, replacement)

	// Assert
	assert.Equal(t, []editor_state_model.ContentRuleRow{replacement}, rules)
}

func TestWhenTheExistingRuleNameDiffersOnlyByCase_ItIsStillReplaced(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	name := gofakeit.Word()
	replacement := editor_state_model.ContentRuleRow{Name: strings.ToUpper(name)}

	// Act
	rules := service.UpsertContentRule([]editor_state_model.ContentRuleRow{{Name: name}}, replacement)

	// Assert
	assert.Equal(t, []editor_state_model.ContentRuleRow{replacement}, rules)
}
