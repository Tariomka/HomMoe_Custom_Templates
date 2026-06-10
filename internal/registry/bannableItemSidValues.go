package registry

type bannableItemSids struct {
	// Movement
	PoleStar          string
	SevenLeagueBoots  string
	SwampBoots        string
	WarlordBoots      string
	MagicKeyRing      string
	LegionsStep       string
	FallenAngelWings  string
	BannerOfFourWinds string
	Spyglass          string

	// Diplomacy
	VoodooshDoll     string
	FlagOfTruce      string
	RingOfNeutrality string

	// Combat
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

	// Magic
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

	// Misc (standalone)
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

	// Set: Resonant Sphere
	ResonantSphereOrbOfTwilight string
	ResonantSphereOrbOfDaylight string
	ResonantSphereOrbOfEternity string
	ResonantSpherePrimalOrb     string

	// Set: Tranquility
	TranquilityBrightmindTiara string
	TranquilityMagicMirror     string
	TranquilityRingOfSerenity  string

	// Set: Shamaniac Soul
	ShamaniacSoulShamanStaff     string
	ShamaniacSoulIridescentCloak string
	ShamaniacSoulGemwoodMask     string
	ShamaniacSoulClutchingRing   string

	// Set: Knight's Honor
	KnightsHonorDrumsOfWar  string
	KnightsHonorLance       string
	KnightsHonorMisericorde string
	KnightsHonorPlateArmor  string
	KnightsHonorArmet       string

	// Set: Ukhtabar Seal
	UkhtabarSealUkhSeal   string
	UkhtabarSealTabarSeal string

	// Set: Milo's Curse
	MilosCurseGoldenPig    string
	MilosCurseGoldenMoth   string
	MilosCurseSkullOfMilos string

	// Set: Pauper's Glory
	PaupersGloryWoodenRing string
	PaupersGloryStrawHat   string
	PaupersGloryRopeBelt   string
	PaupersGloryRags       string
	PaupersGloryDumbClub   string
	PaupersGloryLastCoin   string

	// Set: Angelic Alliance
	AngelicAllianceSwordOfJudgement            string
	AngelicAllianceCelestialSashOfBliss        string
	AngelicAllianceLionsShieldOfCourage        string
	AngelicAllianceArmorOfWonder               string
	AngelicAllianceHelmOfHeavenlyEnlightenment string
	AngelicAllianceSandalsOfTheSaint           string

	// Set: Gifts of Dwarven Lords
	GiftsOfDwarvenLordsAutomatedAntimagicShield      string
	GiftsOfDwarvenLordsAutomatedAntimagicShieldAlt   string
	GiftsOfDwarvenLordsProtectiveBelt                string
	GiftsOfDwarvenLordsProtectiveBeltAlt             string
	GiftsOfDwarvenLordsCrimsonResonanceController    string
	GiftsOfDwarvenLordsCrimsonResonanceControllerAlt string
	GiftsOfDwarvenLordsEmeraldResonanceController    string
	GiftsOfDwarvenLordsEmeraldResonanceControllerAlt string

	// Set: Elixir of Life
	ElixirOfLifeFlaskOfOblivion string
	ElixirOfLifeLifebloodFairy  string
	ElixirOfLifeRingOfLife      string

	// Set: Shadow of Death
	ShadowOfDeathCursedArmor string
	ShadowOfDeathBoneBoots   string
	ShadowOfDeathSecondShade string
	ShadowOfDeathDarkHatchet string

	// Set: Wanderer's Way
	WanderersWayBootsOfTravel string
	WanderersWayBackpack      string

	// Set: Living Arrows
	LivingArrowsShroomwoodBow      string
	LivingArrowsLightAndShadeCloak string
	LivingArrowsQuiveringQuiver    string

	// Set: Duelist's Pride
	DuelistsPrideRapier        string
	DuelistsPrideBuckler       string
	DuelistsPrideBrassKnuckles string

	// Set: Ethereal Knowledge
	EtherealKnowledgeGlassDagger string
	EtherealKnowledgeMirrorShoes string
	EtherealKnowledgeVortexDress string
	EtherealKnowledgeThirdEye    string

	// Set: Inner Song
	InnerSongMusicSheet     string
	InnerSongSingingPanPipe string
	InnerSongFancyMask      string

	// Set: Power of the Dragon Father
	PowerOfTheDragonFatherRedDragonFlameTongue string
	PowerOfTheDragonFatherDragonScaleShield    string
	PowerOfTheDragonFatherDragonScaleArmor     string
	PowerOfTheDragonFatherDragonCrest          string
	PowerOfTheDragonFatherDragonboneGreaves    string
	PowerOfTheDragonFatherSlitheringSash       string
	PowerOfTheDragonFatherDragonWing           string
	PowerOfTheDragonFatherPiercingEyeOfADragon string

	// Set: Beelzebub's Blessing
	BeelzebubsBlessingDemonClaw       string
	BeelzebubsBlessingChitinousShield string
	BeelzebubsBlessingHeartbeat       string
	BeelzebubsBlessingDemonCrest      string

	// Set: Boreolos
	BoreolosHand  string
	BoreolosFoot  string
	BoreolosHeart string
	BoreolosHead  string

	// Set: Holy Sigils
	HolySigilOfRoph         string
	HolySigilOfEridore      string
	HolySigilOfMearea       string
	HolySigilOfInsara       string
	HolySigilOfQuix         string
	HolySigilOfTheSevenMagi string
	HolySigilOfTheSecondMan string
	HolySigilOfUurdt        string

	// Set: Rule of Shadow
	RuleOfShadowLiquidSilence  string
	RuleOfShadowTheTruthmaker  string
	RuleOfShadowTheTruthseeker string
	RuleOfShadowNostriasGaze   string

	// Set: Ambassador's Word
	AmbassadorsWordDiplomaticGifts string
	AmbassadorsWordAmbassadorsSash string

	// Set: Warrior's Strength
	WarriorsStrengthWarriorsBelt     string
	WarriorsStrengthWarriorsOberegus string

	// Set: Keeper's Fortitude
	KeepersFortitudeKeepersRing     string
	KeepersFortitudeKeepersOberegus string

	// Set: Wizard's Might
	WizardsMightWizardsCloak    string
	WizardsMightWizardsOberegus string

	// Set: Scholar's Wisdom
	ScholarsWisdomScholarsTiara    string
	ScholarsWisdomScholarsOberegus string
}

