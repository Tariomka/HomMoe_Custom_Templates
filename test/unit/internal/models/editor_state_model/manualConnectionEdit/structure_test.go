package manualConnectionEdit_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheEditHasStructureAndOverrides_StructureCarriesOnlyTheStructureFields(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := editor_state_model.NewManualConnectionEdit(newDerivedConnection(), newEditedConnection())

	// Act
	structure := edit.Structure()

	// Assert
	assert.Equal(t, template_model.Connection{
		Name:           "Edited",
		From:           "Zone A",
		To:             "Zone C",
		ConnectionType: "Portal",
		IsUserAdded:    true,
	}, structure)
}
