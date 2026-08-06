package zone_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

type IZoneFactory interface {
	CreateSpawnZone(input models.SpawnZoneCreationRequest) entities.Zone
	CreateNeutralZone(input models.NeutralZoneCreationRequest) entities.Zone
	CreateHubZone(input models.HubZoneCreationRequest) entities.Zone
}
