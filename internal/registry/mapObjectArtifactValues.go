package registry

type objectArtifacts struct {
	artifactScrolls
	artifactRandomItems
	artifactMovementItems
	artifactCombatItems
	artifactDiplomacyItems
	artifactMagicItems
	artifactMiscellaneousItems

	ResonantSphereSet         artifactSetResonantSphere
	TranquilitySet            artifactSetTranquility
	ShamaniacSoulSet          artifactSetShamaniacSoul
	KnightsHonorSet           artifactSetKnightsHonor
	UkhtabarSealSet           artifactSetUkhtabarSeal
	MilosCurseSet             artifactSetMilosCurse
	PaupersGlorySet           artifactSetPaupersGlory
	AngelicAllianceSet        artifactSetAngelicAlliance
	GiftsOfDwarvenLordsSet    artifactSetGiftsOfDwarvenLords
	ElixirOfLifeSet           artifactSetElixirOfLife
	ShadowOfDeathSet          artifactSetShadowOfDeath
	WanderersWaySet           artifactSetWanderersWay
	LivingArrowsSet           artifactSetLivingArrows
	DuelistsPrideSet          artifactSetDuelistsPride
	EtherealKnowledgeSet      artifactSetEtherealKnowledge
	InnerSongSet              artifactSetInnerSong
	PowerOfTheDragonFatherSet artifactSetPowerOfTheDragonFather
	BeelzebubsBlessingSet     artifactSetBeelzebubsBlessing
	BoreolosSet               artifactSetBoreolos
	HolySigilsSet             artifactSetHolySigils
	RuleOfShadowSet           artifactSetRuleOfShadow
	AmbassadorsWordSet        artifactSetAmbassadorsWord
	WarriorsStrengthSet       artifactSetWarriorsStrength
	KeepersFortitudeSet       artifactSetKeepersFortitude
	WizardsMightSet           artifactSetWizardsMight
	ScholarsWisdomSet         artifactSetScholarsWisdom
}

func GetMapObjectAllArtifactValues() objectArtifacts {
	return objectArtifacts{
		artifactScrolls:            GetMapObjectScrollValues(),
		artifactRandomItems:        GetMapObjectRandomItemValues(),
		artifactMovementItems:      GetMapObjectMovementArtifactValues(),
		artifactCombatItems:        GetMapObjectCombatArtifactValues(),
		artifactDiplomacyItems:     GetMapObjectDiplomacyArtifactValues(),
		artifactMagicItems:         GetMapObjectMagicArtifactValues(),
		artifactMiscellaneousItems: GetMapObjectMiscellaneousArtifactValues(),

		ResonantSphereSet:         GetMapObjectSetResonantSphereValues(),
		TranquilitySet:            GetMapObjectSetTranquilityValues(),
		ShamaniacSoulSet:          GetMapObjectSetShamaniacSoulValues(),
		KnightsHonorSet:           GetMapObjectSetKnightsHonorValues(),
		UkhtabarSealSet:           GetMapObjectSetUkhtabarSealValues(),
		MilosCurseSet:             GetMapObjectSetMilosCurseValues(),
		PaupersGlorySet:           GetMapObjectSetPaupersGloryValues(),
		AngelicAllianceSet:        GetMapObjectSetAngelicAllianceValues(),
		GiftsOfDwarvenLordsSet:    GetMapObjectSetGiftsOfDwarvenLordsValues(),
		ElixirOfLifeSet:           GetMapObjectSetElixirOfLifeValues(),
		ShadowOfDeathSet:          GetMapObjectSetShadowOfDeathValues(),
		WanderersWaySet:           GetMapObjectSetWanderersWayValues(),
		LivingArrowsSet:           GetMapObjectSetLivingArrowsValues(),
		DuelistsPrideSet:          GetMapObjectSetDuelistsPrideValues(),
		EtherealKnowledgeSet:      GetMapObjectSetEtherealKnowledgeValues(),
		InnerSongSet:              GetMapObjectSetInnerSongValues(),
		PowerOfTheDragonFatherSet: GetMapObjectSetPowerOfTheDragonFatherValues(),
		BeelzebubsBlessingSet:     GetMapObjectSetBeelzebubsBlessingValues(),
		BoreolosSet:               GetMapObjectSetBoreolosValues(),
		HolySigilsSet:             GetMapObjectSetHolySigilsValues(),
		RuleOfShadowSet:           GetMapObjectSetRuleOfShadowValues(),
		AmbassadorsWordSet:        GetMapObjectSetAmbassadorsWordValues(),
		WarriorsStrengthSet:       GetMapObjectSetWarriorsStrengthValues(),
		KeepersFortitudeSet:       GetMapObjectSetKeepersFortitudeValues(),
		WizardsMightSet:           GetMapObjectSetWizardsMightValues(),
		ScholarsWisdomSet:         GetMapObjectSetScholarsWisdomValues(),
	}
}

