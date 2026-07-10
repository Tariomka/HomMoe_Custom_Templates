package registry

type allContentLists struct {
	contentIncludeLists
	basicContentLists
}

// GetAllMandatoryContentIncludeListValues returns all of the content group names used for
//
//	mandatoryContent.content.includeLists
func GetAllMandatoryContentIncludeListValues() allContentLists {
	return allContentLists{
		contentIncludeLists: GetMandatoryContentIncludeListValues(),
		basicContentLists:   GetMandatoryContentBasicIncludeListValues(),
	}
}

type contentIncludeLists struct {
	contentRandomHiresBuildings
	contentResourceBanksBuildings
	contentResources
	contentMines
	contentInteractables
	contentHeroStatsAndSkillsBuildings
	contentScrollBoxes
	contentEnchantedScrollBoxes
	contentMythicScrollBoxes
	contentPandoraBoxes
	contentGuardedUnitBanks
	contentMagicBuildings
	contentMiscellaneous
}

// GetMandatoryContentIncludeListValues returns the "content_list_*" content group names used for
//
//	mandatoryContent.content.includeLists
func GetMandatoryContentIncludeListValues() contentIncludeLists {
	return contentIncludeLists{
		contentRandomHiresBuildings:        GetMandatoryContentRandomHiresBuildingValues(),
		contentResourceBanksBuildings:      GetMandatoryContentResourceBanksBuildingValues(),
		contentResources:                   GetMandatoryContentResourcesValues(),
		contentMines:                       GetMandatoryContentMinesValues(),
		contentInteractables:               GetMandatoryContentInteractablesValues(),
		contentHeroStatsAndSkillsBuildings: GetMandatoryContentHeroStatsAndSkillsBuildingValues(),
		contentScrollBoxes:                 GetMandatoryContentScrollBoxValues(),
		contentEnchantedScrollBoxes:        GetMandatoryContentEnchantedScrollBoxValues(),
		contentMythicScrollBoxes:           GetMandatoryContentMythicScrollBoxValues(),
		contentPandoraBoxes:                GetMandatoryContentPandoraBoxValues(),
		contentGuardedUnitBanks:            GetMandatoryContentGuardedUnitBanksValues(),
		contentMagicBuildings:              GetMandatoryContentMagicBuildingValues(),
		contentMiscellaneous:               GetMandatoryContentMiscellaneousValues(),
	}
}

type contentRandomHiresBuildings struct {
	RandomHires         string
	RandomHiresLowTier  string
	RandomHiresHighTier string
}

func GetMandatoryContentRandomHiresBuildingValues() contentRandomHiresBuildings {
	return contentRandomHiresBuildings{
		RandomHires:         "content_list_building_random_hires",
		RandomHiresLowTier:  "content_list_building_random_hires_low_tier",
		RandomHiresHighTier: "content_list_building_random_hires_high_tier",
	}
}

type contentResourceBanksBuildings struct {
	ResourceBanksCommon           string
	ResourceBanksUncommon         string
	ResourceBanksUncommonFix      string
	GuardedResourceBanksCommonFix string
	GuardedResourceBanksUncommon  string
	GuardedResourceBanksEpic      string
}

func GetMandatoryContentResourceBanksBuildingValues() contentResourceBanksBuildings {
	return contentResourceBanksBuildings{
		ResourceBanksCommon:           "content_list_building_common_resource_banks",
		ResourceBanksUncommon:         "content_list_building_uncommon_resource_banks",
		ResourceBanksUncommonFix:      "content_list_building_uncommon_resource_banks_fix",
		GuardedResourceBanksCommonFix: "content_list_building_common_guarded_resource_bank_fix",
		GuardedResourceBanksUncommon:  "content_list_building_uncommon_guarded_resource_banks",
		GuardedResourceBanksEpic:      "content_list_building_epic_guarded_resource_banks",
	}
}

type contentResources struct {
	ResourcesBasic   string
	ResourcesRare    string
	ResourcesSpecial string
}

