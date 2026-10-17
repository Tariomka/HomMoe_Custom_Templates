package fieldOverrides_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenNothingDiffers_ReturnsNoOverridesAndNoUnset(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newConnection()
	edited := newConnection()

	// Act
	overrides, unset := diffConnection(derived, edited)

	// Assert
	assert.Equal(t, []any{editor_state.ConnectionFieldOverrides{}, []string(nil)}, []any{overrides, unset})
}

func TestWhenAValueFieldDiffers_ItBecomesAPointerOverride(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newConnection()
	edited := newConnection()
	edited.GuardValue = 80000

	// Act
	overrides, _ := diffConnection(derived, edited)

	// Assert
	assert.Equal(t, editor_state.ConnectionFieldOverrides{GuardValue: new(80000)}, overrides)
}

// A zero value is a real choice the user made, so it must be stored rather than
// read as "not overridden".
func TestWhenAValueFieldIsEditedToZero_TheZeroIsStored(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newConnection()
	edited := newConnection()
	edited.GuardValue = 0

	// Act
	overrides, _ := diffConnection(derived, edited)

	// Assert
	assert.Equal(t, editor_state.ConnectionFieldOverrides{GuardValue: new(0)}, overrides)
}

func TestWhenAPointerFieldDiffers_ItKeepsItsOwnType(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newConnection()
	edited := newConnection()
	edited.Road = new(false)

	// Act
	overrides, _ := diffConnection(derived, edited)

	// Assert
	assert.Equal(t, editor_state.ConnectionFieldOverrides{Road: new(false)}, overrides)
}

func TestWhenAFieldIsNilButDerivedHasAValue_ItIsUnset(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newConnection()
	edited := newConnection()
	edited.Road = nil

	// Act
	overrides, unset := diffConnection(derived, edited)

	// Assert
	assert.Equal(t, []any{editor_state.ConnectionFieldOverrides{}, []string{"road"}}, []any{overrides, unset})
}

// nil and [] are distinct: an empty list is an override, not an unset.
func TestWhenASliceIsEmptyButDerivedIsNil_ItIsAnEmptyOverride(t *testing.T) {
	t.Parallel()
	// Arrange
	var derived template_entity.Zone
	edited := derived
	edited.GuardedContentPool = []string{}

	// Act
	overrides, unset := editor_state_model.DiffFieldOverrides[template_entity.Zone, editor_state.ZoneFieldOverrides](
		derived, edited)

	// Assert
	assert.Equal(t,
		[]any{editor_state.ZoneFieldOverrides{GuardedContentPool: []string{}}, []string(nil)},
		[]any{overrides, unset})
}

func TestWhenOnlyStructureDiffers_ItIsNeverCompared(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newConnection()
	edited := newConnection()
	edited.Name, edited.From, edited.To, edited.ConnectionType = "Other", "X", "Y", "Portal"

	// Act
	overrides, unset := diffConnection(derived, edited)

	// Assert
	assert.Equal(t, []any{editor_state.ConnectionFieldOverrides{}, []string(nil)}, []any{overrides, unset})
}

func diffConnection(derived, edited template_entity.Connection) (editor_state.ConnectionFieldOverrides, []string) {
	return editor_state_model.DiffFieldOverrides[template_entity.Connection, editor_state.ConnectionFieldOverrides](
		derived, edited)
}

func newConnection() template_entity.Connection {
	var connection template_entity.Connection
	connection.Name, connection.From, connection.To, connection.ConnectionType = "A-B", "A", "B", "Direct"
	connection.Road = new(true)
	connection.GuardValue = 15000
	return connection
}
