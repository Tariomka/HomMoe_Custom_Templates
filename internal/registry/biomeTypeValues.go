package registry

type biomeTypes struct {
	FromList        string
	Match           string
	MatchMainObject string
	MatchZone       string
}

// GetBiomeTypeValues returns the different biome types used for
//
//	variants.zones.zoneBiome.type
//	variants.zones.contentBiome.type
//	variants.zones.metaObjectBiome.type
func GetBiomeTypeValues() biomeTypes {
	return biomeTypes{
		FromList:        "FromList",
		Match:           "Match",
		MatchMainObject: "MatchMainObject",
		MatchZone:       "MatchZone",
	}
}
