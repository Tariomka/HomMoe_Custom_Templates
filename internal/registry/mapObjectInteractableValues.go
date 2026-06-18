package registry

type objectInteractables struct {
	interactableMines
	interactableStorages
	interactableResourceBanks
	interactableT1GuardedResourceBanks
	interactableT2GuardedResourceBanks
	interactableT3GuardedResourceBanks
	interactableNonContents
	interactableVisionBuildings
	interactableMiscellaneous
	interactableHeroExperienceBuildings
	interactableMagicBuildings
	interactableHeroBuffBuildings
	interactableBuildings
	interactableT1StatsAndSkillsBuildings
	interactableT2StatsAndSkillsBuildings
	interactableT3StatsAndSkillsBuildings
}

var objectInteractableValues = objectInteractables{
	interactableMines:                     interactableMineValues,
	interactableStorages:                  interactableStorageValues,
	interactableResourceBanks:             interactableResourceBankValues,
	interactableT1GuardedResourceBanks:    interactableT1GuardedResourceBankValues,
	interactableT2GuardedResourceBanks:    interactableT2GuardedResourceBankValues,
	interactableT3GuardedResourceBanks:    interactableT3GuardedResourceBankValues,
	interactableNonContents:               interactableNonContentValues,
	interactableVisionBuildings:           interactableVisionBuildingValues,
	interactableMiscellaneous:             interactableMiscellaneousValues,
	interactableHeroExperienceBuildings:   interactableHeroExpValues,
	interactableMagicBuildings:            interactableMagicBuildingValues,
	interactableHeroBuffBuildings:         interactableHeroBuffBuildingValues,
	interactableBuildings:                 interactableBuildingValues,
	interactableT1StatsAndSkillsBuildings: interactableT1StatsAndSkillsValues,
	interactableT2StatsAndSkillsBuildings: interactableT2StatsAndSkillsValues,
	interactableT3StatsAndSkillsBuildings: interactableT3StatsAndSkillsValues,
}

func GetMapObjectAllInteractableValues() objectInteractables {
	return objectInteractableValues
}

type interactableMines struct {
	AlchemyLab   string
	GoldMine     string
	WoodMine     string
	OreMine      string
	CrystalMine  string
	MercuryMine  string
	GemstoneMine string
}

var interactableMineValues = interactableMines{
	AlchemyLab:   "alchemy_lab",
	GoldMine:     "mine_gold",
	WoodMine:     "mine_wood",
	OreMine:      "mine_ore",
	CrystalMine:  "mine_crystals",
	MercuryMine:  "mine_mercury",
	GemstoneMine: "mine_gemstones",
}

type interactableStorages struct {
	WoodStorage     string
	OreStorage      string
	GoldStorage     string
	CrystalStorage  string
	MercuryStorage  string
	GemstoneStorage string
	DustStorage     string
}

var interactableStorageValues = interactableStorages{
	WoodStorage:     "storage_wood",
	OreStorage:      "storage_ore",
	GoldStorage:     "storage_gold",
	CrystalStorage:  "storage_crystals",
	MercuryStorage:  "storage_mercury",
	GemstoneStorage: "storage_gemstones",
	DustStorage:     "storage_dust",
}

type interactableResourceBanks struct {
	Gardener         string
	Windmill         string
	Village          string
	GingerbreadHouse string
	PeasantCart      string
	AbandonedCorpse  string
	CrowNest         string
	GoblinCache      string
	MontyHall        string
	HerosCrypt       string
}

var interactableResourceBankValues = interactableResourceBanks{
	Gardener:         "gardener",
	Windmill:         "windmill",
	Village:          "village",
	GingerbreadHouse: "gingerbread_house",
	PeasantCart:      "peasant_cart",
	AbandonedCorpse:  "abandoned_corpse",
	CrowNest:         "crow_nest",
	GoblinCache:      "goblin_cache",

	MontyHall:  "monty_hall",
	HerosCrypt: "heros_crypt",
}

type interactableT1GuardedResourceBanks struct {
	BlackTower       string
	AbandonedMansion string
	MereasShrine     string
	ShadyDen         string
}

var interactableT1GuardedResourceBankValues = interactableT1GuardedResourceBanks{
	BlackTower:       "black_tower",
	AbandonedMansion: "abandoned_mansion",
	MereasShrine:     "mereas_shrine",
	ShadyDen:         "shady_den",
}

