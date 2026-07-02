package constants

import (
	"image/color"

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

// SpellSchoolOrder gives the canonical display order of the magic schools.
var SpellSchoolOrder = []string{"neutral", "day", "night", "space", "primal"}

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
	return []SpellEntry{
		// High Neutral
		{sids.NeutralPocketDimension, "Pocket Dimension", "neutral", 2},
		{sids.NeutralSecondSight, "Second Sight", "neutral", 2},
		{sids.NeutralShadowForm, "Shadowflight", "neutral", 3},
		{sids.NeutralTownPortal, "Town Portal", "neutral", 3},
		{sids.NeutralDimensionDoor, "Dimension Door", "neutral", 4},
		{sids.NeutralLightGate, "Gate of Light", "neutral", 4},

		// Daylight
		{sids.DaySharpEdge, "Blessing", "day", 1},
		{sids.DayHaste, "Haste", "day", 1},
		{sids.DayHealingWater, "Healing Water", "day", 1},
		{sids.DayShortenShadow, "Shorten Shadow", "day", 1},
		{sids.DayFavorableWind, "Favorable Wind", "day", 2},
		{sids.DayClearView, "From a Bird's Eye", "day", 2},
		{sids.DayInnerLight, "Inner Light", "day", 2},
		{sids.DayCleansingRay, "Weakening Ray", "day", 2},
		{sids.DayArinasHymn, "Arina's Touch", "day", 3},
		{sids.DayMasterfulParry, "Riposte", "day", 3},
		{sids.DaySecondSong, "Song of Power", "day", 3},
		{sids.DayTaunt, "Taunt", "day", 3},
		{sids.DayFarsight, "Clear Fog", "day", 4},
		{sids.DayHolyArms, "Heavenly Blades", "day", 4},
		{sids.DayRadiantArmor, "Radiant Armor", "day", 4},
		{sids.DayVengeance, "Vengeance", "day", 4},
		{sids.DayArinasChosen, "Arina's Chosen", "day", 5},
		{sids.DayJudgement, "Judgement", "day", 5},

		// Nightshade
		{sids.NightDespair, "Despair", "night", 1},
		{sids.NightEnlargeShadow, "Enlarge Shadow", "night", 1},
		{sids.NightFatalDecay, "Fatal Decay", "night", 1},
		{sids.NightUnnaturalCalm, "Unnatural Calm", "night", 1},
		{sids.NightReadMinds, "Read Minds", "night", 2},
		{sids.NightShadeCloak, "Shade Cloak", "night", 2},
		{sids.NightDeathsGrip, "Umbral Grip", "night", 2},
		{sids.NightWeb, "Web", "night", 2},
		{sids.NightNairasVeil, "Naira's Veil", "night", 3},
		{sids.NightSilence, "Silence", "night", 3},
		{sids.NightSleep, "Sleep", "night", 3},
		{sids.NightTwilight, "Twilight", "night", 3},
		{sids.NightBerserker, "Berserk", "night", 4},
		{sids.NightSummonStarchild, "Summon Starchild", "night", 4},
		{sids.NightVulnerability, "Vulnerability", "night", 4},
		{sids.NightDeathsCall, "Coup de Grace", "night", 5},
		{sids.NightNairasKiss, "Naira's Kiss", "night", 5},
		{sids.NightShadowArmy, "Shadow Army", "night", 5},

		// Primal
		{sids.PrimalGroundsight, "Groundsight", "primal", 1},
		{sids.PrimalThunderbolt, "Lightning Bolt", "primal", 1},
		{sids.PrimalThickHide, "Thick Hide", "primal", 1},
		{sids.PrimalCrystalCrown, "Crystal Crown", "primal", 2},
		{sids.PrimalFireGlobe, "Fireball", "primal", 2},
		{sids.PrimalIceBolt, "Ice Bolt", "primal", 2},
		{sids.PrimalWean, "Wean", "primal", 2},
		{sids.PrimalCaveIn, "Cave In", "primal", 3},
		{sids.PrimalEarthsRage, "Earth's Rage", "primal", 3},
		{sids.PrimalWallOfFlame, "Firewall", "primal", 3},
		{sids.PrimalStoneFangs, "Stone Fangs", "primal", 3},
		{sids.PrimalPrimordialPurity, "Anti-Magic", "primal", 4},
		{sids.PrimalChainLightning, "Chain Lightning", "primal", 4},
		{sids.PrimalAvalanche, "Circle of Winter", "primal", 4},
		{sids.PrimalPrimordialChaos, "Primordial Chaos", "primal", 4},
		{sids.PrimalArmageddon, "Armageddon", "primal", 5},
		{sids.PrimalHksmillasRampage, "Hksmilla's Rampage", "primal", 5},
		{sids.PrimalSummonPrimalRemnant, "Summon Primal Remnant", "primal", 5},

		// Arcane
		{sids.SpaceEarlyStart, "Early Start", "space", 1},
		{sids.SpaceEnergyze, "Energize", "space", 1},
		{sids.SpaceDecimate, "Guillotine", "space", 1},
		{sids.SpaceOpticalIllusion, "Optical Illusion", "space", 1},
		{sids.SpaceBlink, "Blink", "space", 2},
		{sids.SpaceCarapace, "Carapace", "space", 2},
		{sids.SpaceEnergyExplosion, "Energy Explosion", "space", 2},
		{sids.SpaceReinforcements, "Reinforcements", "space", 2},
		{sids.SpaceAssemble, "Assemble!", "space", 3},
		{sids.SpaceImpendingFate, "Impending Fate", "space", 3},
		{sids.SpaceShackles, "Shackles", "space", 3},
		{sids.SpaceTrapJaws, "Temporal Spheres", "space", 3},
		{sids.SpaceMirrorCopy, "Mirror Copy", "space", 4},
		{sids.SpaceRewind, "Rewind Life", "space", 4},
		{sids.SpaceTrapSnare, "Spatial Snare", "space", 4},
		{sids.SpaceBlackHole, "Black Hole", "space", 5},
		{sids.SpaceDoreathsTide, "Doreath's Tide", "space", 5},
		{sids.SpaceRealityDistortion, "Reality Distortion", "space", 5},
	}
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

// GetSpellSchoolColorFromDisplayName maps a school display name to its accent color
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

// GetSpellSchoolColor maps a school display name to its accent color
func GetSpellSchoolColor(schoolName string) color.NRGBA {
	switch schoolName {
	case "neutral":
		return themes.ColorSchoolHighNeutral
	case "day":
		return themes.ColorSchoolDaylight
	case "night":
		return themes.ColorSchoolNightshade
	case "space":
		return themes.ColorSchoolArcane
	case "primal":
		return themes.ColorSchoolPrimal
	}
	return themes.ColorAccent
}
