package editorStateDto_test

import (
	"reflect"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenStateIsCloned_CloneEqualsTheSource(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newContentRowState()

	// Act
	clone := state.Clone()

	// Assert
	assert.Equal(t, state, clone)
}

func TestWhenSliceIsNil_CloneSliceStaysNil(t *testing.T) {
	t.Parallel()
	// Arrange - nil must not become empty: the change detection tells them apart.
	state := editor_state_dto.EditorStateDto{}

	// Act
	clone := state.Clone()

	// Assert
	assert.Nil(t, clone.HubZoneContentRows)
}

// TestWhenAContentRowIsMutatedInPlaceOnTheClone_SourceIsUnchanged covers each of
// the six content-row slices; a shallow struct copy would share their storage.
func TestWhenAContentRowIsMutatedInPlaceOnTheClone_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	for fieldName, readRows := range contentRowFields() {
		t.Run("When"+fieldName+"IsMutated_SourceIsUnchanged", func(t *testing.T) {
			t.Parallel()
			// Arrange
			state := newContentRowState()
			clone := state.Clone()
			originalSid := readRows(&state)[0].Sid

			// Act
			readRows(&clone)[0].Sid = originalSid + "-changed"

			// Assert
			assert.Equal(t, originalSid, readRows(&state)[0].Sid)
		})
	}
}

func TestWhenANestedContentRuleIsMutatedInPlaceOnTheClone_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newContentRowState()
	clone := state.Clone()

	// Act
	*clone.PlayerZoneContentRows[0].Rules[0].IsGuarded = false

	// Assert
	require.NotNil(t, state.PlayerZoneContentRows[0].Rules[0].IsGuarded)
	assert.True(t, *state.PlayerZoneContentRows[0].Rules[0].IsGuarded)
}

func TestWhenABonusIsMutatedInPlaceOnTheClone_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newContentRowState()
	clone := state.Clone()
	originalParam := state.Bonuses[0].Param

	// Act
	clone.Bonuses[0].Param = originalParam + "-changed"

	// Assert
	assert.Equal(t, originalParam, state.Bonuses[0].Param)
}

func TestWhenAManualZoneIsMutatedInPlaceOnTheClone_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newFullyPopulatedState(t)
	clone := state.Clone()
	originalName := state.ManualZones[0].Zone.Name

	// Act
	clone.ManualZones[0].Zone.Name = originalName + "-changed"

	// Assert
	assert.Equal(t, originalName, state.ManualZones[0].Zone.Name)
}

func TestWhenAManualConnectionIsMutatedInPlaceOnTheClone_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newFullyPopulatedState(t)
	clone := state.Clone()
	originalName := state.ManualConnections[0].Connection.Name

	// Act
	clone.ManualConnections[0].Connection.Name = originalName + "-changed"

	// Assert
	assert.Equal(t, originalName, state.ManualConnections[0].Connection.Name)
}

// TestWhenEveryReferenceFieldIsWalked_CloneSharesNoStorageWithTheSource is the
// drift guard. It reflects over a state whose every slice and pointer - down
// through the protected entities.Zone and entities.Connection - carries data,
// and fails as soon as a reference field is left aliased. A new slice or
// pointer anywhere in that tree therefore trips this test until Clone covers it.
func TestWhenEveryReferenceFieldIsWalked_CloneSharesNoStorageWithTheSource(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newFullyPopulatedState(t)

	// Act
	clone := state.Clone()

	// Assert
	assertNoSharedStorage(t, reflect.ValueOf(state), reflect.ValueOf(clone), "EditorStateDto")
}

func contentRowFields() map[string]func(state *editor_state_dto.EditorStateDto) []models.ZoneContentRowSave {
	return map[string]func(state *editor_state_dto.EditorStateDto) []models.ZoneContentRowSave{
		"PlayerZoneContentRows": func(state *editor_state_dto.EditorStateDto) []models.ZoneContentRowSave {
			return state.PlayerZoneContentRows
		},
		"LowestNeutralContentRows": func(state *editor_state_dto.EditorStateDto) []models.ZoneContentRowSave {
			return state.LowestNeutralContentRows
		},
		"LowNeutralContentRows": func(state *editor_state_dto.EditorStateDto) []models.ZoneContentRowSave {
			return state.LowNeutralContentRows
		},
		"MediumNeutralContentRows": func(state *editor_state_dto.EditorStateDto) []models.ZoneContentRowSave {
			return state.MediumNeutralContentRows
		},
		"HighNeutralContentRows": func(state *editor_state_dto.EditorStateDto) []models.ZoneContentRowSave {
			return state.HighNeutralContentRows
		},
		"HubZoneContentRows": func(state *editor_state_dto.EditorStateDto) []models.ZoneContentRowSave {
			return state.HubZoneContentRows
		},
	}
}

