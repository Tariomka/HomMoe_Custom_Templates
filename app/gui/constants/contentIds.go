package constants

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

type contentIds struct {
	AbandonedCorpse     models.SidMapping
	AbandonedMansion    models.SidMapping
	AbnormalStructure   models.SidMapping
	AlchemyLab          models.SidMapping
	AltarOfMagic1       models.SidMapping
	AltarOfMagic2       models.SidMapping
	AltarOfMagic3       models.SidMapping
	AltarOfMagic4       models.SidMapping
	AlvarsEye           models.SidMapping
	Arena               models.SidMapping
	ArmoryAutomaton     models.SidMapping
	BeerFountain        models.SidMapping
	BlackTower          models.SidMapping
	BorealCall          models.SidMapping
	CelestialSphere     models.SidMapping
	Chimerologist       models.SidMapping
	CircleOfLife        models.SidMapping
	Circus              models.SidMapping
	CollegeOfWonder     models.SidMapping
	CrowNest            models.SidMapping
	CrystalTrail        models.SidMapping
	CursedOldHouse      models.SidMapping
	DragonUtopia        models.SidMapping
	EnchantedScrollBox  models.SidMapping
	EternalDragon       models.SidMapping
	FickleShrine        models.SidMapping
	FlatteringMirror    models.SidMapping
	Forge               models.SidMapping
	Fort                models.SidMapping
	Fountain            models.SidMapping
	Fountain2           models.SidMapping
	Gardener            models.SidMapping
	GingerbreadHouse    models.SidMapping
	GoblinCache         models.SidMapping
	HerosCrypt          models.SidMapping
	HuntsmansCamp       models.SidMapping
	InfernalCirque      models.SidMapping
	InsarasEye          models.SidMapping
	IridescentAbbey     models.SidMapping
	JoustingRange       models.SidMapping
	KnowledgeGarden     models.SidMapping
	LearningStone       models.SidMapping
	LegionsMemorial     models.SidMapping
	LostLibrary         models.SidMapping
	MagicAmplifier1     models.SidMapping
	MagicAmplifier2     models.SidMapping
	MagicAmplifier3     models.SidMapping
	MagicAmplifier4     models.SidMapping
	MagicWheel          models.SidMapping
	ManaWell            models.SidMapping
	Market              models.SidMapping
	Maze                models.SidMapping
	MercenaryGuild      models.SidMapping
	MereasShrine        models.SidMapping
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
	OvergrownGrave      models.SidMapping
	PandoraBox          models.SidMapping
	PeasantCart         models.SidMapping
	PetrifiedMemorial   models.SidMapping
	PileOfBooks         models.SidMapping
	PointOfBalance      models.SidMapping
	PrismaticLair       models.SidMapping
	Prison              models.SidMapping
	QuixsPath           models.SidMapping
	RaidersCamp         models.SidMapping
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
	ScrollBox           models.SidMapping
	ShadyDen            models.SidMapping
	Stables             models.SidMapping
	StingingSword       models.SidMapping
	StorageCrystals     models.SidMapping
	StorageDust         models.SidMapping
	StorageGemstones    models.SidMapping
	StorageGold         models.SidMapping
	StorageMercury      models.SidMapping
	StorageOre          models.SidMapping
	StorageWood         models.SidMapping
	Tavern              models.SidMapping
	TearOfTruth         models.SidMapping
	TheGorge            models.SidMapping
	TownGate            models.SidMapping
	TreeOfAbundance     models.SidMapping
	TreeOfKnowledge     models.SidMapping
	TrialScales         models.SidMapping
	TroglodyteThrone    models.SidMapping
	TwilightBloom       models.SidMapping
	UncannyRite         models.SidMapping
	UnforgottenGrave    models.SidMapping
	University          models.SidMapping
	UnstableRuins       models.SidMapping
	Village             models.SidMapping
	Watchtower          models.SidMapping
	Windmill            models.SidMapping
	WindRose            models.SidMapping
	WiseOwl             models.SidMapping
}

