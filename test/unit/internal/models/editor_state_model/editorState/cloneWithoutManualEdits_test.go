package editorState_test

import (
	"reflect"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenStateHasManualEdits_SettingsCloneHasNone(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newFullyPopulatedState(t)

	// Act
	settings := state.CloneWithoutManualEdits()

	// Assert
	assert.Zero(t, settings.ManualEditSettings)
}

func TestWhenStateIsClonedWithoutManualEdits_SettingsEqualTheSource(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newFullyPopulatedState(t)
	expected := state
	expected.ManualEditSettings = editor_state_model.ManualEditSettings{}

	// Act
	settings := state.CloneWithoutManualEdits()

	// Assert
	assert.Equal(t, expected, settings)
}

func TestWhenStateIsClonedWithoutManualEdits_SourceKeepsItsManualEdits(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newFullyPopulatedState(t)
	expected := state.Clone()

	// Act
	state.CloneWithoutManualEdits()

	// Assert
	assert.Equal(t, expected.ManualEditSettings, state.ManualEditSettings)
}

// TestWhenEveryReferenceFieldIsWalked_SettingsCloneSharesNoStorageWithTheSource
// is the drift guard for the settings half of Clone: a new slice or pointer in
// any settings group trips it until CloneWithoutManualEdits copies it.
func TestWhenEveryReferenceFieldIsWalked_SettingsCloneSharesNoStorageWithTheSource(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newFullyPopulatedState(t)
	state.ManualEditSettings = editor_state_model.ManualEditSettings{}

	// Act
	settings := state.CloneWithoutManualEdits()

	// Assert
	assertNoSharedStorage(t, reflect.ValueOf(state), reflect.ValueOf(settings), "EditorStateModel")
}
