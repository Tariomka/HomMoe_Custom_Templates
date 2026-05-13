package constants

import "github.com/Tariomka/hommoe_custom_templates/internal/models"

// ContentIds enumerates every world-object SID known to the editor.
//
//nolint:gochecknoglobals // semantic registry
var ContentIds = struct {
	AlchemyLab          models.SidMapping
	Arena               models.SidMapping
	BeerFountain        models.SidMapping
	BorealCall          models.SidMapping
	CelestialSphere     models.SidMapping
	Chimerologist       models.SidMapping
	Circus              models.SidMapping
	CollegeOfWonder     models.SidMapping
	CrystalTrail        models.SidMapping
	DragonUtopia        models.SidMapping
	EternalDragon       models.SidMapping
	FickleShrine        models.SidMapping
	FlatteringMirror    models.SidMapping
	Forge               models.SidMapping
	Fort                models.SidMapping
	Fountain            models.SidMapping
	Fountain2           models.SidMapping
	HuntsmansCamp       models.SidMapping
	InfernalCirque      models.SidMapping
	InsarasEye          models.SidMapping
	JoustingRange       models.SidMapping
	ManaWell            models.SidMapping
	Market              models.SidMapping
	MineCrystals        models.SidMapping
	MineGemstones       models.SidMapping
	MineGold            models.SidMapping
	MineMercury         models.SidMapping
	MineOre             models.SidMapping
	MineWood            models.SidMapping
	Mirage              models.SidMapping
	MontyHall           models.SidMapping
	MysteriousStone     models.SidMapping
	MysticalTower       models.SidMapping
	MythicScrollBox     models.SidMapping
	OrbObservatory      models.SidMapping
	PandoraBox          models.SidMapping
	PetrifiedMemorial   models.SidMapping
	PileOfBooks         models.SidMapping
	PointOfBalance      models.SidMapping
	Prison              models.SidMapping
	QuixsPath           models.SidMapping
	RandomHire1         models.SidMapping
	RandomHire2         models.SidMapping
	RandomHire3         models.SidMapping
	RandomHire4         models.SidMapping
	RandomHire5         models.SidMapping
	RandomHire6         models.SidMapping
	RandomHire7         models.SidMapping
	RandomItemCommon    models.SidMapping
	RandomItemEpic      models.SidMapping
	RandomItemLegendary models.SidMapping
	RandomItemRare      models.SidMapping
	RemoteFoothold      models.SidMapping
	ResearchLaboratory  models.SidMapping
	RitualPyre          models.SidMapping
	SacrificialShrine   models.SidMapping
	ShadyDen            models.SidMapping
	Stables             models.SidMapping
	Tavern              models.SidMapping
	TearOfTruth         models.SidMapping
	TheGorge            models.SidMapping
	TownGate            models.SidMapping
	TreeOfAbundance     models.SidMapping
	TroglodyteThrone    models.SidMapping
	UnforgottenGrave    models.SidMapping
	University          models.SidMapping
	UnstableRuins       models.SidMapping
	Watchtower          models.SidMapping
	WindRose            models.SidMapping
	WiseOwl             models.SidMapping
}{
	AlchemyLab:          models.SidMapping{Sid: "alchemy_lab", Name: "Alchemy Lab"},
	Arena:               models.SidMapping{Sid: "arena", Name: "Arena"},
	BeerFountain:        models.SidMapping{Sid: "beer_fountain", Name: "Beer Fountain"},
	BorealCall:          models.SidMapping{Sid: "boreal_call", Name: "Boreal Call"},
	CelestialSphere:     models.SidMapping{Sid: "celestial_sphere", Name: "Celestial Sphere"},
	Chimerologist:       models.SidMapping{Sid: "chimerologist", Name: "Chimerologist"},
	Circus:              models.SidMapping{Sid: "circus", Name: "Circus"},
	CollegeOfWonder:     models.SidMapping{Sid: "college_of_wonder", Name: "College Of Wonder"},
	CrystalTrail:        models.SidMapping{Sid: "crystal_trail", Name: "Crystal Trail"},
	DragonUtopia:        models.SidMapping{Sid: "dragon_utopia", Name: "Dragon Utopia"},
	EternalDragon:       models.SidMapping{Sid: "eternal_dragon", Name: "Eternal Dragon"},
	FickleShrine:        models.SidMapping{Sid: "fickle_shrine", Name: "Fickle Shrine"},
	FlatteringMirror:    models.SidMapping{Sid: "flattering_mirror", Name: "Flattering Mirror"},
	Forge:               models.SidMapping{Sid: "forge", Name: "Forge"},
	Fort:                models.SidMapping{Sid: "fort", Name: "Fort"},
	Fountain:            models.SidMapping{Sid: "fountain", Name: "Fountain"},
	Fountain2:           models.SidMapping{Sid: "fountain_2", Name: "Fountain 2"},
	HuntsmansCamp:       models.SidMapping{Sid: "huntsmans_camp", Name: "Huntsman's Camp"},
	InfernalCirque:      models.SidMapping{Sid: "infernal_cirque", Name: "Infernal Cirque"},
	InsarasEye:          models.SidMapping{Sid: "insaras_eye", Name: "Insara's Eye"},
	JoustingRange:       models.SidMapping{Sid: "jousting_range", Name: "Jousting Range"},
	ManaWell:            models.SidMapping{Sid: "mana_well", Name: "Mana Well"},
	Market:              models.SidMapping{Sid: "market", Name: "Market"},
	MineCrystals:        models.SidMapping{Sid: "mine_crystals", Name: "Mine Crystals"},
	MineGemstones:       models.SidMapping{Sid: "mine_gemstones", Name: "Mine Gemstones"},
	MineGold:            models.SidMapping{Sid: "mine_gold", Name: "Mine Gold"},
	MineMercury:         models.SidMapping{Sid: "mine_mercury", Name: "Mine Mercury"},
	MineOre:             models.SidMapping{Sid: "mine_ore", Name: "Mine Ore"},
	MineWood:            models.SidMapping{Sid: "mine_wood", Name: "Mine Wood"},
	Mirage:              models.SidMapping{Sid: "mirage", Name: "Mirage"},
	MontyHall:           models.SidMapping{Sid: "monty_hall", Name: "Monty Hall"},
	MysteriousStone:     models.SidMapping{Sid: "mysterious_stone", Name: "Mysterious Stone"},
	MysticalTower:       models.SidMapping{Sid: "mystical_tower", Name: "Mystical Tower"},
	MythicScrollBox:     models.SidMapping{Sid: "mythic_scroll_box", Name: "Mythic Scroll Box"},
	OrbObservatory:      models.SidMapping{Sid: "orb_observatory", Name: "Orb Observatory"},
	PandoraBox:          models.SidMapping{Sid: "pandora_box", Name: "Pandora Box"},
	PetrifiedMemorial:   models.SidMapping{Sid: "petrified_memorial", Name: "Petrified Memorial"},
	PileOfBooks:         models.SidMapping{Sid: "pile_of_books", Name: "Pile Of Books"},
	PointOfBalance:      models.SidMapping{Sid: "point_of_balance", Name: "Point Of Balance"},
	Prison:              models.SidMapping{Sid: "prison", Name: "Prison"},
	QuixsPath:           models.SidMapping{Sid: "quixs_path", Name: "Quix's Path"},
	RandomHire1:         models.SidMapping{Sid: "random_hire_1", Name: "Random Hire 1"},
	RandomHire2:         models.SidMapping{Sid: "random_hire_2", Name: "Random Hire 2"},
	RandomHire3:         models.SidMapping{Sid: "random_hire_3", Name: "Random Hire 3"},
	RandomHire4:         models.SidMapping{Sid: "random_hire_4", Name: "Random Hire 4"},
	RandomHire5:         models.SidMapping{Sid: "random_hire_5", Name: "Random Hire 5"},
	RandomHire6:         models.SidMapping{Sid: "random_hire_6", Name: "Random Hire 6"},
	RandomHire7:         models.SidMapping{Sid: "random_hire_7", Name: "Random Hire 7"},
	RandomItemCommon:    models.SidMapping{Sid: "random_item_common", Name: "Random Item Common"},
	RandomItemEpic:      models.SidMapping{Sid: "random_item_epic", Name: "Random Item Epic"},
	RandomItemLegendary: models.SidMapping{Sid: "random_item_legendary", Name: "Random Item Legendary"},
	RandomItemRare:      models.SidMapping{Sid: "random_item_rare", Name: "Random Item Rare"},
	RemoteFoothold:      models.SidMapping{Sid: "remote_foothold", Name: "Remote Foothold"},
	ResearchLaboratory:  models.SidMapping{Sid: "research_laboratory", Name: "Research Laboratory"},
	RitualPyre:          models.SidMapping{Sid: "ritual_pyre", Name: "Ritual Pyre"},
	SacrificialShrine:   models.SidMapping{Sid: "sacrificial_shrine", Name: "Sacrificial Shrine"},
	ShadyDen:            models.SidMapping{Sid: "shady_den", Name: "Shady Den"},
	Stables:             models.SidMapping{Sid: "stables", Name: "Stables"},
	Tavern:              models.SidMapping{Sid: "tavern", Name: "Tavern"},
	TearOfTruth:         models.SidMapping{Sid: "tear_of_truth", Name: "Tear Of Truth"},
	TheGorge:            models.SidMapping{Sid: "the_gorge", Name: "The Gorge"},
	TownGate:            models.SidMapping{Sid: "town_gate", Name: "Town Gate"},
	TreeOfAbundance:     models.SidMapping{Sid: "tree_of_abundance", Name: "Tree Of Abundance"},
	TroglodyteThrone:    models.SidMapping{Sid: "troglodyte_throne", Name: "Troglodyte Throne"},
	UnforgottenGrave:    models.SidMapping{Sid: "unforgotten_grave", Name: "Unforgotten Grave"},
	University:          models.SidMapping{Sid: "university", Name: "University"},
	UnstableRuins:       models.SidMapping{Sid: "unstable_ruins", Name: "Unstable Ruins"},
	Watchtower:          models.SidMapping{Sid: "watchtower", Name: "Watchtower"},
	WindRose:            models.SidMapping{Sid: "wind_rose", Name: "Wind Rose"},
	WiseOwl:             models.SidMapping{Sid: "wise_owl", Name: "Wise Owl"},
}
