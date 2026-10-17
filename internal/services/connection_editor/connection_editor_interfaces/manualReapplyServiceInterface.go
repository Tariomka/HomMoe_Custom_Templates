package connection_editor_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// IManualReapplyService is the contract for pushing changed castle-count
// options into the manually edited zone snapshot.
type IManualReapplyService interface {
	// ApplyCastleSettingChanges rewrites, in place, the castles of the eligible
	// zones whose castle-count option changed; the tuning counts every zone.
	ApplyCastleSettingChanges(
		zones []template_model.Zone,
		eligible []bool,
		changes editor_state_model.CastleSettingChanges,
		configuration *config.GeneratorConfig)

	// SetNeutralZoneCastleCount rebuilds only the zone's City castles for the
	// new count, keeping everything else untouched.
	SetNeutralZoneCastleCount(zone *template_model.Zone, castleCount int, tuning models.GenerationTuning)
}
