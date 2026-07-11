package constants

import (
	"slices"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

func GetValueOverrideSidsWithExclusions(excluded []string) []string {
	sids := slices.DeleteFunc(
		buildValueOverrideSids(),
		func(sid string) bool { return slices.Contains(excluded, sid) })
	slices.SortStableFunc(sids, strings.Compare)
	return sids
}

func buildValueOverrideSids() []string {
	sids := []string{}
	sids = append(sids, buildHeroBuffAndMagicSids()...)
	sids = append(sids, buildNamedAndRandomUnitSids()...)
	sids = append(sids, buildBuildingAndMiscSids()...)
	sids = append(sids, buildMineAndResourceSids()...)
	sids = append(sids, buildRandomItemAndScrollSids()...)
	sids = append(sids, buildStatsAndSkillsObjectSids()...)
	sids = append(sids, buildGuardedResourceBankSids()...)
	sids = append(sids, buildVisionBuildingSids()...)
	return sids
}

func buildHeroBuffAndMagicSids() []string {
	heroBuffBuildings := registry.GetMapObjectHeroBuffBuildingValues()
	magicBuildings := registry.GetMapObjectMagicBuildingValues()

	return []string{
		heroBuffBuildings.BeerFountain,
		magicBuildings.CelestialSphere,
		heroBuffBuildings.CrystalTrail,
		heroBuffBuildings.Fountain,
		heroBuffBuildings.Fountain2,
		heroBuffBuildings.ManaWell,
		heroBuffBuildings.MysteriousStone,
		magicBuildings.MysticalTower,
		heroBuffBuildings.PileOfBooks,
		heroBuffBuildings.QuixsPath,
		heroBuffBuildings.Stables,
		heroBuffBuildings.TearOfTruth,
	}
}

func buildNamedAndRandomUnitSids() []string {
	namedUnitBanks := registry.GetMapObjectNamedUnitBankValues()
	randomUnitBanks := registry.GetMapObjectRandomUnitBankValues()

	return []string{
		namedUnitBanks.BorealCall,
		namedUnitBanks.JoustingRange,
		namedUnitBanks.PetrifiedMemorial,
		namedUnitBanks.PointOfBalance,
		randomUnitBanks.RandomHireTier1,
		randomUnitBanks.RandomHireTier2,
		randomUnitBanks.RandomHireTier3,
		randomUnitBanks.RandomHireTier4,
		randomUnitBanks.RandomHireTier5,
		randomUnitBanks.RandomHireTier6,
		randomUnitBanks.RandomHireTier7,
		namedUnitBanks.RitualPyre,
		namedUnitBanks.Gorge,
		namedUnitBanks.UnforgottenGrave,
	}
}

func buildBuildingAndMiscSids() []string {
	buildingObjects := registry.GetMapObjectBuildingValues()
	miscObjects := registry.GetMapObjectMiscellaneousValues()
	nonContentObjects := registry.GetMapObjectNonContentValues()

	return []string{
		buildingObjects.Arena,
		buildingObjects.Chimerologist,
		buildingObjects.EternalDragon,
		buildingObjects.FickleShrine,
		nonContentObjects.Forge,
		buildingObjects.HuntsmansCamp,
		nonContentObjects.InsarasEye,
		nonContentObjects.Market,
		nonContentObjects.Mirage,
		miscObjects.Prison,
		nonContentObjects.RemoteFoothold,
		buildingObjects.SacrificialShrine,
		nonContentObjects.Tavern,
		miscObjects.TownGate,
		buildingObjects.TreeOfAbundance,
	}
}

func buildMineAndResourceSids() []string {
	mineObjects := registry.GetMapObjectMineValues()
	resourceBanks := registry.GetMapObjectResourceBankValues()
	resourceObjects := registry.GetMapObjectResourceValues()

	return []string{
		mineObjects.AlchemyLab,
		mineObjects.CrystalMine,
		mineObjects.GemstoneMine,
		mineObjects.GoldMine,
		mineObjects.MercuryMine,
		mineObjects.OreMine,
		mineObjects.WoodMine,
		resourceBanks.MontyHall,
		resourceObjects.PandoraBox,
	}
}

func buildRandomItemAndScrollSids() []string {
	randomItemObjects := registry.GetMapObjectRandomItemValues()
	scrollObjects := registry.GetMapObjectScrollValues()

	return []string{
		scrollObjects.MythicScrollBox,
		randomItemObjects.RandomItemCommon,
		randomItemObjects.RandomItemEpic,
		randomItemObjects.RandomItemLegendary,
		randomItemObjects.RandomItemRare,
	}
}

func buildStatsAndSkillsObjectSids() []string {
	t1StatsAndSkillsObjects := registry.GetMapObjectT1StatsAndSkillsBuildingValues()
	t2StatsAndSkillsObjects := registry.GetMapObjectT2StatsAndSkillsBuildingValues()
	t3StatsAndSkillsObjects := registry.GetMapObjectT3StatsAndSkillsBuildingValues()

	return []string{
		t2StatsAndSkillsObjects.Circus,
		t3StatsAndSkillsObjects.CollegeOfWonder,
		t2StatsAndSkillsObjects.Fort,
		t2StatsAndSkillsObjects.InfernalCirque,
		t2StatsAndSkillsObjects.OrbObservatory,
		t2StatsAndSkillsObjects.University,
		t1StatsAndSkillsObjects.WiseOwl,
	}
}

func buildGuardedResourceBankSids() []string {
	t1GuardedResourceBanks := registry.GetMapObjectT1GuardedResourceBankValues()
	t3GuardedResourceBanks := registry.GetMapObjectT3GuardedResourceBankValues()

	return []string{
		t3GuardedResourceBanks.DragonUtopia,
		t3GuardedResourceBanks.ResearchLaboratory,
		t1GuardedResourceBanks.ShadyDen,
		t3GuardedResourceBanks.TroglodyteThrone,
		t3GuardedResourceBanks.UnstableRuins,
	}
}

func buildVisionBuildingSids() []string {
	visionBuildings := registry.GetMapObjectVisionBuildingValues()

	return []string{
		visionBuildings.FlatteringMirror,
		visionBuildings.Watchtower,
		visionBuildings.WindRose,
	}
}