type interactableT2GuardedResourceBanks struct {
	RaidersCamp       string
	OvergrownGrave    string
	LegionsMemorial   string
	AlvarsEye         string
	CursedOldHouse    string
	AbnormalStructure string
	PrismaticLair     string
	UncannyRite       string
	CircleOfLife      string
	IridescentAbbey   string
}

var interactableT2GuardedResourceBankValues = interactableT2GuardedResourceBanks{
	RaidersCamp:       "raiders_camp",
	OvergrownGrave:    "overgrown_grave",
	LegionsMemorial:   "legions_memorial",
	AlvarsEye:         "alvars_eye",
	CursedOldHouse:    "cursed_old_house",
	AbnormalStructure: "abnormal_structure",
	PrismaticLair:     "prismatic_lair",
	UncannyRite:       "uncanny_rite",
	CircleOfLife:      "circle_of_life",
	IridescentAbbey:   "iridescent_abbey",
}

type interactableT3GuardedResourceBanks struct {
	TroglodyteThrone   string
	TwilightBloom      string
	UnstableRuins      string
	DragonUtopia       string
	ResearchLaboratory string
}

var interactableT3GuardedResourceBankValues = interactableT3GuardedResourceBanks{
	TroglodyteThrone:   "troglodyte_throne",
	TwilightBloom:      "twilight_bloom",
	UnstableRuins:      "unstable_ruins",
	DragonUtopia:       "dragon_utopia",
	ResearchLaboratory: "research_laboratory",
}

type interactableNonContents struct {
	Mirage           string
	InsarasEye       string
	RemoteFoothold   string
	AbandonedOutpost string
	Market           string
	Forge            string
	Tavern           string
}

var interactableNonContentValues = interactableNonContents{
	Mirage:           "mirage",
	InsarasEye:       "insaras_eye",
	RemoteFoothold:   "remote_foothold",
	AbandonedOutpost: "abandoned_outpost",
	Market:           "market",
	Forge:            "forge",
	Tavern:           "tavern",
}

type interactableVisionBuildings struct {
	FlatteringMirror string
	Watchtower       string
	WindRose         string
}

var interactableVisionBuildingValues = interactableVisionBuildings{
	FlatteringMirror: "flattering_mirror",
	Watchtower:       "watchtower",
	WindRose:         "wind_rose",
}

type interactableMiscellaneous struct {
	Prison   string
	TownGate string
}

var interactableMiscellaneousValues = interactableMiscellaneous{
	Prison:   "prison",
	TownGate: "town_gate",
}

type interactableHeroExperienceBuildings struct {
	LearningStone   string
	LostLibrary     string
	TreeOfKnowledge string
}

var interactableHeroExpValues = interactableHeroExperienceBuildings{
	LearningStone:   "learning_stone",
	LostLibrary:     "lost_library",
	TreeOfKnowledge: "tree_of_knowledge",
}

type interactableMagicBuildings struct {
	MysticalTower   string
	CelestialSphere string
	AltarOfMagic1   string
	AltarOfMagic2   string
	AltarOfMagic3   string
	AltarOfMagic4   string
	MagicAmplifier1 string
	MagicAmplifier2 string
	MagicAmplifier3 string
	MagicAmplifier4 string
}

var interactableMagicBuildingValues = interactableMagicBuildings{
	MysticalTower:   "mystical_tower",
	CelestialSphere: "celestial_sphere",

	AltarOfMagic1: "altar_of_magic_1",
	AltarOfMagic2: "altar_of_magic_2",
	AltarOfMagic3: "altar_of_magic_3",
	AltarOfMagic4: "altar_of_magic_4",

	MagicAmplifier1: "magic_amplifier_1",
	MagicAmplifier2: "magic_amplifier_2",
	MagicAmplifier3: "magic_amplifier_3",
	MagicAmplifier4: "magic_amplifier_4",
}

type interactableHeroBuffBuildings struct {
	ManaWell        string
	Fountain        string
	Fountain2       string
	Stables         string
	TearOfTruth     string
	BeerFountain    string
	QuixsPath       string
	PileOfBooks     string
	MysteriousStone string
	CrystalTrail    string
}

var interactableHeroBuffBuildingValues = interactableHeroBuffBuildings{
	ManaWell:        "mana_well",
	Fountain:        "fountain",
	Fountain2:       "fountain_2",
	Stables:         "stables",
	TearOfTruth:     "tear_of_truth",
	BeerFountain:    "beer_fountain",
	QuixsPath:       "quixs_path",
	PileOfBooks:     "pile_of_books",
	MysteriousStone: "mysterious_stone",
	CrystalTrail:    "crystal_trail",
}

