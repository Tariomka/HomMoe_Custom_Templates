package zone_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenAZoneCarriesEncounterHoleSettings_TheyAreFlattened(t *testing.T) {
	t.Parallel()
	// Arrange
	model := newPopulatedZone()

	// Act
	entity := template_model.ToZoneEntity(model)

	// Assert
	expected := &template_entity.EncounterHolesSettings{AffectedEncounters: 1, TwoHoleEncounters: 2}
	assert.Equal(t, expected, entity.EncounterHolesSettings)
}

func TestWhenAMainObjectCarriesAFaction_TheFactionIsFlattened(t *testing.T) {
	t.Parallel()
	// Arrange
	model := newPopulatedZone()

	// Act
	entity := template_model.ToZoneEntity(model)

	// Assert
	expected := &template_entity.TypedRef{Type: "FromList", Args: []string{"factionArg"}}
	assert.Equal(t, expected, entity.MainObjects[0].Faction)
}

func TestWhenAZoneCarriesRoads_TheRoadsAreFlattened(t *testing.T) {
	t.Parallel()
	// Arrange
	model := newPopulatedZone()

	// Act
	entity := template_model.ToZoneEntity(model)

	// Assert
	expected := []template_entity.Road{{
		Type: "road",
		From: template_entity.TypedRef{Type: "from", Args: []string{"fromArg"}},
		To:   template_entity.TypedRef{Type: "to", Args: []string{"toArg"}},
		Road: new(true),
	}}
	assert.Equal(t, expected, entity.Roads)
}