type artifactScrolls struct {
	ScrollBox          string
	EnchantedScrollBox string
	MythicScrollBox    string
}

func GetMapObjectScrollValues() artifactScrolls {
	return artifactScrolls{
		ScrollBox:          "scroll_box",
		EnchantedScrollBox: "enchanted_scroll_box",
		MythicScrollBox:    "mythic_scroll_box",
	}
}

type artifactRandomItems struct {
	RandomItemCommon    string
	RandomItemRare      string
	RandomItemEpic      string
	RandomItemLegendary string
}

func GetMapObjectRandomItemValues() artifactRandomItems {
	return artifactRandomItems{
		RandomItemCommon:    "random_item_common",
		RandomItemRare:      "random_item_rare",
		RandomItemEpic:      "random_item_epic",
		RandomItemLegendary: "random_item_legendary",
	}
}

type artifactMovementItems struct {
	PoleStar          string
	SevenLeagueBoots  string
	SwampBoots        string
	WarlordBoots      string
	MagicKeyRing      string
	LegionsStep       string
	FallenAngelWings  string
	BannerOfFourWinds string
	Spyglass          string
}

func GetMapObjectMovementArtifactValues() artifactMovementItems {
	return artifactMovementItems{
		PoleStar:          "pole_star_artifact",
		SevenLeagueBoots:  "seven_league_boots_artifact",
		SwampBoots:        "swamp_boots_artifact",
		WarlordBoots:      "warlord_boots_artifact",
		MagicKeyRing:      "magic_key_ring_artifact",
		LegionsStep:       "legions_step_artifact",
		FallenAngelWings:  "fallen_angel_wings_artifact",
		BannerOfFourWinds: "banner_of_four_winds_artifact",
		Spyglass:          "spyglass_artifact",
	}
}

type artifactCombatItems struct {
	ShacklesOfWar          string
	OgresClubOfHavoc       string
	TarqOfTheRampagingOgre string
	TunicOfTheCyclopsKing  string
	Garotte                string
	HourglassOfProtection  string
	ShoddyShield           string
	EagleArmor             string
	ChainMail              string
	HeadTorch              string
	FineWand               string
	LordsRing              string
}

func GetMapObjectCombatArtifactValues() artifactCombatItems {
	return artifactCombatItems{
		ShacklesOfWar:          "shackles_of_war_artifact",
		OgresClubOfHavoc:       "ogres_club_of_havoc_artifact",
		TarqOfTheRampagingOgre: "tarq_of_the_rampaging_ogre_artifact",
		TunicOfTheCyclopsKing:  "tunic_of_the_cyclops_king_artifact",
		Garotte:                "garotte_artifact",
		HourglassOfProtection:  "hourglass_of_protection_artifact",
		ShoddyShield:           "shoddy_shield_artifact",
		EagleArmor:             "eagle_armor_artifact",
		ChainMail:              "chain_mail_artifact",
		HeadTorch:              "head_torch_artifact",
		FineWand:               "fine_wand_artifact",
		LordsRing:              "lords_ring_artifact",
	}
}