func GetMandatoryContentResourcesValues() contentResources {
	return contentResources{
		ResourcesBasic:   "content_list_basic_resources",
		ResourcesRare:    "content_list_rare_resources",
		ResourcesSpecial: "content_list_special_resources",
	}
}

type contentMines struct {
	MinesBasic   string
	MinesRare    string
	MinesSpecial string
}

func GetMandatoryContentMinesValues() contentMines {
	return contentMines{
		MinesBasic:   "content_list_basic_mines",
		MinesRare:    "content_list_rare_mines",
		MinesSpecial: "content_list_special_mines",
	}
}

type contentInteractables struct {
	InteractablesCommon         string
	InteractablesCommonBalanced string
	InteractablesUncommon       string
}

func GetMandatoryContentInteractablesValues() contentInteractables {
	return contentInteractables{
		InteractablesCommon:         "content_list_building_common_interact",
		InteractablesCommonBalanced: "content_list_building_common_interact_no_imbalance",
		InteractablesUncommon:       "content_list_building_uncommon_interact",
	}
}

type contentHeroStatsAndSkillsBuildings struct {
	HeroStatsAndSkillsCommon   string
	HeroStatsAndSkillsUncommon string
}

func GetMandatoryContentHeroStatsAndSkillsBuildingValues() contentHeroStatsAndSkillsBuildings {
	return contentHeroStatsAndSkillsBuildings{
		HeroStatsAndSkillsCommon:   "content_list_building_common_hero_stats",
		HeroStatsAndSkillsUncommon: "content_list_building_uncommon_hero_stats",
	}
}

type contentScrollBoxes struct {
	PickupScrollBox      string
	PickupScrollBoxTier1 string
	PickupScrollBoxTier2 string
	PickupScrollBoxTier3 string
	PickupScrollBoxTier4 string
	PickupScrollBoxTier5 string
}

func GetMandatoryContentScrollBoxValues() contentScrollBoxes {
	return contentScrollBoxes{
		PickupScrollBox:      "content_list_pickup_scroll_box",
		PickupScrollBoxTier1: "content_list_pickup_scroll_box_tier_1",
		PickupScrollBoxTier2: "content_list_pickup_scroll_box_tier_2",
		PickupScrollBoxTier3: "content_list_pickup_scroll_box_tier_3",
		PickupScrollBoxTier4: "content_list_pickup_scroll_box_tier_4",
		PickupScrollBoxTier5: "content_list_pickup_scroll_box_tier_5",
	}
}

type contentEnchantedScrollBoxes struct {
	PickupEnchantedScrollBox              string
	PickupEnchantedScrollBoxTier1         string
	PickupEnchantedScrollBoxTier2         string
	PickupEnchantedScrollBoxTier3         string
	PickupEnchantedScrollBoxTier4         string
	PickupEnchantedScrollBoxTier5         string
	PickupEnchantedScrollBoxControlSpells string
}

func GetMandatoryContentEnchantedScrollBoxValues() contentEnchantedScrollBoxes {
	return contentEnchantedScrollBoxes{
		PickupEnchantedScrollBox:              "content_list_pickup_enchanted_scroll_box",
		PickupEnchantedScrollBoxTier1:         "content_list_pickup_enchanted_scroll_box_tier_1",
		PickupEnchantedScrollBoxTier2:         "content_list_pickup_enchanted_scroll_box_tier_2",
		PickupEnchantedScrollBoxTier3:         "content_list_pickup_enchanted_scroll_box_tier_3",
		PickupEnchantedScrollBoxTier4:         "content_list_pickup_enchanted_scroll_box_tier_4",
		PickupEnchantedScrollBoxTier5:         "content_list_pickup_enchanted_scroll_box_tier_5",
		PickupEnchantedScrollBoxControlSpells: "content_list_pickup_enchanted_scroll_box_control_spells",
	}
}

type contentMythicScrollBoxes struct {
	PickupMythicScrollBox string
}

func GetMandatoryContentMythicScrollBoxValues() contentMythicScrollBoxes {
	return contentMythicScrollBoxes{
		PickupMythicScrollBox: "content_list_pickup_mythic_scroll_box",
	}
}

