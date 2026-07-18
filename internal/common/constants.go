package common

const (
	HubZoneName       = "Hub"
	HubZonePrefix     = "Hub-"
	PlayerZonePrefix  = "Spawn-"
	NeutralZonePrefix = "Neutral-"
)

func GetZoneLabels() []string {
	return []string{
		"A", "B", "C", "D", "E", "F", "G", "H",
		"I", "J", "K", "L", "M", "N", "O", "P",
		"Q", "R", "S", "T", "U", "V", "W", "X",
		"Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF",
	}
}
