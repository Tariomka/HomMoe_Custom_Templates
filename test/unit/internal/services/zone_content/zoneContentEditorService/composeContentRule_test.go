package zoneContentEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zone_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTheRuleTypeIsUnknown_TheCompositionIsInvalid(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()

	// Act
	result := service.ComposeContentRule(dtos.ContentRuleCompositionRequestDto{})

	// Assert
	assert.False(t, result.Valid)
}

func TestWhenTheDistanceIndexIsOutOfRange_TheCompositionIsInvalid(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()

	// Act
	result := service.ComposeContentRule(dtos.ContentRuleCompositionRequestDto{
		Option:        dtos.ContentRuleOptionDto{Key: dtos.ContentRuleKeyDistanceToRoad},
		DistanceNames: []string{gofakeit.Word()},
		DistanceIndex: 1,
	})

	// Assert
	assert.False(t, result.Valid)
}

func TestWhenADistanceToRoadRuleIsComposed_ItCarriesTheSelectedDistance(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	name := gofakeit.Word()
	distance := gofakeit.Word()

	// Act
	result := service.ComposeContentRule(dtos.ContentRuleCompositionRequestDto{
		Option:        dtos.ContentRuleOptionDto{Key: dtos.ContentRuleKeyDistanceToRoad, Name: name},
		DistanceNames: []string{gofakeit.Word(), distance},
		DistanceIndex: 1,
	})

	// Assert
	assert.Equal(t, models.ContentRuleRow{Name: name, DistanceName: distance}, result.Rule)
}

func TestWhenADistanceToTownRuleIsComposed_ItCarriesTheSelectedDistance(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	name := gofakeit.Word()
	distance := gofakeit.Word()

	// Act
	result := service.ComposeContentRule(dtos.ContentRuleCompositionRequestDto{
		Option:        dtos.ContentRuleOptionDto{Key: dtos.ContentRuleKeyDistanceToTown, Name: name},
		DistanceNames: []string{distance},
	})

	// Assert
	assert.Equal(t, models.ContentRuleRow{Name: name, DistanceName: distance}, result.Rule)
}

func TestWhenAGuardedRuleIsComposed_ItCarriesTheCheckboxValue(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	name := gofakeit.Word()

	// Act
	result := service.ComposeContentRule(dtos.ContentRuleCompositionRequestDto{
		Option:    dtos.ContentRuleOptionDto{Key: dtos.ContentRuleKeyGuarded, Name: name},
		IsGuarded: true,
	})

	// Assert
	require.NotNil(t, result.Rule.IsGuarded)
	assert.True(t, *result.Rule.IsGuarded)
}

func TestWhenASoloEncounterRuleIsComposed_ItCarriesTheCheckboxValue(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	name := gofakeit.Word()

	// Act
	result := service.ComposeContentRule(dtos.ContentRuleCompositionRequestDto{
		Option:          dtos.ContentRuleOptionDto{Key: dtos.ContentRuleKeySoloEncounter, Name: name},
		IsSoloEncounter: true,
	})

	// Assert
	require.NotNil(t, result.Rule.IsSoloEncounter)
	assert.True(t, *result.Rule.IsSoloEncounter)
}

func TestWhenTheVariantIndexIsOutOfRange_TheCompositionIsInvalid(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()

	// Act
	result := service.ComposeContentRule(dtos.ContentRuleCompositionRequestDto{
		Option:       dtos.ContentRuleOptionDto{Key: dtos.ContentRuleKeyVariant},
		VariantIndex: -1,
	})

	// Assert
	assert.False(t, result.Valid)
}

func TestWhenAVariantRuleIsComposed_ItCarriesTheSelectedVariantId(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	variantID := gofakeit.Number(1, 500)

	// Act
	result := service.ComposeContentRule(dtos.ContentRuleCompositionRequestDto{
		Option:       dtos.ContentRuleOptionDto{Key: dtos.ContentRuleKeyVariant, Name: gofakeit.Word()},
		VariantIDs:   []int{gofakeit.Number(501, 900), variantID},
		VariantIndex: 1,
	})

	// Assert
	require.NotNil(t, result.Rule.VariantID)
	assert.Equal(t, variantID, *result.Rule.VariantID)
}
