package zoneFactory_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenSpawnZoneIsCreated_UsesPlayerLabelName(t *testing.T) {
	t.Parallel()
	// Arrange
	factory := newZoneFactory()

	// Act
	zone := factory.CreateSpawnZone(models.SpawnZoneCreationRequest{
		Label:         "A",
		PlayerName:    "Player1",
		Size:          1,
		GenerateRoads: true,
		Tuning:        newUnitTuning(),
	})

	// Assert
	assert.Equal(t, "Spawn-A", zone.Name)
}

// A spawn zone has no tier; nil is the "infer it" signal the tier service reads.
func TestWhenSpawnZoneIsCreated_LeavesTheQualityUnrecorded(t *testing.T) {
	t.Parallel()
	// Arrange
	factory := newZoneFactory()

	// Act
	zone := factory.CreateSpawnZone(models.SpawnZoneCreationRequest{
		Label:      "A",
		PlayerName: "Player1",
		Size:       1,
		Tuning:     newUnitTuning(),
	})

	// Assert
	assert.Nil(t, zone.Quality)
}