type contentPandoraBoxes struct {
	PickupPandoraBox             string
	PickupPandoraBoxArmy         string
	PickupPandoraBoxArmyLowTier  string
	PickupPandoraBoxArmyHighTier string
	PickupPandoraBoxExperience   string
	PickupPandoraBoxGold         string
}

func GetMandatoryContentPandoraBoxValues() contentPandoraBoxes {
	return contentPandoraBoxes{
		PickupPandoraBox:             "content_list_pickup_pandora_box",
		PickupPandoraBoxArmy:         "content_list_pickup_pandora_box_army",
		PickupPandoraBoxArmyLowTier:  "content_list_pickup_pandora_box_army_low_tier",
		PickupPandoraBoxArmyHighTier: "content_list_pickup_pandora_box_army_high_tier",
		PickupPandoraBoxExperience:   "content_list_pickup_pandora_box_exp",
		PickupPandoraBoxGold:         "content_list_pickup_pandora_box_gold",
	}
}

type contentGuardedUnitBanks struct {
	GuardedUnitsBanksUncommon      string
	GuardedUnitsBanksUncommonEqual string
}

func GetMandatoryContentGuardedUnitBanksValues() contentGuardedUnitBanks {
	return contentGuardedUnitBanks{
		GuardedUnitsBanksUncommon:      "content_list_building_uncommon_guarded_units_banks",
		GuardedUnitsBanksUncommonEqual: "content_list_building_uncommon_guarded_units_banks_equal",
	}
}

type contentMagicBuildings struct {
	MagicBuildingsCommon   string
	MagicBuildingsUncommon string
}

func GetMandatoryContentMagicBuildingValues() contentMagicBuildings {
	return contentMagicBuildings{
		MagicBuildingsCommon:   "content_list_building_commons_magic",
		MagicBuildingsUncommon: "content_list_building_uncommons_magic",
	}
}

type contentMiscellaneous struct {
	Storage                 string
	SpecialBuilding         string
	BasicBuildings          string
	VisionBuildings         string
	PickupRandomItems       string
	PickupPrison            string
	UtopiaBuilding          string
	TownGates               string
	BuffMovePointsBuilding  string
	CommonHeroBuffsBuilding string
	UncommonHeroBanks       string
}

func GetMandatoryContentMiscellaneousValues() contentMiscellaneous {
	return contentMiscellaneous{
		Storage:                 "content_list_basic_storage",
		SpecialBuilding:         "content_list_building_special",
		BasicBuildings:          "content_list_basic_buildings",
		VisionBuildings:         "content_list_vision_buildings",
		PickupRandomItems:       "content_list_pickup_random_items",
		PickupPrison:            "content_list_pickup_prison",
		UtopiaBuilding:          "content_list_building_utopia",
		TownGates:               "content_list_town_gates",
		BuffMovePointsBuilding:  "content_list_building_buff_movepoints",
		CommonHeroBuffsBuilding: "content_list_building_common_hero_buffs",
		UncommonHeroBanks:       "content_list_building_uncommon_hero_banks",
	}
}

type basicContentLists struct {
	basicContentRandomHiresBuildings
	basicContentResourceBanksBuildings
	basicContentResources
	basicContentMines
	basicContentInteractables
	basicContentHeroStatsAndSkillsBuildings
	basicContentScrollBoxes
	basicContentPandoraBoxes
	basicContentGuardedUnitBanks
	basicContentMagicBuildings
	basicContentHeroExperienceBuildings
	basicContentMiscellaneous
}

