package editor_state_v1

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// ManualZoneSave is the frozen v1 snapshot - see editorState.go. Never edit.
// ManualPosition is the [x, y] array v1 wrote; v2 replaced it with an object
// and added the generator stamps v1 never persisted.
type ManualZoneSave struct {
	Zone           entities.Zone `json:"zone"`
	ManualPosition *[2]float64   `json:"manualPosition,omitempty"`
	Quality        *int8         `json:"quality,omitempty"`
}
