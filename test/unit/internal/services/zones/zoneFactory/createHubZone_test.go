package zoneFactory_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenTournamentHubIsCreated_PreservesExplicitName(t *testing.T) {
	t.Parallel()
	// Arrange
	factory := newZoneFactory()
	input := models.HubZoneCreationRequest{
		Name:               "Hub-B",
		Size:               1,
		GuardRandomization: 0.05,
		Tuning:             newUnitTuning(),
	}

	// Act
	zone := factory.CreateHubZone(input)

	// Assert
	assert.Equal(t, "Hub-B", zone.Name)
}