// ContentIds enumerates every world-object SID known to the editor.
var ContentIds = func() contentIds {
	interactableObjects := registry.GetMapObjectAllInteractableValues()
	randomItemObjects := registry.GetMapObjectRandomItemValues()
	resourceObjects := registry.GetMapObjectResourceValues()
	scrollObjects := registry.GetMapObjectScrollValues()

	return contentIds{
		AbandonedCorpse:    models.SidMapping{Sid: interactableObjects.AbandonedCorpse, Name: "Forgotten Remains"},
		AbandonedMansion:   models.SidMapping{Sid: interactableObjects.AbandonedMansion, Name: "Abandoned Mansion"},
		AbnormalStructure:  models.SidMapping{Sid: interactableObjects.AbnormalStructure, Name: "Abnormal Structure"},
		AlchemyLab:         models.SidMapping{Sid: interactableObjects.AlchemyLab, Name: "Alchemical Lab"},
		AltarOfMagic1:      models.SidMapping{Sid: interactableObjects.AltarOfMagic1, Name: "Nightshade Shrine"},
		AltarOfMagic2:      models.SidMapping{Sid: interactableObjects.AltarOfMagic2, Name: "Daylight Shrine"},
		AltarOfMagic3:      models.SidMapping{Sid: interactableObjects.AltarOfMagic3, Name: "Arcane Shrine"},
		AltarOfMagic4:      models.SidMapping{Sid: interactableObjects.AltarOfMagic4, Name: "Primal Shrine"},
		AlvarsEye:          models.SidMapping{Sid: interactableObjects.AlvarsEye, Name: "Alvar Outpost"},
		Arena:              models.SidMapping{Sid: interactableObjects.Arena, Name: "Pit of Glory"},
		ArmoryAutomaton:    models.SidMapping{Sid: interactableObjects.ArmoryAutomaton, Name: "Summit Automaton"},
		BeerFountain:       models.SidMapping{Sid: interactableObjects.BeerFountain, Name: "Beer Fountain"},
		BlackTower:         models.SidMapping{Sid: interactableObjects.BlackTower, Name: "Black Tower"},
		BorealCall:         models.SidMapping{Sid: interactableObjects.BorealCall, Name: "Boreal Call"},
		CelestialSphere:    models.SidMapping{Sid: interactableObjects.CelestialSphere, Name: "Celestial Spire"},
		Chimerologist:      models.SidMapping{Sid: interactableObjects.Chimerologist, Name: "Chimerologist"},
		CircleOfLife:       models.SidMapping{Sid: interactableObjects.CircleOfLife, Name: "Circle of Life"},
		Circus:             models.SidMapping{Sid: interactableObjects.Circus, Name: "Travelling Circus"},
		CollegeOfWonder:    models.SidMapping{Sid: interactableObjects.CollegeOfWonder, Name: "Mountain Monastery"},
		CrowNest:           models.SidMapping{Sid: interactableObjects.CrowNest, Name: "Crow Nest"},
		CrystalTrail:       models.SidMapping{Sid: interactableObjects.CrystalTrail, Name: "Dragon Step"},
		CursedOldHouse:     models.SidMapping{Sid: interactableObjects.CursedOldHouse, Name: "Cursed Old House"},
		DragonUtopia:       models.SidMapping{Sid: interactableObjects.DragonUtopia, Name: "Dragon Utopia"},
		EnchantedScrollBox: models.SidMapping{Sid: scrollObjects.EnchantedScrollBox, Name: "Enchanted Magic Scroll"},
		EternalDragon:      models.SidMapping{Sid: interactableObjects.EternalDragon, Name: "Eternal Dragon"},
		FickleShrine:       models.SidMapping{Sid: interactableObjects.FickleShrine, Name: "Four Scholars Shrine"},
		FlatteringMirror:   models.SidMapping{Sid: interactableObjects.FlatteringMirror, Name: "World Mirror"},
		Forge:              models.SidMapping{Sid: interactableObjects.Forge, Name: "Forge of the Second Man"},
		Fort:               models.SidMapping{Sid: interactableObjects.Fort, Name: "Pauper Knight Order"},
		Fountain:           models.SidMapping{Sid: interactableObjects.Fountain, Name: "Fountain"},
		Fountain2:          models.SidMapping{Sid: interactableObjects.Fountain2, Name: "Fountain 2"},
		Gardener:           models.SidMapping{Sid: interactableObjects.Gardener, Name: "Magic Garden"},
		GingerbreadHouse:   models.SidMapping{Sid: interactableObjects.GingerbreadHouse, Name: "Gingerbread House"},
		GoblinCache:        models.SidMapping{Sid: interactableObjects.GoblinCache, Name: "Goblin Cache"},
		HerosCrypt:         models.SidMapping{Sid: interactableObjects.HerosCrypt, Name: "Tomb of a Nameless Hero"},
		HuntsmansCamp:      models.SidMapping{Sid: interactableObjects.HuntsmansCamp, Name: "Explorer's Camp"},
		InfernalCirque:     models.SidMapping{Sid: interactableObjects.InfernalCirque, Name: "Infernal Cirque"},
		InsarasEye:         models.SidMapping{Sid: interactableObjects.InsarasEye, Name: "Insara's Eye"},
		IridescentAbbey:    models.SidMapping{Sid: interactableObjects.IridescentAbbey, Name: "Iridescent Abbey"},
		JoustingRange:      models.SidMapping{Sid: interactableObjects.JoustingRange, Name: "Colosseum"},
		KnowledgeGarden:    models.SidMapping{Sid: interactableObjects.KnowledgeGarden, Name: "Knowledge Garden"},
		LearningStone:      models.SidMapping{Sid: interactableObjects.LearningStone, Name: "Learning Stone"},
		LegionsMemorial:    models.SidMapping{Sid: interactableObjects.LegionsMemorial, Name: "Legion's Memorial"},
		LostLibrary:        models.SidMapping{Sid: interactableObjects.LostLibrary, Name: "Lost Library"},
		MagicAmplifier1:    models.SidMapping{Sid: interactableObjects.MagicAmplifier1, Name: "Nightshade Amplifier"},
		MagicAmplifier2:    models.SidMapping{Sid: interactableObjects.MagicAmplifier2, Name: "Daylight Amplifier"},
		MagicAmplifier3:    models.SidMapping{Sid: interactableObjects.MagicAmplifier3, Name: "Arcane Amplifier"},
		MagicAmplifier4:    models.SidMapping{Sid: interactableObjects.MagicAmplifier4, Name: "Primal Amplifier"},
		MagicWheel:         models.SidMapping{Sid: interactableObjects.MagicWheel, Name: "Magic Wheel"},
		ManaWell:           models.SidMapping{Sid: interactableObjects.ManaWell, Name: "Well"},
		Market:             models.SidMapping{Sid: interactableObjects.Market, Name: "Marketplace"},
		Maze:               models.SidMapping{Sid: interactableObjects.Maze, Name: "Living Maze"},
		MercenaryGuild:     models.SidMapping{Sid: interactableObjects.MercenaryGuild, Name: "Mercenary Guild"},
		MereasShrine:       models.SidMapping{Sid: interactableObjects.MereasShrine, Name: "Mearea's Altar"},
		MineCrystals:       models.SidMapping{Sid: interactableObjects.CrystalMine, Name: "Crystal Vein"},
		MineGemstones:      models.SidMapping{Sid: interactableObjects.GemstoneMine, Name: "Gem Mound"},
		MineGold:           models.SidMapping{Sid: interactableObjects.GoldMine, Name: "Gold Mine"},
		MineMercury:        models.SidMapping{Sid: interactableObjects.MercuryMine, Name: "Mercury Fissure"},
		MineOre:            models.SidMapping{Sid: interactableObjects.OreMine, Name: "Ore Mine"},
		MineWood:           models.SidMapping{Sid: interactableObjects.WoodMine, Name: "Sawmill"},
		Mirage:             models.SidMapping{Sid: interactableObjects.Mirage, Name: "Mirage"},
		MontyHall:          models.SidMapping{Sid: interactableObjects.MontyHall, Name: "The Monty Hall"},
		MysteriousStone:    models.SidMapping{Sid: interactableObjects.MysteriousStone, Name: "Whispering Stones"},
		MysticalTower:      models.SidMapping{Sid: interactableObjects.MysticalTower, Name: "Altar of Insight"},
		MythicScrollBox:    models.SidMapping{Sid: scrollObjects.MythicScrollBox, Name: "Mythic Magic Scroll"},
		OrbObservatory: models.SidMapping{
			Sid:  interactableObjects.OrbObservatory,
			Name: "Four Scholars Observatory",
		},
		OvergrownGrave:    models.SidMapping{Sid: interactableObjects.OvergrownGrave, Name: "Ancient Crypt"},
		PandoraBox:        models.SidMapping{Sid: resourceObjects.PandoraBox, Name: "Pandora's Box"},
		PeasantCart:       models.SidMapping{Sid: interactableObjects.PeasantCart, Name: "Abandoned Cart"},
		PetrifiedMemorial: models.SidMapping{Sid: interactableObjects.PetrifiedMemorial, Name: "Petrified Memorial"},
		PileOfBooks:       models.SidMapping{Sid: interactableObjects.PileOfBooks, Name: "Pile of Books"},
		PointOfBalance:    models.SidMapping{Sid: interactableObjects.PointOfBalance, Name: "Point of Balance"},
		PrismaticLair:     models.SidMapping{Sid: interactableObjects.PrismaticLair, Name: "Prismatic Nest"},
		Prison:            models.SidMapping{Sid: interactableObjects.Prison, Name: "Hero Cage"},
		QuixsPath:         models.SidMapping{Sid: interactableObjects.QuixsPath, Name: "Quix's Altar"},
		RaidersCamp:       models.SidMapping{Sid: interactableObjects.RaidersCamp, Name: "Raiders' Camp"},
		RandomHire1:       models.SidMapping{Sid: interactableObjects.RandomHireTier1, Name: "Random Hire Tier 1"},
		RandomHire2:       models.SidMapping{Sid: interactableObjects.RandomHireTier2, Name: "Random Hire Tier 2"},
		RandomHire3:       models.SidMapping{Sid: interactableObjects.RandomHireTier3, Name: "Random Hire Tier 3"},
		RandomHire4:       models.SidMapping{Sid: interactableObjects.RandomHireTier4, Name: "Random Hire Tier 4"},
		RandomHire5:       models.SidMapping{Sid: interactableObjects.RandomHireTier5, Name: "Random Hire Tier 5"},
		RandomHire6:       models.SidMapping{Sid: interactableObjects.RandomHireTier6, Name: "Random Hire Tier 6"},
		RandomHire7:       models.SidMapping{Sid: interactableObjects.RandomHireTier7, Name: "Random Hire Tier 7"},
		RandomItemCommon:  models.SidMapping{Sid: randomItemObjects.RandomItemCommon, Name: "Random Item Common"},
		RandomItemEpic:    models.SidMapping{Sid: randomItemObjects.RandomItemEpic, Name: "Random Item Epic"},
		RandomItemLegendary: models.SidMapping{
			Sid:  randomItemObjects.RandomItemLegendary,
			Name: "Random Item Legendary",
		},
		RandomItemRare: models.SidMapping{Sid: randomItemObjects.RandomItemRare, Name: "Random Item Rare"},
		RemoteFoothold: models.SidMapping{Sid: interactableObjects.RemoteFoothold, Name: "Remote Foothold"},
		ResearchLaboratory: models.SidMapping{
			Sid:  interactableObjects.ResearchLaboratory,
			Name: "Research Laboratory",
		},
		RitualPyre:        models.SidMapping{Sid: interactableObjects.RitualPyre, Name: "Ritual Pyre"},
		SacrificialShrine: models.SidMapping{Sid: interactableObjects.SacrificialShrine, Name: "Sacrificial Shrine"},
		ScrollBox:         models.SidMapping{Sid: scrollObjects.ScrollBox, Name: "Magic Scroll"},
		ShadyDen:          models.SidMapping{Sid: interactableObjects.ShadyDen, Name: "Trophy Hunter's Den"},
		Stables:           models.SidMapping{Sid: interactableObjects.Stables, Name: "Stables"},
		StingingSword:     models.SidMapping{Sid: interactableObjects.StingingSword, Name: "Stinging Sword"},
		StorageCrystals:   models.SidMapping{Sid: interactableObjects.CrystalStorage, Name: "Crystal Storage"},
		StorageDust:       models.SidMapping{Sid: interactableObjects.DustStorage, Name: "Alchemical Dust Storage"},
		StorageGemstones:  models.SidMapping{Sid: interactableObjects.GemstoneStorage, Name: "Gem Storage"},
		StorageGold:       models.SidMapping{Sid: interactableObjects.GoldStorage, Name: "Gold Storage"},
		StorageMercury:    models.SidMapping{Sid: interactableObjects.MercuryStorage, Name: "Mercury Storage"},
		StorageOre:        models.SidMapping{Sid: interactableObjects.OreStorage, Name: "Ore Storage"},
		StorageWood:       models.SidMapping{Sid: interactableObjects.WoodStorage, Name: "Wood Storage"},
		Tavern:            models.SidMapping{Sid: interactableObjects.Tavern, Name: "Tavern"},
		TearOfTruth:       models.SidMapping{Sid: interactableObjects.TearOfTruth, Name: "Tear of Truth"},
		TheGorge:          models.SidMapping{Sid: interactableObjects.Gorge, Name: "Carrion Pile"},
		TownGate:          models.SidMapping{Sid: interactableObjects.TownGate, Name: "Town Gate"},
		TreeOfAbundance:   models.SidMapping{Sid: interactableObjects.TreeOfAbundance, Name: "Arborcopia"},
		TreeOfKnowledge:   models.SidMapping{Sid: interactableObjects.TreeOfKnowledge, Name: "Tree of Knowledge"},
		TrialScales:       models.SidMapping{Sid: interactableObjects.TrialScales, Name: "Scales of Worth"},
		TroglodyteThrone:  models.SidMapping{Sid: interactableObjects.TroglodyteThrone, Name: "Troglodyte Throne"},
		TwilightBloom:     models.SidMapping{Sid: interactableObjects.TwilightBloom, Name: "Twilight Bloom"},
		UncannyRite:       models.SidMapping{Sid: interactableObjects.UncannyRite, Name: "Uncanny Rite"},
		UnforgottenGrave:  models.SidMapping{Sid: interactableObjects.UnforgottenGrave, Name: "Tomb of Vigilance"},
		University:        models.SidMapping{Sid: interactableObjects.University, Name: "University"},
		UnstableRuins:     models.SidMapping{Sid: interactableObjects.UnstableRuins, Name: "Overgrown Vori Ruins"},
		Village:           models.SidMapping{Sid: interactableObjects.Village, Name: "Village"},
		Watchtower:        models.SidMapping{Sid: interactableObjects.Watchtower, Name: "Redwood Observatory"},
		Windmill:          models.SidMapping{Sid: interactableObjects.Windmill, Name: "Windmill"},
		WindRose:          models.SidMapping{Sid: interactableObjects.WindRose, Name: "Wind Rose"},
		WiseOwl:           models.SidMapping{Sid: interactableObjects.WiseOwl, Name: "Wise Owl"},
	}
}()
