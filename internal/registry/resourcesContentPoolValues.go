package registry

type resourcesContentPools struct {
	StartZoneVeryPoor string
	StartZonePoor     string
	StartZoneMedium   string
	StartZoneRich     string

	SideZonePoor   string
	SideZoneMedium string
	SideZoneRich   string

	TreasureZoneZero          string
	TreasureZonePoor          string
	TreasureZoneMedium        string
	TreasureZoneRich          string
	TreasureZoneRichNoScrolls string

	DustOnly string
}

// GetResourcesContentPoolValues returns the resource content pool IDs used for
//
//	variants.zones.resourcesContentPool
func GetResourcesContentPoolValues() resourcesContentPools {
	return resourcesContentPools{
		StartZoneVeryPoor: "content_pool_general_resources_start_zone_very_poor",
		StartZonePoor:     "content_pool_general_resources_start_zone_poor",
		StartZoneMedium:   "content_pool_general_resources_start_zone_medium",
		StartZoneRich:     "content_pool_general_resources_start_zone_rich",

		SideZonePoor:   "content_pool_general_resources_side_zone_poor",
		SideZoneMedium: "content_pool_general_resources_side_zone_medium",
		SideZoneRich:   "content_pool_general_resources_side_zone_rich",

		TreasureZoneZero:          "content_pool_general_resources_treasure_zone_zero",
		TreasureZonePoor:          "content_pool_general_resources_treasure_zone_poor",
		TreasureZoneMedium:        "content_pool_general_resources_treasure_zone_medium",
		TreasureZoneRich:          "content_pool_general_resources_treasure_zone_rich",
		TreasureZoneRichNoScrolls: "content_pool_general_resources_treasure_zone_rich_no_scrolls",

		DustOnly: "content_pool_general_resources_dust_only",
	}
}
