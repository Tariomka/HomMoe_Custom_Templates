package regenerationDecisionService_test

import (
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/regeneration"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/editor"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

// debounceWindow mirrors the unexported autoRegenerationDebounce constant. It
// is duplicated deliberately: the tests must fail if the production window
// changes silently.
const debounceWindow = 300 * time.Millisecond

func TestWhenNoPreviousGenerationExists_RegeneratesImmediately(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	request := regeneration.DecisionRequest{
		Previous: nil,
		Current:  defaultState(),
		Now:      gofakeit.Date(),
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, regeneration.Decision{
		Regenerate:      true,
		NextStateAction: regeneration.NextStateLeave,
	}, decision)
}

func TestWhenStateIsUnchangedSinceLastGeneration_CancelsPendingDebounce(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	request := regeneration.DecisionRequest{
		Previous: defaultState(),
		Current:  defaultState(),
		Next:     defaultState(),
		Now:      gofakeit.Date(),
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, regeneration.Decision{
		Regenerate:      false,
		NextStateAction: regeneration.NextStateClear,
	}, decision)
}

func TestWhenLayoutDefiningOptionChanged_RegeneratesImmediately(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	request := regeneration.DecisionRequest{
		Previous: defaultState(),
		Current:  layoutChangedState(),
		Now:      gofakeit.Date(),
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, regeneration.Decision{
		Regenerate:      true,
		NextStateAction: regeneration.NextStateClear,
	}, decision)
}

func TestWhenNonLayoutOptionChangedAndNoDebounceArmed_ArmsDebounce(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	now := gofakeit.Date()
	request := regeneration.DecisionRequest{
		Previous: defaultState(),
		Current:  nonLayoutChangedState(),
		Next:     nil,
		Now:      now,
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, regeneration.Decision{
		Regenerate:      false,
		NextStateAction: regeneration.NextStateSetFromCurrent,
		RedrawAt:        now.Add(debounceWindow),
		ScheduleRedraw:  true,
	}, decision)
}

func TestWhenStateMovedAgainWhileDebounceArmed_RearmsDebounceFromNow(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	now := gofakeit.Date()
	request := regeneration.DecisionRequest{
		Previous:      defaultState(),
		Current:       nonLayoutChangedState(),
		Next:          defaultState(),
		Now:           now,
		DebounceDueAt: now.Add(-time.Hour),
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, regeneration.Decision{
		Regenerate:      false,
		NextStateAction: regeneration.NextStateSetFromCurrent,
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
	request := regeneration.DecisionRequest{
		Previous:      defaultState(),
		Current:       nonLayoutChangedState(),
		Next:          nonLayoutChangedState(),
		Now:           now,
		DebounceDueAt: dueAt,
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, regeneration.Decision{
		Regenerate:      false,
		NextStateAction: regeneration.NextStateLeave,
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
	request := regeneration.DecisionRequest{
		Previous:      defaultState(),
		Current:       nonLayoutChangedState(),
		Next:          nonLayoutChangedState(),
		Now:           now,
		DebounceDueAt: now,
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, regeneration.Decision{
		Regenerate:      true,
		NextStateAction: regeneration.NextStateClear,
	}, decision)
}

func TestWhenStateIsStableAndDebounceOverdue_Regenerates(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	now := gofakeit.Date()
	request := regeneration.DecisionRequest{
		Previous:      defaultState(),
		Current:       nonLayoutChangedState(),
		Next:          nonLayoutChangedState(),
		Now:           now,
		DebounceDueAt: now.Add(-time.Second),
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, regeneration.Decision{
		Regenerate:      true,
		NextStateAction: regeneration.NextStateClear,
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
	request := regeneration.DecisionRequest{
		Previous: defaultState(),
		Current:  current,
		Now:      gofakeit.Date(),
	}

	// Act
	decision := service.DecideRegeneration(request)

	// Assert
	assert.Equal(t, regeneration.Decision{
		Regenerate:      false,
		NextStateAction: regeneration.NextStateClear,
	}, decision)
}

func defaultState() *editor_state_model.EditorState {
	state := editor_state_model.NewDefaultEditorStateModel()
	return &state
}

// layoutChangedState returns a state differing from defaultState in a
// layout-defining option, which invalidates any hand-made zone layout.
func layoutChangedState() *editor_state_model.EditorState {
	state := defaultState()
	state.Topology = config.TopologyChain
	return state
}

// nonLayoutChangedState differs from defaultState only in an option that does
// not alter the zone or connection graph, so it is debounced instead.
func nonLayoutChangedState() *editor_state_model.EditorState {
	state := defaultState()
	state.ResourceDensityPercent = 50
	return state
}
