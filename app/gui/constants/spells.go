package constants

import (
	"cmp"
	"image/color"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

// SpellEntry pairs a learnable spell SID with its display name, magic school
// and tier (used for sorting in the spell picker).
type SpellEntry struct {
	Sid    string
	Name   string
	School string
	Tier   int
}

// SpellSchoolDisplayNames maps a school key to its display label.
var SpellSchoolDisplayNames = map[string]string{
	"neutral": "High Neutral",
	"day":     "Daylight",
	"night":   "Nightshade",
	"space":   "Arcane",
	"primal":  "Primal",
}

// KnownSpells is the catalog of learnable spells. SIDs come from the
// registry; names, schools and tiers are editor-side labels. Grouped
// logically by school; consumers sort by tier within each school as needed.
var KnownSpells = buildKnownSpells()

func buildKnownSpells() []SpellEntry {
	sids := registry.GetSpellSidValues()
	spellSchoolValues := registry.GetSpellSchoolTypeValues()
	return []SpellEntry{
		// High Neutral
		{sids.NeutralPocketDimension, "Pocket Dimension", spellSchoolValues.HighNeutral, 2},
		{sids.NeutralSecondSight, "Second Sight", spellSchoolValues.HighNeutral, 2},
		{sids.NeutralShadowForm, "Shadowflight", spellSchoolValues.HighNeutral, 3},
		{sids.NeutralTownPortal, "Town Portal", spellSchoolValues.HighNeutral, 3},
		{sids.NeutralDimensionDoor, "Dimension Door", spellSchoolValues.HighNeutral, 4},
		{sids.NeutralLightGate, "Gate of Light", spellSchoolValues.HighNeutral, 4},

		// Daylight
		{sids.DaySharpEdge, "Blessing", spellSchoolValues.Daylight, 1},
		{sids.DayHaste, "Haste", spellSchoolValues.Daylight, 1},
		{sids.DayHealingWater, "Healing Water", spellSchoolValues.Daylight, 1},
		{sids.DayShortenShadow, "Shorten Shadow", spellSchoolValues.Daylight, 1},
		{sids.DayFavorableWind, "Favorable Wind", spellSchoolValues.Daylight, 2},
		{sids.DayClearView, "From a Bird's Eye", spellSchoolValues.Daylight, 2},
		{sids.DayInnerLight, "Inner Light", spellSchoolValues.Daylight, 2},
		{sids.DayCleansingRay, "Weakening Ray", spellSchoolValues.Daylight, 2},
		{sids.DayArinasHymn, "Arina's Touch", spellSchoolValues.Daylight, 3},
		{sids.DayMasterfulParry, "Riposte", spellSchoolValues.Daylight, 3},
		{sids.DaySecondSong, "Song of Power", spellSchoolValues.Daylight, 3},
		{sids.DayTaunt, "Taunt", spellSchoolValues.Daylight, 3},
		{sids.DayFarsight, "Clear Fog", spellSchoolValues.Daylight, 4},
		{sids.DayHolyArms, "Heavenly Blades", spellSchoolValues.Daylight, 4},
		{sids.DayRadiantArmor, "Radiant Armor", spellSchoolValues.Daylight, 4},
		{sids.DayVengeance, "Vengeance", spellSchoolValues.Daylight, 4},
		{sids.DayArinasChosen, "Arina's Chosen", spellSchoolValues.Daylight, 5},
		{sids.DayJudgement, "Judgement", spellSchoolValues.Daylight, 5},

		// Nightshade
		{sids.NightDespair, "Despair", spellSchoolValues.Nightshade, 1},
		{sids.NightEnlargeShadow, "Enlarge Shadow", spellSchoolValues.Nightshade, 1},
		{sids.NightFatalDecay, "Fatal Decay", spellSchoolValues.Nightshade, 1},
		{sids.NightUnnaturalCalm, "Unnatural Calm", spellSchoolValues.Nightshade, 1},
		{sids.NightReadMinds, "Read Minds", spellSchoolValues.Nightshade, 2},
		{sids.NightShadeCloak, "Shade Cloak", spellSchoolValues.Nightshade, 2},
		{sids.NightDeathsGrip, "Umbral Grip", spellSchoolValues.Nightshade, 2},
		{sids.NightWeb, "Web", spellSchoolValues.Nightshade, 2},
		{sids.NightNairasVeil, "Naira's Veil", spellSchoolValues.Nightshade, 3},
		{sids.NightSilence, "Silence", spellSchoolValues.Nightshade, 3},
		{sids.NightSleep, "Sleep", spellSchoolValues.Nightshade, 3},
		{sids.NightTwilight, "Twilight", spellSchoolValues.Nightshade, 3},
		{sids.NightBerserker, "Berserk", spellSchoolValues.Nightshade, 4},
		{sids.NightSummonStarchild, "Summon Starchild", spellSchoolValues.Nightshade, 4},
		{sids.NightVulnerability, "Vulnerability", spellSchoolValues.Nightshade, 4},
		{sids.NightDeathsCall, "Coup de Grace", spellSchoolValues.Nightshade, 5},
		{sids.NightNairasKiss, "Naira's Kiss", spellSchoolValues.Nightshade, 5},
		{sids.NightShadowArmy, "Shadow Army", spellSchoolValues.Nightshade, 5},

		// Primal
		{sids.PrimalGroundsight, "Groundsight", spellSchoolValues.Primal, 1},
		{sids.PrimalThunderbolt, "Lightning Bolt", spellSchoolValues.Primal, 1},
		{sids.PrimalThickHide, "Thick Hide", spellSchoolValues.Primal, 1},
		{sids.PrimalCrystalCrown, "Crystal Crown", spellSchoolValues.Primal, 2},
		{sids.PrimalFireGlobe, "Fireball", spellSchoolValues.Primal, 2},
		{sids.PrimalIceBolt, "Ice Bolt", spellSchoolValues.Primal, 2},
		{sids.PrimalWean, "Wean", spellSchoolValues.Primal, 2},
		{sids.PrimalCaveIn, "Cave In", spellSchoolValues.Primal, 3},
		{sids.PrimalEarthsRage, "Earth's Rage", spellSchoolValues.Primal, 3},
		{sids.PrimalWallOfFlame, "Firewall", spellSchoolValues.Primal, 3},
		{sids.PrimalStoneFangs, "Stone Fangs", spellSchoolValues.Primal, 3},
		{sids.PrimalPrimordialPurity, "Anti-Magic", spellSchoolValues.Primal, 4},
		{sids.PrimalChainLightning, "Chain Lightning", spellSchoolValues.Primal, 4},
		{sids.PrimalAvalanche, "Circle of Winter", spellSchoolValues.Primal, 4},
		{sids.PrimalPrimordialChaos, "Primordial Chaos", spellSchoolValues.Primal, 4},
		{sids.PrimalArmageddon, "Armageddon", spellSchoolValues.Primal, 5},
		{sids.PrimalHksmillasRampage, "Hksmilla's Rampage", spellSchoolValues.Primal, 5},
		{sids.PrimalSummonPrimalRemnant, "Summon Primal Remnant", spellSchoolValues.Primal, 5},

		// Arcane
		{sids.SpaceEarlyStart, "Early Start", spellSchoolValues.Arcane, 1},
		{sids.SpaceEnergyze, "Energize", spellSchoolValues.Arcane, 1},
		{sids.SpaceDecimate, "Guillotine", spellSchoolValues.Arcane, 1},
		{sids.SpaceOpticalIllusion, "Optical Illusion", spellSchoolValues.Arcane, 1},
		{sids.SpaceBlink, "Blink", spellSchoolValues.Arcane, 2},
		{sids.SpaceCarapace, "Carapace", spellSchoolValues.Arcane, 2},
		{sids.SpaceEnergyExplosion, "Energy Explosion", spellSchoolValues.Arcane, 2},
		{sids.SpaceReinforcements, "Reinforcements", spellSchoolValues.Arcane, 2},
		{sids.SpaceAssemble, "Assemble!", spellSchoolValues.Arcane, 3},
		{sids.SpaceImpendingFate, "Impending Fate", spellSchoolValues.Arcane, 3},
		{sids.SpaceShackles, "Shackles", spellSchoolValues.Arcane, 3},
		{sids.SpaceTrapJaws, "Temporal Spheres", spellSchoolValues.Arcane, 3},
		{sids.SpaceMirrorCopy, "Mirror Copy", spellSchoolValues.Arcane, 4},
		{sids.SpaceRewind, "Rewind Life", spellSchoolValues.Arcane, 4},
		{sids.SpaceTrapSnare, "Spatial Snare", spellSchoolValues.Arcane, 4},
		{sids.SpaceBlackHole, "Black Hole", spellSchoolValues.Arcane, 5},
		{sids.SpaceDoreathsTide, "Doreath's Tide", spellSchoolValues.Arcane, 5},
		{sids.SpaceRealityDistortion, "Reality Distortion", spellSchoolValues.Arcane, 5},
	}
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
	for _, spell := range KnownSpells {
		if spell.Sid == sid {
			return spell, true
		}
	}
	return SpellEntry{}, false
}

// GetSpellSchoolColorFromDisplayName maps a school display name to its accent color.
func GetSpellSchoolColorFromDisplayName(displayName string) color.NRGBA {
	switch displayName {
	case "High Neutral":
		return themes.ColorSchoolHighNeutral
	case "Daylight":
		return themes.ColorSchoolDaylight
	case "Nightshade":
		return themes.ColorSchoolNightshade
	case "Arcane":
		return themes.ColorSchoolArcane
	case "Primal":
		return themes.ColorSchoolPrimal
	}
	return themes.ColorAccent
}

// GetSpellSchoolColor maps a school display name to its accent color.
func GetSpellSchoolColor(schoolName string) color.NRGBA {
	spellSchoolValues := registry.GetSpellSchoolTypeValues()
	switch schoolName {
	case spellSchoolValues.HighNeutral:
		return themes.ColorSchoolHighNeutral
	case spellSchoolValues.Daylight:
		return themes.ColorSchoolDaylight
	case spellSchoolValues.Nightshade:
		return themes.ColorSchoolNightshade
	case spellSchoolValues.Arcane:
		return themes.ColorSchoolArcane
	case spellSchoolValues.Primal:
		return themes.ColorSchoolPrimal
	}
	return themes.ColorAccent
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