type artifactDiplomacyItems struct {
	VoodooshDoll     string
	FlagOfTruce      string
	RingOfNeutrality string
}

func GetMapObjectDiplomacyArtifactValues() artifactDiplomacyItems {
	return artifactDiplomacyItems{
		VoodooshDoll:     "voodoosh_doll_artifact",
		FlagOfTruce:      "flag_of_truce_artifact",
		RingOfNeutrality: "ring_of_neutrality_artifact",
	}
}

type artifactMagicItems struct {
	CatechismOfNightMagic     string
	CatechismOfDaylightMagic  string
	CatechismOfSpacetimeMagic string
	CatechismOfPrimalMagic    string
	SpellbindersHat           string
	SpellsInABottle           string
	OrbOfInhibition           string
	OrbOfDestruction          string
	SealOfSilence             string
	CrownOfTheSupremeMagi     string
	ClothesOfEnlightenment    string
	CardsDeck                 string
	RunestoneShards           string
}

func GetMapObjectMagicArtifactValues() artifactMagicItems {
	return artifactMagicItems{
		CatechismOfNightMagic:     "catechism_of_night_magic_artifact",
		CatechismOfDaylightMagic:  "catechism_of_daylight_magic_artifact",
		CatechismOfSpacetimeMagic: "catechism_of_spacetime_magic_artifact",
		CatechismOfPrimalMagic:    "catechism_of_primal_magic_artifact",
		SpellbindersHat:           "spellbinders_hat_artifact",
		SpellsInABottle:           "spells_in_a_bottle_artifact",
		OrbOfInhibition:           "orb_of_inhibition_artifact",
		OrbOfDestruction:          "orb_of_destruction_artifact",
		SealOfSilence:             "seal_of_silence_artifact",
		CrownOfTheSupremeMagi:     "crown_of_the_supreme_magi_artifact",
		ClothesOfEnlightenment:    "clothes_of_enlightenment_artifact",
		CardsDeck:                 "cards_deck_artifact",
		RunestoneShards:           "runestone_shards_artifact",
	}
}

type artifactMiscellaneousItems struct {
	GoldenGooseEgg  string
	TacticalGuide   string
	EndlessBag      string
	SoullessSash    string
	MonsterHead     string
	Omencaller      string
	SixthFinger     string
	SoulscallerRing string
	ChainLink       string
	DemonicHeart    string
	TwoFacedMask    string
	AncientIdol     string
	Excalibur       string
	Caduceus        string
}

func GetMapObjectMiscellaneousArtifactValues() artifactMiscellaneousItems {
	return artifactMiscellaneousItems{
		GoldenGooseEgg:  "golden_goose_egg_artifact",
		TacticalGuide:   "tactical_guide_artifact",
		EndlessBag:      "endless_bag_artifact",
		SoullessSash:    "soulless_sash_artifact",
		MonsterHead:     "monster_head_artifact",
		Omencaller:      "omencaller_artifact",
		SixthFinger:     "sixth_finger_artifact",
		SoulscallerRing: "soulscaller_ring_artifact",
		ChainLink:       "chain_link_artifact",
		DemonicHeart:    "demonic_heart_artifact",
		TwoFacedMask:    "two_faced_mask_artifact",
		AncientIdol:     "ancient_idol_artifact",
		Excalibur:       "excalibur_artifact",
		Caduceus:        "caduceus_artifact",
	}
}

type artifactSetResonantSphere struct {
	OrbOfTwilight string
	OrbOfDaylight string
	OrbOfEternity string
	PrimalOrb     string
}

