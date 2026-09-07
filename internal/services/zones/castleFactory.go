package zones

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type CastleFactory struct{}

func NewCastleFactory() zone_interfaces.ICastleFactory {
	return &CastleFactory{}
}

func (this *CastleFactory) CreatePlayerOwnedCastles(
	matchPlayerFaction bool,
	owner string,
	castleCount int) []template_model.MainObject {
	var castles []template_model.MainObject

	for range castleCount {
		objectBuilder := variant_content.NewObjectBuilder().
			WithTypeCity().
			WithOwner(owner).
			WithCastleQualityPoor().
			WithPlacementUniform()
		if matchPlayerFaction {
			objectBuilder = objectBuilder.WithFactionMatch()
		} else {
			objectBuilder = objectBuilder.WithFaction("Random")
		}
		castles = append(castles, objectBuilder.Build())
	}
	return castles
}

func (this *CastleFactory) CreatePlayerUnclaimedCastles(
	matchPlayerFaction bool,
	guardValue, castleCount int) []template_model.MainObject {
	var castles []template_model.MainObject

	for range castleCount {
		objectBuilder := variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardChance(1).
			WithGuardValue(guardValue).
			WithGuardWeeklyIncrement(common_connections.GetGuardWeeklyIncrements().Standard).
			WithCastleQualityMedium().
			WithPlacementUniform().
			WithPlacementArgs("false", "-0.8", "3")
		if matchPlayerFaction {
			objectBuilder = objectBuilder.WithFactionMatch()
		} else {
			objectBuilder = objectBuilder.WithFaction("Random")
		}
		castles = append(castles, objectBuilder.Build())
	}
	return castles
}

func (this *CastleFactory) CreateNeutralZoneCastles(
	profile neutral_zone.Profile,
	tuning models.GenerationTuning,
	castleCount int,
	isHoldCityZone bool) []template_model.MainObject {
	var castles []template_model.MainObject

	if castleCount > 0 {
		objectBuilder := variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardChance(1).
			WithGuardWeeklyIncrement(0.10).
			WithFactionFromList()

		if isHoldCityZone {
			objectBuilder = objectBuilder.
				WithGuardValue(tuning.ScaleByBorderGuardStrength(max(profile.PrimaryCityGuardValue, 20_000))).
				WithCastleQualityUltraRich().
				WithPlacementCenter().
				WithHoldCityWinCon()
		} else {
			objectBuilder = objectBuilder.
				WithGuardValue(tuning.ScaleByBorderGuardStrength(profile.PrimaryCityGuardValue)).
				WithCastleQuality(profile.PrimaryBuildingsSid).
				WithPlacementUniform().
				WithPlacementArgs("true", "0.8", "2")
		}

		castles = append(castles, objectBuilder.Build())
	}

	for range castleCount - 1 {
		castles = append(castles, variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardChance(1).
			WithGuardValue(tuning.ScaleByBorderGuardStrength(profile.ExtraCityGuardValue)).
			WithGuardWeeklyIncrement(0.10).
			WithCastleQuality(profile.ExtraBuildingsSid).
			WithFactionFromList().
			WithPlacementUniform().
			WithPlacementArgs("false", "-0.8", "3").
			Build())
	}

	return castles
}

func (this *CastleFactory) CreateHubZoneCastles(
	tuning models.GenerationTuning,
	castleCount int,
	isHoldCityZone bool) []template_model.MainObject {
	var castles []template_model.MainObject
	newCastleBuilder := func() *variant_content.MainObjectBuilder {
		return variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardWeeklyIncrement(0.10).
			WithFactionFromList()
	}
	buildHoldCityCastle := func(builder *variant_content.MainObjectBuilder) template_model.MainObject {
		return builder.
			WithGuardChance(1).
			WithGuardValue(tuning.ScaleByNeutralGuardStrength(25_000)).
			WithCastleQualityUltraRich().
			WithPlacementCenter().
			WithHoldCityWinCon().
			Build()
	}
	buildCastle := func(builder *variant_content.MainObjectBuilder) template_model.MainObject {
		return builder.
			WithGuardChance(0.5).
			WithGuardValue(tuning.ScaleByNeutralGuardStrength(16_000)).
			WithCastleQualityRich().
			WithPlacementUniform().
			WithPlacementArgs("true", "0.8", "2").
			Build()
	}

	if castleCount > 0 && isHoldCityZone {
		castles = append(castles, buildHoldCityCastle(newCastleBuilder()))
	} else if castleCount > 0 {
		castles = append(castles, buildCastle(newCastleBuilder()))
	}

	for range castleCount - 1 {
		castles = append(castles, buildCastle(newCastleBuilder()))
	}

	return castles
}

func (this *CastleFactory) CreatePlayerSpawnCastle(playerName string, guardValue int) template_model.MainObject {
	return variant_content.NewObjectBuilder().
		WithTypeSpawn().
		WithSpawn(playerName).
		WithNoGuardWhenOwned().
		WithGuardChance(1).
		WithGuardValue(guardValue).
		WithGuardWeeklyIncrement(0.10).
		WithCastleQualityDefault().
		WithPlacementUniform().
		WithPlacementArgs("true", "0.7", "0").
		Build()
}

func (this *CastleFactory) CreateAbandonedOutposts(
	profile neutral_zone.Profile,
	tuning models.GenerationTuning,
	count int) []template_model.MainObject {
	var outposts []template_model.MainObject

	for range count {
		outposts = append(outposts,
			variant_content.NewObjectBuilder().
				WithTypeAbandonedOutpost().
				WithGuardChance(1).
				WithGuardValue(tuning.ScaleByBorderGuardStrength(profile.ExtraCityGuardValue)).
				WithGuardWeeklyIncrement(0.10).
				WithCastleQuality(profile.ExtraBuildingsSid).
				WithFactionFromList().
				WithPlacementUniform().
				WithPlacementArgs("false", "-0.8", "3").
				Build())
	}
	return outposts
}
