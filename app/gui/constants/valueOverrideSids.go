package constants

// ValueOverrideSids are the known object / encounter SIDs offered by the
// guard-value-override picker. SIDs come from the registry; they reference
// world objects whose guard value can be overridden via valueOverrides.
var ValueOverrideSids = buildValueOverrideSids()

func buildValueOverrideSids() []string {
	sids := GetValueOverrideSidValues()
	return []string{
		sids.AlchemyLab,
		sids.Arena,
		sids.BeerFountain,
		sids.BorealCall,
		sids.CelestialSphere,
		sids.Chimerologist,
		sids.Circus,
		sids.CollegeOfWonder,
		sids.CrystalTrail,
		sids.DragonUtopia,
		sids.EternalDragon,
		sids.FickleShrine,
		sids.FlatteringMirror,
		sids.Forge,
		sids.Fort,
		sids.Fountain,
		sids.Fountain2,
		sids.HuntsmansCamp,
		sids.InfernalCirque,
		sids.InsarasEye,
		sids.JoustingRange,
		sids.ManaWell,
		sids.Market,
		sids.MineCrystals,
		sids.MineGemstones,
		sids.MineGold,
		sids.MineMercury,
		sids.MineOre,
		sids.MineWood,
		sids.Mirage,
		sids.MontyHall,
		sids.MysteriousStone,
		sids.MysticalTower,
		sids.MythicScrollBox,
		sids.OrbObservatory,
		sids.PandoraBox,
		sids.PetrifiedMemorial,
		sids.PileOfBooks,
		sids.PointOfBalance,
		sids.Prison,
		sids.QuixsPath,
		sids.RandomHire1,
		sids.RandomHire2,
		sids.RandomHire3,
		sids.RandomHire4,
		sids.RandomHire5,
		sids.RandomHire6,
		sids.RandomHire7,
		sids.RandomItemCommon,
		sids.RandomItemEpic,
		sids.RandomItemLegendary,
		sids.RandomItemRare,
		sids.RemoteFoothold,
		sids.ResearchLaboratory,
		sids.RitualPyre,
		sids.SacrificialShrine,
		sids.ShadyDen,
		sids.Stables,
		sids.Tavern,
		sids.TearOfTruth,
		sids.TheGorge,
		sids.TownGate,
		sids.TreeOfAbundance,
		sids.TroglodyteThrone,
		sids.UnforgottenGrave,
		sids.University,
		sids.UnstableRuins,
		sids.Watchtower,
		sids.WindRose,
		sids.WiseOwl,
	}
}

type valueOverrideSids struct {
	AlchemyLab          string
	Arena               string
	BeerFountain        string
	BorealCall          string
	CelestialSphere     string
	Chimerologist       string
	Circus              string
	CollegeOfWonder     string
	CrystalTrail        string
	DragonUtopia        string
	EternalDragon       string
	FickleShrine        string
	FlatteringMirror    string
	Forge               string
	Fort                string
	Fountain            string
	Fountain2           string
	HuntsmansCamp       string
	InfernalCirque      string
	InsarasEye          string
	JoustingRange       string
	ManaWell            string
	Market              string
	MineCrystals        string
	MineGemstones       string
	MineGold            string
	MineMercury         string
	MineOre             string
	MineWood            string
	Mirage              string
	MontyHall           string
	MysteriousStone     string
	MysticalTower       string
	MythicScrollBox     string
	OrbObservatory      string
	PandoraBox          string
	PetrifiedMemorial   string
	PileOfBooks         string
	PointOfBalance      string
	Prison              string
	QuixsPath           string
	RandomHire1         string
	RandomHire2         string
	RandomHire3         string
	RandomHire4         string
	RandomHire5         string
	RandomHire6         string
	RandomHire7         string
	RandomItemCommon    string
	RandomItemEpic      string
	RandomItemLegendary string
	RandomItemRare      string
	RemoteFoothold      string
	ResearchLaboratory  string
	RitualPyre          string
	SacrificialShrine   string
	ShadyDen            string
	Stables             string
	Tavern              string
	TearOfTruth         string
	TheGorge            string
	TownGate            string
	TreeOfAbundance     string
	TroglodyteThrone    string
	UnforgottenGrave    string
	University          string
	UnstableRuins       string
	Watchtower          string
	WindRose            string
	WiseOwl             string
}

var valueOverrideSidValues = valueOverrideSids{
	AlchemyLab:          "alchemy_lab",
	Arena:               "arena",
	BeerFountain:        "beer_fountain",
	BorealCall:          "boreal_call",
	CelestialSphere:     "celestial_sphere",
	Chimerologist:       "chimerologist",
	Circus:              "circus",
	CollegeOfWonder:     "college_of_wonder",
	CrystalTrail:        "crystal_trail",
	DragonUtopia:        "dragon_utopia",
	EternalDragon:       "eternal_dragon",
	FickleShrine:        "fickle_shrine",
	FlatteringMirror:    "flattering_mirror",
	Forge:               "forge",
	Fort:                "fort",
	Fountain:            "fountain",
	Fountain2:           "fountain_2",
	HuntsmansCamp:       "huntsmans_camp",
	InfernalCirque:      "infernal_cirque",
	InsarasEye:          "insaras_eye",
	JoustingRange:       "jousting_range",
	ManaWell:            "mana_well",
	Market:              "market",
	MineCrystals:        "mine_crystals",
	MineGemstones:       "mine_gemstones",
	MineGold:            "mine_gold",
	MineMercury:         "mine_mercury",
	MineOre:             "mine_ore",
	MineWood:            "mine_wood",
	Mirage:              "mirage",
	MontyHall:           "monty_hall",
	MysteriousStone:     "mysterious_stone",
	MysticalTower:       "mystical_tower",
	MythicScrollBox:     "mythic_scroll_box",
	OrbObservatory:      "orb_observatory",
	PandoraBox:          "pandora_box",
	PetrifiedMemorial:   "petrified_memorial",
	PileOfBooks:         "pile_of_books",
	PointOfBalance:      "point_of_balance",
	Prison:              "prison",
	QuixsPath:           "quixs_path",
	RandomHire1:         "random_hire_1",
	RandomHire2:         "random_hire_2",
	RandomHire3:         "random_hire_3",
	RandomHire4:         "random_hire_4",
	RandomHire5:         "random_hire_5",
	RandomHire6:         "random_hire_6",
	RandomHire7:         "random_hire_7",
	RandomItemCommon:    "random_item_common",
	RandomItemEpic:      "random_item_epic",
	RandomItemLegendary: "random_item_legendary",
	RandomItemRare:      "random_item_rare",
	RemoteFoothold:      "remote_foothold",
	ResearchLaboratory:  "research_laboratory",
	RitualPyre:          "ritual_pyre",
	SacrificialShrine:   "sacrificial_shrine",
	ShadyDen:            "shady_den",
	Stables:             "stables",
	Tavern:              "tavern",
	TearOfTruth:         "tear_of_truth",
	TheGorge:            "the_gorge",
	TownGate:            "town_gate",
	TreeOfAbundance:     "tree_of_abundance",
	TroglodyteThrone:    "troglodyte_throne",
	UnforgottenGrave:    "unforgotten_grave",
	University:          "university",
	UnstableRuins:       "unstable_ruins",
	Watchtower:          "watchtower",
	WindRose:            "wind_rose",
	WiseOwl:             "wise_owl",
}

// GetValueOverrideSidValues returns the object / encounter SIDs used for
//
//	valueOverrides.sid
func GetValueOverrideSidValues() valueOverrideSids {
	return valueOverrideSidValues
}