// GetMandatoryContentBasicIncludeListValues returns the "basic_content_list_*" content group names used for
//
//	mandatoryContent.content.includeLists
func GetMandatoryContentBasicIncludeListValues() basicContentLists {
	return basicContentLists{
		basicContentRandomHiresBuildings:        GetMandatoryContentBasicRandomHiresBuildingValues(),
		basicContentResourceBanksBuildings:      GetMandatoryContentBasicResourceBanksBuildingValues(),
		basicContentResources:                   GetMandatoryContentBasicResourcesValues(),
		basicContentMines:                       GetMandatoryContentBasicMinesValues(),
		basicContentInteractables:               GetMandatoryContentBasicInteractablesValues(),
		basicContentHeroStatsAndSkillsBuildings: GetMandatoryContentBasicHeroStatsAndSkillsBuildingValues(),
		basicContentScrollBoxes:                 GetMandatoryContentBasicScrollBoxValues(),
		basicContentPandoraBoxes:                GetMandatoryContentBasicPandoraBoxValues(),
		basicContentGuardedUnitBanks:            GetMandatoryContentBasicGuardedUnitBankValues(),
		basicContentMagicBuildings:              GetMandatoryContentBasicMagicBuildingValues(),
		basicContentHeroExperienceBuildings:     GetMandatoryContentBasicHeroExperienceBuildingValues(),
		basicContentMiscellaneous:               GetMandatoryContentBasicMiscellaneousValues(),
	}
}

type basicContentRandomHiresBuildings struct {
	BasicRandomHires string
}

func GetMandatoryContentBasicRandomHiresBuildingValues() basicContentRandomHiresBuildings {
	return basicContentRandomHiresBuildings{
		BasicRandomHires: "basic_content_list_building_random_hires",
	}
}

type basicContentResourceBanksBuildings struct {
	BasicResourceBanksTier1                     string
	BasicResourceBanksTier2                     string
	BasicGuardedResourceBanksTier1              string
	BasicGuardedResourceBanksTier2              string
	BasicGuardedResourceBanksTier2NoRestriction string
	BasicGuardedResourceBanksTier3              string
}

func GetMandatoryContentBasicResourceBanksBuildingValues() basicContentResourceBanksBuildings {
	return basicContentResourceBanksBuildings{
		BasicResourceBanksTier1:                     "basic_content_list_building_resource_banks_tier_1",
		BasicResourceBanksTier2:                     "basic_content_list_building_resource_banks_tier_2",
		BasicGuardedResourceBanksTier1:              "basic_content_list_building_guarded_resource_banks_tier_1",
		BasicGuardedResourceBanksTier2:              "basic_content_list_building_guarded_resource_banks_tier_2",
		BasicGuardedResourceBanksTier2NoRestriction: "basic_content_list_building_guarded_resource_banks_tier_2_no_biome_restriction",
		BasicGuardedResourceBanksTier3:              "basic_content_list_building_guarded_resource_banks_tier_3",
	}
}

type basicContentResources struct {
	BasicResourcesBasic   string
	BasicResourcesRare    string
	BasicResourcesSpecial string
}

func GetMandatoryContentBasicResourcesValues() basicContentResources {
	return basicContentResources{
		BasicResourcesBasic:   "basic_content_list_basic_resources",
		BasicResourcesRare:    "basic_content_list_rare_resources",
		BasicResourcesSpecial: "basic_content_list_special_resources",
	}
}

type basicContentMines struct {
	BasicMinesBasic               string
	BasicMinesRare                string
	BasicMinesRareBiomeRestricted string
	BasicMinesSpecial             string
}

func GetMandatoryContentBasicMinesValues() basicContentMines {
	return basicContentMines{
		BasicMinesBasic:               "basic_content_list_basic_mines",
		BasicMinesRare:                "basic_content_list_rare_mines",
		BasicMinesRareBiomeRestricted: "basic_content_list_rare_mines_by_biome",
		BasicMinesSpecial:             "basic_content_list_special_mines",
	}
}

type basicContentInteractables struct {
	BasicInteractablesUncommon string
	BasicInteractablesEpic     string
}

func GetMandatoryContentBasicInteractablesValues() basicContentInteractables {
	return basicContentInteractables{
		BasicInteractablesUncommon: "basic_content_list_building_uncommon_interact",
		BasicInteractablesEpic:     "basic_content_list_building_epic_interact",
	}
}

type basicContentHeroStatsAndSkillsBuildings struct {
	BasicHeroStatsAndSkillsTier1 string
	BasicHeroStatsAndSkillsTier2 string
	BasicHeroStatsAndSkillsTier3 string
}

