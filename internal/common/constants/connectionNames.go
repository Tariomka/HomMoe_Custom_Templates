package constants

const (
	pseudoConnectionPrefix                = "Pseudo-"
	bridgeConnectionPrefix                = "Bridge-"
	fallbackConnectionPrefix              = "Fallback-"
	neutralRingConnectionPrefix           = "NRing-"
	tournamentRingConnectionPrefix        = "TRing-"
	tournamentHubAndSpokeConnectionPrefix = "THubSpoke-"
	tournamentHubRingConnectionPrefix     = "THubRing-"
	tournamentChainConnectionPrefix       = "Tourney-"
	tournamentBalancedConnectionPrefix    = "TBal-"
	webConnectionPrefix                   = "Web-"
	geometricHubConnectionPrefix          = "GeoHub-"
	portalHubConnectionPrefix             = "Portal-Hub-"
	portalConnectionPrefix                = "Portal-"
	manualConnectionPrefix                = "Manual-"
	chainConnectionPrefix                 = "Chain-"
	randomConnectionPrefix                = "Rnd-"
	ringConnectionPrefix                  = "Ring-"
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

func GetNeutralRingConnectionNameFor(labelFrom, labelTo string) string {
	return neutralRingConnectionPrefix + labelFrom + "-" + labelTo
}

func GetTournamentRingConnectionNameFor(labelFrom, labelTo string) string {
	return tournamentRingConnectionPrefix + labelFrom + "-" + labelTo
}

func GetTournamentHubAndSpokeConnectionNameFor(labelFrom, labelTo string) string {
	return tournamentHubAndSpokeConnectionPrefix + labelFrom + "-" + labelTo
}

func GetTournamentHubRingConnectionNameFor(labelFrom, labelCentral, labelTo string) string {
	return tournamentHubRingConnectionPrefix + labelFrom + "-" + labelCentral + "-" + labelTo
}

func GetTournamentChainConnectionNameFor(labelFrom, labelTo string) string {
	return tournamentChainConnectionPrefix + labelFrom + "-" + labelTo
}

func GetTournamentBalancedConnectionNameFor(labelFrom, labelTo string) string {
	return tournamentBalancedConnectionPrefix + labelFrom + "-" + labelTo
}

func GetWebConnectionNameFor(labelFrom, labelTo string) string {
	return webConnectionPrefix + labelFrom + "-" + labelTo
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

func GetChainConnectionNameFor(labelFrom, labelTo string) string {
	return chainConnectionPrefix + labelFrom + "-" + labelTo
}

func GetRandomConnectionNameFor(labelFrom, labelTo string) string {
	return randomConnectionPrefix + labelFrom + "-" + labelTo
}

func GetRingConnectionNameFor(labelFrom, labelTo string) string {
	return ringConnectionPrefix + labelFrom + "-" + labelTo
}
