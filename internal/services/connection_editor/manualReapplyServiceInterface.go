package connection_editor

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// IManualReapplyService is the contract for pushing changed castle-count
// options into the manually edited zone snapshot.
type IManualReapplyService interface {
	// ApplyCastleSettingChanges rewrites, in place, the castles of the manually
	// edited zones whose castle-count option changed.
	ApplyCastleSettingChanges(
		zones []entities.Zone,
		changes editor_state_model.CastleSettingChanges,
		configuration *config.GeneratorConfig)

	// SetNeutralZoneCastleCount rebuilds only the zone's City castles for the
	// new count, keeping everything else untouched.
	SetNeutralZoneCastleCount(zone *entities.Zone, castleCount int, tuning models.GenerationTuning)
}
