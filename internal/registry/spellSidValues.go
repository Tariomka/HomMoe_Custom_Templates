package registry

type spellSids struct {
	HighNeutral highNeutralSpellSids
	Daylight    daylightSpellSids
	Nightshade  nightshadeSpellSids
	Primal      primalSpellSids
	Arcane      arcaneSpellSids
}

// GetSpellSidValues returns the learnable spell SIDs used for
//
//	gameRules.globalBans.magics
//	gameRules.bonuses.parameters
func GetSpellSidValues() spellSids {
	return spellSids{
		HighNeutral: GetHighNeutralSpellSidValues(),
		Daylight:    GetDaylightSpellSidValues(),
		Nightshade:  GetNightshadeSpellSidValues(),
		Primal:      GetPrimalSpellSidValues(),
		Arcane:      GetArcaneSpellSidValues(),
	}
}

type highNeutralSpellSids struct {
	PocketDimension string
	SecondSight     string
	ShadowForm      string
	TownPortal      string
	DimensionDoor   string
	LightGate       string
}

func GetHighNeutralSpellSidValues() highNeutralSpellSids {
	return highNeutralSpellSids{
		PocketDimension: "neutral_magic_pocket_dimension",
		SecondSight:     "neutral_magic_second_sight",
		ShadowForm:      "neutral_magic_shadow_form",
		TownPortal:      "neutral_magic_town_portal",
		DimensionDoor:   "neutral_magic_dimension_door",
		LightGate:       "neutral_magic_light_gate",
	}
}

type daylightSpellSids struct {
	SharpEdge      string
	Haste          string
	HealingWater   string
	ShortenShadow  string
	FavorableWind  string
	ClearView      string
	InnerLight     string
	CleansingRay   string
	ArinasHymn     string
	MasterfulParry string
	SecondSong     string
	Taunt          string
	Farsight       string
	HolyArms       string
	RadiantArmor   string
	Vengeance      string
	ArinasChosen   string
	Judgement      string
}

func GetDaylightSpellSidValues() daylightSpellSids {
	return daylightSpellSids{
		SharpEdge:      "day_2_magic_sharp_edge",
		Haste:          "day_3_magic_haste",
		HealingWater:   "day_1_magic_healing_water",
		ShortenShadow:  "day_5_magic_shorten_shadow",
		FavorableWind:  "day_4_magic_favorable_wind",
		ClearView:      "day_17_magic_clear_view",
		InnerLight:     "day_7_magic_inner_light",
		CleansingRay:   "day_6_magic_cleansing_ray",
		ArinasHymn:     "day_9_magic_arinas_hymn",
		MasterfulParry: "day_11_magic_masterful_parry",
		SecondSong:     "day_10_magic_second_song",
		Taunt:          "day_8_magic_taunt",
		Farsight:       "day_18_magic_farsight",
		HolyArms:       "day_13_magic_holy_arms",
		RadiantArmor:   "day_12_magic_radiant_armor",
		Vengeance:      "day_14_magic_vengeance",
		ArinasChosen:   "day_16_magic_arinas_chosen",
		Judgement:      "day_15_magic_judgement",
	}
}

type nightshadeSpellSids struct {
	Despair         string
	EnlargeShadow   string
	FatalDecay      string
	UnnaturalCalm   string
	ReadMinds       string
	ShadeCloak      string
	DeathsGrip      string
	Web             string
	NairasVeil      string
	Silence         string
	Sleep           string
	Twilight        string
	Berserker       string
	SummonStarchild string
	Vulnerability   string
	DeathsCall      string
	NairasKiss      string
	ShadowArmy      string
}

