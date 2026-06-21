package registry

type allContentLists struct {
	contentLists
	basicContentLists
}

var allContentIncludeListValues = allContentLists{
	contentLists:      contentIncludeListValues,
	basicContentLists: basicContentIncludeListValues,
}

// GetAllMandatoryContentIncludeListValues returns all of the content group names used for
//
//	mandatoryContent.content.includeLists
func GetAllMandatoryContentIncludeListValues() allContentLists {
	return allContentIncludeListValues
}

type contentLists struct {
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

var contentIncludeListValues = contentLists{
	contentRandomHiresBuildings:        contentRandomHiresBuildingValues,
	contentResourceBanksBuildings:      contentResourceBanksBuildingValues,
	contentResources:                   contentResourcesValues,
	contentMines:                       contentMinesValues,
	contentInteractables:               contentInteractablesValues,
	contentHeroStatsAndSkillsBuildings: contentHeroStatsAndSkillsBuildingValues,
	contentScrollBoxes:                 contentScrollBoxValues,
	contentEnchantedScrollBoxes:        contentEnchantedScrollBoxValues,
	contentMythicScrollBoxes:           contentMythicScrollBoxValues,
	contentPandoraBoxes:                contentPandoraBoxValues,
	contentGuardedUnitBanks:            contentGuardedUnitBanksValues,
	contentMagicBuildings:              contentMagicBuildingValues,
	contentMiscellaneous:               contentMiscellaneousValues,
}

// GetMandatoryContentIncludeListValues returns the "content_list_*" content group names used for
//
//	mandatoryContent.content.includeLists
func GetMandatoryContentIncludeListValues() contentLists {
	return contentIncludeListValues
}

type contentRandomHiresBuildings struct {
	RandomHires         string
	RandomHiresLowTier  string
	RandomHiresHighTier string
}

var contentRandomHiresBuildingValues = contentRandomHiresBuildings{
	RandomHires:         "content_list_building_random_hires",
	RandomHiresLowTier:  "content_list_building_random_hires_low_tier",
	RandomHiresHighTier: "content_list_building_random_hires_high_tier",
}

func GetMandatoryContentRandomHiresBuildingValues() contentRandomHiresBuildings {
	return contentRandomHiresBuildingValues
}

type contentResourceBanksBuildings struct {
	ResourceBanksCommon           string
	ResourceBanksUncommon         string
	ResourceBanksUncommonFix      string
	GuardedResourceBanksCommonFix string
	GuardedResourceBanksUncommon  string
	GuardedResourceBanksEpic      string
}

var contentResourceBanksBuildingValues = contentResourceBanksBuildings{
	ResourceBanksCommon:           "content_list_building_common_resource_banks",
	ResourceBanksUncommon:         "content_list_building_uncommon_resource_banks",
	ResourceBanksUncommonFix:      "content_list_building_uncommon_resource_banks_fix",
	GuardedResourceBanksCommonFix: "content_list_building_common_guarded_resource_bank_fix",
	GuardedResourceBanksUncommon:  "content_list_building_uncommon_guarded_resource_banks",
	GuardedResourceBanksEpic:      "content_list_building_epic_guarded_resource_banks",
}

func GetMandatoryContentResourceBanksBuildingValues() contentResourceBanksBuildings {
	return contentResourceBanksBuildingValues
}

type contentResources struct {
	ResourcesBasic   string
	ResourcesRare    string
	ResourcesSpecial string
}

var contentResourcesValues = contentResources{
	ResourcesBasic:   "content_list_basic_resources",
	ResourcesRare:    "content_list_rare_resources",
	ResourcesSpecial: "content_list_special_resources",
}

type contentMines struct {
	MinesBasic   string
	MinesRare    string
	MinesSpecial string
}

var contentMinesValues = contentMines{
	MinesBasic:   "content_list_basic_mines",
	MinesRare:    "content_list_rare_mines",
	MinesSpecial: "content_list_special_mines",
}

type contentInteractables struct {
	InteractablesCommon         string
	InteractablesCommonBalanced string
	InteractablesUncommon       string
}

var contentInteractablesValues = contentInteractables{
	InteractablesCommon:         "content_list_building_common_interact",
	InteractablesCommonBalanced: "content_list_building_common_interact_no_imbalance",
	InteractablesUncommon:       "content_list_building_uncommon_interact",
}

type contentHeroStatsAndSkillsBuildings struct {
	HeroStatsAndSkillsCommon   string
	HeroStatsAndSkillsUncommon string
}

var contentHeroStatsAndSkillsBuildingValues = contentHeroStatsAndSkillsBuildings{
	HeroStatsAndSkillsCommon:   "content_list_building_common_hero_stats",
	HeroStatsAndSkillsUncommon: "content_list_building_uncommon_hero_stats",
}

type contentScrollBoxes struct {
	PickupScrollBox      string
	PickupScrollBoxTier1 string
	PickupScrollBoxTier2 string
	PickupScrollBoxTier3 string
	PickupScrollBoxTier4 string
	PickupScrollBoxTier5 string
}

var contentScrollBoxValues = contentScrollBoxes{
	PickupScrollBox:      "content_list_pickup_scroll_box",
	PickupScrollBoxTier1: "content_list_pickup_scroll_box_tier_1",
	PickupScrollBoxTier2: "content_list_pickup_scroll_box_tier_2",
	PickupScrollBoxTier3: "content_list_pickup_scroll_box_tier_3",
	PickupScrollBoxTier4: "content_list_pickup_scroll_box_tier_4",
	PickupScrollBoxTier5: "content_list_pickup_scroll_box_tier_5",
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

var contentEnchantedScrollBoxValues = contentEnchantedScrollBoxes{
	PickupEnchantedScrollBox:              "content_list_pickup_enchanted_scroll_box",
	PickupEnchantedScrollBoxTier1:         "content_list_pickup_enchanted_scroll_box_tier_1",
	PickupEnchantedScrollBoxTier2:         "content_list_pickup_enchanted_scroll_box_tier_2",
	PickupEnchantedScrollBoxTier3:         "content_list_pickup_enchanted_scroll_box_tier_3",
	PickupEnchantedScrollBoxTier4:         "content_list_pickup_enchanted_scroll_box_tier_4",
	PickupEnchantedScrollBoxTier5:         "content_list_pickup_enchanted_scroll_box_tier_5",
	PickupEnchantedScrollBoxControlSpells: "content_list_pickup_enchanted_scroll_box_control_spells",
}

type contentMythicScrollBoxes struct {
	PickupMythicScrollBox string
}

var contentMythicScrollBoxValues = contentMythicScrollBoxes{
	PickupMythicScrollBox: "content_list_pickup_mythic_scroll_box",
}

type contentPandoraBoxes struct {
	PickupPandoraBox             string
	PickupPandoraBoxArmy         string
	PickupPandoraBoxArmyLowTier  string
	PickupPandoraBoxArmyHighTier string
	PickupPandoraBoxExperience   string
	PickupPandoraBoxGold         string
}

var contentPandoraBoxValues = contentPandoraBoxes{
	PickupPandoraBox:             "content_list_pickup_pandora_box",
	PickupPandoraBoxArmy:         "content_list_pickup_pandora_box_army",
	PickupPandoraBoxArmyLowTier:  "content_list_pickup_pandora_box_army_low_tier",
	PickupPandoraBoxArmyHighTier: "content_list_pickup_pandora_box_army_high_tier",
	PickupPandoraBoxExperience:   "content_list_pickup_pandora_box_exp",
	PickupPandoraBoxGold:         "content_list_pickup_pandora_box_gold",
}

func GetMandatoryContentPandoraBoxValues() contentPandoraBoxes {
	return contentPandoraBoxValues
}

type contentGuardedUnitBanks struct {
	GuardedUnitsBanksUncommon      string
	GuardedUnitsBanksUncommonEqual string
}

var contentGuardedUnitBanksValues = contentGuardedUnitBanks{
	GuardedUnitsBanksUncommon:      "content_list_building_uncommon_guarded_units_banks",
	GuardedUnitsBanksUncommonEqual: "content_list_building_uncommon_guarded_units_banks_equal",
}

type contentMagicBuildings struct {
	MagicBuildingsCommon   string
	MagicBuildingsUncommon string
}

var contentMagicBuildingValues = contentMagicBuildings{
	MagicBuildingsCommon:   "content_list_building_commons_magic",
	MagicBuildingsUncommon: "content_list_building_uncommons_magic",
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

var contentMiscellaneousValues = contentMiscellaneous{
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

func GetMandatoryContentMiscellaneousValues() contentMiscellaneous {
	return contentMiscellaneousValues
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

var basicContentIncludeListValues = basicContentLists{
	basicContentRandomHiresBuildings:        basicContentRandomHiresBuildingValues,
	basicContentResourceBanksBuildings:      basicContentResourceBanksBuildingValues,
	basicContentResources:                   basicContentResourcesValues,
	basicContentMines:                       basicContentMinesValues,
	basicContentInteractables:               basicContentInteractablesValues,
	basicContentHeroStatsAndSkillsBuildings: basicContentHeroStatsAndSkillsBuildingValues,
	basicContentScrollBoxes:                 basicContentScrollBoxValues,
	basicContentPandoraBoxes:                basicContentPandoraBoxValues,
	basicContentGuardedUnitBanks:            basicContentGuardedUnitBanksValues,
	basicContentMagicBuildings:              basicContentMagicBuildingValues,
	basicContentHeroExperienceBuildings:     basicContentHeroExperienceBuildingValues,
	basicContentMiscellaneous:               basicContentMiscellaneousValues,
}

// GetMandatoryContentBasicIncludeListValues returns the "basic_content_list_*" content group names used for
//
//	mandatoryContent.content.includeLists
func GetMandatoryContentBasicIncludeListValues() basicContentLists {
	return basicContentIncludeListValues
}

type basicContentRandomHiresBuildings struct {
	BasicRandomHires string
}

var basicContentRandomHiresBuildingValues = basicContentRandomHiresBuildings{
	BasicRandomHires: "basic_content_list_building_random_hires",
}

func GetMandatoryContentBasicRandomHiresBuildingValues() basicContentRandomHiresBuildings {
	return basicContentRandomHiresBuildingValues
}

type basicContentResourceBanksBuildings struct {
	BasicResourceBanksTier1                     string
	BasicResourceBanksTier2                     string
	BasicGuardedResourceBanksTier1              string
	BasicGuardedResourceBanksTier2              string
	BasicGuardedResourceBanksTier2NoRestriction string
	BasicGuardedResourceBanksTier3              string
}

var basicContentResourceBanksBuildingValues = basicContentResourceBanksBuildings{
	BasicResourceBanksTier1:                     "basic_content_list_building_resource_banks_tier_1",
	BasicResourceBanksTier2:                     "basic_content_list_building_resource_banks_tier_2",
	BasicGuardedResourceBanksTier1:              "basic_content_list_building_guarded_resource_banks_tier_1",
	BasicGuardedResourceBanksTier2:              "basic_content_list_building_guarded_resource_banks_tier_2",
	BasicGuardedResourceBanksTier2NoRestriction: "basic_content_list_building_guarded_resource_banks_tier_2_no_biome_restriction",
	BasicGuardedResourceBanksTier3:              "basic_content_list_building_guarded_resource_banks_tier_3",
}

func GetMandatoryContentBasicResourceBanksBuildingValues() basicContentResourceBanksBuildings {
	return basicContentResourceBanksBuildingValues
}

type basicContentResources struct {
	BasicResourcesBasic   string
	BasicResourcesRare    string
	BasicResourcesSpecial string
}

var basicContentResourcesValues = basicContentResources{
	BasicResourcesBasic:   "basic_content_list_basic_resources",
	BasicResourcesRare:    "basic_content_list_rare_resources",
	BasicResourcesSpecial: "basic_content_list_special_resources",
}

type basicContentMines struct {
	BasicMinesBasic               string
	BasicMinesRare                string
	BasicMinesRareBiomeRestricted string
	BasicMinesSpecial             string
}

var basicContentMinesValues = basicContentMines{
	BasicMinesBasic:               "basic_content_list_basic_mines",
	BasicMinesRare:                "basic_content_list_rare_mines",
	BasicMinesRareBiomeRestricted: "basic_content_list_rare_mines_by_biome",
	BasicMinesSpecial:             "basic_content_list_special_mines",
}

func GetMandatoryContentBasicMinesValues() basicContentMines {
	return basicContentMinesValues
}

type basicContentInteractables struct {
	BasicInteractablesUncommon string
	BasicInteractablesEpic     string
}

var basicContentInteractablesValues = basicContentInteractables{
	BasicInteractablesUncommon: "basic_content_list_building_uncommon_interact",
	BasicInteractablesEpic:     "basic_content_list_building_epic_interact",
}

type basicContentHeroStatsAndSkillsBuildings struct {
	BasicHeroStatsAndSkillsTier1 string
	BasicHeroStatsAndSkillsTier2 string
	BasicHeroStatsAndSkillsTier3 string
}

var basicContentHeroStatsAndSkillsBuildingValues = basicContentHeroStatsAndSkillsBuildings{
	BasicHeroStatsAndSkillsTier1: "basic_content_list_building_hero_stats_and_skills_tier_1",
	BasicHeroStatsAndSkillsTier2: "basic_content_list_building_hero_stats_and_skills_tier_2",
	BasicHeroStatsAndSkillsTier3: "basic_content_list_building_hero_stats_and_skills_tier_3",
}

func GetMandatoryContentBasicHeroStatsAndSkillsBuildingValues() basicContentHeroStatsAndSkillsBuildings {
	return basicContentHeroStatsAndSkillsBuildingValues
}

type basicContentScrollBoxes struct {
	BasicPickupScrollBox          string
	BasicPickupEnchantedScrollBox string
	BasicPickupMythicScrollBox    string
}

var basicContentScrollBoxValues = basicContentScrollBoxes{
	BasicPickupScrollBox:          "basic_content_list_pickup_scroll_box",
	BasicPickupEnchantedScrollBox: "basic_content_list_pickup_enchanted_scroll_box",
	BasicPickupMythicScrollBox:    "basic_content_list_pickup_mythic_scroll_box",
}

func GetMandatoryContentBasicScrollBoxValues() basicContentScrollBoxes {
	return basicContentScrollBoxValues
}

type basicContentPandoraBoxes struct {
	BasicPickupPandoraBox string
}

var basicContentPandoraBoxValues = basicContentPandoraBoxes{
	BasicPickupPandoraBox: "basic_content_list_pickup_pandora_box",
}

type basicContentGuardedUnitBanks struct {
	BasicGuardedUnitBanks                string
	BasicGuardedUnitBanksNotRestricted   string
	BasicGuardedUnitBanksBiomeRestricted string
}

var basicContentGuardedUnitBanksValues = basicContentGuardedUnitBanks{
	BasicGuardedUnitBanks:                "basic_content_list_building_guarded_units_banks",
	BasicGuardedUnitBanksNotRestricted:   "basic_content_list_building_guarded_units_banks_no_biome_restriction",
	BasicGuardedUnitBanksBiomeRestricted: "basic_content_list_building_guarded_units_banks_only_biome_restriction",
}

func GetMandatoryContentBasicGuardedUnitBankValues() basicContentGuardedUnitBanks {
	return basicContentGuardedUnitBanksValues
}

type basicContentMagicBuildings struct {
	BasicMagicBuildingsTier1 string
	BasicMagicBuildingsTier2 string
}

var basicContentMagicBuildingValues = basicContentMagicBuildings{
	BasicMagicBuildingsTier1: "basic_content_list_building_magic_tier_1",
	BasicMagicBuildingsTier2: "basic_content_list_building_magic_tier_2",
}

func GetMandatoryContentBasicMagicBuildingValues() basicContentMagicBuildings {
	return basicContentMagicBuildingValues
}

type basicContentHeroExperienceBuildings struct {
	BasicHeroExperienceBuildingTier1 string
	BasicHeroExperienceBuildingTier2 string
}

var basicContentHeroExperienceBuildingValues = basicContentHeroExperienceBuildings{
	BasicHeroExperienceBuildingTier1: "basic_content_list_building_hero_exp_tier_1",
	BasicHeroExperienceBuildingTier2: "basic_content_list_building_hero_exp_tier_2",
}

func GetMandatoryContentBasicHeroExperienceBuildingValues() basicContentHeroExperienceBuildings {
	return basicContentHeroExperienceBuildingValues
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

var basicContentMiscellaneousValues = basicContentMiscellaneous{
	BasicStorage:               "basic_content_list_basic_storage",
	BasicNonContent:            "basic_content_list_non_content",
	BasicVisionBuildingsTier1:  "basic_content_list_vision_buildings_tier_1",
	BasicVisionBuildingsTier2:  "basic_content_list_vision_buildings_tier_2",
	BasicPickupRandomItems:     "basic_content_list_pickup_random_items",
	BasicPickupPrison:          "basic_content_list_pickup_prison",
	BasicHeroBuffBuildingTier1: "basic_content_list_building_hero_buff_tier_1",
}

func GetMandatoryContentBasicMiscellaneousValues() basicContentMiscellaneous {
	return basicContentMiscellaneousValues
}
