package registry

type roadConnectionTypes struct {
	MainObject       string
	Connection       string
	Crossroads       string
	MandatoryContent string
}

var roadConnectionTypeValues = roadConnectionTypes{
	MainObject:       "MainObject",
	Connection:       "Connection",
	Crossroads:       "Crossroads",
	MandatoryContent: "MandatoryContent",
}

// GetRoadConnectionTypeValues returns the road connection types used for
//
//	variants.zones.mainObjects.roads.from.type
//	variants.zones.mainObjects.roads.to.type
func GetRoadConnectionTypeValues() roadConnectionTypes {
	return roadConnectionTypeValues
}
