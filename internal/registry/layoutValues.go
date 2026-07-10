package registry

type layouts struct {
	AiSpawn           string
	Back              string
	Center            string
	CenterZone        string
	Leaf              string
	PlayerSpawn       string
	SecondSpawn       string
	SideSpawnZone     string
	SideZone          string
	Sides             string
	Spawn             string
	Spawns            string
	StartZone         string
	SuperTreasureZone string
	Treasure          string
	TreasureZone      string
	Treasures         string
	WinConditionZone  string
}

// GetLayoutValues returns the layout values used for
//
//	variants.zones.layout
func GetLayoutValues() layouts {
	return layouts{
		AiSpawn:    "zone_layout_ai_spawn",
		Back:       "zone_layout_back",
		Center:     "zone_layout_center",
		CenterZone: "zone_layout_center_zone",
		// Default: "zone_layout_default",
		Leaf:              "zone_layout_leaf",
		PlayerSpawn:       "zone_layout_player_spawn",
		SecondSpawn:       "zone_layout_second_spawn",
		SideSpawnZone:     "zone_layout_side_spawn_zone",
		SideZone:          "zone_layout_side_zone",
		Sides:             "zone_layout_sides",
		Spawn:             "zone_layout_spawn",
		Spawns:            "zone_layout_spawns",
		StartZone:         "zone_layout_start_zone",
		SuperTreasureZone: "zone_layout_supertreasure_zone",
		Treasure:          "zone_layout_treasure",
		TreasureZone:      "zone_layout_treasure_zone",
		Treasures:         "zone_layout_treasures",
		WinConditionZone:  "zone_layout_wincondition_zone",
	}
}
