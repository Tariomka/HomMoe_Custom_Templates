package registry

type mainObjectTypes struct {
	AbandonedOutpost string
	City             string
	GladiatorArena   string
	Spawn            string
}

var mainObjectTypeValues = mainObjectTypes{
	AbandonedOutpost: "AbandonedOutpost",
	City:             "City",
	GladiatorArena:   "GladiatorArena",
	Spawn:            "Spawn",
}

// GetMainObjectTypeValues returns the main object types used for
//
//	variants.zones.mainObjects.type
func GetMainObjectTypeValues() mainObjectTypes {
	return mainObjectTypeValues
}