func GetMandatoryContentBasicHeroStatsAndSkillsBuildingValues() basicContentHeroStatsAndSkillsBuildings {
	return basicContentHeroStatsAndSkillsBuildings{
		BasicHeroStatsAndSkillsTier1: "basic_content_list_building_hero_stats_and_skills_tier_1",
		BasicHeroStatsAndSkillsTier2: "basic_content_list_building_hero_stats_and_skills_tier_2",
		BasicHeroStatsAndSkillsTier3: "basic_content_list_building_hero_stats_and_skills_tier_3",
	}
}

type basicContentScrollBoxes struct {
	BasicPickupScrollBox          string
	BasicPickupEnchantedScrollBox string
	BasicPickupMythicScrollBox    string
}

func GetMandatoryContentBasicScrollBoxValues() basicContentScrollBoxes {
	return basicContentScrollBoxes{
		BasicPickupScrollBox:          "basic_content_list_pickup_scroll_box",
		BasicPickupEnchantedScrollBox: "basic_content_list_pickup_enchanted_scroll_box",
		BasicPickupMythicScrollBox:    "basic_content_list_pickup_mythic_scroll_box",
	}
}

type basicContentPandoraBoxes struct {
	BasicPickupPandoraBox string
}

func GetMandatoryContentBasicPandoraBoxValues() basicContentPandoraBoxes {
	return basicContentPandoraBoxes{
		BasicPickupPandoraBox: "basic_content_list_pickup_pandora_box",
	}
}

type basicContentGuardedUnitBanks struct {
	BasicGuardedUnitBanks                string
	BasicGuardedUnitBanksNotRestricted   string
	BasicGuardedUnitBanksBiomeRestricted string
}

func GetMandatoryContentBasicGuardedUnitBankValues() basicContentGuardedUnitBanks {
	return basicContentGuardedUnitBanks{
		BasicGuardedUnitBanks:                "basic_content_list_building_guarded_units_banks",
		BasicGuardedUnitBanksNotRestricted:   "basic_content_list_building_guarded_units_banks_no_biome_restriction",
		BasicGuardedUnitBanksBiomeRestricted: "basic_content_list_building_guarded_units_banks_only_biome_restriction",
	}
}

type basicContentMagicBuildings struct {
	BasicMagicBuildingsTier1 string
	BasicMagicBuildingsTier2 string
}

func GetMandatoryContentBasicMagicBuildingValues() basicContentMagicBuildings {
	return basicContentMagicBuildings{
		BasicMagicBuildingsTier1: "basic_content_list_building_magic_tier_1",
		BasicMagicBuildingsTier2: "basic_content_list_building_magic_tier_2",
	}
}

type basicContentHeroExperienceBuildings struct {
	BasicHeroExperienceBuildingTier1 string
	BasicHeroExperienceBuildingTier2 string
}

func GetMandatoryContentBasicHeroExperienceBuildingValues() basicContentHeroExperienceBuildings {
	return basicContentHeroExperienceBuildings{
		BasicHeroExperienceBuildingTier1: "basic_content_list_building_hero_exp_tier_1",
		BasicHeroExperienceBuildingTier2: "basic_content_list_building_hero_exp_tier_2",
	}
}

type basicContentMiscellaneous struct {
	BasicStorage               string
	BasicNonContent            string
	BasicVisionBuildingsTier1  string
	BasicVisionBuildingsTier2  string
	BasicPickupRandomItems     string
	BasicPickupPrison          string
	BasicHeroBuffBuildingTier1 string
}

func GetMandatoryContentBasicMiscellaneousValues() basicContentMiscellaneous {
	return basicContentMiscellaneous{
		BasicStorage:               "basic_content_list_basic_storage",
		BasicNonContent:            "basic_content_list_non_content",
		BasicVisionBuildingsTier1:  "basic_content_list_vision_buildings_tier_1",
		BasicVisionBuildingsTier2:  "basic_content_list_vision_buildings_tier_2",
		BasicPickupRandomItems:     "basic_content_list_pickup_random_items",
		BasicPickupPrison:          "basic_content_list_pickup_prison",
		BasicHeroBuffBuildingTier1: "basic_content_list_building_hero_buff_tier_1",
	}
}
