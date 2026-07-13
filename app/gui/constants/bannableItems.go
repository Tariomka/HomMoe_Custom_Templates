package constants

import (
	"cmp"
	"slices"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

const (
	categoryMovement      = "Movement"
	categoryDiplomacy     = "Diplomacy"
	categoryCombat        = "Combat"
	categoryMagic         = "Magic"
	categoryMiscellaneous = "Misc"
	categorySet           = "Set"
)

// BannableItemEntry pairs an artifact SID with its human-readable name and the
// UI grouping category used by the item picker.
type BannableItemEntry struct {
	Sid      string
	Name     string
	Category string
}

func buildBannableItems() []BannableItemEntry {
	items := []BannableItemEntry{}
	items = append(items, buildMovementItems()...)
	items = append(items, buildDiplomacyItems()...)
	items = append(items, buildCombatItems()...)
	items = append(items, buildMagicItems()...)
	items = append(items, buildMiscellaneousItems()...)
	items = append(items, buildResonantSphereItems()...)
	items = append(items, buildTranquilityItems()...)
	items = append(items, buildShamaniacSoulItems()...)
	items = append(items, buildKnightsHonorItems()...)
	items = append(items, buildUkhtabarSealItems()...)
	items = append(items, buildMilosCurseItems()...)
	items = append(items, buildPaupersGloryItems()...)
	items = append(items, buildAngelicAllianceItems()...)
	items = append(items, buildGiftsOfDwarvenLordsItems()...)
	items = append(items, buildElixirOfLifeItems()...)
	items = append(items, buildShadowOfDeathItems()...)
	items = append(items, buildWanderersWayItems()...)
	items = append(items, buildLivingArrowsItems()...)
	items = append(items, buildDuelistsPrideItems()...)
	items = append(items, buildEtherealKnowledgeItems()...)
	items = append(items, buildInnerSongItems()...)
	items = append(items, buildPowerOfTheDragonFatherItems()...)
	items = append(items, buildBeelzebubsBlessingItems()...)
	items = append(items, buildBoreolosItems()...)
	items = append(items, buildHolySigilItems()...)
	items = append(items, buildRuleOfShadowItems()...)
	items = append(items, buildAmbassadorsWordItems()...)
	items = append(items, buildWarriorsStrengthItems()...)
	items = append(items, buildKeepersFortitudeItems()...)
	items = append(items, buildWizardsMightItems()...)
	items = append(items, buildScholarsWisdomItems()...)
	return items
}

func buildMovementItems() []BannableItemEntry {
	sids := registry.GetMapObjectMovementArtifactValues()
	return []BannableItemEntry{
		{sids.PoleStar, "Pole Star", categoryMovement},
		{sids.SevenLeagueBoots, "Seven League Boots", categoryMovement},
		{sids.SwampBoots, "Swamp Boots", categoryMovement},
		{sids.WarlordBoots, "Warlord Boots", categoryMovement},
		{sids.MagicKeyRing, "Magic Key Ring", categoryMovement},
		{sids.LegionsStep, "Legion's Step", categoryMovement},
		{sids.FallenAngelWings, "Fallen Angel Wings", categoryMovement},
		{sids.BannerOfFourWinds, "Banner of Four Winds", categoryMovement},
		{sids.Spyglass, "Spyglass", categoryMovement},
	}
}

func buildDiplomacyItems() []BannableItemEntry {
	sids := registry.GetMapObjectDiplomacyArtifactValues()
	return []BannableItemEntry{
		{sids.VoodooshDoll, "Voodoosh Doll", categoryDiplomacy},
		{sids.FlagOfTruce, "Flag of Truce", categoryDiplomacy},
		{sids.RingOfNeutrality, "Ring of Neutrality", categoryDiplomacy},
	}
}

func buildCombatItems() []BannableItemEntry {
	sids := registry.GetMapObjectCombatArtifactValues()
	return []BannableItemEntry{
		{sids.ShacklesOfWar, "Shackles of War", categoryCombat},
		{sids.OgresClubOfHavoc, "Ogre's Club of Havoc", categoryCombat},
		{sids.TarqOfTheRampagingOgre, "Tarq of the Rampaging Ogre", categoryCombat},
		{sids.TunicOfTheCyclopsKing, "Tunic of the Cyclops King", categoryCombat},
		{sids.Garotte, "Garotte", categoryCombat},
		{sids.HourglassOfProtection, "Hourglass of Protection", categoryCombat},
		{sids.ShoddyShield, "Shoddy Shield", categoryCombat},
		{sids.EagleArmor, "Eagle Armor", categoryCombat},
		{sids.ChainMail, "Chain Mail", categoryCombat},
		{sids.HeadTorch, "Head Torch", categoryCombat},
		{sids.FineWand, "Fine Wand", categoryCombat},
		{sids.LordsRing, "Lord's Ring", categoryCombat},
	}
}

func buildMagicItems() []BannableItemEntry {
	sids := registry.GetMapObjectMagicArtifactValues()
	return []BannableItemEntry{
		{sids.CatechismOfNightMagic, "Catechism of Night Magic", categoryMagic},
		{sids.CatechismOfDaylightMagic, "Catechism of Daylight Magic", categoryMagic},
		{sids.CatechismOfSpacetimeMagic, "Catechism of Spacetime Magic", categoryMagic},
		{sids.CatechismOfPrimalMagic, "Catechism of Primal Magic", categoryMagic},
		{sids.SpellbindersHat, "Spellbinder's Hat", categoryMagic},
		{sids.SpellsInABottle, "Spells in a Bottle", categoryMagic},
		{sids.OrbOfInhibition, "Orb of Inhibition", categoryMagic},
		{sids.OrbOfDestruction, "Orb of Destruction", categoryMagic},
		{sids.SealOfSilence, "Seal of Silence", categoryMagic},
		{sids.CrownOfTheSupremeMagi, "Crown of the Supreme Magi", categoryMagic},
		{sids.ClothesOfEnlightenment, "Clothes of Enlightenment", categoryMagic},
		{sids.CardsDeck, "Cards Deck", categoryMagic},
		{sids.RunestoneShards, "Runestone Shards", categoryMagic},
	}
}

func buildMiscellaneousItems() []BannableItemEntry {
	sids := registry.GetMapObjectMiscellaneousArtifactValues()
	return []BannableItemEntry{
		{sids.GoldenGooseEgg, "Golden Goose Egg", categoryMiscellaneous},
		{sids.TacticalGuide, "Tactical Guide", categoryMiscellaneous},
		{sids.EndlessBag, "Endless Bag", categoryMiscellaneous},
		{sids.SoullessSash, "Soulless Sash", categoryMiscellaneous},
		{sids.MonsterHead, "Monster Head", categoryMiscellaneous},
		{sids.Omencaller, "Omencaller", categoryMiscellaneous},
		{sids.SixthFinger, "Sixth Finger", categoryMiscellaneous},
		{sids.SoulscallerRing, "Soulscaller Ring", categoryMiscellaneous},
		{sids.ChainLink, "Chain Link", categoryMiscellaneous},
		{sids.DemonicHeart, "Demonic Heart", categoryMiscellaneous},
		{sids.TwoFacedMask, "Two-Faced Mask", categoryMiscellaneous},
		{sids.AncientIdol, "Ancient Idol", categoryMiscellaneous},
		{sids.Excalibur, "Excalibur", categoryMiscellaneous},
		{sids.Caduceus, "Caduceus", categoryMiscellaneous},
	}
}

func buildResonantSphereItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetResonantSphereValues()
	return []BannableItemEntry{
		{sids.OrbOfTwilight, "Resonant Sphere: Orb of Twilight", categorySet},
		{sids.OrbOfDaylight, "Resonant Sphere: Orb of Daylight", categorySet},
		{sids.OrbOfEternity, "Resonant Sphere: Orb of Eternity", categorySet},
		{sids.PrimalOrb, "Resonant Sphere: Primal Orb", categorySet},
	}
}

func buildTranquilityItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetTranquilityValues()
	return []BannableItemEntry{
		{sids.BrightmindTiara, "Tranquility: Brightmind Tiara", categorySet},
		{sids.MagicMirror, "Tranquility: Magic Mirror", categorySet},
		{sids.RingOfSerenity, "Tranquility: Ring of Serenity", categorySet},
	}
}

func buildShamaniacSoulItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetShamaniacSoulValues()
	return []BannableItemEntry{
		{sids.ShamanStaff, "Shamaniac Soul: Shaman Staff", categorySet},
		{sids.IridescentCloak, "Shamaniac Soul: Iridescent Cloak", categorySet},
		{sids.GemwoodMask, "Shamaniac Soul: Gemwood Mask", categorySet},
		{sids.ClutchingRing, "Shamaniac Soul: Clutching Ring", categorySet},
	}
}

func buildKnightsHonorItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetKnightsHonorValues()
	return []BannableItemEntry{
		{sids.DrumsOfWar, "Knight's Honor: Drums of War", categorySet},
		{sids.Lance, "Knight's Honor: Lance", categorySet},
		{sids.Misericorde, "Knight's Honor: Misericorde", categorySet},
		{sids.PlateArmor, "Knight's Honor: Plate Armor", categorySet},
		{sids.Armet, "Knight's Honor: Armet", categorySet},
	}
}

func buildUkhtabarSealItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetUkhtabarSealValues()
	return []BannableItemEntry{
		{sids.UkhSeal, "Ukhtabar Seal: Ukh Seal", categorySet},
		{sids.TabarSeal, "Ukhtabar Seal: Tabar Seal", categorySet},
	}
}

func buildMilosCurseItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetMilosCurseValues()
	return []BannableItemEntry{
		{sids.GoldenPig, "Milo's Curse: Golden Pig", categorySet},
		{sids.GoldenMoth, "Milo's Curse: Golden Moth", categorySet},
		{sids.SkullOfMilos, "Milo's Curse: Skull of Milos", categorySet},
	}
}

func buildPaupersGloryItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetPaupersGloryValues()
	return []BannableItemEntry{
		{sids.WoodenRing, "Pauper's Glory: Wooden Ring", categorySet},
		{sids.StrawHat, "Pauper's Glory: Straw Hat", categorySet},
		{sids.RopeBelt, "Pauper's Glory: Rope Belt", categorySet},
		{sids.Rags, "Pauper's Glory: Rags", categorySet},
		{sids.DumbClub, "Pauper's Glory: Dumb Club", categorySet},
		{sids.LastCoin, "Pauper's Glory: Last Coin", categorySet},
	}
}

func buildAngelicAllianceItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetAngelicAllianceValues()
	return []BannableItemEntry{
		{sids.SwordOfJudgement, "Angelic Alliance: Sword of Judgement", categorySet},
		{sids.CelestialSashOfBliss, "Angelic Alliance: Celestial Sash of Bliss", categorySet},
		{sids.LionsShieldOfCourage, "Angelic Alliance: Lion's Shield of Courage", categorySet},
		{sids.ArmorOfWonder, "Angelic Alliance: Armor of Wonder", categorySet},
		{sids.HelmOfHeavenlyEnlightenment, "Angelic Alliance: Helm of Heavenly Enlightenment", categorySet},
		{sids.SandalsOfTheSaint, "Angelic Alliance: Sandals of the Saint", categorySet},
	}
}

func buildGiftsOfDwarvenLordsItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetGiftsOfDwarvenLordsValues()
	return []BannableItemEntry{
		{sids.AutomatedAntimagicShield, "Dwarven Gifts: Automated Antimagic Shield", categorySet},
		{sids.ProtectiveBelt, "Dwarven Gifts: Protective Belt", categorySet},
		{sids.CrimsonResonanceController, "Dwarven Gifts: Crimson Resonance Controller", categorySet},
		{sids.EmeraldResonanceController, "Dwarven Gifts: Emerald Resonance Controller", categorySet},
	}
}

func buildElixirOfLifeItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetElixirOfLifeValues()
	return []BannableItemEntry{
		{sids.FlaskOfOblivion, "Elixir of Life: Flask of Oblivion", categorySet},
		{sids.LifebloodFairy, "Elixir of Life: Lifeblood Fairy", categorySet},
		{sids.RingOfLife, "Elixir of Life: Ring of Life", categorySet},
	}
}

func buildShadowOfDeathItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetShadowOfDeathValues()
	return []BannableItemEntry{
		{sids.CursedArmor, "Shadow of Death: Cursed Armor", categorySet},
		{sids.BoneBoots, "Shadow of Death: Bone Boots", categorySet},
		{sids.SecondShade, "Shadow of Death: Second Shade", categorySet},
		{sids.DarkHatchet, "Shadow of Death: Dark Hatchet", categorySet},
	}
}

func buildWanderersWayItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetWanderersWayValues()
	return []BannableItemEntry{
		{sids.BootsOfTravel, "Wanderer's Way: Boots of Travel", categorySet},
		{sids.Backpack, "Wanderer's Way: Backpack", categorySet},
	}
}

func buildLivingArrowsItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetLivingArrowsValues()
	return []BannableItemEntry{
		{sids.ShroomwoodBow, "Living Arrows: Shroomwood Bow", categorySet},
		{sids.LightAndShadeCloak, "Living Arrows: Light and Shade Cloak", categorySet},
		{sids.QuiveringQuiver, "Living Arrows: Quivering Quiver", categorySet},
	}
}

func buildDuelistsPrideItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetDuelistsPrideValues()
	return []BannableItemEntry{
		{sids.Rapier, "Duelist's Pride: Rapier", categorySet},
		{sids.Buckler, "Duelist's Pride: Buckler", categorySet},
		{sids.BrassKnuckles, "Duelist's Pride: Brass Knuckles", categorySet},
	}
}

func buildEtherealKnowledgeItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetEtherealKnowledgeValues()
	return []BannableItemEntry{
		{sids.GlassDagger, "Ethereal Knowledge: Glass Dagger", categorySet},
		{sids.MirrorShoes, "Ethereal Knowledge: Mirror Shoes", categorySet},
		{sids.VortexDress, "Ethereal Knowledge: Vortex Dress", categorySet},
		{sids.ThirdEye, "Ethereal Knowledge: Third Eye", categorySet},
	}
}

func buildInnerSongItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetInnerSongValues()
	return []BannableItemEntry{
		{sids.MusicSheet, "Inner Song: Music Sheet", categorySet},
		{sids.SingingPanPipe, "Inner Song: Singing Pan Pipe", categorySet},
		{sids.FancyMask, "Inner Song: Fancy Mask", categorySet},
	}
}

func buildPowerOfTheDragonFatherItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetPowerOfTheDragonFatherValues()
	return []BannableItemEntry{
		{sids.RedDragonFlameTongue, "Dragon Father: Red Dragon Flame Tongue", categorySet},
		{sids.DragonScaleShield, "Dragon Father: Dragon Scale Shield", categorySet},
		{sids.DragonScaleArmor, "Dragon Father: Dragon Scale Armor", categorySet},
		{sids.DragonCrest, "Dragon Father: Dragon Crest", categorySet},
		{sids.DragonBoneGreaves, "Dragon Father: Dragonbone Greaves", categorySet},
		{sids.SlitheringSash, "Dragon Father: Slithering Sash", categorySet},
		{sids.DragonWing, "Dragon Father: Dragon Wing", categorySet},
		{sids.PiercingEyeOfADragon, "Dragon Father: Piercing Eye of a Dragon", categorySet},
	}
}

func buildBeelzebubsBlessingItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetBeelzebubsBlessingValues()
	return []BannableItemEntry{
		{sids.DemonClaw, "Beelzebub's Blessing: Demon Claw", categorySet},
		{sids.ChitinousShield, "Beelzebub's Blessing: Chitinous Shield", categorySet},
		{sids.Heartbeat, "Beelzebub's Blessing: Heartbeat", categorySet},
		{sids.DemonCrest, "Beelzebub's Blessing: Demon Crest", categorySet},
	}
}

func buildBoreolosItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetBoreolosValues()
	return []BannableItemEntry{
		{sids.Hand, "Boreolos: Hand", categorySet},
		{sids.Foot, "Boreolos: Foot", categorySet},
		{sids.Heart, "Boreolos: Heart", categorySet},
		{sids.Head, "Boreolos: Head", categorySet},
	}
}

func buildHolySigilItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetHolySigilsValues()
	return []BannableItemEntry{
		{sids.SigilOfRoph, "Holy Sigil of Roph", categorySet},
		{sids.SigilOfEridore, "Holy Sigil of Eridore", categorySet},
		{sids.SigilOfMearea, "Holy Sigil of Mearea", categorySet},
		{sids.SigilOfInsara, "Holy Sigil of Insara", categorySet},
		{sids.SigilOfQuix, "Holy Sigil of Quix", categorySet},
		{sids.SigilOfTheSevenMagi, "Holy Sigil of the Seven Magi", categorySet},
		{sids.SigilOfTheSecondMan, "Holy Sigil of the Second Man", categorySet},
		{sids.SigilOfUurdt, "Holy Sigil of Uurdt", categorySet},
	}
}

func buildRuleOfShadowItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetRuleOfShadowValues()
	return []BannableItemEntry{
		{sids.LiquidSilence, "Rule of Shadow: Liquid Silence", categorySet},
		{sids.TheTruthmaker, "Rule of Shadow: The Truthmaker", categorySet},
		{sids.TheTruthseeker, "Rule of Shadow: The Truthseeker", categorySet},
		{sids.NostriasGaze, "Rule of Shadow: Nostria's Gaze", categorySet},
	}
}

func buildAmbassadorsWordItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetAmbassadorsWordValues()
	return []BannableItemEntry{
		{sids.DiplomaticGifts, "Ambassador's Word: Diplomatic Gifts", categorySet},
		{sids.AmbassadorsSash, "Ambassador's Word: Ambassador's Sash", categorySet},
	}
}

func buildWarriorsStrengthItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetWarriorsStrengthValues()
	return []BannableItemEntry{
		{sids.WarriorsBelt, "Warrior's Strength: Warrior's Belt", categorySet},
		{sids.WarriorsOberegus, "Warrior's Strength: Warrior's Oberegus", categorySet},
	}
}

func buildKeepersFortitudeItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetKeepersFortitudeValues()
	return []BannableItemEntry{
		{sids.KeepersRing, "Keeper's Fortitude: Keeper's Ring", categorySet},
		{sids.KeepersOberegus, "Keeper's Fortitude: Keeper's Oberegus", categorySet},
	}
}

func buildWizardsMightItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetWizardsMightValues()
	return []BannableItemEntry{
		{sids.WizardsCloak, "Wizard's Might: Wizard's Cloak", categorySet},
		{sids.WizardsOberegus, "Wizard's Might: Wizard's Oberegus", categorySet},
	}
}

func buildScholarsWisdomItems() []BannableItemEntry {
	sids := registry.GetMapObjectSetScholarsWisdomValues()
	return []BannableItemEntry{
		{sids.ScholarsTiara, "Scholar's Wisdom: Scholar's Tiara", categorySet},
		{sids.ScholarsOberegus, "Scholar's Wisdom: Scholar's Oberegus", categorySet},
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
	for _, item := range buildBannableItems() {
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
