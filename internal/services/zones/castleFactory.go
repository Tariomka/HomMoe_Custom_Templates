package zones

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
)

type CastleFactory struct{}

func NewCastleFactory() *CastleFactory {
	return &CastleFactory{}
}

func (this *CastleFactory) CreatePlayerOwnedCastles(
	matchPlayerFaction bool,
	owner string,
	castleCount int) []entities.MainObject {
	var castles []entities.MainObject

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
	guardValue, castleCount int) []entities.MainObject {
	var castles []entities.MainObject

	for range castleCount {
		objectBuilder := variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardChance(1).
			WithGuardValue(guardValue).
			WithGuardWeeklyIncrement(0.15).
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
	isHoldCityZone bool) []entities.MainObject {
	var castles []entities.MainObject

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
	isHoldCityZone bool) []entities.MainObject {
	var castles []entities.MainObject
	newCastleBuilder := func() *variant_content.MainObjectBuilder {
		return variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardWeeklyIncrement(0.10).
			WithFactionFromList()
	}
	buildHoldCityCastle := func(builder *variant_content.MainObjectBuilder) entities.MainObject {
		return builder.
			WithGuardChance(1).
			WithGuardValue(tuning.ScaleByNeutralGuardStrength(25_000)).
			WithCastleQualityUltraRich().
			WithPlacementCenter().
			WithHoldCityWinCon().
			Build()
	}
	buildCastle := func(builder *variant_content.MainObjectBuilder) entities.MainObject {
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

func (this *CastleFactory) createPlayerSpawnCastle(playerName string, guardValue int) entities.MainObject {
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

func (this *CastleFactory) createAbandonedOutposts(
	profile neutral_zone.Profile,
	tuning models.GenerationTuning,
	count int) []entities.MainObject {
	var outposts []entities.MainObject

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