func GetNightshadeSpellSidValues() nightshadeSpellSids {
	return nightshadeSpellSids{
		Despair:         "night_4_magic_despair",
		EnlargeShadow:   "night_3_magic_enlarge_shadow",
		FatalDecay:      "night_7_magic_fatal_decay",
		UnnaturalCalm:   "night_1_magic_unnatural_calm",
		ReadMinds:       "night_17_magic_read_minds",
		ShadeCloak:      "night_5_magic_shade_cloak",
		DeathsGrip:      "night_6_magic_deaths_grip",
		Web:             "night_2_magic_web",
		NairasVeil:      "night_18_magic_nairas_veil",
		Silence:         "night_10_magic_silence",
		Sleep:           "night_8_magic_sleep",
		Twilight:        "night_9_magic_twilight",
		Berserker:       "night_13_magic_berserker",
		SummonStarchild: "night_12_magic_summon_starchild",
		Vulnerability:   "night_11_magic_vulnerability",
		DeathsCall:      "night_15_magic_deaths_call",
		NairasKiss:      "night_14_magic_nairas_kiss",
		ShadowArmy:      "night_16_magic_shadow_army",
	}
}

type primalSpellSids struct {
	Groundsight         string
	Thunderbolt         string
	ThickHide           string
	CrystalCrown        string
	FireGlobe           string
	IceBolt             string
	Wean                string
	CaveIn              string
	EarthsRage          string
	WallOfFlame         string
	StoneFangs          string
	PrimordialPurity    string
	ChainLightning      string
	Avalanche           string
	PrimordialChaos     string
	Armageddon          string
	HksmillasRampage    string
	SummonPrimalRemnant string
}

func GetPrimalSpellSidValues() primalSpellSids {
	return primalSpellSids{
		Groundsight:         "primal_17_magic_groundsight",
		Thunderbolt:         "primal_1_magic_thunderbolt",
		ThickHide:           "primal_2_magic_thick_hide",
		CrystalCrown:        "primal_5_magic_crystal_crown",
		FireGlobe:           "primal_4_magic_fire_globe",
		IceBolt:             "primal_6_magic_ice_bolt",
		Wean:                "primal_3_magic_wean",
		CaveIn:              "primal_8_magic_cave_in",
		EarthsRage:          "primal_9_magic_earths_rage",
		WallOfFlame:         "primal_7_magic_wall_of_flame",
		StoneFangs:          "primal_16_magic_stone_fangs",
		PrimordialPurity:    "primal_10_magic_primordial_purity",
		ChainLightning:      "primal_12_magic_chain_lightning",
		Avalanche:           "primal_13_magic_avalanche",
		PrimordialChaos:     "primal_18_magic_primordial_chaos",
		Armageddon:          "primal_11_magic_armageddon",
		HksmillasRampage:    "primal_14_magic_hksmillas_rampage",
		SummonPrimalRemnant: "primal_15_magic_summon_primal_remnant",
	}
}

type arcaneSpellSids struct {
	EarlyStart        string
	Energyze          string
	Decimate          string
	OpticalIllusion   string
	Blink             string
	Carapace          string
	EnergyExplosion   string
	Reinforcements    string
	Assemble          string
	ImpendingFate     string
	Shackles          string
	TrapJaws          string
	MirrorCopy        string
	Rewind            string
	TrapSnare         string
	BlackHole         string
	DoreathsTide      string
	RealityDistortion string
}

func GetArcaneSpellSidValues() arcaneSpellSids {
	return arcaneSpellSids{
		EarlyStart:        "space_1_magic_early_start",
		Energyze:          "space_3_magic_energyze",
		Decimate:          "space_11_magic_decimate",
		OpticalIllusion:   "space_4_magic_optical_illusion",
		Blink:             "space_6_magic_blink",
		Carapace:          "space_8_magic_carapace",
		EnergyExplosion:   "space_2_magic_energy_explosion",
		Reinforcements:    "space_17_magic_reinforcements",
		Assemble:          "space_18_magic_assemble",
		ImpendingFate:     "space_9_magic_impending_fate",
		Shackles:          "space_7_magic_shackles",
		TrapJaws:          "space_5_magic_trap_jaws",
		MirrorCopy:        "space_10_magic_mirror_copy",
		Rewind:            "space_12_magic_rewind",
		TrapSnare:         "space_15_magic_trap_snare",
		BlackHole:         "space_13_magic_black_hole",
		DoreathsTide:      "space_14_magic_doreaths_tide",
		RealityDistortion: "space_16_magic_reality_distortion",
	}
}
