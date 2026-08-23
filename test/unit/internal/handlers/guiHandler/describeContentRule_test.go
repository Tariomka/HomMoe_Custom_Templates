package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/stretchr/testify/assert"
)

func TestWhenSavedRuleIsValid_ReturnsDisplayTextAndMarker(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	savedRule := models.ContentRuleRow{Name: "Distance to road", DistanceName: "Far"}
	expected := dtos.ContentRuleDescriptionDto{
		Key:         dtos.ContentRuleKeyDistanceToRoad,
		DisplayText: "Distance to road: Far",
		Marker:      "R",
		Valid:       true,
		SavedRule:   savedRule,
	}

	// Act
	result := handler.DescribeContentRule(models.SidMapping{}, savedRule)

	// Assert
	assert.Equal(t, expected, result)
}

func TestWhenBooleanRuleIsFalse_ReturnsNegatedMarker(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	value := false
	savedRule := models.ContentRuleRow{Name: "Guarded", IsGuarded: &value}

	// Act
	result := handler.DescribeContentRule(models.SidMapping{}, savedRule)

	// Assert
	assert.Equal(t, "!G", result.Marker)
}

func TestWhenVariantRuleIsValid_ReturnsVariantLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	variantID := 2
	content := models.SidMapping{
		Sid:  registry.GetMapObjectT3GuardedResourceBankValues().DragonUtopia,
		Name: "Dragon Utopia",
	}
	savedRule := models.ContentRuleRow{Name: "Variant", VariantID: &variantID}
	expected := dtos.ContentRuleDescriptionDto{
		Key:          dtos.ContentRuleKeyVariant,
		DisplayText:  "Variant: Large Guard",
		VariantLabel: "Large Guard",
		Valid:        true,
		SavedRule:    savedRule,
	}

	// Act
	result := handler.DescribeContentRule(content, savedRule)

	// Assert
	assert.Equal(t, expected, result)
}

func TestWhenSavedRuleIsInvalid_ReturnsFallbackDescription(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	savedRule := models.ContentRuleRow{Name: "Unknown rule"}
	expected := dtos.ContentRuleDescriptionDto{
		DisplayText: "Unknown rule",
		SavedRule:   savedRule,
	}

	// Act
	result := handler.DescribeContentRule(models.SidMapping{}, savedRule)

	// Assert
	assert.Equal(t, expected, result)
}
