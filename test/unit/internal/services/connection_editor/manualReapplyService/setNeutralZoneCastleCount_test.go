package manualReapplyService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenCountIncreases_RebuildsRequestedCastleCount(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := makeNeutralZone("G", neutral_zone.QualityHigh, 1)

	// Act
	newManualReapplyService().SetNeutralZoneCastleCount(&zone, 3, defaultTuning())

	// Assert
	assert.Equal(t, 3, test_helpers.NewZoneEditorService().CountZoneCastles(zone))
}

func TestWhenCastlesAreRebuilt_KeepsQualityProfile(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := makeNeutralZone("G", neutral_zone.QualityHigh, 1)

	// Act
	newManualReapplyService().SetNeutralZoneCastleCount(&zone, 3, defaultTuning())

	// Assert
	assert.Equal(t, neutral_zone.QualityHigh, zone_services.NewZoneTierService().GetQuality(zone))
}

func TestWhenCastlesAreRebuilt_KeepsGuardMultiplier(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := makeNeutralZone("G", neutral_zone.QualityHigh, 1)
	originalGuardMultiplier := zone.GuardMultiplier

	// Act
	newManualReapplyService().SetNeutralZoneCastleCount(&zone, 3, defaultTuning())

	// Assert
	assert.InDelta(t, originalGuardMultiplier, zone.GuardMultiplier, test_helpers.Delta,
		"guard multiplier must not be re-profiled")
}

func TestWhenCastlesAreRebuilt_KeepsGuardedContentPool(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := makeNeutralZone("G", neutral_zone.QualityHigh, 1)
	originalPool := zone.GuardedContentPool

	// Act
	newManualReapplyService().SetNeutralZoneCastleCount(&zone, 3, defaultTuning())

	// Assert
	assert.Equal(t, originalPool, zone.GuardedContentPool, "content pools must not be re-profiled")
}

func TestWhenCastlesAreRebuilt_KeepsGuardedContentValue(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := makeNeutralZone("G", neutral_zone.QualityHigh, 1)
	originalGuardedValue := zone.GuardedContentValue

	// Act
	newManualReapplyService().SetNeutralZoneCastleCount(&zone, 3, defaultTuning())

	// Assert
	assert.Equal(t, originalGuardedValue, zone.GuardedContentValue, "content values must not be re-profiled")
}

func TestWhenZoneSizeWasEditedManually_KeepsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := makeNeutralZone("G", neutral_zone.QualityHigh, 1)
	manualSize := gofakeit.Float64Range(0.5, 3.0)
	zone.Size = manualSize

	// Act
	newManualReapplyService().SetNeutralZoneCastleCount(&zone, 3, defaultTuning())

	// Assert
	assert.InDelta(t, manualSize, zone.Size, test_helpers.Delta)
}

func TestWhenZoneHasAbandonedOutpost_PreservesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := makeNeutralZone("G", neutral_zone.QualityMedium, 1)
	zone.MainObjects = append(zone.MainObjects, entities.MainObject{Type: "AbandonedOutpost"})

	// Act
	newManualReapplyService().SetNeutralZoneCastleCount(&zone, 2, defaultTuning())

	// Assert
	outpostCount := 0
	for _, mainObject := range zone.MainObjects {
		if mainObject.Type == "AbandonedOutpost" {
			outpostCount++
		}
	}
	assert.Equal(t, 1, outpostCount, "abandoned outposts must survive a castle rebuild")
}

func TestWhenPrimaryCastleHoldsWinCondition_PreservesHoldCityFlag(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := makeNeutralZone("G", neutral_zone.QualityHigh, 1)
	zone.MainObjects[0].HoldCityWinCon = true

	// Act
	newManualReapplyService().SetNeutralZoneCastleCount(&zone, 2, defaultTuning())

	// Assert
	require.Equal(t, 2, test_helpers.NewZoneEditorService().CountZoneCastles(zone))
	assert.True(t, zone.MainObjects[0].HoldCityWinCon, "hold-city win condition was lost by the rebuild")
}
