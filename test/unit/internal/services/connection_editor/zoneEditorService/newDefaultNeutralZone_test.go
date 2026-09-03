package zoneEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenLabelIsGiven_NamesZoneNeutralLabel(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	zone := test_helpers.NewZoneEditorService().
		NewDefaultNeutralZone("Q", neutral_zone.QualityMedium, 1, false, defaultTuning())

	// Assert
	assert.Equal(t, "Neutral-Q", zone.Name)
}

func TestWhenZoneIsCreatedManually_ClearsMandatoryContentReference(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	zone := test_helpers.NewZoneEditorService().
		NewDefaultNeutralZone("Q", neutral_zone.QualityMedium, 1, false, defaultTuning())

	// Assert
	assert.Nil(t, zone.MandatoryContent)
}

func TestWhenCastleCountIsOne_CreatesOneCastle(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	zone := test_helpers.NewZoneEditorService().
		NewDefaultNeutralZone("Q", neutral_zone.QualityMedium, 1, false, defaultTuning())

	// Assert
	assert.Equal(t, 1, test_helpers.NewZoneEditorService().CountZoneCastles(zone))
}

func TestWhenCastleCountIsZero_CreatesNoCastles(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	zone := test_helpers.NewZoneEditorService().
		NewDefaultNeutralZone("R", neutral_zone.QualityLow, 0, false, defaultTuning())

	// Assert
	assert.Equal(t, 0, test_helpers.NewZoneEditorService().CountZoneCastles(zone))
}

func TestWhenQualityIsRequested_ProfilesZoneWithThatQuality(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		quality neutral_zone.Quality
	}{
		{"WhenQualityIsLow_ProfilesZoneAsLow", neutral_zone.QualityLow},
		{"WhenQualityIsMedium_ProfilesZoneAsMedium", neutral_zone.QualityMedium},
		{"WhenQualityIsHigh_ProfilesZoneAsHigh", neutral_zone.QualityHigh},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			zone := test_helpers.NewZoneEditorService().
				NewDefaultNeutralZone("Z", testCase.quality, 1, false, defaultTuning())

			// Assert
			assert.Equal(t,
				testCase.quality,
				zone_services.NewZoneTierService().GetQuality(template_model.ToZoneEntity(zone)))
		})
	}
}

func TestWhenQualityIsRequested_RecordsThatTierOnTheZone(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	zone := test_helpers.NewZoneEditorService().
		NewDefaultNeutralZone("Z", neutral_zone.QualityHigh, 1, false, defaultTuning())

	// Assert
	assert.Equal(t, neutral_zone.QualityHigh, *zone.Quality)
}