func GetMapObjectSetResonantSphereValues() artifactSetResonantSphere {
	return artifactSetResonantSphere{
		OrbOfTwilight: "resonant_sphere_orb_of_twilight_artifact",
		OrbOfDaylight: "resonant_sphere_orb_of_daylight_artifact",
		OrbOfEternity: "resonant_sphere_orb_of_eternity_artifact",
		PrimalOrb:     "resonant_sphere_primal_orb_artifact",
	}
}

type artifactSetTranquility struct {
	BrightmindTiara string
	MagicMirror     string
	RingOfSerenity  string
}

func GetMapObjectSetTranquilityValues() artifactSetTranquility {
	return artifactSetTranquility{
		BrightmindTiara: "tranquility_brightmind_tiara_artifact",
		MagicMirror:     "tranquility_magic_mirror_artifact",
		RingOfSerenity:  "tranquility_ring_of_serenity_artifact",
	}
}

type artifactSetShamaniacSoul struct {
	ShamanStaff     string
	IridescentCloak string
	GemwoodMask     string
	ClutchingRing   string
}

func GetMapObjectSetShamaniacSoulValues() artifactSetShamaniacSoul {
	return artifactSetShamaniacSoul{
		ShamanStaff:     "shamaniac_soul_shaman_staff_artifact",
		IridescentCloak: "shamaniac_soul_iridescent_cloak_artifact",
		GemwoodMask:     "shamaniac_soul_gemwood_mask_artifact",
		ClutchingRing:   "shamaniac_soul_clutching_ring_artifact",
	}
}

type artifactSetKnightsHonor struct {
	DrumsOfWar  string
	Lance       string
	Misericorde string
	PlateArmor  string
	Armet       string
}

func GetMapObjectSetKnightsHonorValues() artifactSetKnightsHonor {
	return artifactSetKnightsHonor{
		DrumsOfWar:  "knights_honor_drums_of_war_artifact",
		Lance:       "knights_honor_lance_artifact",
		Misericorde: "knights_honor_misericorde_artifact",
		PlateArmor:  "knights_honor_plate_armor_artifact",
		Armet:       "knights_honor_armet_artifact",
	}
}

type artifactSetUkhtabarSeal struct {
	UkhSeal   string
	TabarSeal string
}

func GetMapObjectSetUkhtabarSealValues() artifactSetUkhtabarSeal {
	return artifactSetUkhtabarSeal{
		UkhSeal:   "ukhtabar_seal_ukh_seal_artifact",
		TabarSeal: "ukhtabar_seal_tabar_seal_artifact",
	}
}

type artifactSetMilosCurse struct {
	GoldenPig    string
	GoldenMoth   string
	SkullOfMilos string
}

func GetMapObjectSetMilosCurseValues() artifactSetMilosCurse {
	return artifactSetMilosCurse{
		GoldenPig:    "milos_curse_golden_pig_artifact",
		GoldenMoth:   "milos_curse_golden_moth_artifact",
		SkullOfMilos: "milos_curse_skull_of_milos_artifact",
	}
}

type artifactSetPaupersGlory struct {
	WoodenRing string
	StrawHat   string
	RopeBelt   string
	Rags       string
	DumbClub   string
	LastCoin   string
}

func GetMapObjectSetPaupersGloryValues() artifactSetPaupersGlory {
	return artifactSetPaupersGlory{
		WoodenRing: "paupers_glory_wooden_ring_artifact",
		StrawHat:   "paupers_glory_straw_hat_artifact",
		RopeBelt:   "paupers_glory_rope_belt_artifact",
		Rags:       "paupers_glory_rags_artifact",
		DumbClub:   "paupers_glory_dumb_club_artifact",
		LastCoin:   "paupers_glory_last_coin_artifact",
	}
}

type artifactSetAngelicAlliance struct {
	SwordOfJudgement            string
	CelestialSashOfBliss        string
	LionsShieldOfCourage        string
	ArmorOfWonder               string
	HelmOfHeavenlyEnlightenment string
	SandalsOfTheSaint           string
}

