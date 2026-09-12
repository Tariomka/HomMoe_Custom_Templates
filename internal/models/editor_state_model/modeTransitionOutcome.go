package editor_state_model

type ModeTransitionOutcome struct {
	ManualEditsDiscarded     bool // true when a manual snapshot existed and was dropped.
	TournamentCountCorrected bool // true when the requested state carried a != 2 tournament player count
}

func (this ModeTransitionOutcome) MergedWith(other ModeTransitionOutcome) ModeTransitionOutcome {
	return ModeTransitionOutcome{
		ManualEditsDiscarded:     this.ManualEditsDiscarded || other.ManualEditsDiscarded,
		TournamentCountCorrected: this.TournamentCountCorrected || other.TournamentCountCorrected,
	}
}
