package zone

// Zone represents a single zone in the map.
type Zone struct {
	Name string `json:"name"`

	// GeneratorPosition is an optional normalised [0,1]×[0,1] hint stamped by
	// the generator for Random / Balanced layouts. It is NOT serialised — the
	// preview renderer uses it to reproduce the same geometry that drove
	// connection generation. Mirrors the C# editor's Zone.GeneratorPosition.
	GeneratorPosition *[2]float64 `json:"-"`

	// GeneratorRing is the concentric-ring index (0 = outermost player ring,
	// increasing toward centre) stamped by the generator for Balanced layouts.
	// Not serialised — used only by the preview renderer to snap zones onto
	// the correct ring. Mirrors C# Zone.GeneratorRing.
	GeneratorRing *int `json:"-"`

	Size   float64 `json:"size"`
	Layout string  `json:"layout"` // references a ZoneLayoutDef.Name in the template's `zoneLayouts`

	GuardCutoffValue          int     `json:"guardCutoffValue"`
	GuardRandomization        float64 `json:"guardRandomization"`
	GuardMultiplier           float64 `json:"guardMultiplier,omitempty"`
	GuardWeeklyIncrement      float64 `json:"guardWeeklyIncrement"`
	GuardReactionDistribution []int   `json:"guardReactionDistribution"`
	DiplomacyModifier         float64 `json:"diplomacyModifier,omitempty"`

	EncounterHolesSettings *EncounterHolesSettings `json:"encounterHolesSettings,omitempty"`

	RandomHireEnableWeeklyUnitIncrement []bool `json:"randomHireEnableWeeklyUnitIncrement,omitempty"`
	RandomHireInitialUnitIncrement      []int  `json:"randomHireInitialUnitIncrement,omitempty"`

	GuardedContentPool   []string `json:"guardedContentPool"`
	UnguardedContentPool []string `json:"unguardedContentPool"`
	ResourcesContentPool []string `json:"resourcesContentPool"`

	MandatoryContent   StringList `json:"mandatoryContent,omitempty"`
	ContentCountLimits StringList `json:"contentCountLimits,omitempty"`

	GuardedContentValue          int `json:"guardedContentValue"`
	GuardedContentValuePerArea   int `json:"guardedContentValuePerArea"`
	UnguardedContentValue        int `json:"unguardedContentValue"`
	UnguardedContentValuePerArea int `json:"unguardedContentValuePerArea"`
	ResourcesValue               int `json:"resourcesValue"`
	ResourcesValuePerArea        int `json:"resourcesValuePerArea"`

	MainObjects []MainObject `json:"mainObjects"`

	ZoneBiome        TypedRef `json:"zoneBiome"`
	ContentBiome     TypedRef `json:"contentBiome"`
	MetaObjectsBiome TypedRef `json:"metaObjectsBiome"`

	CrossroadsPosition *int   `json:"crossroadsPosition,omitempty"`
	Roads              []Road `json:"roads,omitempty"`
}
