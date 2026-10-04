package constants

const (
	pseudoConnectionPrefix             = "Pseudo-"
	bridgeConnectionPrefix             = "Bridge-"
	fallbackConnectionPrefix           = "Fallback-"
	tournamentBalancedConnectionPrefix = "TBal-"
	geometricHubConnectionPrefix       = "GeoHub-"
	portalHubConnectionPrefix          = "Portal-Hub-"
	portalConnectionPrefix             = "Portal-"
	manualConnectionPrefix             = "Manual-"
	randomConnectionPrefix             = "Rnd-"
)

func GetPseudoConnectionNameFor(labelFrom, labelTo string) string {
	return pseudoConnectionPrefix + labelFrom + "-" + labelTo
}

func GetBridgeConnectionNameFor(labelFrom, labelTo string) string {
	return bridgeConnectionPrefix + labelFrom + "-" + labelTo
}

func GetFallbackConnectionNameFor(labelFrom, labelTo string) string {
	return fallbackConnectionPrefix + labelFrom + "-" + labelTo
}

func GetTournamentBalancedConnectionNameFor(labelFrom, labelTo string) string {
	return tournamentBalancedConnectionPrefix + labelFrom + "-" + labelTo
}

func GetGeometricHubConnectionNameFor(labelFrom, labelTo string) string {
	return geometricHubConnectionPrefix + labelFrom + "-" + labelTo
}

func GetPortalHubConnectionNameFor(label string) string {
	return portalHubConnectionPrefix + label
}

func GetPortalConnectionNameFor(labelFrom, labelTo string) string {
	return portalConnectionPrefix + labelFrom + "-" + labelTo
}

func GetManualConnectionNameFor(labelFrom, labelTo string) string {
	return manualConnectionPrefix + labelFrom + "-" + labelTo
}

func GetRandomConnectionNameFor(labelFrom, labelTo string) string {
	return randomConnectionPrefix + labelFrom + "-" + labelTo
}
