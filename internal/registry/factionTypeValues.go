package registry

type factionTypes struct {
	FromList string
	Match    string
}

var factionTypeValues = factionTypes{
	FromList: "FromList",
	Match:    "Match",
}

// GetFactionTypeValues returns the faction generation types used for
//
//	variants.zones.mainObjects.faction.type
func GetFactionTypeValues() factionTypes {
	return factionTypeValues
}
