package constants

import (
	"slices"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

// ValueOverrideSids are the known object / encounter SIDs offered by the
// guard-value-override picker. SIDs come from the registry; they reference
// world objects whose guard value can be overridden via valueOverrides.
var ValueOverrideSids = buildValueOverrideSids()

func buildValueOverrideSids() []string {
	buildingObjects := registry.GetMapObjectBuildingValues()
	heroBuffBuildings := registry.GetMapObjectHeroBuffBuildingValues()
	magicBuildings := registry.GetMapObjectMagicBuildingValues()
	mineObjects := registry.GetMapObjectMineValues()
	miscObjects := registry.GetMapObjectMiscellaneousValues()
	namedUnitBanks := registry.GetMapObjectNamedUnitBankValues()
	nonContentObjects := registry.GetMapObjectNonContentValues()
	randomUnitBanks := registry.GetMapObjectRandomUnitBankValues()
	randomItemObjects := registry.GetMapObjectRandomItemValues()
	resourceBanks := registry.GetMapObjectResourceBankValues()
	resourceObjects := registry.GetMapObjectResourceValues()
	scrollObjects := registry.GetMapObjectScrollValues()
	t1StatsAndSkillsObjects := registry.GetMapObjectT1StatsAndSkillsBuildingValues()
	t2StatsAndSkillsObjects := registry.GetMapObjectT2StatsAndSkillsBuildingValues()
	t3StatsAndSkillsObjects := registry.GetMapObjectT3StatsAndSkillsBuildingValues()
	t1GuardedResourceBanks := registry.GetMapObjectT1GuardedResourceBankValues()
	t3GuardedResourceBanks := registry.GetMapObjectT3GuardedResourceBankValues()
	visionBuildings := registry.GetMapObjectVisionBuildingValues()

	return []string{
		mineObjects.AlchemyLab,
		buildingObjects.Arena,
		heroBuffBuildings.BeerFountain,
		namedUnitBanks.BorealCall,
		magicBuildings.CelestialSphere,
		buildingObjects.Chimerologist,
		t2StatsAndSkillsObjects.Circus,
		t3StatsAndSkillsObjects.CollegeOfWonder,
		heroBuffBuildings.CrystalTrail,
		t3GuardedResourceBanks.DragonUtopia,
		buildingObjects.EternalDragon,
		buildingObjects.FickleShrine,
		visionBuildings.FlatteringMirror,
		nonContentObjects.Forge,
		t2StatsAndSkillsObjects.Fort,
		heroBuffBuildings.Fountain,
		heroBuffBuildings.Fountain2,
		buildingObjects.HuntsmansCamp,
		t2StatsAndSkillsObjects.InfernalCirque,
		nonContentObjects.InsarasEye,
		namedUnitBanks.JoustingRange,
		heroBuffBuildings.ManaWell,
		nonContentObjects.Market,
		mineObjects.CrystalMine,
		mineObjects.GemstoneMine,
		mineObjects.GoldMine,
		mineObjects.MercuryMine,
		mineObjects.OreMine,
		mineObjects.WoodMine,
		nonContentObjects.Mirage,
		resourceBanks.MontyHall,
		heroBuffBuildings.MysteriousStone,
		magicBuildings.MysticalTower,
		scrollObjects.MythicScrollBox,
		t2StatsAndSkillsObjects.OrbObservatory,
		resourceObjects.PandoraBox,
		namedUnitBanks.PetrifiedMemorial,
		heroBuffBuildings.PileOfBooks,
		namedUnitBanks.PointOfBalance,
		miscObjects.Prison,
		heroBuffBuildings.QuixsPath,
		randomUnitBanks.RandomHireTier1,
		randomUnitBanks.RandomHireTier2,
		randomUnitBanks.RandomHireTier3,
		randomUnitBanks.RandomHireTier4,
		randomUnitBanks.RandomHireTier5,
		randomUnitBanks.RandomHireTier6,
		randomUnitBanks.RandomHireTier7,
		randomItemObjects.RandomItemCommon,
		randomItemObjects.RandomItemEpic,
		randomItemObjects.RandomItemLegendary,
		randomItemObjects.RandomItemRare,
		nonContentObjects.RemoteFoothold,
		t3GuardedResourceBanks.ResearchLaboratory,
		namedUnitBanks.RitualPyre,
		buildingObjects.SacrificialShrine,
		t1GuardedResourceBanks.ShadyDen,
		heroBuffBuildings.Stables,
		nonContentObjects.Tavern,
		heroBuffBuildings.TearOfTruth,
		namedUnitBanks.Gorge,
		miscObjects.TownGate,
		buildingObjects.TreeOfAbundance,
		t3GuardedResourceBanks.TroglodyteThrone,
		namedUnitBanks.UnforgottenGrave,
		t2StatsAndSkillsObjects.University,
		t3GuardedResourceBanks.UnstableRuins,
		visionBuildings.Watchtower,
		visionBuildings.WindRose,
		t1StatsAndSkillsObjects.WiseOwl,
	}
}

func GetValueOverrideSidsWithExclusions(excluded []string) []string {
	sids := slices.DeleteFunc(
		buildValueOverrideSids(),
		func(sid string) bool { return slices.Contains(excluded, sid) })
	slices.SortStableFunc(sids, strings.Compare)
	return sids
}
