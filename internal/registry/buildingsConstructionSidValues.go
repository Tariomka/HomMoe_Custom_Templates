package registry

type buildingsConstructionSids struct {
	Arcade       string
	Army         string
	ChosenOne    string
	ChosenOneUp1 string
	ChosenOneUp2 string
	ChosenOneUp3 string
	Default      string
	ExtraPoor    string
	ExtraRich    string
	Full         string
	Massacre     string
	MassacreUp1  string
	MassacreUp2  string
	MassacreUp3  string
	Medium       string
	Poor         string
	Rich         string
	Siege        string
	UltraRich    string
}

var buildingsConstructionSidValues = buildingsConstructionSids{
	Arcade:       "arcade_buildings_construction",
	Army:         "army_buildings_construction",
	ChosenOne:    "chosen_one_buildings_construction",
	ChosenOneUp1: "chosen_one_buildings_construction_up_1",
	ChosenOneUp2: "chosen_one_buildings_construction_up_2",
	ChosenOneUp3: "chosen_one_buildings_construction_up_3",
	Default:      "default_buildings_construction",
	ExtraPoor:    "extra_poor_buildings_construction",
	ExtraRich:    "extra_rich_buildings_construction",
	Full:         "full_buildings_construction",
	Massacre:     "massacre_buildings_construction",
	MassacreUp1:  "massacre_buildings_construction_up_1",
	MassacreUp2:  "massacre_buildings_construction_up_2",
	MassacreUp3:  "massacre_buildings_construction_up_3",
	Medium:       "medium_buildings_construction",
	Poor:         "poor_buildings_construction",
	Rich:         "rich_buildings_construction",
	Siege:        "siege_buildings_construction",
	UltraRich:    "ultra_rich_buildings_construction",
}

// GetBuildingsConstructionSidValues returns the SIDs used for
//
//	variants.zones.mainObjects.buildingsConstructionSid
func GetBuildingsConstructionSidValues() buildingsConstructionSids {
	return buildingsConstructionSidValues
}