func GetMapObjectSetAngelicAllianceValues() artifactSetAngelicAlliance {
	return artifactSetAngelicAlliance{
		SwordOfJudgement:            "angelic_alliance_sword_of_judgement_artifact",
		CelestialSashOfBliss:        "angelic_alliance_celestial_sash_of_bliss_artifact",
		LionsShieldOfCourage:        "angelic_alliance_lions_shield_of_courage_artifact",
		ArmorOfWonder:               "angelic_alliance_armor_of_wonder_artifact",
		HelmOfHeavenlyEnlightenment: "angelic_alliance_helm_of_heavenly_enlightenment_artifact",
		SandalsOfTheSaint:           "angelic_alliance_sandals_of_the_saint_artifact",
	}
}

type artifactSetGiftsOfDwarvenLords struct {
	AutomatedAntimagicShield   string
	ProtectiveBelt             string
	CrimsonResonanceController string
	EmeraldResonanceController string
}

func GetMapObjectSetGiftsOfDwarvenLordsValues() artifactSetGiftsOfDwarvenLords {
	return artifactSetGiftsOfDwarvenLords{
		AutomatedAntimagicShield: "gifts_of_dwarven_lords_automated_antimagic_shield_artifact",
		// AutomatedAntimagicShieldAlt:   "gifts_of_dwarven_lords_automated_antimagic_shield_artifact_alt",
		ProtectiveBelt: "gifts_of_dwarven_lords_protective_belt_artifact",
		// ProtectiveBeltAlt:             "gifts_of_dwarven_lords_protective_belt_artifact_alt",
		CrimsonResonanceController: "gifts_of_dwarven_lords_crimson_resonance_controller_artifact",
		// CrimsonResonanceControllerAlt: "gifts_of_dwarven_lords_crimson_resonance_controller_artifact_alt",
		EmeraldResonanceController: "gifts_of_dwarven_lords_emerald_resonance_controller_artifact",
		// EmeraldResonanceControllerAlt: "gifts_of_dwarven_lords_emerald_resonance_controller_artifact_alt",
	}
}

type artifactSetElixirOfLife struct {
	FlaskOfOblivion string
	LifebloodFairy  string
	RingOfLife      string
}

func GetMapObjectSetElixirOfLifeValues() artifactSetElixirOfLife {
	return artifactSetElixirOfLife{
		FlaskOfOblivion: "elixir_of_life_flask_of_oblivion_artifact",
		LifebloodFairy:  "elixir_of_life_lifeblood_fairy_artifact",
		RingOfLife:      "elixir_of_life_ring_of_life_artifact",
	}
}

type artifactSetShadowOfDeath struct {
	CursedArmor string
	BoneBoots   string
	SecondShade string
	DarkHatchet string
}

func GetMapObjectSetShadowOfDeathValues() artifactSetShadowOfDeath {
	return artifactSetShadowOfDeath{
		CursedArmor: "shadow_of_death_cursed_armor_artifact",
		BoneBoots:   "shadow_of_death_bone_boots_artifact",
		SecondShade: "shadow_of_death_second_shade_artifact",
		DarkHatchet: "shadow_of_death_dark_hatchet_artifact",
	}
}

type artifactSetWanderersWay struct {
	BootsOfTravel string
	Backpack      string
}

func GetMapObjectSetWanderersWayValues() artifactSetWanderersWay {
	return artifactSetWanderersWay{
		BootsOfTravel: "wanderers_way_boots_of_travel_artifact",
		Backpack:      "wanderers_way_backpack_artifact",
	}
}

type artifactSetLivingArrows struct {
	ShroomwoodBow      string
	LightAndShadeCloak string
	QuiveringQuiver    string
}

func GetMapObjectSetLivingArrowsValues() artifactSetLivingArrows {
	return artifactSetLivingArrows{
		ShroomwoodBow:      "living_arrows_shroomwood_bow_artifact",
		LightAndShadeCloak: "living_arrows_light_and_shade_cloak_artifact",
		QuiveringQuiver:    "living_arrows_quivering_quiver_artifact",
	}
}

