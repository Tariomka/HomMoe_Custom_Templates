package templateMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A raw .rmg.json records no tier, so every lifted zone must come back with a
// nil Quality - "not recorded, infer it". A value field would claim Plastic.
func TestWhenATemplateIsLifted_NoZoneCarriesARecordedTier(t *testing.T) {
	t.Parallel()
	// Arrange
	mapper := mappers.NewTemplateMapper()
	entity := test_helpers.NewAllFieldsTemplate()
	require.NotEmpty(t, entity.Variants)
	require.NotEmpty(t, entity.Variants[0].Zones)

	// Act
	model := mapper.ToModel(entity)

	// Assert
	var tieredZoneNames []string
	for _, zone := range model.Variants[0].Zones {
		if zone.Quality != nil {
			tieredZoneNames = append(tieredZoneNames, zone.Name)
		}
	}
	assert.Empty(t, tieredZoneNames)
}

func TestWhenATemplateIsLifted_TheZoneNamesAreCarriedAcross(t *testing.T) {
	t.Parallel()
	// Arrange
	mapper := mappers.NewTemplateMapper()
	entity := test_helpers.NewAllFieldsTemplate()
	require.NotEmpty(t, entity.Variants)
	var expected []string
	for _, zone := range entity.Variants[0].Zones {
		expected = append(expected, zone.Name)
	}

	// Act
	model := mapper.ToModel(entity)

	// Assert
	var actual []string
	for _, zone := range model.Variants[0].Zones {
		actual = append(actual, zone.Name)
	}
	assert.Equal(t, expected, actual)
}
