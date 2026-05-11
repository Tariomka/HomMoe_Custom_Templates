package models

// SidMapping pairs a string ID (sid) with a human-readable name. Mirrors
// Services/ContentManagement/SidMapping.cs from the C# project.
type SidMapping struct {
	Sid  string `json:"sid"`
	Name string `json:"name"`
}

// ContentIds enumerates every world-object SID known to the editor.
// Mirrors Services/ContentManagement/SidMapping.cs ContentIds.
//
//nolint:gochecknoglobals // semantic registry
var ContentIds = struct {
	AlchemyLab          SidMapping
	Arena               SidMapping
	BeerFountain        SidMapping
	BorealCall          SidMapping
	CelestialSphere     SidMapping
	Chimerologist       SidMapping
	Circus              SidMapping
	CollegeOfWonder     SidMapping
	CrystalTrail        SidMapping
	DragonUtopia        SidMapping
	EternalDragon       SidMapping
	FickleShrine        SidMapping
	FlatteringMirror    SidMapping
	Forge               SidMapping
	Fort                SidMapping
	Fountain            SidMapping
	Fountain2           SidMapping
	HuntsmansCamp       SidMapping
	InfernalCirque      SidMapping
	InsarasEye          SidMapping
	JoustingRange       SidMapping
	ManaWell            SidMapping
	Market              SidMapping
	MineCrystals        SidMapping
	MineGemstones       SidMapping
	MineGold            SidMapping
	MineMercury         SidMapping
	MineOre             SidMapping
	MineWood            SidMapping
	Mirage              SidMapping
	MontyHall           SidMapping
	MysteriousStone     SidMapping
	MysticalTower       SidMapping
	MythicScrollBox     SidMapping
	OrbObservatory      SidMapping
	PandoraBox          SidMapping
	PetrifiedMemorial   SidMapping
	PileOfBooks         SidMapping
	PointOfBalance      SidMapping
	Prison              SidMapping
	QuixsPath           SidMapping
	RandomHire1         SidMapping
	RandomHire2         SidMapping
	RandomHire3         SidMapping
	RandomHire4         SidMapping
	RandomHire5         SidMapping
	RandomHire6         SidMapping
	RandomHire7         SidMapping
	RandomItemCommon    SidMapping
	RandomItemEpic      SidMapping
	RandomItemLegendary SidMapping
	RandomItemRare      SidMapping
	RemoteFoothold      SidMapping
	ResearchLaboratory  SidMapping
	RitualPyre          SidMapping
	SacrificialShrine   SidMapping
	ShadyDen            SidMapping
	Stables             SidMapping
	Tavern              SidMapping
	TearOfTruth         SidMapping
	TheGorge            SidMapping
	TownGate            SidMapping
	TreeOfAbundance     SidMapping
	TroglodyteThrone    SidMapping
	UnforgottenGrave    SidMapping
	University          SidMapping
	UnstableRuins       SidMapping
	Watchtower          SidMapping
	WindRose            SidMapping
	WiseOwl             SidMapping
}{
	AlchemyLab:          SidMapping{"alchemy_lab", "Alchemy Lab"},
	Arena:               SidMapping{"arena", "Arena"},
	BeerFountain:        SidMapping{"beer_fountain", "Beer Fountain"},
	BorealCall:          SidMapping{"boreal_call", "Boreal Call"},
	CelestialSphere:     SidMapping{"celestial_sphere", "Celestial Sphere"},
	Chimerologist:       SidMapping{"chimerologist", "Chimerologist"},
	Circus:              SidMapping{"circus", "Circus"},
	CollegeOfWonder:     SidMapping{"college_of_wonder", "College Of Wonder"},
	CrystalTrail:        SidMapping{"crystal_trail", "Crystal Trail"},
	DragonUtopia:        SidMapping{"dragon_utopia", "Dragon Utopia"},
	EternalDragon:       SidMapping{"eternal_dragon", "Eternal Dragon"},
	FickleShrine:        SidMapping{"fickle_shrine", "Fickle Shrine"},
	FlatteringMirror:    SidMapping{"flattering_mirror", "Flattering Mirror"},
	Forge:               SidMapping{"forge", "Forge"},
	Fort:                SidMapping{"fort", "Fort"},
	Fountain:            SidMapping{"fountain", "Fountain"},
	Fountain2:           SidMapping{"fountain_2", "Fountain 2"},
	HuntsmansCamp:       SidMapping{"huntsmans_camp", "Huntsman's Camp"},
	InfernalCirque:      SidMapping{"infernal_cirque", "Infernal Cirque"},
	InsarasEye:          SidMapping{"insaras_eye", "Insara's Eye"},
	JoustingRange:       SidMapping{"jousting_range", "Jousting Range"},
	ManaWell:            SidMapping{"mana_well", "Mana Well"},
	Market:              SidMapping{"market", "Market"},
	MineCrystals:        SidMapping{"mine_crystals", "Mine Crystals"},
	MineGemstones:       SidMapping{"mine_gemstones", "Mine Gemstones"},
	MineGold:            SidMapping{"mine_gold", "Mine Gold"},
	MineMercury:         SidMapping{"mine_mercury", "Mine Mercury"},
	MineOre:             SidMapping{"mine_ore", "Mine Ore"},
	MineWood:            SidMapping{"mine_wood", "Mine Wood"},
	Mirage:              SidMapping{"mirage", "Mirage"},
	MontyHall:           SidMapping{"monty_hall", "Monty Hall"},
	MysteriousStone:     SidMapping{"mysterious_stone", "Mysterious Stone"},
	MysticalTower:       SidMapping{"mystical_tower", "Mystical Tower"},
	MythicScrollBox:     SidMapping{"mythic_scroll_box", "Mythic Scroll Box"},
	OrbObservatory:      SidMapping{"orb_observatory", "Orb Observatory"},
	PandoraBox:          SidMapping{"pandora_box", "Pandora Box"},
	PetrifiedMemorial:   SidMapping{"petrified_memorial", "Petrified Memorial"},
	PileOfBooks:         SidMapping{"pile_of_books", "Pile Of Books"},
	PointOfBalance:      SidMapping{"point_of_balance", "Point Of Balance"},
	Prison:              SidMapping{"prison", "Prison"},
	QuixsPath:           SidMapping{"quixs_path", "Quix's Path"},
	RandomHire1:         SidMapping{"random_hire_1", "Random Hire 1"},
	RandomHire2:         SidMapping{"random_hire_2", "Random Hire 2"},
	RandomHire3:         SidMapping{"random_hire_3", "Random Hire 3"},
	RandomHire4:         SidMapping{"random_hire_4", "Random Hire 4"},
	RandomHire5:         SidMapping{"random_hire_5", "Random Hire 5"},
	RandomHire6:         SidMapping{"random_hire_6", "Random Hire 6"},
	RandomHire7:         SidMapping{"random_hire_7", "Random Hire 7"},
	RandomItemCommon:    SidMapping{"random_item_common", "Random Item Common"},
	RandomItemEpic:      SidMapping{"random_item_epic", "Random Item Epic"},
	RandomItemLegendary: SidMapping{"random_item_legendary", "Random Item Legendary"},
	RandomItemRare:      SidMapping{"random_item_rare", "Random Item Rare"},
	RemoteFoothold:      SidMapping{"remote_foothold", "Remote Foothold"},
	ResearchLaboratory:  SidMapping{"research_laboratory", "Research Laboratory"},
	RitualPyre:          SidMapping{"ritual_pyre", "Ritual Pyre"},
	SacrificialShrine:   SidMapping{"sacrificial_shrine", "Sacrificial Shrine"},
	ShadyDen:            SidMapping{"shady_den", "Shady Den"},
	Stables:             SidMapping{"stables", "Stables"},
	Tavern:              SidMapping{"tavern", "Tavern"},
	TearOfTruth:         SidMapping{"tear_of_truth", "Tear Of Truth"},
	TheGorge:            SidMapping{"the_gorge", "The Gorge"},
	TownGate:            SidMapping{"town_gate", "Town Gate"},
	TreeOfAbundance:     SidMapping{"tree_of_abundance", "Tree Of Abundance"},
	TroglodyteThrone:    SidMapping{"troglodyte_throne", "Troglodyte Throne"},
	UnforgottenGrave:    SidMapping{"unforgotten_grave", "Unforgotten Grave"},
	University:          SidMapping{"university", "University"},
	UnstableRuins:       SidMapping{"unstable_ruins", "Unstable Ruins"},
	Watchtower:          SidMapping{"watchtower", "Watchtower"},
	WindRose:            SidMapping{"wind_rose", "Wind Rose"},
	WiseOwl:             SidMapping{"wise_owl", "Wise Owl"},
}

