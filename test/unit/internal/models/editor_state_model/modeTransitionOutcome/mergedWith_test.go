package modeTransitionOutcome_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenOutcomesReportDifferentCorrections_BothSurvive(t *testing.T) {
	t.Parallel()
	// Arrange
	pending := editor_state_model.ModeTransitionOutcome{ManualEditsDiscarded: true}
	incoming := editor_state_model.ModeTransitionOutcome{TournamentCountCorrected: true}

	// Act
	merged := pending.MergedWith(incoming)

	// Assert
	assert.Equal(t, editor_state_model.ModeTransitionOutcome{
		ManualEditsDiscarded:     true,
		TournamentCountCorrected: true,
	}, merged)
}

func TestWhenTheIncomingOutcomeIsEmpty_ThePendingOneIsKept(t *testing.T) {
	t.Parallel()
	// Arrange
	pending := editor_state_model.ModeTransitionOutcome{ManualEditsDiscarded: true}

	// Act
	merged := pending.MergedWith(editor_state_model.ModeTransitionOutcome{})

	// Assert
	assert.Equal(t, pending, merged)
}

func TestWhenBothOutcomesAreEmpty_NothingIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	pending := editor_state_model.ModeTransitionOutcome{}

	// Act
	merged := pending.MergedWith(editor_state_model.ModeTransitionOutcome{})

	// Assert
	assert.Equal(t, editor_state_model.ModeTransitionOutcome{}, merged)
}
