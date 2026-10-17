// Package manualConnectionEdit_test contains shared arrangement helpers for the
// manualConnectionEdit.go unit tests.
package manualConnectionEdit_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// newDerivedConnection builds a generator-shaped connection to diff against.
func newDerivedConnection() template_model.Connection {
	return template_model.Connection{
		Name:           "Derived",
		From:           "Zone A",
		To:             "Zone B",
		ConnectionType: "Proximity",
		Length:         1.0,
		GuardValue:     1000,
		Road:           new(true),
		GuardZone:      "Zone A",
		GatePlacement:  "Center",
		PortalPlacementRulesFrom: []template_model.PlacementRule{
			newPlacementRule("derived"),
		},
	}
}

// newEditedConnection builds a connection differing from newDerivedConnection
// in structure and in several fields.
func newEditedConnection() template_model.Connection {
	connection := newDerivedConnection()
	connection.Name = "Edited"
	connection.To = "Zone C"
	connection.ConnectionType = "Portal"
	connection.Length = 2.5
	connection.GuardValue = 0
	connection.Road = new(false)
	connection.GuardEscape = true
	connection.PortalPlacementRulesFrom = []template_model.PlacementRule{newPlacementRule("edited")}
	connection.PortalPlacementRulesTo = []template_model.PlacementRule{newPlacementRule("to")}
	connection.IsUserAdded = true
	return connection
}

func newPlacementRule(ruleType string) template_model.PlacementRule {
	var rule template_model.PlacementRule
	rule.Type = ruleType
	rule.Args = []any{"arg"}
	return rule
}

// newPopulatedEdit builds an edit whose every reference field carries data.
func newPopulatedEdit() editor_state_model.ManualConnectionEdit {
	edit := editor_state_model.NewManualConnectionEdit(newDerivedConnection(), newEditedConnection())
	edit.Unset = []string{"guardZone"}
	return edit
}