// IncludeListIds enumerates the named include-list SIDs.
// Mirrors Services/ContentManagement/SidMapping.cs IncludeListIds.
//
//nolint:gochecknoglobals // semantic registry
var IncludeListIds = struct {
	RandomHiresLowTier  SidMapping
	RandomHiresHighTier SidMapping
	RandomHiresAllTier  SidMapping
	ResourceBanksTier1  SidMapping
	ResourceBanksTier2  SidMapping
}{
	RandomHiresLowTier:  SidMapping{"content_list_building_random_hires_low_tier", "Random Hires Low Tier"},
	RandomHiresHighTier: SidMapping{"content_list_building_random_hires_high_tier", "Random Hires High Tier"},
	RandomHiresAllTier:  SidMapping{"basic_content_list_building_random_hires", "Random Hires All Tier"},
	ResourceBanksTier1:  SidMapping{"basic_content_list_building_guarded_resource_banks_tier_1", "Resource Banks T1"},
	ResourceBanksTier2:  SidMapping{"basic_content_list_building_guarded_resource_banks_tier_2", "Resource Banks T2"},
}

// ContentItemGroup categorizes SidMappings by zone-content kind.
// Mirrors Services/ContentManagement/ContentItemGroup.cs.
//
//nolint:gochecknoglobals // semantic registry
var ContentItemGroup = struct {
	Mines         []SidMapping
	Treasures     []SidMapping
	HireBuildings []SidMapping
	ResourceBanks []SidMapping
}{
	Mines: []SidMapping{
		ContentIds.MineWood, ContentIds.MineOre, ContentIds.MineGold,
		ContentIds.MineMercury, ContentIds.MineCrystals, ContentIds.MineGemstones,
		ContentIds.AlchemyLab,
	},
	Treasures: []SidMapping{
		ContentIds.MythicScrollBox, ContentIds.PandoraBox,
		ContentIds.RandomItemCommon, ContentIds.RandomItemEpic, ContentIds.RandomItemLegendary,
	},
	HireBuildings: []SidMapping{
		ContentIds.RandomHire1, ContentIds.RandomHire2, ContentIds.RandomHire3,
		ContentIds.RandomHire4, ContentIds.RandomHire5, ContentIds.RandomHire6, ContentIds.RandomHire7,
		IncludeListIds.RandomHiresLowTier, IncludeListIds.RandomHiresHighTier, IncludeListIds.RandomHiresAllTier,
	},
	ResourceBanks: []SidMapping{
		IncludeListIds.ResourceBanksTier1, IncludeListIds.ResourceBanksTier2,
	},
}

