package constants

import (
	"cmp"
	"slices"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

// BannableItemEntry pairs an artifact SID with its human-readable name and the
// UI grouping category used by the item picker.
type BannableItemEntry struct {
	Sid      string
	Name     string
	Category string
}

// BannableItems is the catalog of artifacts that can appear in
// globalBans.items. SIDs come from the registry; names and categories
// (Movement, Diplomacy, Combat, Magic, Misc, Set) are editor-side labels
// used to group the picker.
//
//nolint:gochecknoglobals // semantic catalog
var BannableItems = buildBannableItems()

func buildBannableItems() []BannableItemEntry {
	sids := registry.GetBannableItemSidValues()
	return []BannableItemEntry{
		// Movement
		{sids.PoleStar, "Pole Star", "Movement"},
		{sids.SevenLeagueBoots, "Seven League Boots", "Movement"},
		{sids.SwampBoots, "Swamp Boots", "Movement"},
		{sids.WarlordBoots, "Warlord Boots", "Movement"},
		{sids.MagicKeyRing, "Magic Key Ring", "Movement"},
		{sids.LegionsStep, "Legion's Step", "Movement"},
		{sids.FallenAngelWings, "Fallen Angel Wings", "Movement"},
		{sids.BannerOfFourWinds, "Banner of Four Winds", "Movement"},
		{sids.Spyglass, "Spyglass", "Movement"},

		// Diplomacy
		{sids.VoodooshDoll, "Voodoosh Doll", "Diplomacy"},
		{sids.FlagOfTruce, "Flag of Truce", "Diplomacy"},
		{sids.RingOfNeutrality, "Ring of Neutrality", "Diplomacy"},

		// Combat
		{sids.ShacklesOfWar, "Shackles of War", "Combat"},
		{sids.OgresClubOfHavoc, "Ogre's Club of Havoc", "Combat"},
		{sids.TarqOfTheRampagingOgre, "Tarq of the Rampaging Ogre", "Combat"},
		{sids.TunicOfTheCyclopsKing, "Tunic of the Cyclops King", "Combat"},
		{sids.Garotte, "Garotte", "Combat"},
		{sids.HourglassOfProtection, "Hourglass of Protection", "Combat"},
		{sids.ShoddyShield, "Shoddy Shield", "Combat"},
		{sids.EagleArmor, "Eagle Armor", "Combat"},
		{sids.ChainMail, "Chain Mail", "Combat"},
		{sids.HeadTorch, "Head Torch", "Combat"},
		{sids.FineWand, "Fine Wand", "Combat"},
		{sids.LordsRing, "Lord's Ring", "Combat"},

		// Magic
		{sids.CatechismOfNightMagic, "Catechism of Night Magic", "Magic"},
		{sids.CatechismOfDaylightMagic, "Catechism of Daylight Magic", "Magic"},
		{sids.CatechismOfSpacetimeMagic, "Catechism of Spacetime Magic", "Magic"},
		{sids.CatechismOfPrimalMagic, "Catechism of Primal Magic", "Magic"},
		{sids.SpellbindersHat, "Spellbinder's Hat", "Magic"},
		{sids.SpellsInABottle, "Spells in a Bottle", "Magic"},
		{sids.OrbOfInhibition, "Orb of Inhibition", "Magic"},
		{sids.OrbOfDestruction, "Orb of Destruction", "Magic"},
		{sids.SealOfSilence, "Seal of Silence", "Magic"},
		{sids.CrownOfTheSupremeMagi, "Crown of the Supreme Magi", "Magic"},
		{sids.ClothesOfEnlightenment, "Clothes of Enlightenment", "Magic"},
		{sids.CardsDeck, "Cards Deck", "Magic"},
		{sids.RunestoneShards, "Runestone Shards", "Magic"},

		// Misc (standalone)
		{sids.GoldenGooseEgg, "Golden Goose Egg", "Misc"},
		{sids.TacticalGuide, "Tactical Guide", "Misc"},
		{sids.EndlessBag, "Endless Bag", "Misc"},
		{sids.SoullessSash, "Soulless Sash", "Misc"},
		{sids.MonsterHead, "Monster Head", "Misc"},
		{sids.Omencaller, "Omencaller", "Misc"},
		{sids.SixthFinger, "Sixth Finger", "Misc"},
		{sids.SoulscallerRing, "Soulscaller Ring", "Misc"},
		{sids.ChainLink, "Chain Link", "Misc"},
		{sids.DemonicHeart, "Demonic Heart", "Misc"},
		{sids.TwoFacedMask, "Two-Faced Mask", "Misc"},
		{sids.AncientIdol, "Ancient Idol", "Misc"},
		{sids.Excalibur, "Excalibur", "Misc"},
		{sids.Caduceus, "Caduceus", "Misc"},

		// Set: Resonant Sphere
		{sids.ResonantSphereOrbOfTwilight, "Resonant Sphere: Orb of Twilight", "Set"},
		{sids.ResonantSphereOrbOfDaylight, "Resonant Sphere: Orb of Daylight", "Set"},
		{sids.ResonantSphereOrbOfEternity, "Resonant Sphere: Orb of Eternity", "Set"},
		{sids.ResonantSpherePrimalOrb, "Resonant Sphere: Primal Orb", "Set"},

		// Set: Tranquility
		{sids.TranquilityBrightmindTiara, "Tranquility: Brightmind Tiara", "Set"},
		{sids.TranquilityMagicMirror, "Tranquility: Magic Mirror", "Set"},
		{sids.TranquilityRingOfSerenity, "Tranquility: Ring of Serenity", "Set"},

		// Set: Shamaniac Soul
		{sids.ShamaniacSoulShamanStaff, "Shamaniac Soul: Shaman Staff", "Set"},
		{sids.ShamaniacSoulIridescentCloak, "Shamaniac Soul: Iridescent Cloak", "Set"},
		{sids.ShamaniacSoulGemwoodMask, "Shamaniac Soul: Gemwood Mask", "Set"},
		{sids.ShamaniacSoulClutchingRing, "Shamaniac Soul: Clutching Ring", "Set"},

		// Set: Knight's Honor
		{sids.KnightsHonorDrumsOfWar, "Knight's Honor: Drums of War", "Set"},
		{sids.KnightsHonorLance, "Knight's Honor: Lance", "Set"},
		{sids.KnightsHonorMisericorde, "Knight's Honor: Misericorde", "Set"},
		{sids.KnightsHonorPlateArmor, "Knight's Honor: Plate Armor", "Set"},
		{sids.KnightsHonorArmet, "Knight's Honor: Armet", "Set"},

		// Set: Ukhtabar Seal
		{sids.UkhtabarSealUkhSeal, "Ukhtabar Seal: Ukh Seal", "Set"},
		{sids.UkhtabarSealTabarSeal, "Ukhtabar Seal: Tabar Seal", "Set"},

		// Set: Milo's Curse
		{sids.MilosCurseGoldenPig, "Milo's Curse: Golden Pig", "Set"},
		{sids.MilosCurseGoldenMoth, "Milo's Curse: Golden Moth", "Set"},
		{sids.MilosCurseSkullOfMilos, "Milo's Curse: Skull of Milos", "Set"},

		// Set: Pauper's Glory
		{sids.PaupersGloryWoodenRing, "Pauper's Glory: Wooden Ring", "Set"},
		{sids.PaupersGloryStrawHat, "Pauper's Glory: Straw Hat", "Set"},
		{sids.PaupersGloryRopeBelt, "Pauper's Glory: Rope Belt", "Set"},
		{sids.PaupersGloryRags, "Pauper's Glory: Rags", "Set"},
		{sids.PaupersGloryDumbClub, "Pauper's Glory: Dumb Club", "Set"},
		{sids.PaupersGloryLastCoin, "Pauper's Glory: Last Coin", "Set"},

		// Set: Angelic Alliance
		{sids.AngelicAllianceSwordOfJudgement, "Angelic Alliance: Sword of Judgement", "Set"},
		{sids.AngelicAllianceCelestialSashOfBliss, "Angelic Alliance: Celestial Sash of Bliss", "Set"},
		{sids.AngelicAllianceLionsShieldOfCourage, "Angelic Alliance: Lion's Shield of Courage", "Set"},
		{sids.AngelicAllianceArmorOfWonder, "Angelic Alliance: Armor of Wonder", "Set"},
		{sids.AngelicAllianceHelmOfHeavenlyEnlightenment, "Angelic Alliance: Helm of Heavenly Enlightenment", "Set"},
		{sids.AngelicAllianceSandalsOfTheSaint, "Angelic Alliance: Sandals of the Saint", "Set"},

		// Set: Gifts of Dwarven Lords
		{sids.GiftsOfDwarvenLordsAutomatedAntimagicShield, "Dwarven Gifts: Automated Antimagic Shield", "Set"},
		{sids.GiftsOfDwarvenLordsAutomatedAntimagicShieldAlt, "Dwarven Gifts: Automated Antimagic Shield (Alt)", "Set"},
		{sids.GiftsOfDwarvenLordsProtectiveBelt, "Dwarven Gifts: Protective Belt", "Set"},
		{sids.GiftsOfDwarvenLordsProtectiveBeltAlt, "Dwarven Gifts: Protective Belt (Alt)", "Set"},
		{sids.GiftsOfDwarvenLordsCrimsonResonanceController, "Dwarven Gifts: Crimson Resonance Controller", "Set"},
		{sids.GiftsOfDwarvenLordsCrimsonResonanceControllerAlt, "Dwarven Gifts: Crimson Resonance Controller (Alt)", "Set"},
		{sids.GiftsOfDwarvenLordsEmeraldResonanceController, "Dwarven Gifts: Emerald Resonance Controller", "Set"},
		{sids.GiftsOfDwarvenLordsEmeraldResonanceControllerAlt, "Dwarven Gifts: Emerald Resonance Controller (Alt)", "Set"},

		// Set: Elixir of Life
		{sids.ElixirOfLifeFlaskOfOblivion, "Elixir of Life: Flask of Oblivion", "Set"},
		{sids.ElixirOfLifeLifebloodFairy, "Elixir of Life: Lifeblood Fairy", "Set"},
		{sids.ElixirOfLifeRingOfLife, "Elixir of Life: Ring of Life", "Set"},

		// Set: Shadow of Death
		{sids.ShadowOfDeathCursedArmor, "Shadow of Death: Cursed Armor", "Set"},
		{sids.ShadowOfDeathBoneBoots, "Shadow of Death: Bone Boots", "Set"},
		{sids.ShadowOfDeathSecondShade, "Shadow of Death: Second Shade", "Set"},
		{sids.ShadowOfDeathDarkHatchet, "Shadow of Death: Dark Hatchet", "Set"},

		// Set: Wanderer's Way
		{sids.WanderersWayBootsOfTravel, "Wanderer's Way: Boots of Travel", "Set"},
		{sids.WanderersWayBackpack, "Wanderer's Way: Backpack", "Set"},

		// Set: Living Arrows
		{sids.LivingArrowsShroomwoodBow, "Living Arrows: Shroomwood Bow", "Set"},
		{sids.LivingArrowsLightAndShadeCloak, "Living Arrows: Light and Shade Cloak", "Set"},
		{sids.LivingArrowsQuiveringQuiver, "Living Arrows: Quivering Quiver", "Set"},

		// Set: Duelist's Pride
		{sids.DuelistsPrideRapier, "Duelist's Pride: Rapier", "Set"},
		{sids.DuelistsPrideBuckler, "Duelist's Pride: Buckler", "Set"},
		{sids.DuelistsPrideBrassKnuckles, "Duelist's Pride: Brass Knuckles", "Set"},

		// Set: Ethereal Knowledge
		{sids.EtherealKnowledgeGlassDagger, "Ethereal Knowledge: Glass Dagger", "Set"},
		{sids.EtherealKnowledgeMirrorShoes, "Ethereal Knowledge: Mirror Shoes", "Set"},
		{sids.EtherealKnowledgeVortexDress, "Ethereal Knowledge: Vortex Dress", "Set"},
		{sids.EtherealKnowledgeThirdEye, "Ethereal Knowledge: Third Eye", "Set"},

		// Set: Inner Song
		{sids.InnerSongMusicSheet, "Inner Song: Music Sheet", "Set"},
		{sids.InnerSongSingingPanPipe, "Inner Song: Singing Pan Pipe", "Set"},
		{sids.InnerSongFancyMask, "Inner Song: Fancy Mask", "Set"},

		// Set: Power of the Dragon Father
		{sids.PowerOfTheDragonFatherRedDragonFlameTongue, "Dragon Father: Red Dragon Flame Tongue", "Set"},
		{sids.PowerOfTheDragonFatherDragonScaleShield, "Dragon Father: Dragon Scale Shield", "Set"},
		{sids.PowerOfTheDragonFatherDragonScaleArmor, "Dragon Father: Dragon Scale Armor", "Set"},
		{sids.PowerOfTheDragonFatherDragonCrest, "Dragon Father: Dragon Crest", "Set"},
		{sids.PowerOfTheDragonFatherDragonboneGreaves, "Dragon Father: Dragonbone Greaves", "Set"},
		{sids.PowerOfTheDragonFatherSlitheringSash, "Dragon Father: Slithering Sash", "Set"},
		{sids.PowerOfTheDragonFatherDragonWing, "Dragon Father: Dragon Wing", "Set"},
		{sids.PowerOfTheDragonFatherPiercingEyeOfADragon, "Dragon Father: Piercing Eye of a Dragon", "Set"},

		// Set: Beelzebub's Blessing
		{sids.BeelzebubsBlessingDemonClaw, "Beelzebub's Blessing: Demon Claw", "Set"},
		{sids.BeelzebubsBlessingChitinousShield, "Beelzebub's Blessing: Chitinous Shield", "Set"},
		{sids.BeelzebubsBlessingHeartbeat, "Beelzebub's Blessing: Heartbeat", "Set"},
		{sids.BeelzebubsBlessingDemonCrest, "Beelzebub's Blessing: Demon Crest", "Set"},

		// Set: Boreolos
		{sids.BoreolosHand, "Boreolos: Hand", "Set"},
		{sids.BoreolosFoot, "Boreolos: Foot", "Set"},
		{sids.BoreolosHeart, "Boreolos: Heart", "Set"},
		{sids.BoreolosHead, "Boreolos: Head", "Set"},

		// Set: Holy Sigils
		{sids.HolySigilOfRoph, "Holy Sigil of Roph", "Set"},
		{sids.HolySigilOfEridore, "Holy Sigil of Eridore", "Set"},
		{sids.HolySigilOfMearea, "Holy Sigil of Mearea", "Set"},
		{sids.HolySigilOfInsara, "Holy Sigil of Insara", "Set"},
		{sids.HolySigilOfQuix, "Holy Sigil of Quix", "Set"},
		{sids.HolySigilOfTheSevenMagi, "Holy Sigil of the Seven Magi", "Set"},
		{sids.HolySigilOfTheSecondMan, "Holy Sigil of the Second Man", "Set"},
		{sids.HolySigilOfUurdt, "Holy Sigil of Uurdt", "Set"},

		// Set: Rule of Shadow
		{sids.RuleOfShadowLiquidSilence, "Rule of Shadow: Liquid Silence", "Set"},
		{sids.RuleOfShadowTheTruthmaker, "Rule of Shadow: The Truthmaker", "Set"},
		{sids.RuleOfShadowTheTruthseeker, "Rule of Shadow: The Truthseeker", "Set"},
		{sids.RuleOfShadowNostriasGaze, "Rule of Shadow: Nostria's Gaze", "Set"},

		// Set: Ambassador's Word
		{sids.AmbassadorsWordDiplomaticGifts, "Ambassador's Word: Diplomatic Gifts", "Set"},
		{sids.AmbassadorsWordAmbassadorsSash, "Ambassador's Word: Ambassador's Sash", "Set"},

		// Set: Warrior's Strength
		{sids.WarriorsStrengthWarriorsBelt, "Warrior's Strength: Warrior's Belt", "Set"},
		{sids.WarriorsStrengthWarriorsOberegus, "Warrior's Strength: Warrior's Oberegus", "Set"},

		// Set: Keeper's Fortitude
		{sids.KeepersFortitudeKeepersRing, "Keeper's Fortitude: Keeper's Ring", "Set"},
		{sids.KeepersFortitudeKeepersOberegus, "Keeper's Fortitude: Keeper's Oberegus", "Set"},

		// Set: Wizard's Might
		{sids.WizardsMightWizardsCloak, "Wizard's Might: Wizard's Cloak", "Set"},
		{sids.WizardsMightWizardsOberegus, "Wizard's Might: Wizard's Oberegus", "Set"},

		// Set: Scholar's Wisdom
		{sids.ScholarsWisdomScholarsTiara, "Scholar's Wisdom: Scholar's Tiara", "Set"},
		{sids.ScholarsWisdomScholarsOberegus, "Scholar's Wisdom: Scholar's Oberegus", "Set"},
	}
}

func GetBannableItemsWithExclusions(excluded []string) []BannableItemEntry {
	items := slices.DeleteFunc(
		buildBannableItems(),
		func(item BannableItemEntry) bool { return slices.Contains(excluded, item.Sid) })
	slices.SortStableFunc(items, CompareBannableItems)
	return items
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

func CompareBannableItems(a, b BannableItemEntry) int {
	if comparison := cmp.Compare(a.Category, b.Category); comparison != 0 {
		return comparison
	}

	return cmp.Compare(a.Name, b.Name)
}
