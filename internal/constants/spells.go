package constants

// SpellEntry pairs a learnable spell SID with its display name, magic school
// and tier (used for sorting in the spell picker).
type SpellEntry struct {
	Sid    string
	Name   string
	School string
	Tier   int
}

// SpellSchoolOrder gives the canonical display order of the magic schools.
//
//nolint:gochecknoglobals // semantic registry
var SpellSchoolOrder = []string{"neutral", "day", "night", "space", "primal"}

// SpellSchoolDisplayNames maps a school key to its display label.
//
//nolint:gochecknoglobals // semantic registry
var SpellSchoolDisplayNames = map[string]string{
	"neutral": "Neutral",
	"day":     "Day",
	"night":   "Night",
	"space":   "Space",
	"primal":  "Primal",
}

// KnownSpells is the catalog of learnable spells. Grouped logically by school;
// consumers sort by tier within each school as needed.
//
//nolint:gochecknoglobals // semantic registry
var KnownSpells = []SpellEntry{
	// Neutral
	{"neutral_magic_pocket_dimension", "Pocket Dimension", "neutral", 2},
	{"neutral_magic_second_sight", "Second Sight", "neutral", 2},
	{"neutral_magic_shadow_form", "Shadowflight", "neutral", 3},
	{"neutral_magic_town_portal", "Town Portal", "neutral", 3},
	{"neutral_magic_dimension_door", "Dimension Door", "neutral", 4},
	{"neutral_magic_light_gate", "Gate of Light", "neutral", 4},

	// Day
	{"day_2_magic_sharp_edge", "Blessing", "day", 1},
	{"day_3_magic_haste", "Haste", "day", 1},
	{"day_1_magic_healing_water", "Healing Water", "day", 1},
	{"day_5_magic_shorten_shadow", "Shorten Shadow", "day", 1},
	{"day_4_magic_favorable_wind", "Favorable Wind", "day", 2},
	{"day_17_magic_clear_view", "From a Bird's Eye", "day", 2},
	{"day_7_magic_inner_light", "Inner Light", "day", 2},
	{"day_6_magic_cleansing_ray", "Weakening Ray", "day", 2},
	{"day_9_magic_arinas_hymn", "Arina's Touch", "day", 3},
	{"day_11_magic_masterful_parry", "Riposte", "day", 3},
	{"day_10_magic_second_song", "Song of Power", "day", 3},
	{"day_8_magic_taunt", "Taunt", "day", 3},
	{"day_18_magic_farsight", "Clear Fog", "day", 4},
	{"day_13_magic_holy_arms", "Heavenly Blades", "day", 4},
	{"day_12_magic_radiant_armor", "Radiant Armor", "day", 4},
	{"day_14_magic_vengeance", "Vengeance", "day", 4},
	{"day_16_magic_arinas_chosen", "Arina's Chosen", "day", 5},
	{"day_15_magic_judgement", "Judgement", "day", 5},

	// Night
	{"night_4_magic_despair", "Despair", "night", 1},
	{"night_3_magic_enlarge_shadow", "Enlarge Shadow", "night", 1},
	{"night_7_magic_fatal_decay", "Fatal Decay", "night", 1},
	{"night_1_magic_unnatural_calm", "Unnatural Calm", "night", 1},
	{"night_17_magic_read_minds", "Read Minds", "night", 2},
	{"night_5_magic_shade_cloak", "Shade Cloak", "night", 2},
	{"night_6_magic_deaths_grip", "Umbral Grip", "night", 2},
	{"night_2_magic_web", "Web", "night", 2},
	{"night_18_magic_nairas_veil", "Naira's Veil", "night", 3},
	{"night_10_magic_silence", "Silence", "night", 3},
	{"night_8_magic_sleep", "Sleep", "night", 3},
	{"night_9_magic_twilight", "Twilight", "night", 3},
	{"night_13_magic_berserker", "Berserk", "night", 4},
	{"night_12_magic_summon_starchild", "Summon Starchild", "night", 4},
	{"night_11_magic_vulnerability", "Vulnerability", "night", 4},
	{"night_15_magic_deaths_call", "Coup de Grace", "night", 5},
	{"night_14_magic_nairas_kiss", "Naira's Kiss", "night", 5},
	{"night_16_magic_shadow_army", "Shadow Army", "night", 5},

	// Primal
	{"primal_17_magic_groundsight", "Groundsight", "primal", 1},
	{"primal_1_magic_thunderbolt", "Lightning Bolt", "primal", 1},
	{"primal_2_magic_thick_hide", "Thick Hide", "primal", 1},
	{"primal_5_magic_crystal_crown", "Crystal Crown", "primal", 2},
	{"primal_4_magic_fire_globe", "Fireball", "primal", 2},
	{"primal_6_magic_ice_bolt", "Ice Bolt", "primal", 2},
	{"primal_3_magic_wean", "Wean", "primal", 2},
	{"primal_8_magic_cave_in", "Cave In", "primal", 3},
	{"primal_9_magic_earths_rage", "Earth's Rage", "primal", 3},
	{"primal_7_magic_wall_of_flame", "Firewall", "primal", 3},
	{"primal_16_magic_stone_fangs", "Stone Fangs", "primal", 3},
	{"primal_10_magic_primordial_purity", "Anti-Magic", "primal", 4},
	{"primal_12_magic_chain_lightning", "Chain Lightning", "primal", 4},
	{"primal_13_magic_avalanche", "Circle of Winter", "primal", 4},
	{"primal_18_magic_primordial_chaos", "Primordial Chaos", "primal", 4},
	{"primal_11_magic_armageddon", "Armageddon", "primal", 5},
	{"primal_14_magic_hksmillas_rampage", "Hksmilla's Rampage", "primal", 5},
	{"primal_15_magic_summon_primal_remnant", "Summon Primal Remnant", "primal", 5},

	// Space
	{"space_1_magic_early_start", "Early Start", "space", 1},
	{"space_3_magic_energyze", "Energize", "space", 1},
	{"space_11_magic_decimate", "Guillotine", "space", 1},
	{"space_4_magic_optical_illusion", "Optical Illusion", "space", 1},
	{"space_6_magic_blink", "Blink", "space", 2},
	{"space_8_magic_carapace", "Carapace", "space", 2},
	{"space_2_magic_energy_explosion", "Energy Explosion", "space", 2},
	{"space_17_magic_reinforcements", "Reinforcements", "space", 2},
	{"space_18_magic_assemble", "Assemble!", "space", 3},
	{"space_9_magic_impending_fate", "Impending Fate", "space", 3},
	{"space_7_magic_shackles", "Shackles", "space", 3},
	{"space_5_magic_trap_jaws", "Temporal Spheres", "space", 3},
	{"space_10_magic_mirror_copy", "Mirror Copy", "space", 4},
	{"space_12_magic_rewind", "Rewind Life", "space", 4},
	{"space_15_magic_trap_snare", "Spatial Snare", "space", 4},
	{"space_13_magic_black_hole", "Black Hole", "space", 5},
	{"space_14_magic_doreaths_tide", "Doreath's Tide", "space", 5},
	{"space_16_magic_reality_distortion", "Reality Distortion", "space", 5},
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
