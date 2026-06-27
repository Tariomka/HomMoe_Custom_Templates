package providers

import "github.com/Tariomka/hommoe_custom_templates/internal/registry"

var (
	buildingObjects           = registry.GetMapObjectBuildingValues()
	championSelectRules       = registry.GetChampionSelectValues()
	gameModes                 = registry.GetGameModeValues()
	heroBuffBuildings         = registry.GetMapObjectHeroBuffBuildingValues()
	magicBuildings            = registry.GetMapObjectMagicBuildingValues()
	mapBonuses                = registry.GetMapBonusesValues()
	nonContentObjects         = registry.GetMapObjectNonContentValues()
	randomUnitBanks           = registry.GetMapObjectRandomUnitBankValues()
	receiversFilters          = registry.GetReceiversFilterValues()
	resourceObjects           = registry.GetMapObjectResourceValues()
	ruleTypes                 = registry.GetRuleTypeValues()
	spellSids                 = registry.GetSpellSidValues()
	t1GuardedResourceBanks    = registry.GetMapObjectT1GuardedResourceBankValues()
	t1StatsAndSkillsBuildings = registry.GetMapObjectT1StatsAndSkillsValues()
	t2StatsAndSkillsBuildings = registry.GetMapObjectT2StatsAndSkillsBuildingValues()
	unitBanks                 = registry.GetMapObjectNamedUnitBankValues()
	visionBuildings           = registry.GetMapObjectVisionBuildingValues()
	winConditionValues        = registry.GetWinningConditionValues()
	zoneLayouts               = registry.GetLayoutValues()
)
