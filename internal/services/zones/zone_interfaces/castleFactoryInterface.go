package zone_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ICastleFactory interface {
	CreatePlayerOwnedCastles(
		matchPlayerFaction bool,
		owner string,
		castleCount int) []template_model.MainObject

	CreatePlayerUnclaimedCastles(
		matchPlayerFaction bool,
		guardValue, castleCount int) []template_model.MainObject

	CreateNeutralZoneCastles(
		profile neutral_zone.Profile,
		tuning models.GenerationTuning,
		castleCount int,
		isHoldCityZone bool) []template_model.MainObject

	CreateHubZoneCastles(
		tuning models.GenerationTuning,
		castleCount int,
		isHoldCityZone bool) []template_model.MainObject

	CreatePlayerSpawnCastle(playerName string, guardValue int) template_model.MainObject

	CreateAbandonedOutposts(
		profile neutral_zone.Profile,
		tuning models.GenerationTuning,
		count int) []template_model.MainObject
}