var bannableItemSidValues = bannableItemSids{
	// Movement
	PoleStar:          "pole_star_artifact",
	SevenLeagueBoots:  "seven_league_boots_artifact",
	SwampBoots:        "swamp_boots_artifact",
	WarlordBoots:      "warlord_boots_artifact",
	MagicKeyRing:      "magic_key_ring_artifact",
	LegionsStep:       "legions_step_artifact",
	FallenAngelWings:  "fallen_angel_wings_artifact",
	BannerOfFourWinds: "banner_of_four_winds_artifact",
	Spyglass:          "spyglass_artifact",

	// Diplomacy
	VoodooshDoll:     "voodoosh_doll_artifact",
	FlagOfTruce:      "flag_of_truce_artifact",
	RingOfNeutrality: "ring_of_neutrality_artifact",

	// Combat
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

	// Magic
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

	// Misc (standalone)
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

	// Set: Resonant Sphere
	ResonantSphereOrbOfTwilight: "resonant_sphere_orb_of_twilight_artifact",
	ResonantSphereOrbOfDaylight: "resonant_sphere_orb_of_daylight_artifact",
	ResonantSphereOrbOfEternity: "resonant_sphere_orb_of_eternity_artifact",
	ResonantSpherePrimalOrb:     "resonant_sphere_primal_orb_artifact",

	// Set: Tranquility
	TranquilityBrightmindTiara: "tranquility_brightmind_tiara_artifact",
	TranquilityMagicMirror:     "tranquility_magic_mirror_artifact",
	TranquilityRingOfSerenity:  "tranquility_ring_of_serenity_artifact",

	// Set: Shamaniac Soul
	ShamaniacSoulShamanStaff:     "shamaniac_soul_shaman_staff_artifact",
	ShamaniacSoulIridescentCloak: "shamaniac_soul_iridescent_cloak_artifact",
	ShamaniacSoulGemwoodMask:     "shamaniac_soul_gemwood_mask_artifact",
	ShamaniacSoulClutchingRing:   "shamaniac_soul_clutching_ring_artifact",

	// Set: Knight's Honor
	KnightsHonorDrumsOfWar:  "knights_honor_drums_of_war_artifact",
	KnightsHonorLance:       "knights_honor_lance_artifact",
	KnightsHonorMisericorde: "knights_honor_misericorde_artifact",
	KnightsHonorPlateArmor:  "knights_honor_plate_armor_artifact",
	KnightsHonorArmet:       "knights_honor_armet_artifact",

	// Set: Ukhtabar Seal
	UkhtabarSealUkhSeal:   "ukhtabar_seal_ukh_seal_artifact",
	UkhtabarSealTabarSeal: "ukhtabar_seal_tabar_seal_artifact",

	// Set: Milo's Curse
	MilosCurseGoldenPig:    "milos_curse_golden_pig_artifact",
	MilosCurseGoldenMoth:   "milos_curse_golden_moth_artifact",
	MilosCurseSkullOfMilos: "milos_curse_skull_of_milos_artifact",

	// Set: Pauper's Glory
	PaupersGloryWoodenRing: "paupers_glory_wooden_ring_artifact",
	PaupersGloryStrawHat:   "paupers_glory_straw_hat_artifact",
	PaupersGloryRopeBelt:   "paupers_glory_rope_belt_artifact",
	PaupersGloryRags:       "paupers_glory_rags_artifact",
	PaupersGloryDumbClub:   "paupers_glory_dumb_club_artifact",
	PaupersGloryLastCoin:   "paupers_glory_last_coin_artifact",

	// Set: Angelic Alliance
	AngelicAllianceSwordOfJudgement:            "angelic_alliance_sword_of_judgement_artifact",
	AngelicAllianceCelestialSashOfBliss:        "angelic_alliance_celestial_sash_of_bliss_artifact",
	AngelicAllianceLionsShieldOfCourage:        "angelic_alliance_lions_shield_of_courage_artifact",
	AngelicAllianceArmorOfWonder:               "angelic_alliance_armor_of_wonder_artifact",
	AngelicAllianceHelmOfHeavenlyEnlightenment: "angelic_alliance_helm_of_heavenly_enlightenment_artifact",
	AngelicAllianceSandalsOfTheSaint:           "angelic_alliance_sandals_of_the_saint_artifact",

	// Set: Gifts of Dwarven Lords
	GiftsOfDwarvenLordsAutomatedAntimagicShield:      "gifts_of_dwarven_lords_automated_antimagic_shield_artifact",
	GiftsOfDwarvenLordsAutomatedAntimagicShieldAlt:   "gifts_of_dwarven_lords_automated_antimagic_shield_artifact_alt",
	GiftsOfDwarvenLordsProtectiveBelt:                "gifts_of_dwarven_lords_protective_belt_artifact",
	GiftsOfDwarvenLordsProtectiveBeltAlt:             "gifts_of_dwarven_lords_protective_belt_artifact_alt",
	GiftsOfDwarvenLordsCrimsonResonanceController:    "gifts_of_dwarven_lords_crimson_resonance_controller_artifact",
	GiftsOfDwarvenLordsCrimsonResonanceControllerAlt: "gifts_of_dwarven_lords_crimson_resonance_controller_artifact_alt",
	GiftsOfDwarvenLordsEmeraldResonanceController:    "gifts_of_dwarven_lords_emerald_resonance_controller_artifact",
	GiftsOfDwarvenLordsEmeraldResonanceControllerAlt: "gifts_of_dwarven_lords_emerald_resonance_controller_artifact_alt",

	// Set: Elixir of Life
	ElixirOfLifeFlaskOfOblivion: "elixir_of_life_flask_of_oblivion_artifact",
	ElixirOfLifeLifebloodFairy:  "elixir_of_life_lifeblood_fairy_artifact",
	ElixirOfLifeRingOfLife:      "elixir_of_life_ring_of_life_artifact",

	// Set: Shadow of Death
	ShadowOfDeathCursedArmor: "shadow_of_death_cursed_armor_artifact",
	ShadowOfDeathBoneBoots:   "shadow_of_death_bone_boots_artifact",
	ShadowOfDeathSecondShade: "shadow_of_death_second_shade_artifact",
	ShadowOfDeathDarkHatchet: "shadow_of_death_dark_hatchet_artifact",

	// Set: Wanderer's Way
	WanderersWayBootsOfTravel: "wanderers_way_boots_of_travel_artifact",
	WanderersWayBackpack:      "wanderers_way_backpack_artifact",

	// Set: Living Arrows
	LivingArrowsShroomwoodBow:      "living_arrows_shroomwood_bow_artifact",
	LivingArrowsLightAndShadeCloak: "living_arrows_light_and_shade_cloak_artifact",
	LivingArrowsQuiveringQuiver:    "living_arrows_quivering_quiver_artifact",

	// Set: Duelist's Pride
	DuelistsPrideRapier:        "duelists_pride_rapier_artifact",
	DuelistsPrideBuckler:       "duelists_pride_buckler_artifact",
	DuelistsPrideBrassKnuckles: "duelists_pride_brass_knuckles_artifact",

	// Set: Ethereal Knowledge
	EtherealKnowledgeGlassDagger: "ethereal_knowledge_glass_dagger_artifact",
	EtherealKnowledgeMirrorShoes: "ethereal_knowledge_mirror_shoes_artifact",
	EtherealKnowledgeVortexDress: "ethereal_knowledge_vortex_dress_artifact",
	EtherealKnowledgeThirdEye:    "ethereal_knowledge_third_eye_artifact",

	// Set: Inner Song
	InnerSongMusicSheet:     "inner_song_music_sheet_artifact",
	InnerSongSingingPanPipe: "inner_song_singing_pan_pipe_artifact",
	InnerSongFancyMask:      "inner_song_fancy_mask_artifact",

	// Set: Power of the Dragon Father
	PowerOfTheDragonFatherRedDragonFlameTongue: "power_of_the_dragon_father_red_dragon_flame_tongue_artifact",
	PowerOfTheDragonFatherDragonScaleShield:    "power_of_the_dragon_father_dragon_scale_shield_artifact",
	PowerOfTheDragonFatherDragonScaleArmor:     "power_of_the_dragon_father_dragon_scale_armor_artifact",
	PowerOfTheDragonFatherDragonCrest:          "power_of_the_dragon_father_dragon_crest_artifact",
	PowerOfTheDragonFatherDragonboneGreaves:    "power_of_the_dragon_father_dragonbone_greaves_artifact",
	PowerOfTheDragonFatherSlitheringSash:       "power_of_the_dragon_father_slithering_sash_artifact",
	PowerOfTheDragonFatherDragonWing:           "power_of_the_dragon_father_dragon_wing_artifact",
	PowerOfTheDragonFatherPiercingEyeOfADragon: "power_of_the_dragon_father_piercing_eye_of_a_dragon_artifact",

	// Set: Beelzebub's Blessing
	BeelzebubsBlessingDemonClaw:       "beelzebubs_blessing_demon_claw_artifact",
	BeelzebubsBlessingChitinousShield: "beelzebubs_blessing_chitinous_shield_artifact",
	BeelzebubsBlessingHeartbeat:       "beelzebubs_blessing_heartbeat_artifact",
	BeelzebubsBlessingDemonCrest:      "beelzebubs_blessing_demon_crest_artifact",

	// Set: Boreolos
	BoreolosHand:  "boreolos_hand_artifact",
	BoreolosFoot:  "boreolos_foot_artifact",
	BoreolosHeart: "boreolos_heart_artifact",
	BoreolosHead:  "boreolos_head_artifact",

	// Set: Holy Sigils
	HolySigilOfRoph:         "holy_sigil_of_roph_artifact",
	HolySigilOfEridore:      "holy_sigil_of_eridore_artifact",
	HolySigilOfMearea:       "holy_sigil_of_mearea_artifact",
	HolySigilOfInsara:       "holy_sigil_of_insara_artifact",
	HolySigilOfQuix:         "holy_sigil_of_quix_artifact",
	HolySigilOfTheSevenMagi: "holy_sigil_of_the_seven_magi_artifact",
	HolySigilOfTheSecondMan: "holy_sigil_of_the_second_man_artifact",
	HolySigilOfUurdt:        "holy_sigil_of_uurdt_artifact",

	// Set: Rule of Shadow
	RuleOfShadowLiquidSilence:  "rule_of_shadow_liquid_silence_artifact",
	RuleOfShadowTheTruthmaker:  "rule_of_shadow_the_truthmaker_artifact",
	RuleOfShadowTheTruthseeker: "rule_of_shadow_the_truthseeker_artifact",
	RuleOfShadowNostriasGaze:   "rule_of_shadow_nostrias_gaze_artifact",

	// Set: Ambassador's Word
	AmbassadorsWordDiplomaticGifts: "ambassadors_word_diplomatic_gifts_artifact",
	AmbassadorsWordAmbassadorsSash: "ambassadors_word_ambassadors_sash_artifact",

	// Set: Warrior's Strength
	WarriorsStrengthWarriorsBelt:     "warriors_strength_warriors_belt_artifact",
	WarriorsStrengthWarriorsOberegus: "warriors_strength_warriors_oberegus_artifact",

	// Set: Keeper's Fortitude
	KeepersFortitudeKeepersRing:     "keepers_fortitude_keepers_ring_artifact",
	KeepersFortitudeKeepersOberegus: "keepers_fortitude_keepers_oberegus_artifact",

	// Set: Wizard's Might
	WizardsMightWizardsCloak:    "wizards_might_wizards_cloak_artifact",
	WizardsMightWizardsOberegus: "wizards_might_wizards_oberegus_artifact",

	// Set: Scholar's Wisdom
	ScholarsWisdomScholarsTiara:    "scholars_wisdom_scholars_tiara_artifact",
	ScholarsWisdomScholarsOberegus: "scholars_wisdom_scholars_oberegus_artifact",
}

// GetBannableItemSidValues returns the artifact SIDs used for
//
//	globalBans.items
func GetBannableItemSidValues() bannableItemSids {
	return bannableItemSidValues
}
