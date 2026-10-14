package connection_editor_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// IManualEditDeltaService records manual edits as structure plus overrides on
// top of a fresh generation, and rebuilds them on a later one.
type IManualEditDeltaService interface {
	// Capture records the edited zones and connections against the fresh generation.
	Capture(request models.ManualEditCaptureRequest) editor_state_model.ManualEditSettings

	// Expand rebuilds every manual zone and connection on the fresh generation,
	// reapplying changed castle options to the zones that own their castles.
	Expand(request models.ManualEditExpandRequest) models.ManualEditExpansion
}
