package constants

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

var (
	randomHiresBuildings             = registry.GetMandatoryContentRandomHiresBuildingValues()
	resourceBanksBuildings           = registry.GetMandatoryContentResourceBanksBuildingValues()
	pandoraBoxes                     = registry.GetMandatoryContentPandoraBoxValues()
	miscellaneous                    = registry.GetMandatoryContentMiscellaneousValues()
	basicRandomHiresBuildings        = registry.GetMandatoryContentBasicRandomHiresBuildingValues()
	basicResourceBanksBuildings      = registry.GetMandatoryContentBasicResourceBanksBuildingValues()
	basicMines                       = registry.GetMandatoryContentBasicMinesValues()
	basicHeroStatsAndSkillsBuildings = registry.GetMandatoryContentBasicHeroStatsAndSkillsBuildingValues()
	basicScrollBoxes                 = registry.GetMandatoryContentBasicScrollBoxValues()
	basicGuardedUnitBanks            = registry.GetMandatoryContentBasicGuardedUnitBankValues()
	basicMagicBuildings              = registry.GetMandatoryContentBasicMagicBuildingValues()
	basicHeroExpBuildings            = registry.GetMandatoryContentBasicHeroExperienceBuildingValues()
	basicMiscellaneous               = registry.GetMandatoryContentBasicMiscellaneousValues()
)

