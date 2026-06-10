package constants

import "strings"

// BannableItemEntry pairs an artifact SID with its human-readable name and the
// UI grouping category used by the item picker.
type BannableItemEntry struct {
	Sid      string
	Name     string
	Category string
}

// BannableItems is the catalog of artifacts that can appear in
// globalBans.items, sourced from official game templates and content lists.
// Categories (Movement, Diplomacy, Combat, Magic, Misc, Set) group the picker.
//
//nolint:gochecknoglobals // semantic registry
var BannableItems = []BannableItemEntry{
	// Movement
	{"pole_star_artifact", "Pole Star", "Movement"},
	{"seven_league_boots_artifact", "Seven League Boots", "Movement"},
	{"swamp_boots_artifact", "Swamp Boots", "Movement"},
	{"warlord_boots_artifact", "Warlord Boots", "Movement"},
	{"magic_key_ring_artifact", "Magic Key Ring", "Movement"},
	{"legions_step_artifact", "Legion's Step", "Movement"},
	{"fallen_angel_wings_artifact", "Fallen Angel Wings", "Movement"},
	{"banner_of_four_winds_artifact", "Banner of Four Winds", "Movement"},
	{"spyglass_artifact", "Spyglass", "Movement"},

	// Diplomacy
	{"voodoosh_doll_artifact", "Voodoosh Doll", "Diplomacy"},
	{"flag_of_truce_artifact", "Flag of Truce", "Diplomacy"},
	{"ring_of_neutrality_artifact", "Ring of Neutrality", "Diplomacy"},

	// Combat
	{"shackles_of_war_artifact", "Shackles of War", "Combat"},
	{"ogres_club_of_havoc_artifact", "Ogre's Club of Havoc", "Combat"},
	{"tarq_of_the_rampaging_ogre_artifact", "Tarq of the Rampaging Ogre", "Combat"},
	{"tunic_of_the_cyclops_king_artifact", "Tunic of the Cyclops King", "Combat"},
	{"garotte_artifact", "Garotte", "Combat"},
	{"hourglass_of_protection_artifact", "Hourglass of Protection", "Combat"},
	{"shoddy_shield_artifact", "Shoddy Shield", "Combat"},
	{"eagle_armor_artifact", "Eagle Armor", "Combat"},
	{"chain_mail_artifact", "Chain Mail", "Combat"},
	{"head_torch_artifact", "Head Torch", "Combat"},
	{"fine_wand_artifact", "Fine Wand", "Combat"},
	{"lords_ring_artifact", "Lord's Ring", "Combat"},

	// Magic
	{"catechism_of_night_magic_artifact", "Catechism of Night Magic", "Magic"},
	{"catechism_of_daylight_magic_artifact", "Catechism of Daylight Magic", "Magic"},
	{"catechism_of_spacetime_magic_artifact", "Catechism of Spacetime Magic", "Magic"},
	{"catechism_of_primal_magic_artifact", "Catechism of Primal Magic", "Magic"},
	{"spellbinders_hat_artifact", "Spellbinder's Hat", "Magic"},
	{"spells_in_a_bottle_artifact", "Spells in a Bottle", "Magic"},
	{"orb_of_inhibition_artifact", "Orb of Inhibition", "Magic"},
	{"orb_of_destruction_artifact", "Orb of Destruction", "Magic"},
	{"seal_of_silence_artifact", "Seal of Silence", "Magic"},
	{"crown_of_the_supreme_magi_artifact", "Crown of the Supreme Magi", "Magic"},
	{"clothes_of_enlightenment_artifact", "Clothes of Enlightenment", "Magic"},
	{"cards_deck_artifact", "Cards Deck", "Magic"},
	{"runestone_shards_artifact", "Runestone Shards", "Magic"},

	// Misc (standalone)
	{"golden_goose_egg_artifact", "Golden Goose Egg", "Misc"},
	{"tactical_guide_artifact", "Tactical Guide", "Misc"},
	{"endless_bag_artifact", "Endless Bag", "Misc"},
	{"soulless_sash_artifact", "Soulless Sash", "Misc"},
	{"monster_head_artifact", "Monster Head", "Misc"},
	{"omencaller_artifact", "Omencaller", "Misc"},
	{"sixth_finger_artifact", "Sixth Finger", "Misc"},
	{"soulscaller_ring_artifact", "Soulscaller Ring", "Misc"},
	{"chain_link_artifact", "Chain Link", "Misc"},
	{"demonic_heart_artifact", "Demonic Heart", "Misc"},
	{"two_faced_mask_artifact", "Two-Faced Mask", "Misc"},
	{"ancient_idol_artifact", "Ancient Idol", "Misc"},
	{"excalibur_artifact", "Excalibur", "Misc"},
	{"caduceus_artifact", "Caduceus", "Misc"},

	// Set: Resonant Sphere
	{"resonant_sphere_orb_of_twilight_artifact", "Resonant Sphere: Orb of Twilight", "Set"},
	{"resonant_sphere_orb_of_daylight_artifact", "Resonant Sphere: Orb of Daylight", "Set"},
	{"resonant_sphere_orb_of_eternity_artifact", "Resonant Sphere: Orb of Eternity", "Set"},
	{"resonant_sphere_primal_orb_artifact", "Resonant Sphere: Primal Orb", "Set"},

	// Set: Tranquility
	{"tranquility_brightmind_tiara_artifact", "Tranquility: Brightmind Tiara", "Set"},
	{"tranquility_magic_mirror_artifact", "Tranquility: Magic Mirror", "Set"},
	{"tranquility_ring_of_serenity_artifact", "Tranquility: Ring of Serenity", "Set"},

	// Set: Shamaniac Soul
	{"shamaniac_soul_shaman_staff_artifact", "Shamaniac Soul: Shaman Staff", "Set"},
	{"shamaniac_soul_iridescent_cloak_artifact", "Shamaniac Soul: Iridescent Cloak", "Set"},
	{"shamaniac_soul_gemwood_mask_artifact", "Shamaniac Soul: Gemwood Mask", "Set"},
	{"shamaniac_soul_clutching_ring_artifact", "Shamaniac Soul: Clutching Ring", "Set"},

	// Set: Knight's Honor
	{"knights_honor_drums_of_war_artifact", "Knight's Honor: Drums of War", "Set"},
	{"knights_honor_lance_artifact", "Knight's Honor: Lance", "Set"},
	{"knights_honor_misericorde_artifact", "Knight's Honor: Misericorde", "Set"},
	{"knights_honor_plate_armor_artifact", "Knight's Honor: Plate Armor", "Set"},
	{"knights_honor_armet_artifact", "Knight's Honor: Armet", "Set"},

	// Set: Ukhtabar Seal
	{"ukhtabar_seal_ukh_seal_artifact", "Ukhtabar Seal: Ukh Seal", "Set"},
	{"ukhtabar_seal_tabar_seal_artifact", "Ukhtabar Seal: Tabar Seal", "Set"},

	// Set: Milo's Curse
	{"milos_curse_golden_pig_artifact", "Milo's Curse: Golden Pig", "Set"},
	{"milos_curse_golden_moth_artifact", "Milo's Curse: Golden Moth", "Set"},
	{"milos_curse_skull_of_milos_artifact", "Milo's Curse: Skull of Milos", "Set"},

	// Set: Pauper's Glory
	{"paupers_glory_wooden_ring_artifact", "Pauper's Glory: Wooden Ring", "Set"},
	{"paupers_glory_straw_hat_artifact", "Pauper's Glory: Straw Hat", "Set"},
	{"paupers_glory_rope_belt_artifact", "Pauper's Glory: Rope Belt", "Set"},
	{"paupers_glory_rags_artifact", "Pauper's Glory: Rags", "Set"},
	{"paupers_glory_dumb_club_artifact", "Pauper's Glory: Dumb Club", "Set"},
	{"paupers_glory_last_coin_artifact", "Pauper's Glory: Last Coin", "Set"},

	// Set: Angelic Alliance
	{"angelic_alliance_sword_of_judgement_artifact", "Angelic Alliance: Sword of Judgement", "Set"},
	{"angelic_alliance_celestial_sash_of_bliss_artifact", "Angelic Alliance: Celestial Sash of Bliss", "Set"},
	{"angelic_alliance_lions_shield_of_courage_artifact", "Angelic Alliance: Lion's Shield of Courage", "Set"},
	{"angelic_alliance_armor_of_wonder_artifact", "Angelic Alliance: Armor of Wonder", "Set"},
	{"angelic_alliance_helm_of_heavenly_enlightenment_artifact", "Angelic Alliance: Helm of Heavenly Enlightenment", "Set"},
	{"angelic_alliance_sandals_of_the_saint_artifact", "Angelic Alliance: Sandals of the Saint", "Set"},

	// Set: Gifts of Dwarven Lords
	{"gifts_of_dwarven_lords_automated_antimagic_shield_artifact", "Dwarven Gifts: Automated Antimagic Shield", "Set"},
	{"gifts_of_dwarven_lords_automated_antimagic_shield_artifact_alt", "Dwarven Gifts: Automated Antimagic Shield (Alt)", "Set"},
	{"gifts_of_dwarven_lords_protective_belt_artifact", "Dwarven Gifts: Protective Belt", "Set"},
	{"gifts_of_dwarven_lords_protective_belt_artifact_alt", "Dwarven Gifts: Protective Belt (Alt)", "Set"},
	{"gifts_of_dwarven_lords_crimson_resonance_controller_artifact", "Dwarven Gifts: Crimson Resonance Controller", "Set"},
	{"gifts_of_dwarven_lords_crimson_resonance_controller_artifact_alt", "Dwarven Gifts: Crimson Resonance Controller (Alt)", "Set"},
	{"gifts_of_dwarven_lords_emerald_resonance_controller_artifact", "Dwarven Gifts: Emerald Resonance Controller", "Set"},
	{"gifts_of_dwarven_lords_emerald_resonance_controller_artifact_alt", "Dwarven Gifts: Emerald Resonance Controller (Alt)", "Set"},

	// Set: Elixir of Life
	{"elixir_of_life_flask_of_oblivion_artifact", "Elixir of Life: Flask of Oblivion", "Set"},
	{"elixir_of_life_lifeblood_fairy_artifact", "Elixir of Life: Lifeblood Fairy", "Set"},
	{"elixir_of_life_ring_of_life_artifact", "Elixir of Life: Ring of Life", "Set"},

	// Set: Shadow of Death
	{"shadow_of_death_cursed_armor_artifact", "Shadow of Death: Cursed Armor", "Set"},
	{"shadow_of_death_bone_boots_artifact", "Shadow of Death: Bone Boots", "Set"},
	{"shadow_of_death_second_shade_artifact", "Shadow of Death: Second Shade", "Set"},
	{"shadow_of_death_dark_hatchet_artifact", "Shadow of Death: Dark Hatchet", "Set"},

	// Set: Wanderer's Way
	{"wanderers_way_boots_of_travel_artifact", "Wanderer's Way: Boots of Travel", "Set"},
	{"wanderers_way_backpack_artifact", "Wanderer's Way: Backpack", "Set"},

	// Set: Living Arrows
	{"living_arrows_shroomwood_bow_artifact", "Living Arrows: Shroomwood Bow", "Set"},
	{"living_arrows_light_and_shade_cloak_artifact", "Living Arrows: Light and Shade Cloak", "Set"},
	{"living_arrows_quivering_quiver_artifact", "Living Arrows: Quivering Quiver", "Set"},

	// Set: Duelist's Pride
	{"duelists_pride_rapier_artifact", "Duelist's Pride: Rapier", "Set"},
	{"duelists_pride_buckler_artifact", "Duelist's Pride: Buckler", "Set"},
	{"duelists_pride_brass_knuckles_artifact", "Duelist's Pride: Brass Knuckles", "Set"},

	// Set: Ethereal Knowledge
	{"ethereal_knowledge_glass_dagger_artifact", "Ethereal Knowledge: Glass Dagger", "Set"},
	{"ethereal_knowledge_mirror_shoes_artifact", "Ethereal Knowledge: Mirror Shoes", "Set"},
	{"ethereal_knowledge_vortex_dress_artifact", "Ethereal Knowledge: Vortex Dress", "Set"},
	{"ethereal_knowledge_third_eye_artifact", "Ethereal Knowledge: Third Eye", "Set"},

	// Set: Inner Song
	{"inner_song_music_sheet_artifact", "Inner Song: Music Sheet", "Set"},
	{"inner_song_singing_pan_pipe_artifact", "Inner Song: Singing Pan Pipe", "Set"},
	{"inner_song_fancy_mask_artifact", "Inner Song: Fancy Mask", "Set"},

	// Set: Power of the Dragon Father
	{"power_of_the_dragon_father_red_dragon_flame_tongue_artifact", "Dragon Father: Red Dragon Flame Tongue", "Set"},
	{"power_of_the_dragon_father_dragon_scale_shield_artifact", "Dragon Father: Dragon Scale Shield", "Set"},
	{"power_of_the_dragon_father_dragon_scale_armor_artifact", "Dragon Father: Dragon Scale Armor", "Set"},
	{"power_of_the_dragon_father_dragon_crest_artifact", "Dragon Father: Dragon Crest", "Set"},
	{"power_of_the_dragon_father_dragonbone_greaves_artifact", "Dragon Father: Dragonbone Greaves", "Set"},
	{"power_of_the_dragon_father_slithering_sash_artifact", "Dragon Father: Slithering Sash", "Set"},
	{"power_of_the_dragon_father_dragon_wing_artifact", "Dragon Father: Dragon Wing", "Set"},
	{"power_of_the_dragon_father_piercing_eye_of_a_dragon_artifact", "Dragon Father: Piercing Eye of a Dragon", "Set"},

	// Set: Beelzebub's Blessing
	{"beelzebubs_blessing_demon_claw_artifact", "Beelzebub's Blessing: Demon Claw", "Set"},
	{"beelzebubs_blessing_chitinous_shield_artifact", "Beelzebub's Blessing: Chitinous Shield", "Set"},
	{"beelzebubs_blessing_heartbeat_artifact", "Beelzebub's Blessing: Heartbeat", "Set"},
	{"beelzebubs_blessing_demon_crest_artifact", "Beelzebub's Blessing: Demon Crest", "Set"},

	// Set: Boreolos
	{"boreolos_hand_artifact", "Boreolos: Hand", "Set"},
	{"boreolos_foot_artifact", "Boreolos: Foot", "Set"},
	{"boreolos_heart_artifact", "Boreolos: Heart", "Set"},
	{"boreolos_head_artifact", "Boreolos: Head", "Set"},

	// Set: Holy Sigils
	{"holy_sigil_of_roph_artifact", "Holy Sigil of Roph", "Set"},
	{"holy_sigil_of_eridore_artifact", "Holy Sigil of Eridore", "Set"},
	{"holy_sigil_of_mearea_artifact", "Holy Sigil of Mearea", "Set"},
	{"holy_sigil_of_insara_artifact", "Holy Sigil of Insara", "Set"},
	{"holy_sigil_of_quix_artifact", "Holy Sigil of Quix", "Set"},
	{"holy_sigil_of_the_seven_magi_artifact", "Holy Sigil of the Seven Magi", "Set"},
	{"holy_sigil_of_the_second_man_artifact", "Holy Sigil of the Second Man", "Set"},
	{"holy_sigil_of_uurdt_artifact", "Holy Sigil of Uurdt", "Set"},

	// Set: Rule of Shadow
	{"rule_of_shadow_liquid_silence_artifact", "Rule of Shadow: Liquid Silence", "Set"},
	{"rule_of_shadow_the_truthmaker_artifact", "Rule of Shadow: The Truthmaker", "Set"},
	{"rule_of_shadow_the_truthseeker_artifact", "Rule of Shadow: The Truthseeker", "Set"},
	{"rule_of_shadow_nostrias_gaze_artifact", "Rule of Shadow: Nostria's Gaze", "Set"},

	// Set: Ambassador's Word
	{"ambassadors_word_diplomatic_gifts_artifact", "Ambassador's Word: Diplomatic Gifts", "Set"},
	{"ambassadors_word_ambassadors_sash_artifact", "Ambassador's Word: Ambassador's Sash", "Set"},

	// Set: Warrior's Strength
	{"warriors_strength_warriors_belt_artifact", "Warrior's Strength: Warrior's Belt", "Set"},
	{"warriors_strength_warriors_oberegus_artifact", "Warrior's Strength: Warrior's Oberegus", "Set"},

	// Set: Keeper's Fortitude
	{"keepers_fortitude_keepers_ring_artifact", "Keeper's Fortitude: Keeper's Ring", "Set"},
	{"keepers_fortitude_keepers_oberegus_artifact", "Keeper's Fortitude: Keeper's Oberegus", "Set"},

	// Set: Wizard's Might
	{"wizards_might_wizards_cloak_artifact", "Wizard's Might: Wizard's Cloak", "Set"},
	{"wizards_might_wizards_oberegus_artifact", "Wizard's Might: Wizard's Oberegus", "Set"},

	// Set: Scholar's Wisdom
	{"scholars_wisdom_scholars_tiara_artifact", "Scholar's Wisdom: Scholar's Tiara", "Set"},
	{"scholars_wisdom_scholars_oberegus_artifact", "Scholar's Wisdom: Scholar's Oberegus", "Set"},
}

// FindBannableItem returns the catalog entry for an artifact SID, or ok=false
// when the SID is not in the catalog (e.g. a custom/modded artifact).
func FindBannableItem(sid string) (BannableItemEntry, bool) {
	for _, item := range BannableItems {
		if item.Sid == sid {
			return item, true
		}
	}
	return BannableItemEntry{}, false
}

// SidToDisplayName converts a snake_case SID (with optional _artifact suffix)
// to a sentence-case display name. Used as a fallback for IDs not present in
// any catalog.
func SidToDisplayName(sid string) string {
	s := strings.ReplaceAll(strings.ReplaceAll(sid, "_artifact", ""), "_", " ")
	if len(s) == 0 {
		return sid
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
