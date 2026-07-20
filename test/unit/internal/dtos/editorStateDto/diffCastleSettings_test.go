package editorStateDto_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/stretchr/testify/assert"
)

func TestWhenNothingChanged_ReportsNoChanges(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := dtos.NewDefaultEditorStateDto()
	current := previous

	// Act
	changes := previous.DiffCastleSettings(&current)

	// Assert
	assert.Equal(t, editor_state_dto.CastleSettingChanges{}, changes)
}

func TestWhenSimpleModeNeutralCountChanges_FlagsNeutralSimpleOnly(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := dtos.NewDefaultEditorStateDto()
	current := previous
	current.NeutralZoneCastles = 3
	current.NeutralHighCastlesPerZone = 4 // advanced-only option, must be ignored

	// Act
	changes := previous.DiffCastleSettings(&current)

	// Assert
	assert.Equal(t, editor_state_dto.CastleSettingChanges{NeutralSimple: true}, changes)
}

func TestWhenAdvancedModeHighCountChanges_FlagsNeutralHighOnly(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := dtos.NewDefaultEditorStateDto()
	previous.AdvancedMode = true
	current := previous
	current.NeutralHighCastlesPerZone = 4
	current.NeutralZoneCastles = 3 // simple-only option, must be ignored

	// Act
	changes := previous.DiffCastleSettings(&current)

	// Assert
	assert.Equal(t, editor_state_dto.CastleSettingChanges{NeutralHigh: true}, changes)
}

func TestWhenAdvancedModeLowestCountChanges_FlagsNeutralLowestOnly(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := dtos.NewDefaultEditorStateDto()
	previous.AdvancedMode = true
	current := previous
	current.NeutralLowestCastlesPerZone = 4

	// Act
	changes := previous.DiffCastleSettings(&current)

	// Assert
	assert.Equal(t, editor_state_dto.CastleSettingChanges{NeutralLowest: true}, changes)
}

func TestWhenPlayerAndHubCountsChange_FlagsPlayerCastlesAndHub(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := dtos.NewDefaultEditorStateDto()
	current := previous
	current.PlayerOwnedCastles = 2
	current.HubZoneCastles = 3

	// Act
	changes := previous.DiffCastleSettings(&current)

	// Assert
	assert.Equal(t, editor_state_dto.CastleSettingChanges{PlayerCastles: true, Hub: true}, changes)
}