// IncludeListIds enumerates the named include-list SIDs (group references
// that resolve to many concrete content items at generation time).
var IncludeListIds = struct {
	RandomHiresLowTier              models.SidMapping
	RandomHiresHighTier             models.SidMapping
	RandomHiresAllTier              models.SidMapping
	RandomHiresAllTierWeighted      models.SidMapping
	ResourceBanksTier1              models.SidMapping
	ResourceBanksTier2              models.SidMapping
	GuardedBanksTier1               models.SidMapping
	GuardedBanksTier2               models.SidMapping
	GuardedBanksTier3               models.SidMapping
	BasicStorageBanks               models.SidMapping
	RandomRareMines                 models.SidMapping
	RandomRareMinesBiomeRestricted  models.SidMapping
	RandomGuardedUnitBank           models.SidMapping
	HeroBuffTier1                   models.SidMapping
	HeroExpTier2                    models.SidMapping
	HeroStatsAndSkillsTier1         models.SidMapping
	HeroStatsAndSkillsTier2         models.SidMapping
	HeroStatsAndSkillsTier3         models.SidMapping
	MagicBuildingsTier1             models.SidMapping
	MagicBuildingsTier2             models.SidMapping
	HeroImprovementUncommon         models.SidMapping
	VisionBuildingsTier1            models.SidMapping
	RandomPickupItems               models.SidMapping
	MythicScrollBoxPickup           models.SidMapping
	PandoraBoxArmyLowTier           models.SidMapping
	PandoraBoxArmyHighTier          models.SidMapping
	UtopiaBuildings                 models.SidMapping
	EpicGuardedResourceBanks        models.SidMapping
	GuardedUnitBanksBiomeRestricted models.SidMapping
	GuardedUnitBanksNoBiome         models.SidMapping
}{
	RandomHiresLowTier:              models.SidMapping{Sid: randomHiresBuildings.RandomHiresLowTier, Name: "Random Hires Low Tier"},
	RandomHiresHighTier:             models.SidMapping{Sid: randomHiresBuildings.RandomHiresHighTier, Name: "Random Hires High Tier"},
	RandomHiresAllTier:              models.SidMapping{Sid: basicRandomHiresBuildings.BasicRandomHires, Name: "Random Hires Any Tier"},
	RandomHiresAllTierWeighted:      models.SidMapping{Sid: randomHiresBuildings.RandomHires, Name: "Random Hires Any Tier (Weighted)"},
	ResourceBanksTier1:              models.SidMapping{Sid: basicResourceBanksBuildings.BasicResourceBanksTier1, Name: "Resource Banks T1"},
	ResourceBanksTier2:              models.SidMapping{Sid: basicResourceBanksBuildings.BasicResourceBanksTier2, Name: "Resource Banks T2"},
	GuardedBanksTier1:               models.SidMapping{Sid: basicResourceBanksBuildings.BasicGuardedResourceBanksTier1, Name: "Guarded Banks T1"},
	GuardedBanksTier2:               models.SidMapping{Sid: basicResourceBanksBuildings.BasicGuardedResourceBanksTier2, Name: "Guarded Banks T2"},
	GuardedBanksTier3:               models.SidMapping{Sid: basicResourceBanksBuildings.BasicGuardedResourceBanksTier3, Name: "Guarded Banks T3"},
	BasicStorageBanks:               models.SidMapping{Sid: basicMiscellaneous.BasicStorage, Name: "Random Basic Storage"},
	RandomRareMines:                 models.SidMapping{Sid: basicMines.BasicMinesRare, Name: "Random Rare Mine"},
	RandomRareMinesBiomeRestricted:  models.SidMapping{Sid: basicMines.BasicMinesRareBiomeRestricted, Name: "Random Rare Mine (Biome Restricted)"},
	RandomGuardedUnitBank:           models.SidMapping{Sid: basicGuardedUnitBanks.BasicGuardedUnitBanks, Name: "Random Guarded Unit Bank"},
	HeroBuffTier1:                   models.SidMapping{Sid: basicMiscellaneous.BasicHeroBuffBuildingTier1, Name: "Random Hero Buff Tier 1"},
	HeroExpTier2:                    models.SidMapping{Sid: basicHeroExpBuildings.BasicHeroExperienceBuildingTier2, Name: "Random Hero Exp Tier 2"},
	HeroStatsAndSkillsTier1:         models.SidMapping{Sid: basicHeroStatsAndSkillsBuildings.BasicHeroStatsAndSkillsTier1, Name: "Random Hero Stat/Skill Tier 1"},
	HeroStatsAndSkillsTier2:         models.SidMapping{Sid: basicHeroStatsAndSkillsBuildings.BasicHeroStatsAndSkillsTier2, Name: "Random Hero Stat/Skill Tier 2"},
	HeroStatsAndSkillsTier3:         models.SidMapping{Sid: basicHeroStatsAndSkillsBuildings.BasicHeroStatsAndSkillsTier3, Name: "Random Hero Stat/Skill Tier 3"},
	MagicBuildingsTier1:             models.SidMapping{Sid: basicMagicBuildings.BasicMagicBuildingsTier1, Name: "Random Magic Building Tier 1"},
	MagicBuildingsTier2:             models.SidMapping{Sid: basicMagicBuildings.BasicMagicBuildingsTier2, Name: "Random Magic Building Tier 2"},
	HeroImprovementUncommon:         models.SidMapping{Sid: miscellaneous.UncommonHeroBanks, Name: "Uncommon Hero Improvement"},
	VisionBuildingsTier1:            models.SidMapping{Sid: basicMiscellaneous.BasicVisionBuildingsTier1, Name: "Random Vision Building Tier 1"},
	RandomPickupItems:               models.SidMapping{Sid: basicMiscellaneous.BasicPickupRandomItems, Name: "Random Pickup Items"},
	MythicScrollBoxPickup:           models.SidMapping{Sid: basicScrollBoxes.BasicPickupMythicScrollBox, Name: "Random Mythic Scroll Box"},
	PandoraBoxArmyLowTier:           models.SidMapping{Sid: pandoraBoxes.PickupPandoraBoxArmyLowTier, Name: "Pandora Box Army (Low Tier)"},
	PandoraBoxArmyHighTier:          models.SidMapping{Sid: pandoraBoxes.PickupPandoraBoxArmyHighTier, Name: "Pandora Box Army (High Tier)"},
	UtopiaBuildings:                 models.SidMapping{Sid: miscellaneous.UtopiaBuilding, Name: "Utopia (Dragon/Unstable/Lab)"},
	EpicGuardedResourceBanks:        models.SidMapping{Sid: resourceBanksBuildings.GuardedResourceBanksEpic, Name: "Epic Guarded Resource Banks"},
	GuardedUnitBanksBiomeRestricted: models.SidMapping{Sid: basicGuardedUnitBanks.BasicGuardedUnitBanksBiomeRestricted, Name: "Guarded Unit Bank (Biome Restricted)"},
	GuardedUnitBanksNoBiome:         models.SidMapping{Sid: basicGuardedUnitBanks.BasicGuardedUnitBanksNotRestricted, Name: "Guarded Unit Bank (No Biome Restriction)"},
}
