package neutral_zone

type Profile struct {
	Layout                       string
	GuardReactionDistribution    []int
	GuardMultiplier              float64
	GuardedContentPool           []string
	UnguardedContentPool         []string
	ResourcesContentPool         []string
	GuardedContentValue          int
	GuardedContentValuePerArea   int
	UnguardedContentValue        int
	UnguardedContentValuePerArea int
	ResourcesValue               int
	ResourcesValuePerArea        int
	PrimaryCityGuardValue        int
	ExtraCityGuardValue          int
	PrimaryBuildingsSid          string
	ExtraBuildingsSid            string
}
