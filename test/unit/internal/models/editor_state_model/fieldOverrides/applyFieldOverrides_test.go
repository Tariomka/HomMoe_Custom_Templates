package fieldOverrides_test

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenOverridesAreApplied_TheyReplaceTheDerivedValues(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newConnection()
	overrides := editor_state.ConnectionFieldOverrides{GuardValue: new(80000), Road: new(false)}
	expected := newConnection()
	expected.GuardValue = 80000
	expected.Road = new(false)

	// Act
	applied := applyConnection(derived, overrides, nil)

	// Assert
	assert.Equal(t, expected, applied)
}

func TestWhenAFieldIsUnset_ItIsCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newConnection()
	expected := newConnection()
	expected.Road = nil

	// Act
	applied := applyConnection(derived, editor_state.ConnectionFieldOverrides{}, []string{"road"})

	// Assert
	assert.Equal(t, expected, applied)
}

func TestWhenAnUnsetNameIsUnknown_ItIsIgnored(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newConnection()

	// Act
	applied := applyConnection(derived, editor_state.ConnectionFieldOverrides{}, []string{"notAField"})

	// Assert
	assert.Equal(t, newConnection(), applied)
}

func TestWhenOverridesAreApplied_TheDerivedValueIsNotModified(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newConnection()

	// Act
	applyConnection(derived, editor_state.ConnectionFieldOverrides{GuardValue: new(1)}, []string{"road"})

	// Assert
	assert.Equal(t, newConnection(), derived)
}

// Diff then apply must give back exactly the edited item, whatever the
// derivation was; that is the whole contract of the override form.
func TestWhenADiffIsAppliedToItsDerivation_TheEditedItemComesBack(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newConnection()
	edited := newConnection()
	edited.GuardValue = 0
	edited.Road = nil
	edited.GatePlacement = "Center"
	overrides, unset := diffConnection(derived, edited)

	// Act
	applied := applyConnection(derived, overrides, unset)

	// Assert
	assert.Equal(t, edited, applied)
}

// The parity guards: a new field on the RMG zone or connection that the
// override structs do not mirror would silently stop being persisted.
func TestWhenTheZoneOverridesAreCompared_TheyMirrorEveryNonStructureZoneField(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	mismatches := overrideMismatches(
		reflect.TypeFor[template_entity.Zone](), reflect.TypeFor[editor_state.ZoneFieldOverrides](), "Name")

	// Assert
	assert.Empty(t, mismatches)
}

func TestWhenTheConnectionOverridesAreCompared_TheyMirrorEveryNonStructureConnectionField(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	mismatches := overrideMismatches(
		reflect.TypeFor[template_entity.Connection](), reflect.TypeFor[editor_state.ConnectionFieldOverrides](),
		"Name", "From", "To", "ConnectionType")

	// Assert
	assert.Empty(t, mismatches)
}

func applyConnection(
	derived template_entity.Connection,
	overrides editor_state.ConnectionFieldOverrides,
	unset []string) template_entity.Connection {
	return editor_state_model.ApplyFieldOverrides(derived, overrides, unset)
}

// overrideMismatches lists every item field that has no same-named, same-JSON
// override field of the expected type, and every override field with no item
// field behind it.
func overrideMismatches(itemType, overridesType reflect.Type, structure ...string) []string {
	var mismatches []string
	for field := range itemType.Fields() {
		if slices.Contains(structure, field.Name) {
			continue
		}

		override, found := overridesType.FieldByName(field.Name)
		switch {
		case !found:
			mismatches = append(mismatches, field.Name+": missing")
		case override.Type != expectedOverrideType(field.Type):
			mismatches = append(mismatches, field.Name+": type "+override.Type.String())
		case jsonName(override) != jsonName(field) || !strings.Contains(override.Tag.Get("json"), "omitzero"):
			mismatches = append(mismatches, field.Name+": tag "+override.Tag.Get("json"))
		}
	}
	for field := range overridesType.Fields() {
		if _, found := itemType.FieldByName(field.Name); !found {
			mismatches = append(mismatches, field.Name+": no item field")
		}
	}
	return mismatches
}

func expectedOverrideType(fieldType reflect.Type) reflect.Type {
	if fieldType.Kind() == reflect.Pointer || fieldType.Kind() == reflect.Slice {
		return fieldType
	}
	return reflect.PointerTo(fieldType)
}

func jsonName(field reflect.StructField) string {
	name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
	return name
}
