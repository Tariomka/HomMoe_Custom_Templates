package registry

import "strings"

// ObjectSids are known object / encounter SIDs used in ValueOverride and as
// mandatory-content / content-pool references across example templates. Ported
// from the C# KnownValues.ObjectSids.
var ObjectSids = []string{
	"alchemy_lab",
	"arena",
	"beer_fountain",
	"boreal_call",
	"celestial_sphere",
	"chimerologist",
	"circus",
	"college_of_wonder",
	"crystal_trail",
	"dragon_utopia",
	"eternal_dragon",
	"fickle_shrine",
	"flattering_mirror",
	"forge",
	"fort",
	"fountain",
	"fountain_2",
	"huntsmans_camp",
	"infernal_cirque",
	"insaras_eye",
	"jousting_range",
	"mana_well",
	"market",
	"mine_crystals",
	"mine_gemstones",
	"mine_gold",
	"mine_mercury",
	"mine_ore",
	"mine_wood",
	"mirage",
	"monty_hall",
	"mysterious_stone",
	"mystical_tower",
	"mythic_scroll_box",
	"orb_observatory",
	"pandora_box",
	"petrified_memorial",
	"pile_of_books",
	"point_of_balance",
	"prison",
	"quixs_path",
	"random_hire_1",
	"random_hire_2",
	"random_hire_3",
	"random_hire_4",
	"random_hire_5",
	"random_hire_6",
	"random_hire_7",
	"random_item_common",
	"random_item_epic",
	"random_item_legendary",
	"random_item_rare",
	"remote_foothold",
	"research_laboratory",
	"ritual_pyre",
	"sacrificial_shrine",
	"shady_den",
	"stables",
	"tavern",
	"tear_of_truth",
	"the_gorge",
	"town_gate",
	"tree_of_abundance",
	"troglodyte_throne",
	"unforgotten_grave",
	"university",
	"unstable_ruins",
	"watchtower",
	"wind_rose",
	"wise_owl",
}

// SidToDisplayName converts a snake_case SID (with optional _artifact suffix)
// to a Title Case display name. Used as a fallback for IDs not in a catalog.
func SidToDisplayName(sid string) string {
	s := strings.ReplaceAll(strings.ReplaceAll(sid, "_artifact", ""), "_", " ")
	if len(s) == 0 {
		return sid
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
