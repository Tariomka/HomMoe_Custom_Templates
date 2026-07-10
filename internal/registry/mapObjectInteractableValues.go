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
	interactableNamedUnitBanks
	interactableRandomUnitBanks
}

func GetMapObjectAllInteractableValues() objectInteractables {
	return objectInteractables{
		interactableMines:                     GetMapObjectMineValues(),
		interactableStorages:                  GetMapObjectStorageValues(),
		interactableResourceBanks:             GetMapObjectResourceBankValues(),
		interactableT1GuardedResourceBanks:    GetMapObjectT1GuardedResourceBankValues(),
		interactableT2GuardedResourceBanks:    GetMapObjectT2GuardedResourceBankValues(),
		interactableT3GuardedResourceBanks:    GetMapObjectT3GuardedResourceBankValues(),
		interactableNonContents:               GetMapObjectNonContentValues(),
		interactableVisionBuildings:           GetMapObjectVisionBuildingValues(),
		interactableMiscellaneous:             GetMapObjectMiscellaneousValues(),
		interactableHeroExperienceBuildings:   GetMapObjectHeroExperienceBuildingValues(),
		interactableMagicBuildings:            GetMapObjectMagicBuildingValues(),
		interactableHeroBuffBuildings:         GetMapObjectHeroBuffBuildingValues(),
		interactableBuildings:                 GetMapObjectBuildingValues(),
		interactableT1StatsAndSkillsBuildings: GetMapObjectT1StatsAndSkillsBuildingValues(),
		interactableT2StatsAndSkillsBuildings: GetMapObjectT2StatsAndSkillsBuildingValues(),
		interactableT3StatsAndSkillsBuildings: GetMapObjectT3StatsAndSkillsBuildingValues(),
		interactableNamedUnitBanks:            GetMapObjectNamedUnitBankValues(),
		interactableRandomUnitBanks:           GetMapObjectRandomUnitBankValues(),
	}
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

func GetMapObjectMineValues() interactableMines {
	return interactableMines{
		AlchemyLab:   "alchemy_lab",
		GoldMine:     "mine_gold",
		WoodMine:     "mine_wood",
		OreMine:      "mine_ore",
		CrystalMine:  "mine_crystals",
		MercuryMine:  "mine_mercury",
		GemstoneMine: "mine_gemstones",
	}
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

func GetMapObjectStorageValues() interactableStorages {
	return interactableStorages{
		WoodStorage:     "storage_wood",
		OreStorage:      "storage_ore",
		GoldStorage:     "storage_gold",
		CrystalStorage:  "storage_crystals",
		MercuryStorage:  "storage_mercury",
		GemstoneStorage: "storage_gemstones",
		DustStorage:     "storage_dust",
	}
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

func GetMapObjectResourceBankValues() interactableResourceBanks {
	return interactableResourceBanks{
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
}

type interactableT1GuardedResourceBanks struct {
	BlackTower       string
	AbandonedMansion string
	MereasShrine     string
	ShadyDen         string
}

func GetMapObjectT1GuardedResourceBankValues() interactableT1GuardedResourceBanks {
	return interactableT1GuardedResourceBanks{
		BlackTower:       "black_tower",
		AbandonedMansion: "abandoned_mansion",
		MereasShrine:     "mereas_shrine",
		ShadyDen:         "shady_den",
	}
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

func GetMapObjectT2GuardedResourceBankValues() interactableT2GuardedResourceBanks {
	return interactableT2GuardedResourceBanks{
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
}

type interactableT3GuardedResourceBanks struct {
	TroglodyteThrone   string
	TwilightBloom      string
	UnstableRuins      string
	DragonUtopia       string
	ResearchLaboratory string
}

func GetMapObjectT3GuardedResourceBankValues() interactableT3GuardedResourceBanks {
	return interactableT3GuardedResourceBanks{
		TroglodyteThrone:   "troglodyte_throne",
		TwilightBloom:      "twilight_bloom",
		UnstableRuins:      "unstable_ruins",
		DragonUtopia:       "dragon_utopia",
		ResearchLaboratory: "research_laboratory",
	}
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

func GetMapObjectNonContentValues() interactableNonContents {
	return interactableNonContents{
		Mirage:           "mirage",
		InsarasEye:       "insaras_eye",
		RemoteFoothold:   "remote_foothold",
		AbandonedOutpost: "abandoned_outpost",
		Market:           "market",
		Forge:            "forge",
		Tavern:           "tavern",
	}
}

type interactableVisionBuildings struct {
	FlatteringMirror string
	Watchtower       string
	WindRose         string
}

func GetMapObjectVisionBuildingValues() interactableVisionBuildings {
	return interactableVisionBuildings{
		FlatteringMirror: "flattering_mirror",
		Watchtower:       "watchtower",
		WindRose:         "wind_rose",
	}
}

type interactableMiscellaneous struct {
	Prison   string
	TownGate string
}

func GetMapObjectMiscellaneousValues() interactableMiscellaneous {
	return interactableMiscellaneous{
		Prison:   "prison",
		TownGate: "town_gate",
	}
}

type interactableHeroExperienceBuildings struct {
	LearningStone   string
	LostLibrary     string
	TreeOfKnowledge string
}

func GetMapObjectHeroExperienceBuildingValues() interactableHeroExperienceBuildings {
	return interactableHeroExperienceBuildings{
		LearningStone:   "learning_stone",
		LostLibrary:     "lost_library",
		TreeOfKnowledge: "tree_of_knowledge",
	}
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

func GetMapObjectMagicBuildingValues() interactableMagicBuildings {
	return interactableMagicBuildings{
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

func GetMapObjectHeroBuffBuildingValues() interactableHeroBuffBuildings {
	return interactableHeroBuffBuildings{
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
}

type interactableBuildings struct {
	HuntsmansCamp     string
	MercenaryGuild    string
	SacrificialShrine string
	Chimerologist     string
	Arena             string
	EternalDragon     string
	FickleShrine      string
	TreeOfAbundance   string
}

func GetMapObjectBuildingValues() interactableBuildings {
	return interactableBuildings{
		HuntsmansCamp: "huntsmans_camp",

		MercenaryGuild:    "mercenary_guild",
		SacrificialShrine: "sacrificial_shrine",
		Chimerologist:     "chimerologist",
		Arena:             "arena",

		EternalDragon:   "eternal_dragon",
		FickleShrine:    "fickle_shrine",
		TreeOfAbundance: "tree_of_abundance",
	}
}

type interactableT1StatsAndSkillsBuildings struct {
	StingingSword   string
	ArmoryAutomaton string
	MagicWheel      string
	KnowledgeGarden string
	WiseOwl         string
}

func GetMapObjectT1StatsAndSkillsBuildingValues() interactableT1StatsAndSkillsBuildings {
	return interactableT1StatsAndSkillsBuildings{
		StingingSword:   "stinging_sword",
		ArmoryAutomaton: "armory_automaton",
		MagicWheel:      "magic_wheel",
		KnowledgeGarden: "knowledge_garden",
		WiseOwl:         "wise_owl",
	}
}

type interactableT2StatsAndSkillsBuildings struct {
	Fort           string
	OrbObservatory string
	University     string
	Circus         string
	InfernalCirque string
}

func GetMapObjectT2StatsAndSkillsBuildingValues() interactableT2StatsAndSkillsBuildings {
	return interactableT2StatsAndSkillsBuildings{
		Fort:           "fort",
		OrbObservatory: "orb_observatory",
		University:     "university",
		Circus:         "circus",
		InfernalCirque: "infernal_cirque",
	}
}

type interactableT3StatsAndSkillsBuildings struct {
	Maze            string
	TrialScales     string
	CollegeOfWonder string
}

func GetMapObjectT3StatsAndSkillsBuildingValues() interactableT3StatsAndSkillsBuildings {
	return interactableT3StatsAndSkillsBuildings{
		Maze:            "maze",
		TrialScales:     "trial_scales",
		CollegeOfWonder: "college_of_wonder",
	}
}

type interactableNamedUnitBanks struct {
	JoustingRange     string
	UnforgottenGrave  string
	PetrifiedMemorial string
	RitualPyre        string
	BorealCall        string
	Gorge             string
	PointOfBalance    string
}

func GetMapObjectNamedUnitBankValues() interactableNamedUnitBanks {
	return interactableNamedUnitBanks{
		JoustingRange:     "jousting_range",
		UnforgottenGrave:  "unforgotten_grave",
		PetrifiedMemorial: "petrified_memorial",
		RitualPyre:        "ritual_pyre",
		BorealCall:        "boreal_call",
		Gorge:             "the_gorge",
		PointOfBalance:    "point_of_balance",
	}
}

type interactableRandomUnitBanks struct {
	RandomHireTier1 string
	RandomHireTier2 string
	RandomHireTier3 string
	RandomHireTier4 string
	RandomHireTier5 string
	RandomHireTier6 string
	RandomHireTier7 string
}

func GetMapObjectRandomUnitBankValues() interactableRandomUnitBanks {
	return interactableRandomUnitBanks{
		RandomHireTier1: "random_hire_1",
		RandomHireTier2: "random_hire_2",
		RandomHireTier3: "random_hire_3",
		RandomHireTier4: "random_hire_4",
		RandomHireTier5: "random_hire_5",
		RandomHireTier6: "random_hire_6",
		RandomHireTier7: "random_hire_7",
	}
}

// Below are SIDs currently unused
// TODO: make them also registry entries
// also need correct names as the parameters for each of these building ids

// "id": "underground_lair",
// "id": "vanguard",
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
