package zone_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type IZoneFactory interface {
	CreateSpawnZone(input models.SpawnZoneCreationRequest) template_model.Zone
	CreateNeutralZone(input models.NeutralZoneCreationRequest) template_model.Zone
	CreateHubZone(input models.HubZoneCreationRequest) template_model.Zone
}