type interactableBuildings struct {
	HuntsmanCamp      string
	MercenaryGuild    string
	SacrificialShrine string
	Chimerologist     string
	Arena             string
	EternalDragon     string
	FickleShrine      string
	TreeOfAbundance   string
}

var interactableBuildingValues = interactableBuildings{
	HuntsmanCamp: "huntsmans_camp",

	MercenaryGuild:    "mercenary_guild",
	SacrificialShrine: "sacrificial_shrine",
	Chimerologist:     "chimerologist",
	Arena:             "arena",

	EternalDragon:   "eternal_dragon",
	FickleShrine:    "fickle_shrine",
	TreeOfAbundance: "tree_of_abundance",
}

type interactableT1StatsAndSkillsBuildings struct {
	StingingSword   string
	ArmoryAutomaton string
	MagicWheel      string
	KnowledgeGarden string
	WiseOwl         string
}

var interactableT1StatsAndSkillsValues = interactableT1StatsAndSkillsBuildings{
	StingingSword:   "stinging_sword",
	ArmoryAutomaton: "armory_automaton",
	MagicWheel:      "magic_wheel",
	KnowledgeGarden: "knowledge_garden",
	WiseOwl:         "wise_owl",
}

type interactableT2StatsAndSkillsBuildings struct {
	Fort           string
	OrbObservatory string
	University     string
	Circus         string
	InfernalCirque string
}

var interactableT2StatsAndSkillsValues = interactableT2StatsAndSkillsBuildings{
	Fort:           "fort",
	OrbObservatory: "orb_observatory",
	University:     "university",
	Circus:         "circus",
	InfernalCirque: "infernal_cirque",
}

type interactableT3StatsAndSkillsBuildings struct {
	Maze            string
	TrialScales     string
	CollegeOfWonder string
}

var interactableT3StatsAndSkillsValues = interactableT3StatsAndSkillsBuildings{
	Maze:            "maze",
	TrialScales:     "trial_scales",
	CollegeOfWonder: "college_of_wonder",
}

// "id": "temporary_camp",

// "id": "portal_magic",

// "id": "testing_grounds",
// "id": "custom_testing_grounds",

// "id": "underground_lair",
// "id": "custom_underground_lair",
// "id": "point_of_balance",
// "id": "custom_point_of_balance",
// "id": "jousting_range",
// "id": "custom_jousting_range",
// "id": "vanguard",
// "id": "custom_vanguard",
// "id": "ritual_pyre",
// "id": "custom_ritual_pyre",
// "id": "unforgotten_grave",
// "id": "custom_unforgotten_grave",
// "id": "boreal_call",
// "id": "custom_boreal_call",
// "id": "the_gorge",
// "id": "petrified_memorial",
// "id": "custom_petrified_memorial",
// "id": "pocket_dimension",
// "id": "human_city",
// "id": "undead_city",
// "id": "dungeon_city",
// "id": "nature_city",
// "id": "unfrozen_city",
// "id": "demon_city",
// "id": "barracks_neutral_dragon_lich",
// "id": "barracks_human_1",
// "id": "barracks_human_2",
// "id": "barracks_human_3",
// "id": "barracks_human_4",
// "id": "barracks_human_5",
// "id": "barracks_human_6",
// "id": "barracks_human_7",
// "id": "barracks_necropolis_1",
// "id": "barracks_necropolis_2",
// "id": "barracks_necropolis_3",
// "id": "barracks_necropolis_4",
// "id": "barracks_necropolis_5",
// "id": "barracks_necropolis_6",
// "id": "barracks_necropolis_7",
// "id": "barracks_nature_1",
// "id": "barracks_nature_2",
// "id": "barracks_nature_3",
// "id": "barracks_nature_4",
// "id": "barracks_nature_5",
// "id": "barracks_nature_6",
// "id": "barracks_nature_7",
// "id": "barracks_dungeon_1",
// "id": "barracks_dungeon_2",
// "id": "barracks_dungeon_3",
// "id": "barracks_dungeon_4",
// "id": "barracks_dungeon_5",
// "id": "barracks_dungeon_6",
// "id": "barracks_dungeon_7",
// "id": "barracks_unfrozen_1",
// "id": "barracks_unfrozen_2",
// "id": "barracks_unfrozen_3",
// "id": "barracks_unfrozen_4",
// "id": "barracks_unfrozen_5",
// "id": "barracks_unfrozen_6",
// "id": "barracks_unfrozen_7",
// "id": "barracks_demon_1",
// "id": "barracks_demon_2",
// "id": "barracks_demon_3",
// "id": "barracks_demon_4",
// "id": "barracks_demon_5",
// "id": "barracks_demon_6",
// "id": "barracks_demon_7",
// "id": "pvp_promo_barracks",
// "id": "pvp_promo_barracks_necropolis",
// "id": "barracks_neutral_1",
// "id": "barracks_neutral_2",
// "id": "barracks_neutral_3",
// "id": "barracks_neutral_4",
// "id": "barracks_neutral_5",
// "id": "barracks_neutral_6",
// "id": "barracks_neutral_7",
// "id": "barracks_neutral_8",
// "id": "barracks_neutral_9",
// "id": "barracks_neutral_10",
// "id": "barracks_neutral_11",
// "id": "barracks_neutral_12",
// "id": "barracks_neutral_13",
// "id": "barracks_neutral_14",
// "id": "barracks_neutral_15",
// "id": "barracks_neutral_16",
// "id": "barracks_neutral_17",
// "id": "unit_trade_lab_kitten_horn",
// "id": "unit_trade_lab_gnat",
// "id": "block",
// "id": "shroom_of_growth",
// "id": "fairy_ring",
// "id": "learning_stone_old",
// "id": "gladiator_arena",

