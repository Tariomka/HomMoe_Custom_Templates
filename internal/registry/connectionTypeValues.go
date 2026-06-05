package registry

type connectionTypes struct {
	Default        string
	Direct         string
	GladiatorArena string
	Portal         string
	Proximity      string
}

var connectionTypeValues = connectionTypes{
	Default:        "Default",
	Direct:         "Direct",
	GladiatorArena: "GladiatorArena",
	Portal:         "Portal",
	Proximity:      "Proximity",
}

// GetConnectionTypeValues returns the connection type values used for
//
//	variants.connections.connectionType
func GetConnectionTypeValues() connectionTypes {
	return connectionTypeValues
}
