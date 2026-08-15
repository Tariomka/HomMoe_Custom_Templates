package stateHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/validators"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenStateHasNoIssues_ReturnsNoWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewStateHandler(&test_helpers.FileServiceMock{}, newPassingValidator())

	// Act
	validation := handler.ValidateEditorState(dtos.NewDefaultEditorStateDto(), true)

	// Assert
	assert.Empty(t, validation.Warnings)
}

func TestWhenIssuesAreNotFixed_ReturnsTheStateUnmodified(t *testing.T) {
	t.Parallel()
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.AdvancedMode = true
	state.NeutralZoneCount = gofakeit.IntRange(1, 16)
	handler := handlers.NewStateHandler(
		&test_helpers.FileServiceMock{},
		newValidatorReporting(gofakeit.Sentence(3)))

	// Act
	validation := handler.ValidateEditorState(state, false)

	// Assert
	assert.Equal(t, state, validation.State)
}

func TestWhenIssuesAreFixed_AppliesEachIssueFix(t *testing.T) {
	t.Parallel()
	// Arrange - the real validator is used because ValidationIssue's fix is
	// package-private and cannot be supplied by a mock.
	state := dtos.NewDefaultEditorStateDto()
	state.PlayerCount = 99
	handler := handlers.NewStateHandler(&test_helpers.FileServiceMock{}, validators.NewEditorStateValidator())

	// Act
	validation := handler.ValidateEditorState(state, true)

	// Assert
	assert.Equal(t, 8, validation.State.PlayerCount)
}

func TestWhenAdvancedModeIsOn_ZeroesTheSimpleNeutralZoneCount(t *testing.T) {
	t.Parallel()
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.AdvancedMode = true
	state.NeutralZoneCount = gofakeit.IntRange(1, 16)
	handler := handlers.NewStateHandler(&test_helpers.FileServiceMock{}, newPassingValidator())

	// Act
	validation := handler.ValidateEditorState(state, true)

	// Assert
	assert.Zero(t, validation.State.NeutralZoneCount)
}

// TestWhenValidationFixesAContentRow_TheCallersSliceIsUnchanged pins the clone
// on entry: the DTO arrives by value, so without it the validated state would
// hand every caller a slice that still points at the caller's own storage.
func TestWhenValidationFixesAContentRow_TheCallersSliceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.PlayerCount = 99
	state.PlayerZoneContentRows = []models.ZoneContentRowSave{{Sid: "sawmill", Count: 1}}
	handler := handlers.NewStateHandler(&test_helpers.FileServiceMock{}, validators.NewEditorStateValidator())
	validation := handler.ValidateEditorState(state, true)

	// Act
	validation.State.PlayerZoneContentRows[0].Count = 5

	// Assert
	assert.Equal(t, 1, state.PlayerZoneContentRows[0].Count)
}

func TestWhenAdvancedModeIsOn_KeepsThePerTierCounts(t *testing.T) {
	t.Parallel()
	// Arrange
	state := advancedTierCountsState()
	state.AdvancedMode = true
	handler := handlers.NewStateHandler(&test_helpers.FileServiceMock{}, newPassingValidator())
	expected := tierCountsOf(state)

	// Act
	validation := handler.ValidateEditorState(state, true)

	// Assert
	assert.Equal(t, expected, tierCountsOf(validation.State))
}

func TestWhenAdvancedModeIsOff_ZeroesThePerTierCounts(t *testing.T) {
	t.Parallel()
	// Arrange
	state := advancedTierCountsState()
	state.AdvancedMode = false
	handler := handlers.NewStateHandler(&test_helpers.FileServiceMock{}, newPassingValidator())

	// Act
	validation := handler.ValidateEditorState(state, true)

	// Assert
	assert.Equal(t, make([]int, 8), tierCountsOf(validation.State))
}

func TestWhenAdvancedModeIsOff_KeepsTheSimpleNeutralZoneCount(t *testing.T) {
	t.Parallel()
	// Arrange
	state := advancedTierCountsState()
	state.AdvancedMode = false
	state.NeutralZoneCount = gofakeit.IntRange(1, 16)
	handler := handlers.NewStateHandler(&test_helpers.FileServiceMock{}, newPassingValidator())

	// Act
	validation := handler.ValidateEditorState(state, true)

	// Assert
	assert.Equal(t, state.NeutralZoneCount, validation.State.NeutralZoneCount)
}

func TestWhenIssuesAreNotFixed_SkipsTheInactiveCountNormalization(t *testing.T) {
	t.Parallel()
	// Arrange
	state := advancedTierCountsState()
	state.AdvancedMode = false
	handler := handlers.NewStateHandler(&test_helpers.FileServiceMock{}, newPassingValidator())
	expected := tierCountsOf(state)

	// Act
	validation := handler.ValidateEditorState(state, false)

	// Assert
	assert.Equal(t, expected, tierCountsOf(validation.State))
}

func TestWhenStateIsValidated_LeavesTheCallersStateUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	state := advancedTierCountsState()
	state.AdvancedMode = false
	handler := handlers.NewStateHandler(&test_helpers.FileServiceMock{}, newPassingValidator())
	expected := tierCountsOf(state)

	// Act
	_ = handler.ValidateEditorState(state, true)

	// Assert
	assert.Equal(t, expected, tierCountsOf(state))
}

// advancedTierCountsState returns a state whose eight per-tier neutral counts
// are all non-zero, so a normalization to zero is observable.
func advancedTierCountsState() dtos.EditorStateDto {
	state := dtos.NewDefaultEditorStateDto()
	state.NeutralLowestNoCastleCount = gofakeit.IntRange(1, 8)
	state.NeutralLowestCastleCount = gofakeit.IntRange(1, 8)
	state.NeutralLowNoCastleCount = gofakeit.IntRange(1, 8)
	state.NeutralLowCastleCount = gofakeit.IntRange(1, 8)
	state.NeutralMediumNoCastleCount = gofakeit.IntRange(1, 8)
	state.NeutralMediumCastleCount = gofakeit.IntRange(1, 8)
	state.NeutralHighNoCastleCount = gofakeit.IntRange(1, 8)
	state.NeutralHighCastleCount = gofakeit.IntRange(1, 8)
	return state
}

func tierCountsOf(state dtos.EditorStateDto) []int {
	return []int{
		state.NeutralLowestNoCastleCount,
		state.NeutralLowestCastleCount,
		state.NeutralLowNoCastleCount,
		state.NeutralLowCastleCount,
		state.NeutralMediumNoCastleCount,
		state.NeutralMediumCastleCount,
		state.NeutralHighNoCastleCount,
		state.NeutralHighCastleCount,
	}
}