type artifactSetDuelistsPride struct {
	Rapier        string
	Buckler       string
	BrassKnuckles string
}

func GetMapObjectSetDuelistsPrideValues() artifactSetDuelistsPride {
	return artifactSetDuelistsPride{
		Rapier:        "duelists_pride_rapier_artifact",
		Buckler:       "duelists_pride_buckler_artifact",
		BrassKnuckles: "duelists_pride_brass_knuckles_artifact",
	}
}

type artifactSetEtherealKnowledge struct {
	GlassDagger string
	MirrorShoes string
	VortexDress string
	ThirdEye    string
}

func GetMapObjectSetEtherealKnowledgeValues() artifactSetEtherealKnowledge {
	return artifactSetEtherealKnowledge{
		GlassDagger: "ethereal_knowledge_glass_dagger_artifact",
		MirrorShoes: "ethereal_knowledge_mirror_shoes_artifact",
		VortexDress: "ethereal_knowledge_vortex_dress_artifact",
		ThirdEye:    "ethereal_knowledge_third_eye_artifact",
	}
}

type artifactSetInnerSong struct {
	MusicSheet     string
	SingingPanPipe string
	FancyMask      string
}

func GetMapObjectSetInnerSongValues() artifactSetInnerSong {
	return artifactSetInnerSong{
		MusicSheet:     "inner_song_music_sheet_artifact",
		SingingPanPipe: "inner_song_singing_pan_pipe_artifact",
		FancyMask:      "inner_song_fancy_mask_artifact",
	}
}

type artifactSetPowerOfTheDragonFather struct {
	RedDragonFlameTongue string
	DragonScaleShield    string
	DragonScaleArmor     string
	DragonCrest          string
	DragonBoneGreaves    string
	SlitheringSash       string
	DragonWing           string
	PiercingEyeOfADragon string
}

func GetMapObjectSetPowerOfTheDragonFatherValues() artifactSetPowerOfTheDragonFather {
	return artifactSetPowerOfTheDragonFather{
		RedDragonFlameTongue: "power_of_the_dragon_father_red_dragon_flame_tongue_artifact",
		DragonScaleShield:    "power_of_the_dragon_father_dragon_scale_shield_artifact",
		DragonScaleArmor:     "power_of_the_dragon_father_dragon_scale_armor_artifact",
		DragonCrest:          "power_of_the_dragon_father_dragon_crest_artifact",
		DragonBoneGreaves:    "power_of_the_dragon_father_dragonbone_greaves_artifact",
		SlitheringSash:       "power_of_the_dragon_father_slithering_sash_artifact",
		DragonWing:           "power_of_the_dragon_father_dragon_wing_artifact",
		PiercingEyeOfADragon: "power_of_the_dragon_father_piercing_eye_of_a_dragon_artifact",
	}
}

type artifactSetBeelzebubsBlessing struct {
	DemonClaw       string
	ChitinousShield string
	Heartbeat       string
	DemonCrest      string
}

func GetMapObjectSetBeelzebubsBlessingValues() artifactSetBeelzebubsBlessing {
	return artifactSetBeelzebubsBlessing{
		DemonClaw:       "beelzebubs_blessing_demon_claw_artifact",
		ChitinousShield: "beelzebubs_blessing_chitinous_shield_artifact",
		Heartbeat:       "beelzebubs_blessing_heartbeat_artifact",
		DemonCrest:      "beelzebubs_blessing_demon_crest_artifact",
	}
}

type artifactSetBoreolos struct {
	Hand  string
	Foot  string
	Heart string
	Head  string
}

func GetMapObjectSetBoreolosValues() artifactSetBoreolos {
	return artifactSetBoreolos{
		Hand:  "boreolos_hand_artifact",
		Foot:  "boreolos_foot_artifact",
		Heart: "boreolos_heart_artifact",
		Head:  "boreolos_head_artifact",
	}
}

