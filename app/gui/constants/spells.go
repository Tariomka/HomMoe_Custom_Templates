package constants

import (
	"cmp"
	"image/color"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

const (
	schoolNameNeutral = "High Neutral"
	schoolNameDay     = "Daylight"
	schoolNameNight   = "Nightshade"
	schoolNameSpace   = "Arcane"
	schoolNamePrimal  = "Primal"
)

// SpellEntry pairs a learnable spell SID with its display name, magic school
// and tier (used for sorting in the spell picker).
type SpellEntry struct {
	Sid    string
	Name   string
	School string
	Tier   int
}

func GetKnownSpellsWithExclusions(excluded []string) []SpellEntry {
	spells := slices.DeleteFunc(
		buildKnownSpells(),
		func(spell SpellEntry) bool { return slices.Contains(excluded, spell.Sid) })
	slices.SortStableFunc(spells, CompareSpellEntries)
	return spells
}

// FindSpell returns the catalog entry for a spell SID, or ok=false when the
// SID is not in the catalog.
func FindSpell(sid string) (SpellEntry, bool) {
	for _, spell := range buildKnownSpells() {
		if spell.Sid == sid {
			return spell, true
		}
	}

	return SpellEntry{}, false
}

func GetSpellSchoolDisplayName(schoolType string) string {
	spellSchoolValues := registry.GetSpellSchoolTypeValues()
	switch schoolType {
	case spellSchoolValues.HighNeutral:
		return schoolNameNeutral
	case spellSchoolValues.Daylight:
		return schoolNameDay
	case spellSchoolValues.Nightshade:
		return schoolNameNight
	case spellSchoolValues.Arcane:
		return schoolNameSpace
	case spellSchoolValues.Primal:
		return schoolNamePrimal
	default:
		return schoolType
	}
}

// GetSpellSchoolColorFromDisplayName maps a school display name to its accent color.
func GetSpellSchoolColorFromDisplayName(displayName string) color.NRGBA {
	switch displayName {
	case schoolNameNeutral:
		return themes.ColorsSpellSchools.HighNeutral
	case schoolNameDay:
		return themes.ColorsSpellSchools.Daylight
	case schoolNameNight:
		return themes.ColorsSpellSchools.Nightshade
	case schoolNameSpace:
		return themes.ColorsSpellSchools.Arcane
	case schoolNamePrimal:
		return themes.ColorsSpellSchools.Primal
	default:
		return themes.ColorsBase.Accent
	}
}

// GetSpellSchoolColor maps a school display name to its accent color.
func GetSpellSchoolColor(schoolName string) color.NRGBA {
	spellSchoolValues := registry.GetSpellSchoolTypeValues()
	switch schoolName {
	case spellSchoolValues.HighNeutral:
		return themes.ColorsSpellSchools.HighNeutral
	case spellSchoolValues.Daylight:
		return themes.ColorsSpellSchools.Daylight
	case spellSchoolValues.Nightshade:
		return themes.ColorsSpellSchools.Nightshade
	case spellSchoolValues.Arcane:
		return themes.ColorsSpellSchools.Arcane
	case spellSchoolValues.Primal:
		return themes.ColorsSpellSchools.Primal
	default:
		return themes.ColorsBase.Accent
	}
}

func CompareSpellEntries(a, b SpellEntry) int {
	schoolIndexA, schoolIndexB := 99, 99
	for i, school := range registry.GetSpellSchoolTypeList() {
		if a.School == school {
			schoolIndexA = i
		}
		if b.School == school {
			schoolIndexB = i
		}
	}
	if comparison := cmp.Compare(schoolIndexA, schoolIndexB); comparison != 0 {
		return comparison
	}

	if comparison := cmp.Compare(a.Tier, b.Tier); comparison != 0 {
		return comparison
	}

	return cmp.Compare(a.Name, b.Name)
}

func buildKnownSpells() []SpellEntry {
	spells := []SpellEntry{}
	spells = append(spells, buildHighNeutralSpells()...)
	spells = append(spells, buildDaylightSpells()...)
	spells = append(spells, buildNightshadeSpells()...)
	spells = append(spells, buildArcaneSpells()...)
	spells = append(spells, buildPrimalSpells()...)
	return spells
}

func buildHighNeutralSpells() []SpellEntry {
	spells := registry.GetHighNeutralSpellSidValues()
	spellSchoolValues := registry.GetSpellSchoolTypeValues()
	return []SpellEntry{
		{spells.PocketDimension, "Pocket Dimension", spellSchoolValues.HighNeutral, 2},
		{spells.SecondSight, "Second Sight", spellSchoolValues.HighNeutral, 2},
		{spells.ShadowForm, "Shadowflight", spellSchoolValues.HighNeutral, 3},
		{spells.TownPortal, "Town Portal", spellSchoolValues.HighNeutral, 3},
		{spells.DimensionDoor, "Dimension Door", spellSchoolValues.HighNeutral, 4},
		{spells.LightGate, "Gate of Light", spellSchoolValues.HighNeutral, 4},
	}
}

func buildDaylightSpells() []SpellEntry {
	spells := registry.GetDaylightSpellSidValues()
	spellSchoolValues := registry.GetSpellSchoolTypeValues()
	return []SpellEntry{
		{spells.SharpEdge, "Blessing", spellSchoolValues.Daylight, 1},
		{spells.Haste, "Haste", spellSchoolValues.Daylight, 1},
		{spells.HealingWater, "Healing Water", spellSchoolValues.Daylight, 1},
		{spells.ShortenShadow, "Shorten Shadow", spellSchoolValues.Daylight, 1},
		{spells.FavorableWind, "Favorable Wind", spellSchoolValues.Daylight, 2},
		{spells.ClearView, "From a Bird's Eye", spellSchoolValues.Daylight, 2},
		{spells.InnerLight, "Inner Light", spellSchoolValues.Daylight, 2},
		{spells.CleansingRay, "Weakening Ray", spellSchoolValues.Daylight, 2},
		{spells.ArinasHymn, "Arina's Touch", spellSchoolValues.Daylight, 3},
		{spells.MasterfulParry, "Riposte", spellSchoolValues.Daylight, 3},
		{spells.SecondSong, "Song of Power", spellSchoolValues.Daylight, 3},
		{spells.Taunt, "Taunt", spellSchoolValues.Daylight, 3},
		{spells.Farsight, "Clear Fog", spellSchoolValues.Daylight, 4},
		{spells.HolyArms, "Heavenly Blades", spellSchoolValues.Daylight, 4},
		{spells.RadiantArmor, "Radiant Armor", spellSchoolValues.Daylight, 4},
		{spells.Vengeance, "Vengeance", spellSchoolValues.Daylight, 4},
		{spells.ArinasChosen, "Arina's Chosen", spellSchoolValues.Daylight, 5},
		{spells.Judgement, "Judgement", spellSchoolValues.Daylight, 5},
	}
}

func buildNightshadeSpells() []SpellEntry {
	spells := registry.GetNightshadeSpellSidValues()
	spellSchoolValues := registry.GetSpellSchoolTypeValues()
	return []SpellEntry{
		{spells.Despair, "Despair", spellSchoolValues.Nightshade, 1},
		{spells.EnlargeShadow, "Enlarge Shadow", spellSchoolValues.Nightshade, 1},
		{spells.FatalDecay, "Fatal Decay", spellSchoolValues.Nightshade, 1},
		{spells.UnnaturalCalm, "Unnatural Calm", spellSchoolValues.Nightshade, 1},
		{spells.ReadMinds, "Read Minds", spellSchoolValues.Nightshade, 2},
		{spells.ShadeCloak, "Shade Cloak", spellSchoolValues.Nightshade, 2},
		{spells.DeathsGrip, "Umbral Grip", spellSchoolValues.Nightshade, 2},
		{spells.Web, "Web", spellSchoolValues.Nightshade, 2},
		{spells.NairasVeil, "Naira's Veil", spellSchoolValues.Nightshade, 3},
		{spells.Silence, "Silence", spellSchoolValues.Nightshade, 3},
		{spells.Sleep, "Sleep", spellSchoolValues.Nightshade, 3},
		{spells.Twilight, "Twilight", spellSchoolValues.Nightshade, 3},
		{spells.Berserker, "Berserk", spellSchoolValues.Nightshade, 4},
		{spells.SummonStarchild, "Summon Starchild", spellSchoolValues.Nightshade, 4},
		{spells.Vulnerability, "Vulnerability", spellSchoolValues.Nightshade, 4},
		{spells.DeathsCall, "Coup de Grace", spellSchoolValues.Nightshade, 5},
		{spells.NairasKiss, "Naira's Kiss", spellSchoolValues.Nightshade, 5},
		{spells.ShadowArmy, "Shadow Army", spellSchoolValues.Nightshade, 5},
	}
}

func buildArcaneSpells() []SpellEntry {
	spells := registry.GetArcaneSpellSidValues()
	spellSchoolValues := registry.GetSpellSchoolTypeValues()
	return []SpellEntry{
		{spells.EarlyStart, "Early Start", spellSchoolValues.Arcane, 1},
		{spells.Energyze, "Energize", spellSchoolValues.Arcane, 1},
		{spells.Decimate, "Guillotine", spellSchoolValues.Arcane, 1},
		{spells.OpticalIllusion, "Optical Illusion", spellSchoolValues.Arcane, 1},
		{spells.Blink, "Blink", spellSchoolValues.Arcane, 2},
		{spells.Carapace, "Carapace", spellSchoolValues.Arcane, 2},
		{spells.EnergyExplosion, "Energy Explosion", spellSchoolValues.Arcane, 2},
		{spells.Reinforcements, "Reinforcements", spellSchoolValues.Arcane, 2},
		{spells.Assemble, "Assemble!", spellSchoolValues.Arcane, 3},
		{spells.ImpendingFate, "Impending Fate", spellSchoolValues.Arcane, 3},
		{spells.Shackles, "Shackles", spellSchoolValues.Arcane, 3},
		{spells.TrapJaws, "Temporal Spheres", spellSchoolValues.Arcane, 3},
		{spells.MirrorCopy, "Mirror Copy", spellSchoolValues.Arcane, 4},
		{spells.Rewind, "Rewind Life", spellSchoolValues.Arcane, 4},
		{spells.TrapSnare, "Spatial Snare", spellSchoolValues.Arcane, 4},
		{spells.BlackHole, "Black Hole", spellSchoolValues.Arcane, 5},
		{spells.DoreathsTide, "Doreath's Tide", spellSchoolValues.Arcane, 5},
		{spells.RealityDistortion, "Reality Distortion", spellSchoolValues.Arcane, 5},
	}
}

func buildPrimalSpells() []SpellEntry {
	spells := registry.GetPrimalSpellSidValues()
	spellSchoolValues := registry.GetSpellSchoolTypeValues()
	return []SpellEntry{
		{spells.Groundsight, "Groundsight", spellSchoolValues.Primal, 1},
		{spells.Thunderbolt, "Lightning Bolt", spellSchoolValues.Primal, 1},
		{spells.ThickHide, "Thick Hide", spellSchoolValues.Primal, 1},
		{spells.CrystalCrown, "Crystal Crown", spellSchoolValues.Primal, 2},
		{spells.FireGlobe, "Fireball", spellSchoolValues.Primal, 2},
		{spells.IceBolt, "Ice Bolt", spellSchoolValues.Primal, 2},
		{spells.Wean, "Wean", spellSchoolValues.Primal, 2},
		{spells.CaveIn, "Cave In", spellSchoolValues.Primal, 3},
		{spells.EarthsRage, "Earth's Rage", spellSchoolValues.Primal, 3},
		{spells.WallOfFlame, "Firewall", spellSchoolValues.Primal, 3},
		{spells.StoneFangs, "Stone Fangs", spellSchoolValues.Primal, 3},
		{spells.PrimordialPurity, "Anti-Magic", spellSchoolValues.Primal, 4},
		{spells.ChainLightning, "Chain Lightning", spellSchoolValues.Primal, 4},
		{spells.Avalanche, "Circle of Winter", spellSchoolValues.Primal, 4},
		{spells.PrimordialChaos, "Primordial Chaos", spellSchoolValues.Primal, 4},
		{spells.Armageddon, "Armageddon", spellSchoolValues.Primal, 5},
		{spells.HksmillasRampage, "Hksmilla's Rampage", spellSchoolValues.Primal, 5},
		{spells.SummonPrimalRemnant, "Summon Primal Remnant", spellSchoolValues.Primal, 5},
	}
}
