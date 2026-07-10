package registry

type objectArtifacts struct {
	artifactScrolls
	artifactRandomItems
}

func GetMapObjectAllArtifactValues() objectArtifacts {
	return objectArtifacts{
		artifactScrolls:     GetMapObjectScrollValues(),
		artifactRandomItems: GetMapObjectRandomItemValues(),
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

// 			"id": "resonant_sphere_orb_of_twilight_artifact",

// 			"id": "resonant_sphere_orb_of_daylight_artifact",
// 			"id": "resonant_sphere_orb_of_eternity_artifact",
// 			"id": "resonant_sphere_primal_orb_artifact",
// 			"id": "tranquility_brightmind_tiara_artifact",
// 			"id": "tranquility_magic_mirror_artifact",
// 			"id": "tranquility_ring_of_serenity_artifact",
// 			"id": "shamaniac_soul_shaman_staff_artifact",
// 			"id": "shamaniac_soul_iridescent_cloak_artifact",
// 			"id": "shamaniac_soul_gemwood_mask_artifact",
// 			"id": "shamaniac_soul_clutching_ring_artifact",
// 			"id": "knights_honor_drums_of_war_artifact",
// 			"id": "knights_honor_lance_artifact",
// 			"id": "knights_honor_misericorde_artifact",
// 			"id": "knights_honor_plate_armor_artifact",
// 			"id": "knights_honor_armet_artifact",
// 			"id": "ukhtabar_seal_ukh_seal_artifact",
// 			"id": "ukhtabar_seal_tabar_seal_artifact",
// 			"id": "milos_curse_golden_pig_artifact",
// 			"id": "milos_curse_golden_moth_artifact",
// 			"id": "milos_curse_skull_of_milos_artifact",

// 			"id": "paupers_glory_wooden_ring_artifact",
// 			"id": "paupers_glory_straw_hat_artifact",
// 			"id": "paupers_glory_rope_belt_artifact",
// 			"id": "paupers_glory_rags_artifact",
// 			"id": "paupers_glory_dumb_club_artifact",
// 			"id": "paupers_glory_last_coin_artifact",
// 			"id": "angelic_alliance_sword_of_judgement_artifact",
// 			"id": "angelic_alliance_celestial_sash_of_bliss_artifact",
// 			"id": "angelic_alliance_lions_shield_of_courage_artifact",
// 			"id": "angelic_alliance_armor_of_wonder_artifact",
// 			"id": "angelic_alliance_helm_of_heavenly_enlightenment_artifact",
// 			"id": "angelic_alliance_sandals_of_the_saint_artifact",

// 			"id": "gifts_of_dwarven_lords_automated_antimagic_shield_artifact",
// 			"id": "gifts_of_dwarven_lords_protective_belt_artifact",
// 			"id": "gifts_of_dwarven_lords_crimson_resonance_controller_artifact",
// 			"id": "gifts_of_dwarven_lords_emerald_resonance_controller_artifact",
// 			"id": "elixir_of_life_lifeblood_fairy_artifact",
// 			"id": "elixir_of_life_flask_of_oblivion_artifact",
// 			"id": "elixir_of_life_ring_of_life_artifact",
// 			"id": "shadow_of_death_bone_boots_artifact",
// 			"id": "shadow_of_death_cursed_armor_artifact",
// 			"id": "shadow_of_death_dark_hatchet_artifact",
// 			"id": "shadow_of_death_second_shade_artifact",
// 			"id": "wanderers_way_boots_of_travel_artifact",
// 			"id": "wanderers_way_backpack_artifact",
// 			"id": "living_arrows_shroomwood_bow_artifact",
// 			"id": "living_arrows_light_and_shade_cloak_artifact",
// 			"id": "living_arrows_quivering_quiver_artifact",
// 			"id": "duelists_pride_rapier_artifact",
// 			"id": "duelists_pride_buckler_artifact",
// 			"id": "duelists_pride_brass_knuckles_artifact",
// 			"id": "ethereal_knowledge_glass_dagger_artifact",
// 			"id": "ethereal_knowledge_mirror_shoes_artifact",
// 			"id": "ethereal_knowledge_vortex_dress_artifact",
// 			"id": "ethereal_knowledge_third_eye_artifact",
// 			"id": "inner_song_music_sheet_artifact",
// 			"id": "inner_song_singing_pan_pipe_artifact",
// 			"id": "inner_song_fancy_mask_artifact",
// 			"id": "power_of_the_dragon_father_red_dragon_flame_tongue_artifact",
// 			"id": "power_of_the_dragon_father_dragon_scale_shield_artifact",
// 			"id": "power_of_the_dragon_father_dragon_scale_armor_artifact",
// 			"id": "power_of_the_dragon_father_dragon_crest_artifact",
// 			"id": "power_of_the_dragon_father_dragonbone_greaves_artifact",
// 			"id": "power_of_the_dragon_father_slithering_sash_artifact",
// 			"id": "power_of_the_dragon_father_dragon_wing_artifact",
// 			"id": "power_of_the_dragon_father_piercing_eye_of_a_dragon_artifact",
// 			"id": "beelzebubs_blessing_demon_claw_artifact",
// 			"id": "beelzebubs_blessing_chitinous_shield_artifact",
// 			"id": "beelzebubs_blessing_heartbeat_artifact",
// 			"id": "beelzebubs_blessing_demon_crest_artifact",
// 			"id": "rule_of_shadow_liquid_silence_artifact",
// 			"id": "rule_of_shadow_the_truthmaker_artifact",
// 			"id": "rule_of_shadow_the_truthseeker_artifact",
// 			"id": "rule_of_shadow_nostrias_gaze_artifact",
// 			"id": "warriors_strength_warriors_oberegus_artifact",
// 			"id": "warriors_strength_warriors_belt_artifact",
// 			"id": "keepers_fortitude_keepers_oberegus_artifact",

// 			"id": "keepers_fortitude_keepers_ring_artifact",

// 			"id": "wizards_might_wizards_oberegus_artifact",
// 			"id": "wizards_might_wizards_cloak_artifact",
// 			"id": "scholars_wisdom_scholars_oberegus_artifact",
// 			"id": "scholars_wisdom_scholars_tiara_artifact",
// 			"id": "cards_deck_artifact",
// 			"id": "runestone_shards_artifact",
// 			"id": "clothes_of_enlightenment_artifact",
// 			"id": "spells_in_a_bottle_artifact",
// 			"id": "spellbinders_hat_artifact",
// 			"id": "ambassadors_word_diplomatic_gifts_artifact",
// 			"id": "ambassadors_word_ambassadors_sash_artifact",
// 			"id": "orb_of_inhibition_artifact",
// 			"id": "boreolos_hand_artifact",
// 			"id": "boreolos_foot_artifact",
// 			"id": "boreolos_heart_artifact",
// 			"id": "boreolos_head_artifact",
// 			"id": "catechism_of_night_magic_artifact",
// 			"id": "catechism_of_daylight_magic_artifact",
// 			"id": "catechism_of_spacetime_magic_artifact",
// 			"id": "catechism_of_primal_magic_artifact",
// 			"id": "golden_goose_egg_artifact",
// 			"id": "tactical_guide_artifact",
// 			"id": "pole_star_artifact",
// 			"id": "ogres_club_of_havoc_artifact",
// 			"id": "garotte_artifact",
// 			"id": "tarq_of_the_rampaging_ogre_artifact",
// 			"id": "excalibur_artifact",
// 			"id": "caduceus_artifact",
// 			"id": "two_faced_mask_artifact",
// 			"id": "ancient_idol_artifact",
// 			"id": "orb_of_destruction_artifact",
// 			"id": "seal_of_silence_artifact",
// 			"id": "hourglass_of_protection_artifact",
// 			"id": "fine_wand_artifact",
// 			"id": "shoddy_shield_artifact",
// 			"id": "tunic_of_the_cyclops_king_artifact",
// 			"id": "eagle_armor_artifact",
// 			"id": "chain_mail_artifact",
// 			"id": "crown_of_the_supreme_magi_artifact",
// 			"id": "head_torch_artifact",
// 			"id": "legions_step_artifact",
// 			"id": "seven_league_boots_artifact",
// 			"id": "shackles_of_war_artifact",
// 			"id": "swamp_boots_artifact",
// 			"id": "warlord_boots_artifact",
// 			"id": "magic_key_ring_artifact",
// 			"id": "endless_bag_artifact",
// 			"id": "spyglass_artifact",
// 			"id": "soulless_sash_artifact",
// 			"id": "banner_of_four_winds_artifact",
// 			"id": "fallen_angel_wings_artifact",
// 			"id": "flag_of_truce_artifact",
// 			"id": "monster_head_artifact",
// 			"id": "omencaller_artifact",
// 			"id": "sixth_finger_artifact",
// 			"id": "soulscaller_ring_artifact",
// 			"id": "chain_link_artifact",
// 			"id": "ring_of_neutrality_artifact",
// 			"id": "lords_ring_artifact",
// 			"id": "holy_sigil_of_roph_artifact",
// 			"id": "holy_sigil_of_eridore_artifact",
// 			"id": "holy_sigil_of_mearea_artifact",
// 			"id": "holy_sigil_of_quix_artifact",
// 			"id": "holy_sigil_of_the_seven_magi_artifact",
// 			"id": "holy_sigil_of_the_second_man_artifact",
// 			"id": "holy_sigil_of_insara_artifact",
// 			"id": "holy_sigil_of_uurdt_artifact",
// 			"id": "voodoosh_doll_artifact",
// 			"id": "demonic_heart_artifact",
