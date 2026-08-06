package zone_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

type ICastleFactory interface {
	CreatePlayerOwnedCastles(
		matchPlayerFaction bool,
		owner string,
		castleCount int) []entities.MainObject

	CreatePlayerUnclaimedCastles(
		matchPlayerFaction bool,
		guardValue, castleCount int) []entities.MainObject

	CreateNeutralZoneCastles(
		profile neutral_zone.Profile,
		tuning models.GenerationTuning,
		castleCount int,
		isHoldCityZone bool) []entities.MainObject

	CreateHubZoneCastles(
		tuning models.GenerationTuning,
		castleCount int,
		isHoldCityZone bool) []entities.MainObject

	CreatePlayerSpawnCastle(playerName string, guardValue int) entities.MainObject

	CreateAbandonedOutposts(
		profile neutral_zone.Profile,
		tuning models.GenerationTuning,
		count int) []entities.MainObject
}
