package registry

type spellSids struct {
	// High Neutral
	NeutralPocketDimension string
	NeutralSecondSight     string
	NeutralShadowForm      string
	NeutralTownPortal      string
	NeutralDimensionDoor   string
	NeutralLightGate       string

	// Daylight
	DaySharpEdge      string
	DayHaste          string
	DayHealingWater   string
	DayShortenShadow  string
	DayFavorableWind  string
	DayClearView      string
	DayInnerLight     string
	DayCleansingRay   string
	DayArinasHymn     string
	DayMasterfulParry string
	DaySecondSong     string
	DayTaunt          string
	DayFarsight       string
	DayHolyArms       string
	DayRadiantArmor   string
	DayVengeance      string
	DayArinasChosen   string
	DayJudgement      string

	// Nightshade
	NightDespair         string
	NightEnlargeShadow   string
	NightFatalDecay      string
	NightUnnaturalCalm   string
	NightReadMinds       string
	NightShadeCloak      string
	NightDeathsGrip      string
	NightWeb             string
	NightNairasVeil      string
	NightSilence         string
	NightSleep           string
	NightTwilight        string
	NightBerserker       string
	NightSummonStarchild string
	NightVulnerability   string
	NightDeathsCall      string
	NightNairasKiss      string
	NightShadowArmy      string

	// Primal
	PrimalGroundsight         string
	PrimalThunderbolt         string
	PrimalThickHide           string
	PrimalCrystalCrown        string
	PrimalFireGlobe           string
	PrimalIceBolt             string
	PrimalWean                string
	PrimalCaveIn              string
	PrimalEarthsRage          string
	PrimalWallOfFlame         string
	PrimalStoneFangs          string
	PrimalPrimordialPurity    string
	PrimalChainLightning      string
	PrimalAvalanche           string
	PrimalPrimordialChaos     string
	PrimalArmageddon          string
	PrimalHksmillasRampage    string
	PrimalSummonPrimalRemnant string

	// Arcane
	SpaceEarlyStart        string
	SpaceEnergyze          string
	SpaceDecimate          string
	SpaceOpticalIllusion   string
	SpaceBlink             string
	SpaceCarapace          string
	SpaceEnergyExplosion   string
	SpaceReinforcements    string
	SpaceAssemble          string
	SpaceImpendingFate     string
	SpaceShackles          string
	SpaceTrapJaws          string
	SpaceMirrorCopy        string
	SpaceRewind            string
	SpaceTrapSnare         string
	SpaceBlackHole         string
	SpaceDoreathsTide      string
	SpaceRealityDistortion string
}