//
//
//
// "id": "stinging_sword_campaign",
// "id": "armory_automaton_campaign",
// "id": "magic_wheel_campaign",
// "id": "knowledge_garden_campaign",
// "id": "maze_campaign",
// "id": "fort_campaign",
// "id": "orb_observatory_campaign",
// "id": "college_of_wonder_campaign",
// "id": "goblin_cache_campaign",
// "id": "peasant_cart_campaign",
// "id": "crow_nest_campaign",
// "id": "shady_den_campaign",

// "id": "block_campaign_tree_grass",
// "id": "block_campaign_tree_autumn",
// "id": "block_campaign_tree_dirt",
// "id": "block_campaign_tree_dead",
// "id": "block_campaign_tree_withered",
// "id": "block_campaign_tree_snow",
// "id": "block_campaign_mountain_dirt",
// "id": "block_campaign",

// "id": "campaign_stinging_sword",
// "id": "campaign_armory_automaton",
// "id": "campaign_magic_wheel",
// "id": "campaign_knowledge_garden",
// "id": "campaign_college_of_wonder",
// "id": "campaign_fort",
// "id": "campaign_orb_observatory",
// "id": "campaign_maze",
// "id": "campaign_trial_scales",
// "id": "campaign_shady_den",
// "id": "campaign_flattering_mirror",
// "id": "campaign_gingerbread_house",
// "id": "campaign_lost_library_empty",
// "id": "campaign_M1_broken_bridge",
// "id": "campaign_M2_shroom_tether",
// "id": "campaign_M2_rockfall_1",
// "id": "campaign_M2_rockfall_2",
// "id": "campaign_M2_empty_mine",
// "id": "campaign_M3_necro_altar",
// "id": "campaign_M3_cage",
// "id": "campaign_M4_burning_man",
// "id": "campaign_M4_construction_site",
// "id": "campaign_M4_diary",
// "id": "campaign_M4_stargazer_tower",
// "id": "campaign_M5_death_knight_house_empty",
// "id": "campaign_M5_death_knight_house",
// "id": "campaign_M6_mearea_shrine",
// "id": "campaign_M6_cage",
// "id": "campaign_M9_scientist_house",
// "id": "campaign_M9_necromancy_amplifier",
// "id": "campaign_M9_necromancy_amplifier_foundation",
// "id": "campaign_M9_necromancy_amplifier_destroyed",
// "id": "campaign_M9_block_angel1",
// "id": "campaign_M9_block_angel2",
// "id": "campaign_M9_block_angel3",
// "id": "campaign_M9_block_angel4",
// "id": "campaign_M9_sylvan_altar",
// "id": "campaign_M9_phoenix_nest_shroom",
// "id": "campaign_M9_phoenix_nest",
// "id": "campaign_M9_phoenix_egg",
// "id": "campaign_M9_tomb",
