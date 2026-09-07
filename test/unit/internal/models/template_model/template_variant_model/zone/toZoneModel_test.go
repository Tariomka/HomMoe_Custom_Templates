package zone_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

// Encounter-hole settings hang off a pointer, so lifting has to go through the
// pointer converter rather than copying a zero-valued block.
func TestWhenAZoneEntityCarriesEncounterHoleSettings_TheyAreLifted(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := newPopulatedZoneEntity()

	// Act
	model := template_model.ToZoneModel(entity)

	// Assert
	expected := &template_model.EncounterHolesSettings{AffectedEncounters: 1, TwoHoleEncounters: 2}
	assert.Equal(t, expected, model.EncounterHolesSettings)
}

// A main object's faction is a TypedRef behind a pointer, which is the deepest
// converter the zone reaches.
func TestWhenAMainObjectCarriesAFaction_TheFactionIsLifted(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := newPopulatedZoneEntity()

	// Act
	model := template_model.ToZoneModel(entity)

	// Assert
	expected := &template_model.TypedRef{Type: "FromList", Args: []string{"factionArg"}}
	assert.Equal(t, expected, model.MainObjects[0].Faction)
}

func TestWhenAZoneEntityCarriesRoads_TheRoadsAreLifted(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := newPopulatedZoneEntity()

	// Act
	model := template_model.ToZoneModel(entity)

	// Assert
	expected := []template_model.Road{{
		Type: "road",
		From: template_model.TypedRef{Type: "from", Args: []string{"fromArg"}},
		To:   template_model.TypedRef{Type: "to", Args: []string{"toArg"}},
		Road: new(true),
	}}
	assert.Equal(t, expected, model.Roads)
}

// A zone without encounter holes must lift to nil, not to a zero-valued block:
// the schema writes the key only when it is present.
func TestWhenAZoneEntityHasNoEncounterHoleSettings_TheModelKeepsThemAbsent(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := newPopulatedZoneEntity()
	entity.EncounterHolesSettings = nil

	// Act
	model := template_model.ToZoneModel(entity)

	// Assert
	assert.Nil(t, model.EncounterHolesSettings)
}
