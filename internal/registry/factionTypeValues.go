package registry

type factionTypes struct {
	FromList string
	Match    string
}

// GetFactionTypeValues returns the faction generation types used for
//
//	variants.zones.mainObjects.faction.type
func GetFactionTypeValues() factionTypes {
	return factionTypes{
		FromList: "FromList",
		Match:    "Match",
	}
}
