package zoneContentEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zone_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoVariantRuleApplies_TheRowKeepsItsPlainName(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	name := gofakeit.Word()

	// Act
	displayName := service.GetContentRowDisplayName(name, []dtos.ContentRuleDescriptionDto{
		{Key: dtos.ContentRuleKeyGuarded, Valid: true},
	})

	// Assert
	assert.Equal(t, name, displayName)
}

func TestWhenTheVariantRuleIsInvalid_TheRowKeepsItsPlainName(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	name := gofakeit.Word()

	// Act
	displayName := service.GetContentRowDisplayName(name, []dtos.ContentRuleDescriptionDto{
		{Key: dtos.ContentRuleKeyVariant, VariantLabel: gofakeit.Word()},
	})

	// Assert
	assert.Equal(t, name, displayName)
}

func TestWhenAVariantRuleApplies_TheRowNameCarriesTheVariantLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	name := gofakeit.Word()
	label := gofakeit.Word()

	// Act
	displayName := service.GetContentRowDisplayName(name, []dtos.ContentRuleDescriptionDto{
		{Key: dtos.ContentRuleKeyVariant, Valid: true, VariantLabel: label},
	})

	// Assert
	assert.Equal(t, name+" ("+label+")", displayName)
}
