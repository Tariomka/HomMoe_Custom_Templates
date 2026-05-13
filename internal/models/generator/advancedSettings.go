package generator

type AdvancedSettings struct {
	Enabled                    bool
	NeutralLowNoCastleCount    int
	NeutralLowCastleCount      int
	NeutralMediumNoCastleCount int
	NeutralMediumCastleCount   int
	NeutralHighNoCastleCount   int
	NeutralHighCastleCount     int
	PlayerZoneSize             float64
	NeutralZoneSize            float64
	GuardRandomization         float64
}