var spellSidValues = spellSids{
	// High Neutral
	NeutralPocketDimension: "neutral_magic_pocket_dimension",
	NeutralSecondSight:     "neutral_magic_second_sight",
	NeutralShadowForm:      "neutral_magic_shadow_form",
	NeutralTownPortal:      "neutral_magic_town_portal",
	NeutralDimensionDoor:   "neutral_magic_dimension_door",
	NeutralLightGate:       "neutral_magic_light_gate",

	// Daylight
	DaySharpEdge:      "day_2_magic_sharp_edge",
	DayHaste:          "day_3_magic_haste",
	DayHealingWater:   "day_1_magic_healing_water",
	DayShortenShadow:  "day_5_magic_shorten_shadow",
	DayFavorableWind:  "day_4_magic_favorable_wind",
	DayClearView:      "day_17_magic_clear_view",
	DayInnerLight:     "day_7_magic_inner_light",
	DayCleansingRay:   "day_6_magic_cleansing_ray",
	DayArinasHymn:     "day_9_magic_arinas_hymn",
	DayMasterfulParry: "day_11_magic_masterful_parry",
	DaySecondSong:     "day_10_magic_second_song",
	DayTaunt:          "day_8_magic_taunt",
	DayFarsight:       "day_18_magic_farsight",
	DayHolyArms:       "day_13_magic_holy_arms",
	DayRadiantArmor:   "day_12_magic_radiant_armor",
	DayVengeance:      "day_14_magic_vengeance",
	DayArinasChosen:   "day_16_magic_arinas_chosen",
	DayJudgement:      "day_15_magic_judgement",

	// Nightshade
	NightDespair:         "night_4_magic_despair",
	NightEnlargeShadow:   "night_3_magic_enlarge_shadow",
	NightFatalDecay:      "night_7_magic_fatal_decay",
	NightUnnaturalCalm:   "night_1_magic_unnatural_calm",
	NightReadMinds:       "night_17_magic_read_minds",
	NightShadeCloak:      "night_5_magic_shade_cloak",
	NightDeathsGrip:      "night_6_magic_deaths_grip",
	NightWeb:             "night_2_magic_web",
	NightNairasVeil:      "night_18_magic_nairas_veil",
	NightSilence:         "night_10_magic_silence",
	NightSleep:           "night_8_magic_sleep",
	NightTwilight:        "night_9_magic_twilight",
	NightBerserker:       "night_13_magic_berserker",
	NightSummonStarchild: "night_12_magic_summon_starchild",
	NightVulnerability:   "night_11_magic_vulnerability",
	NightDeathsCall:      "night_15_magic_deaths_call",
	NightNairasKiss:      "night_14_magic_nairas_kiss",
	NightShadowArmy:      "night_16_magic_shadow_army",

	// Primal
	PrimalGroundsight:         "primal_17_magic_groundsight",
	PrimalThunderbolt:         "primal_1_magic_thunderbolt",
	PrimalThickHide:           "primal_2_magic_thick_hide",
	PrimalCrystalCrown:        "primal_5_magic_crystal_crown",
	PrimalFireGlobe:           "primal_4_magic_fire_globe",
	PrimalIceBolt:             "primal_6_magic_ice_bolt",
	PrimalWean:                "primal_3_magic_wean",
	PrimalCaveIn:              "primal_8_magic_cave_in",
	PrimalEarthsRage:          "primal_9_magic_earths_rage",
	PrimalWallOfFlame:         "primal_7_magic_wall_of_flame",
	PrimalStoneFangs:          "primal_16_magic_stone_fangs",
	PrimalPrimordialPurity:    "primal_10_magic_primordial_purity",
	PrimalChainLightning:      "primal_12_magic_chain_lightning",
	PrimalAvalanche:           "primal_13_magic_avalanche",
	PrimalPrimordialChaos:     "primal_18_magic_primordial_chaos",
	PrimalArmageddon:          "primal_11_magic_armageddon",
	PrimalHksmillasRampage:    "primal_14_magic_hksmillas_rampage",
	PrimalSummonPrimalRemnant: "primal_15_magic_summon_primal_remnant",

	// Arcane
	SpaceEarlyStart:        "space_1_magic_early_start",
	SpaceEnergyze:          "space_3_magic_energyze",
	SpaceDecimate:          "space_11_magic_decimate",
	SpaceOpticalIllusion:   "space_4_magic_optical_illusion",
	SpaceBlink:             "space_6_magic_blink",
	SpaceCarapace:          "space_8_magic_carapace",
	SpaceEnergyExplosion:   "space_2_magic_energy_explosion",
	SpaceReinforcements:    "space_17_magic_reinforcements",
	SpaceAssemble:          "space_18_magic_assemble",
	SpaceImpendingFate:     "space_9_magic_impending_fate",
	SpaceShackles:          "space_7_magic_shackles",
	SpaceTrapJaws:          "space_5_magic_trap_jaws",
	SpaceMirrorCopy:        "space_10_magic_mirror_copy",
	SpaceRewind:            "space_12_magic_rewind",
	SpaceTrapSnare:         "space_15_magic_trap_snare",
	SpaceBlackHole:         "space_13_magic_black_hole",
	SpaceDoreathsTide:      "space_14_magic_doreaths_tide",
	SpaceRealityDistortion: "space_16_magic_reality_distortion",
}

// GetSpellSidValues returns the learnable spell SIDs used for
//
//	gameRules.globalBans.magics
//	gameRules.bonuses.parameters
func GetSpellSidValues() spellSids {
	return spellSidValues
}
