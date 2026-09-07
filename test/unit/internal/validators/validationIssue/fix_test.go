package validationIssue_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/validators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenPlayerCountAboveMaximum_FixClampsToMaximum(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.PlayerCount = 50

	// Act
	applyAllFixes(t, &state)

	// Assert
	assert.Equal(t, 8, state.PlayerCount)
}

func TestWhenPlayerCountBelowMinimum_FixClampsToMinimum(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.PlayerCount = 0

	// Act
	applyAllFixes(t, &state)

	// Assert
	assert.Equal(t, 2, state.PlayerCount)
}

func TestWhenPercentFieldBelowMinimum_FixClampsToMinimum(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.ResourceDensityPercent = 0

	// Act
	applyAllFixes(t, &state)

	// Assert
	assert.Equal(t, 25, state.ResourceDensityPercent)
}

func TestWhenCountFieldIsNegative_FixSetsToZero(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.NeutralZoneCount = -5

	// Act
	applyAllFixes(t, &state)

	// Assert
	assert.Equal(t, 0, state.NeutralZoneCount)
}

func TestWhenHeroMaxIsLessThanHeroMin_FixRaisesHeroMaxToHeroMin(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.HeroCountMin = 6
	state.HeroCountMax = 4

	// Act
	applyAllFixes(t, &state)

	// Assert
	assert.Equal(t, 6, state.HeroCountMax)
}

func TestWhenHeroMinAboveRangeAndHeroMaxBelowHeroMin_FixesRestoreOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.HeroCountMin = 50
	state.HeroCountMax = 5

	// Act
	applyAllFixes(t, &state)

	// Assert
	assert.Equal(t, 12, state.HeroCountMax)
}

func TestWhenMapSizeIsUnknown_FixSnapsToNearestSize(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.MapSize = 100

	// Act
	applyAllFixes(t, &state)

	// Assert
	assert.Equal(t, 96, state.MapSize)
}

func TestWhenGameModeIsUnknown_FixResetsToClassic(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.GameMode = "NotARealGameMode"

	// Act
	applyAllFixes(t, &state)

	// Assert
	assert.Equal(t, registry.GetGameModeValues().Classic, state.GameMode)
}

func TestWhenVictoryConditionIsUnknown_FixResetsToStandard(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.VictoryCondition = "NotARealCondition"

	// Act
	applyAllFixes(t, &state)

	// Assert
	assert.Equal(t, registry.GetWinningConditionValues().Standard, state.VictoryCondition)
}

// applyAllFixes validates the state and applies every returned fix in order,
// requiring that at least one issue was found.
func applyAllFixes(t *testing.T, state *editor_state_model.EditorState) {
	t.Helper()
	issues := validators.NewEditorStateValidator().Validate(state)
	require.NotEmpty(t, issues)
	for _, issue := range issues {
		issue.Fix(state)
	}
}
