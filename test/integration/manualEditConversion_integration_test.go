package integration_test

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Converting a legacy snapshot against the settings it was made with and
// expanding it onto a fresh generation from the same settings must rebuild
// exactly the zones the snapshot held.
func TestWhenACuratedLegacySnapshotIsConvertedAndExpanded_TheZonesAreUnchanged(t *testing.T) {
	t.Parallel()
	for _, fixture := range curatedManualEditFixtures() {
		t.Run(fixture+"_ZonesUnchanged", func(t *testing.T) {
			t.Parallel()
			// Arrange
			handler := composition.InitializeGuiHandler()
			state, snapshotZones, _ := loadLegacySnapshot(t, handler, fixture)

			// Act
			expansion := convertAndExpand(t, handler, state)

			// Assert
			assert.Equal(t, snapshotZones, expansion.Zones)
		})
	}
}

func TestWhenACuratedLegacySnapshotIsConvertedAndExpanded_TheConnectionsAreUnchanged(t *testing.T) {
	t.Parallel()
	for _, fixture := range curatedManualEditFixtures() {
		t.Run(fixture+"_ConnectionsUnchanged", func(t *testing.T) {
			t.Parallel()
			// Arrange
			handler := composition.InitializeGuiHandler()
			state, _, snapshotConnections := loadLegacySnapshot(t, handler, fixture)

			// Act
			expansion := convertAndExpand(t, handler, state)

			// Assert
			assert.Equal(t, snapshotConnections, expansion.Connections)
		})
	}
}

// The full production cycle: the converted compact edits are saved, loaded
// and expanded on a new generation, and still rebuild the snapshot exactly.
func TestWhenAConvertedCuratedDocumentIsSavedAndReloaded_TheLayoutIsUnchanged(t *testing.T) {
	t.Parallel()
	for _, fixture := range curatedManualEditFixtures() {
		t.Run(fixture+"_LayoutUnchanged", func(t *testing.T) {
			t.Parallel()
			// Arrange
			handler := composition.InitializeGuiHandler()
			state, snapshotZones, snapshotConnections := loadLegacySnapshot(t, handler, fixture)
			converted, err := handler.ConvertManualEdits(dtos.ManualEditConvertRequestDto{
				Baseline:    editor_state_dto.EditorStateDto{EditorState: state.CloneWithoutManualEdits()},
				ManualEdits: state.ManualEditSettings,
			})
			require.NoError(t, err)
			state.ManualEditSettings = converted.ManualEdits
			state.GeneratorVersion = 1
			reloaded, err := handler.LoadState(saveStateTo(t, state, t.TempDir()), false)
			require.NoError(t, err)
			require.False(t, reloaded.State.NeedsDerivation())

			// Act
			expansion := expandOnFreshGeneration(t, handler, reloaded.State)

			// Assert
			assert.Equal(t,
				template_model.Variant{Zones: snapshotZones, Connections: snapshotConnections},
				template_model.Variant{Zones: expansion.Zones, Connections: expansion.Connections})
		})
	}
}

// The compact form must be a fraction of the old snapshot for a large layout
// whose only edit was a move (report §13 targets 30% or less).
func TestWhenTheLargeHubMoveIsConverted_TheManualSectionShrinksToAtMostThirtyPercent(t *testing.T) {
	t.Parallel()
	// Arrange
	const fixture = "editorState_v2_geohub6p_move.gen.json"
	handler := composition.InitializeGuiHandler()
	state, _, _ := loadLegacySnapshot(t, handler, fixture)
	converted, err := handler.ConvertManualEdits(dtos.ManualEditConvertRequestDto{
		Baseline:    editor_state_dto.EditorStateDto{EditorState: state.CloneWithoutManualEdits()},
		ManualEdits: state.ManualEditSettings,
	})
	require.NoError(t, err)
	state.ManualEditSettings = converted.ManualEdits
	state.GeneratorVersion = 1
	legacySize := compactSize(t, readFile(t, filepath.Join("..", "test_helpers", "testdata", fixture)),
		"manualZones", "manualConnections")

	// Act
	compactedSize := compactSize(t, readFile(t, saveStateTo(t, state, t.TempDir())), "manualEdits")

	// Assert
	t.Logf("manual section: legacy %d B, compact %d B", legacySize, compactedSize)
	assert.LessOrEqual(t, float64(compactedSize), 0.3*float64(legacySize))
}

// compactSize is the byte count of the named top-level values with all
// insignificant whitespace removed.
func compactSize(t *testing.T, document string, keys ...string) int {
	t.Helper()
	var values map[string]json.RawMessage
	require.NoError(t, json.Unmarshal([]byte(document), &values))
	size := 0
	for _, key := range keys {
		var compacted bytes.Buffer
		require.NoError(t, json.Compact(&compacted, values[key]))
		size += compacted.Len()
	}

	return size
}

func loadLegacySnapshot(
	t *testing.T,
	handler handler_interfaces.IGuiHandler,
	fixture string) (editor_state_model.EditorState, []template_model.Zone, []template_model.Connection) {
	t.Helper()
	loaded, err := handler.LoadState(filepath.Join("..", "test_helpers", "testdata", fixture), false)
	require.NoError(t, err)
	require.True(t, loaded.State.NeedsDerivation())
	zones, connections := loaded.State.ToSnapshot()

	return loaded.State, zones, connections
}

func convertAndExpand(
	t *testing.T,
	handler handler_interfaces.IGuiHandler,
	state editor_state_model.EditorState) dtos.ManualEditExpansionDto {
	t.Helper()
	settings := editor_state_dto.EditorStateDto{EditorState: state.CloneWithoutManualEdits()}
	converted, err := handler.ConvertManualEdits(dtos.ManualEditConvertRequestDto{
		Baseline:    settings,
		ManualEdits: state.ManualEditSettings,
	})
	require.NoError(t, err)
	state.ManualEditSettings = converted.ManualEdits

	return expandOnFreshGeneration(t, handler, state)
}

// expandOnFreshGeneration generates from the state's settings and expands its
// compact manual edits on that generation.
func expandOnFreshGeneration(
	t *testing.T,
	handler handler_interfaces.IGuiHandler,
	state editor_state_model.EditorState) dtos.ManualEditExpansionDto {
	t.Helper()
	generated, err := handler.GenerateTemplate(
		editor_state_dto.EditorStateDto{EditorState: state.CloneWithoutManualEdits()})
	require.NoError(t, err)

	return handler.ExpandManualEdits(dtos.ManualEditExpandRequestDto{
		Template: generated.Template,
		State:    editor_state_dto.EditorStateDto{EditorState: state},
	})
}
