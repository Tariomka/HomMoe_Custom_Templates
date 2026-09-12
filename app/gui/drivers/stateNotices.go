package drivers

import "github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"

// noteStateTransition shows an outcome's message and keeps it pending, because
// the regeneration that follows writes its own status.
func (this *State) noteStateTransition(outcome editor_state_model.ModeTransitionOutcome) {
	if outcome == (editor_state_model.ModeTransitionOutcome{}) {
		return
	}

	this.pendingOutcome = this.pendingOutcome.MergedWith(outcome)
	this.SetStatus(stateTransitionNotice(this.pendingOutcome), false)
}

func (this *State) takePendingNotice() string {
	notice := stateTransitionNotice(this.pendingOutcome)
	this.pendingOutcome = editor_state_model.ModeTransitionOutcome{}
	return notice
}

func stateTransitionNotice(outcome editor_state_model.ModeTransitionOutcome) string {
	switch {
	case outcome.TournamentCountCorrected && outcome.ManualEditsDiscarded:
		return "Tournament mode requires 2 players - the player count was corrected " +
			"and the manual zone layout was discarded."
	case outcome.TournamentCountCorrected:
		return "Tournament mode requires 2 players - the player count was corrected."
	case outcome.ManualEditsDiscarded:
		return "The manual zone layout was discarded because the game mode change regenerates the map."
	default:
		return ""
	}
}
