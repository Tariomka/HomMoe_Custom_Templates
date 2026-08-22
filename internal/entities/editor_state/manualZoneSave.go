package editor_state

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// ManualZoneSave persists a zone edited in the manual zone editor. The zone's
// ManualPosition is captured separately because entities.Zone deliberately
// omits it from JSON (json:"-"), yet that normalized position is the essential
// piece of a hand-made layout and must survive a save/load round-trip.
type ManualZoneSave struct {
	Zone           entities.Zone `json:"zone"`
	ManualPosition *[2]float64   `json:"manualPosition,omitempty"`
}