// newContentRowState returns a state whose six content-row slices and bonus
// list each hold one row, with a rule pointer on the player rows.
func newContentRowState() editor_state_dto.EditorStateDto {
	newRow := func() []models.ZoneContentRowSave {
		return []models.ZoneContentRowSave{{
			Sid:   gofakeit.LetterN(10),
			Count: 1,
			Rules: []models.ContentRuleRowSave{{Name: "Guarded", IsGuarded: new(true)}},
		}}
	}

	state := editor_state_dto.NewDefaultEditorStateDto()
	state.PlayerZoneContentRows = newRow()
	state.LowestNeutralContentRows = newRow()
	state.LowNeutralContentRows = newRow()
	state.MediumNeutralContentRows = newRow()
	state.HighNeutralContentRows = newRow()
	state.HubZoneContentRows = newRow()
	state.Bonuses = []config.BonusEntry{{ReceiverFilter: "start_hero", Param: gofakeit.LetterN(6)}}
	return state
}

// newFullyPopulatedState builds a state in which every slice holds an element
// and every pointer is set, recursively, so the drift guard has something to
// compare on every branch of the tree.
func newFullyPopulatedState(t *testing.T) editor_state_dto.EditorStateDto {
	t.Helper()
	var state editor_state_dto.EditorStateDto
	fillReferenceFields(reflect.ValueOf(&state).Elem())
	require.NotEmpty(t, state.ManualZones, "the populator must reach the manual-edit slices")
	return state
}

// fillReferenceFields recursively assigns non-zero data to every settable field,
// allocating one element per slice and one value per pointer.
func fillReferenceFields(value reflect.Value) {
	//nolint:exhaustive // the kinds left to the default carry no storage to populate
	switch value.Kind() {
	case reflect.Bool:
		value.SetBool(true)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value.SetInt(7)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value.SetUint(7)
	case reflect.Float32, reflect.Float64:
		value.SetFloat(1.5)
	case reflect.String:
		value.SetString("filled")
	case reflect.Interface:
		value.Set(reflect.ValueOf("filled"))
	case reflect.Pointer:
		allocated := reflect.New(value.Type().Elem())
		fillReferenceFields(allocated.Elem())
		value.Set(allocated)
	case reflect.Slice:
		created := reflect.MakeSlice(value.Type(), 1, 1)
		fillReferenceFields(created.Index(0))
		value.Set(created)
	case reflect.Array:
		for index := range value.Len() {
			fillReferenceFields(value.Index(index))
		}
	case reflect.Struct:
		for _, field := range value.Fields() {
			if field.CanSet() {
				fillReferenceFields(field)
			}
		}
	default:
	}
}

// assertNoSharedStorage fails when the clone reaches the same backing array or
// pointed-to value as the source anywhere in the tree.
func assertNoSharedStorage(t *testing.T, source, clone reflect.Value, path string) {
	t.Helper()
	//nolint:exhaustive // the kinds left to the default cannot share storage
	switch source.Kind() {
	case reflect.Pointer:
		if source.IsNil() {
			return
		}
		require.False(t, clone.IsNil(), "%s: the clone dropped a pointer the source has", path)
		assert.NotEqual(t, source.Pointer(), clone.Pointer(), "%s: clone shares a pointer with the source", path)
		assertNoSharedStorage(t, source.Elem(), clone.Elem(), path)
	case reflect.Slice:
		if source.Len() == 0 {
			return
		}
		require.Equal(t, source.Len(), clone.Len(), "%s: the clone has a different length", path)
		assert.NotEqual(t, source.Pointer(), clone.Pointer(), "%s: clone shares a backing array with the source", path)
		for index := range source.Len() {
			assertNoSharedStorage(t, source.Index(index), clone.Index(index), path+"[0]")
		}
	case reflect.Array:
		for index := range source.Len() {
			assertNoSharedStorage(t, source.Index(index), clone.Index(index), path+"[0]")
		}
	case reflect.Struct:
		for index := range source.NumField() {
			field := source.Type().Field(index)
			if !field.IsExported() {
				continue
			}
			assertNoSharedStorage(t, source.Field(index), clone.Field(index), path+"."+field.Name)
		}
	case reflect.Interface:
		// PlacementRule.Args holds boxed scalars decoded from JSON; they are
		// immutable in practice, so cloneConnection copies the slice only.
	default:
	}
}