// LookupSidByName returns the SidMapping whose Name matches name (case-insensitive).
func LookupSidByName(name string) (SidMapping, bool) {
	for _, m := range allSidMappings() {
		if equalFold(m.Name, name) {
			return m, true
		}
	}
	return SidMapping{}, false
}

// LookupSid returns the SidMapping whose Sid matches sid.
func LookupSid(sid string) (SidMapping, bool) {
	for _, m := range allSidMappings() {
		if m.Sid == sid {
			return m, true
		}
	}
	return SidMapping{}, false
}

func allSidMappings() []SidMapping {
	all := make([]SidMapping, 0, 80)
	all = append(all,
		ContentIds.AlchemyLab, ContentIds.Arena, ContentIds.BeerFountain, ContentIds.BorealCall,
		ContentIds.CelestialSphere, ContentIds.Chimerologist, ContentIds.Circus, ContentIds.CollegeOfWonder,
		ContentIds.CrystalTrail, ContentIds.DragonUtopia, ContentIds.EternalDragon, ContentIds.FickleShrine,
		ContentIds.FlatteringMirror, ContentIds.Forge, ContentIds.Fort, ContentIds.Fountain,
		ContentIds.Fountain2, ContentIds.HuntsmansCamp, ContentIds.InfernalCirque, ContentIds.InsarasEye,
		ContentIds.JoustingRange, ContentIds.ManaWell, ContentIds.Market, ContentIds.MineCrystals,
		ContentIds.MineGemstones, ContentIds.MineGold, ContentIds.MineMercury, ContentIds.MineOre,
		ContentIds.MineWood, ContentIds.Mirage, ContentIds.MontyHall, ContentIds.MysteriousStone,
		ContentIds.MysticalTower, ContentIds.MythicScrollBox, ContentIds.OrbObservatory, ContentIds.PandoraBox,
		ContentIds.PetrifiedMemorial, ContentIds.PileOfBooks, ContentIds.PointOfBalance, ContentIds.Prison,
		ContentIds.QuixsPath, ContentIds.RandomHire1, ContentIds.RandomHire2, ContentIds.RandomHire3,
		ContentIds.RandomHire4, ContentIds.RandomHire5, ContentIds.RandomHire6, ContentIds.RandomHire7,
		ContentIds.RandomItemCommon, ContentIds.RandomItemEpic, ContentIds.RandomItemLegendary,
		ContentIds.RandomItemRare, ContentIds.RemoteFoothold, ContentIds.ResearchLaboratory,
		ContentIds.RitualPyre, ContentIds.SacrificialShrine, ContentIds.ShadyDen, ContentIds.Stables,
		ContentIds.Tavern, ContentIds.TearOfTruth, ContentIds.TheGorge, ContentIds.TownGate,
		ContentIds.TreeOfAbundance, ContentIds.TroglodyteThrone, ContentIds.UnforgottenGrave,
		ContentIds.University, ContentIds.UnstableRuins, ContentIds.Watchtower, ContentIds.WindRose,
		ContentIds.WiseOwl,
		IncludeListIds.RandomHiresLowTier, IncludeListIds.RandomHiresHighTier, IncludeListIds.RandomHiresAllTier,
		IncludeListIds.ResourceBanksTier1, IncludeListIds.ResourceBanksTier2,
	)
	return all
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}
