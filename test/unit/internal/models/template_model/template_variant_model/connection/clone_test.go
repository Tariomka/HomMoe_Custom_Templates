package connection_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_common_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_variant_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenConnectionIsCloned_ValueFieldsArePreserved(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := newPopulatedConnection()

	// Act
	clone := connection.Clone()

	// Assert
	assert.Equal(t, connection, clone)
}

func TestWhenARoadFlagIsMutatedOnAConnectionClone_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := newPopulatedConnection()
	clone := connection.Clone()
	expected := *connection.Road

	// Act
	*clone.Road = false

	// Assert
	assert.Equal(t, expected, *connection.Road)
}

func TestWhenAPlacementRuleIsMutatedOnAConnectionClone_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := newPopulatedConnection()
	clone := connection.Clone()
	expected := connection.PortalPlacementRulesFrom[0].Type

	// Act
	clone.PortalPlacementRulesFrom[0].Type = "changed"

	// Assert
	assert.Equal(t, expected, connection.PortalPlacementRulesFrom[0].Type)
}

func TestWhenAPlacementRuleArgumentIsMutatedOnAConnectionClone_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := newPopulatedConnection()
	clone := connection.Clone()
	expected := connection.PortalPlacementRulesTo[0].Args[0]

	// Act
	clone.PortalPlacementRulesTo[0].Args[0] = "changed"

	// Assert
	assert.Equal(t, expected, connection.PortalPlacementRulesTo[0].Args[0])
}

func TestWhenConnectionWithNilReferenceFieldsIsCloned_ReferencesRemainNil(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_variant_model.Connection{}

	// Act
	clone := connection.Clone()

	// Assert
	assert.Nil(t, clone.Road)
}

func newPopulatedConnection() template_variant_model.Connection {
	return template_variant_model.Connection{
		Name:                     "connection",
		From:                     "from",
		To:                       "to",
		ConnectionType:           "Portal",
		Length:                   1.5,
		SimTurnSquad:             true,
		Road:                     new(true),
		GuardZone:                "guardZone",
		GuardEscape:              true,
		GuardValue:               100,
		GuardRandomization:       0.5,
		GuardWeeklyIncrement:     0.25,
		GuardMatchGroup:          "group",
		GatePlacement:            "gate",
		PortalPlacementRulesFrom: []template_common_model.PlacementRule{{Type: "from", Args: []any{"fromArg"}}},
		PortalPlacementRulesTo:   []template_common_model.PlacementRule{{Type: "to", Args: []any{"toArg"}}},
		IsUserAdded:              true,
	}
}
