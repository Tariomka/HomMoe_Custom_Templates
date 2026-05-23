package constants

import "github.com/Tariomka/hommoe_custom_templates/internal/models"

// IncludeListIds enumerates the named include-list SIDs (group references
// that resolve to many concrete content items at generation time).
//
//nolint:gochecknoglobals // semantic registry
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
	RandomHiresLowTier:              models.SidMapping{Sid: "content_list_building_random_hires_low_tier", Name: "Random Hires Low Tier"},
	RandomHiresHighTier:             models.SidMapping{Sid: "content_list_building_random_hires_high_tier", Name: "Random Hires High Tier"},
	RandomHiresAllTier:              models.SidMapping{Sid: "basic_content_list_building_random_hires", Name: "Random Hires Any Tier"},
	RandomHiresAllTierWeighted:      models.SidMapping{Sid: "content_list_building_random_hires", Name: "Random Hires Any Tier (Weighted)"},
	ResourceBanksTier1:              models.SidMapping{Sid: "basic_content_list_building_guarded_resource_banks_tier_1", Name: "Resource Banks T1"},
	ResourceBanksTier2:              models.SidMapping{Sid: "basic_content_list_building_guarded_resource_banks_tier_2", Name: "Resource Banks T2"},
	GuardedBanksTier1:               models.SidMapping{Sid: "basic_content_list_building_guarded_resource_banks_tier_1", Name: "Guarded Banks T1"},
	GuardedBanksTier2:               models.SidMapping{Sid: "basic_content_list_building_guarded_resource_banks_tier_2", Name: "Guarded Banks T2"},
	GuardedBanksTier3:               models.SidMapping{Sid: "basic_content_list_building_guarded_resource_banks_tier_3", Name: "Guarded Banks T3"},
	BasicStorageBanks:               models.SidMapping{Sid: "basic_content_list_basic_storage", Name: "Random Basic Storage"},
	RandomRareMines:                 models.SidMapping{Sid: "basic_content_list_rare_mines", Name: "Random Rare Mine"},
	RandomRareMinesBiomeRestricted:  models.SidMapping{Sid: "basic_content_list_rare_mines_by_biome", Name: "Random Rare Mine (Biome Restricted)"},
	RandomGuardedUnitBank:           models.SidMapping{Sid: "basic_content_list_building_guarded_units_banks", Name: "Random Guarded Unit Bank"},
	HeroBuffTier1:                   models.SidMapping{Sid: "basic_content_list_building_hero_buff_tier_1", Name: "Random Hero Buff Tier 1"},
	HeroExpTier2:                    models.SidMapping{Sid: "basic_content_list_building_hero_exp_tier_2", Name: "Random Hero Exp Tier 2"},
	HeroStatsAndSkillsTier1:         models.SidMapping{Sid: "basic_content_list_building_hero_stats_and_skills_tier_1", Name: "Random Hero Stat/Skill Tier 1"},
	HeroStatsAndSkillsTier2:         models.SidMapping{Sid: "basic_content_list_building_hero_stats_and_skills_tier_2", Name: "Random Hero Stat/Skill Tier 2"},
	HeroStatsAndSkillsTier3:         models.SidMapping{Sid: "basic_content_list_building_hero_stats_and_skills_tier_3", Name: "Random Hero Stat/Skill Tier 3"},
	MagicBuildingsTier1:             models.SidMapping{Sid: "basic_content_list_building_magic_tier_1", Name: "Random Magic Building Tier 1"},
	MagicBuildingsTier2:             models.SidMapping{Sid: "basic_content_list_building_magic_tier_2", Name: "Random Magic Building Tier 2"},
	HeroImprovementUncommon:         models.SidMapping{Sid: "content_list_building_uncommon_hero_banks", Name: "Uncommon Hero Improvement"},
	VisionBuildingsTier1:            models.SidMapping{Sid: "basic_content_list_vision_buildings_tier_1", Name: "Random Vision Building Tier 1"},
	RandomPickupItems:               models.SidMapping{Sid: "basic_content_list_pickup_random_items", Name: "Random Pickup Items"},
	MythicScrollBoxPickup:           models.SidMapping{Sid: "basic_content_list_pickup_mythic_scroll_box", Name: "Random Mythic Scroll Box"},
	PandoraBoxArmyLowTier:           models.SidMapping{Sid: "content_list_pickup_pandora_box_army_low_tier", Name: "Pandora Box Army (Low Tier)"},
	PandoraBoxArmyHighTier:          models.SidMapping{Sid: "content_list_pickup_pandora_box_army_high_tier", Name: "Pandora Box Army (High Tier)"},
	UtopiaBuildings:                 models.SidMapping{Sid: "content_list_building_utopia", Name: "Utopia (Dragon/Unstable/Lab)"},
	EpicGuardedResourceBanks:        models.SidMapping{Sid: "content_list_building_epic_guarded_resource_banks", Name: "Epic Guarded Resource Banks"},
	GuardedUnitBanksBiomeRestricted: models.SidMapping{Sid: "basic_content_list_building_guarded_units_banks_only_biome_restriction", Name: "Guarded Unit Bank (Biome Restricted)"},
	GuardedUnitBanksNoBiome:         models.SidMapping{Sid: "basic_content_list_building_guarded_units_banks_no_biome_restriction", Name: "Guarded Unit Bank (No Biome Restriction)"},
}
