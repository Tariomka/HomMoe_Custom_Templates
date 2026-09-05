package constants

const (
	HubZoneName       = "Hub"
	HubZonePrefix     = "Hub-"
	PlayerZonePrefix  = "Spawn-"
	NeutralZonePrefix = "Neutral-"
)

func GetHubZoneNameFor(label string) string {
	return HubZonePrefix + label
}

func GetPlayerZoneNameFor(label string) string {
	return PlayerZonePrefix + label
}

func GetNeutralZoneNameFor(label string) string {
	return NeutralZonePrefix + label
}
