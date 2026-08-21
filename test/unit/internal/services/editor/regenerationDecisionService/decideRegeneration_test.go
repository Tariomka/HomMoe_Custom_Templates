package regenerationDecisionService_test

import (
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/editor"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

// debounceWindow mirrors the unexported autoRegenerationDebounce constant. It
// is duplicated deliberately: the tests must fail if the production window
// changes silently.
const debounceWindow = 300 * time.Millisecond

func defaultState() *editor_state_dto.EditorStateDto {
	state := editor_state_dto.NewDefaultEditorStateDto()
	return &state
}

// layoutChangedState returns a state differing from defaultState in a
// layout-defining option, which invalidates any hand-made zone layout.
func layoutChangedState() *editor_state_dto.EditorStateDto {
	state := defaultState()
	state.Topology = config_inner.TopologyChain
	return state
}

// nonLayoutChangedState differs from defaultState only in an option that does
// not alter the zone or connection graph, so it is debounced instead.
func nonLayoutChangedState() *editor_state_dto.EditorStateDto {
	state := defaultState()
	state.ResourceDensityPercent = 50
	return state
}

func TestWhenNoPreviousGenerationExists_RegeneratesImmediately(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	request := dtos.RegenerationDecisionRequestDto{
		Previous: nil,
		Current:  defaultState(),
		Now:      gofakeit.Date(),
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, dtos.RegenerationDecisionDto{
		Regenerate:      true,
		NextStateAction: dtos.NextStateLeave,
	}, decision)
}

func TestWhenStateIsUnchangedSinceLastGeneration_CancelsPendingDebounce(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	request := dtos.RegenerationDecisionRequestDto{
		Previous: defaultState(),
		Current:  defaultState(),
		Next:     defaultState(),
		Now:      gofakeit.Date(),
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, dtos.RegenerationDecisionDto{
		Regenerate:      false,
		NextStateAction: dtos.NextStateClear,
	}, decision)
}

func TestWhenLayoutDefiningOptionChanged_RegeneratesImmediately(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	request := dtos.RegenerationDecisionRequestDto{
		Previous: defaultState(),
		Current:  layoutChangedState(),
		Now:      gofakeit.Date(),
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, dtos.RegenerationDecisionDto{
		Regenerate:      true,
		NextStateAction: dtos.NextStateClear,
	}, decision)
}

func TestWhenNonLayoutOptionChangedAndNoDebounceArmed_ArmsDebounce(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	now := gofakeit.Date()
	request := dtos.RegenerationDecisionRequestDto{
		Previous: defaultState(),
		Current:  nonLayoutChangedState(),
		Next:     nil,
		Now:      now,
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, dtos.RegenerationDecisionDto{
		Regenerate:      false,
		NextStateAction: dtos.NextStateSetFromCurrent,
		RedrawAt:        now.Add(debounceWindow),
		ScheduleRedraw:  true,
	}, decision)
}

func TestWhenStateMovedAgainWhileDebounceArmed_RearmsDebounceFromNow(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	now := gofakeit.Date()
	request := dtos.RegenerationDecisionRequestDto{
		Previous:      defaultState(),
		Current:       nonLayoutChangedState(),
		Next:          defaultState(),
		Now:           now,
		DebounceDueAt: now.Add(-time.Hour),
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, dtos.RegenerationDecisionDto{
		Regenerate:      false,
		NextStateAction: dtos.NextStateSetFromCurrent,
		RedrawAt:        now.Add(debounceWindow),
		ScheduleRedraw:  true,
	}, decision)
}

func TestWhenStateIsStableAndDebounceNotDue_KeepsWaiting(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	now := gofakeit.Date()
	dueAt := now.Add(time.Nanosecond)
	request := dtos.RegenerationDecisionRequestDto{
		Previous:      defaultState(),
		Current:       nonLayoutChangedState(),
		Next:          nonLayoutChangedState(),
		Now:           now,
		DebounceDueAt: dueAt,
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, dtos.RegenerationDecisionDto{
		Regenerate:      false,
		NextStateAction: dtos.NextStateLeave,
		RedrawAt:        dueAt,
		ScheduleRedraw:  true,
	}, decision)
}

// The boundary is exclusive: at exactly the due instant the debounce has
// elapsed and regeneration happens. Untestable before the window was injected.
func TestWhenStateIsStableAndDebounceExactlyDue_Regenerates(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	now := gofakeit.Date()
	request := dtos.RegenerationDecisionRequestDto{
		Previous:      defaultState(),
		Current:       nonLayoutChangedState(),
		Next:          nonLayoutChangedState(),
		Now:           now,
		DebounceDueAt: now,
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, dtos.RegenerationDecisionDto{
		Regenerate:      true,
		NextStateAction: dtos.NextStateClear,
	}, decision)
}

func TestWhenStateIsStableAndDebounceOverdue_Regenerates(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	now := gofakeit.Date()
	request := dtos.RegenerationDecisionRequestDto{
		Previous:      defaultState(),
		Current:       nonLayoutChangedState(),
		Next:          nonLayoutChangedState(),
		Now:           now,
		DebounceDueAt: now.Add(-time.Second),
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, dtos.RegenerationDecisionDto{
		Regenerate:      true,
		NextStateAction: dtos.NextStateClear,
	}, decision)
}

// Manual edits alone must never trigger a regeneration: they are reapplied
// through a separate path.
func TestWhenOnlyManualEditsDiffer_CancelsPendingDebounce(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	current := defaultState()
	current.ManualZones = manualZoneSaves()
	request := dtos.RegenerationDecisionRequestDto{
		Previous: defaultState(),
		Current:  current,
		Now:      gofakeit.Date(),
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, dtos.RegenerationDecisionDto{
		Regenerate:      false,
		NextStateAction: dtos.NextStateClear,
	}, decision)
}
