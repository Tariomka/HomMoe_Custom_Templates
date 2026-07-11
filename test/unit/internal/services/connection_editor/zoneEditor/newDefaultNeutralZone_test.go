package zoneEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenLabelIsGiven_NamesZoneNeutralLabel(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	zone := connection_editor.NewDefaultNeutralZone("Q", models.QualityMedium, 1, false, defaultTuning())

	// Assert
	assert.Equal(t, "Neutral-Q", zone.Name)
}

func TestWhenZoneIsCreatedManually_ClearsMandatoryContentReference(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	zone := connection_editor.NewDefaultNeutralZone("Q", models.QualityMedium, 1, false, defaultTuning())

	// Assert
	assert.Nil(t, zone.MandatoryContent)
}

func TestWhenCastleCountIsOne_CreatesOneCastle(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	zone := connection_editor.NewDefaultNeutralZone("Q", models.QualityMedium, 1, false, defaultTuning())

	// Assert
	assert.Equal(t, 1, connection_editor.CountZoneCastles(zone))
}

func TestWhenCastleCountIsZero_CreatesNoCastles(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	zone := connection_editor.NewDefaultNeutralZone("R", models.QualityLow, 0, false, defaultTuning())

	// Assert
	assert.Equal(t, 0, connection_editor.CountZoneCastles(zone))
}

func TestWhenQualityIsRequested_ProfilesZoneWithThatQuality(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		quality models.NeutralZoneQuality
	}{
		{"WhenQualityIsLow_ProfilesZoneAsLow", models.QualityLow},
		{"WhenQualityIsMedium_ProfilesZoneAsMedium", models.QualityMedium},
		{"WhenQualityIsHigh_ProfilesZoneAsHigh", models.QualityHigh},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			zone := connection_editor.NewDefaultNeutralZone("Z", testCase.quality, 1, false, defaultTuning())

			// Assert
			assert.Equal(t, testCase.quality, connection_editor.QualityOfZone(zone))
		})
	}
}