type artifactSetHolySigils struct {
	SigilOfRoph         string
	SigilOfEridore      string
	SigilOfMearea       string
	SigilOfInsara       string
	SigilOfQuix         string
	SigilOfTheSevenMagi string
	SigilOfTheSecondMan string
	SigilOfUurdt        string
}

func GetMapObjectSetHolySigilsValues() artifactSetHolySigils {
	return artifactSetHolySigils{
		SigilOfRoph:         "holy_sigil_of_roph_artifact",
		SigilOfEridore:      "holy_sigil_of_eridore_artifact",
		SigilOfMearea:       "holy_sigil_of_mearea_artifact",
		SigilOfInsara:       "holy_sigil_of_insara_artifact",
		SigilOfQuix:         "holy_sigil_of_quix_artifact",
		SigilOfTheSevenMagi: "holy_sigil_of_the_seven_magi_artifact",
		SigilOfTheSecondMan: "holy_sigil_of_the_second_man_artifact",
		SigilOfUurdt:        "holy_sigil_of_uurdt_artifact",
	}
}

type artifactSetRuleOfShadow struct {
	LiquidSilence  string
	TheTruthmaker  string
	TheTruthseeker string
	NostriasGaze   string
}

func GetMapObjectSetRuleOfShadowValues() artifactSetRuleOfShadow {
	return artifactSetRuleOfShadow{
		LiquidSilence:  "rule_of_shadow_liquid_silence_artifact",
		TheTruthmaker:  "rule_of_shadow_the_truthmaker_artifact",
		TheTruthseeker: "rule_of_shadow_the_truthseeker_artifact",
		NostriasGaze:   "rule_of_shadow_nostrias_gaze_artifact",
	}
}

type artifactSetAmbassadorsWord struct {
	DiplomaticGifts string
	AmbassadorsSash string
}

func GetMapObjectSetAmbassadorsWordValues() artifactSetAmbassadorsWord {
	return artifactSetAmbassadorsWord{
		DiplomaticGifts: "ambassadors_word_diplomatic_gifts_artifact",
		AmbassadorsSash: "ambassadors_word_ambassadors_sash_artifact",
	}
}

type artifactSetWarriorsStrength struct {
	WarriorsBelt     string
	WarriorsOberegus string
}

func GetMapObjectSetWarriorsStrengthValues() artifactSetWarriorsStrength {
	return artifactSetWarriorsStrength{
		WarriorsBelt:     "warriors_strength_warriors_belt_artifact",
		WarriorsOberegus: "warriors_strength_warriors_oberegus_artifact",
	}
}

type artifactSetKeepersFortitude struct {
	KeepersRing     string
	KeepersOberegus string
}

func GetMapObjectSetKeepersFortitudeValues() artifactSetKeepersFortitude {
	return artifactSetKeepersFortitude{
		KeepersRing:     "keepers_fortitude_keepers_ring_artifact",
		KeepersOberegus: "keepers_fortitude_keepers_oberegus_artifact",
	}
}

type artifactSetWizardsMight struct {
	WizardsCloak    string
	WizardsOberegus string
}

func GetMapObjectSetWizardsMightValues() artifactSetWizardsMight {
	return artifactSetWizardsMight{
		WizardsCloak:    "wizards_might_wizards_cloak_artifact",
		WizardsOberegus: "wizards_might_wizards_oberegus_artifact",
	}
}

type artifactSetScholarsWisdom struct {
	ScholarsTiara    string
	ScholarsOberegus string
}

func GetMapObjectSetScholarsWisdomValues() artifactSetScholarsWisdom {
	return artifactSetScholarsWisdom{
		ScholarsTiara:    "scholars_wisdom_scholars_tiara_artifact",
		ScholarsOberegus: "scholars_wisdom_scholars_oberegus_artifact",
	}
}
