package zone

import "github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner"

// Zone represents a single zone in the map.
type Zone struct {
	Name string `json:"name"`

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

	MandatoryContent   inner.StringList `json:"mandatoryContent,omitempty"`
	ContentCountLimits inner.StringList `json:"contentCountLimits,omitempty"`

	GuardedContentValue          int `json:"guardedContentValue"`
	GuardedContentValuePerArea   int `json:"guardedContentValuePerArea"`
	UnguardedContentValue        int `json:"unguardedContentValue"`
	UnguardedContentValuePerArea int `json:"unguardedContentValuePerArea"`
	ResourcesValue               int `json:"resourcesValue"`
	ResourcesValuePerArea        int `json:"resourcesValuePerArea"`

	MainObjects []inner.MainObject `json:"mainObjects"`

	ZoneBiome        inner.TypedRef `json:"zoneBiome"`
	ContentBiome     inner.TypedRef `json:"contentBiome"`
	MetaObjectsBiome inner.TypedRef `json:"metaObjectsBiome"`

	CrossroadsPosition *int         `json:"crossroadsPosition,omitempty"`
	Roads              []inner.Road `json:"roads,omitempty"`
}
