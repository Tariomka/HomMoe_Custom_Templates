package template_variant

// Zone represents a single zone in the map.
type Zone struct {
	Name string `json:"name"`

	// GeneratorPosition is an optional normalized [0,1]×[0,1] hint stamped by
	// the generator for the position-driven layouts (Random, Circles, Square,
	// Geometric, Cross). It is NOT serialized - the
	// preview renderer uses it to reproduce the same geometry that drove
	// connection generation
	GeneratorPosition *[2]float64 `json:"-"` //TODO: make it vec2[float64]. As it is not serialized to json, it should be safe to change

	// GeneratorRing is the concentric-ring index (0 = outermost player ring,
	// increasing toward centre) stamped by the generator for Circles layouts.
	// Not serialized - used only by the preview renderer to snap zones onto
	// the correct ring
	GeneratorRing *int `json:"-"`

	// ManualPosition is a normalized [0,1]×[0,1] position assigned when the
	// user moves or adds zones in the manual zone editor. It is NOT serialized
	// - when every zone has one, the preview renderer honours these positions
	// verbatim instead of computing a topology layout
	ManualPosition *[2]float64 `json:"-"`

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
