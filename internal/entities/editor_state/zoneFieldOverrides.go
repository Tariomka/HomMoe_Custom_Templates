package editor_state

import "github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"

type ZoneFieldOverrides struct {
	Size   *float64 `json:"size,omitzero"`
	Layout *string  `json:"layout,omitzero"`

	GuardCutoffValue          *int     `json:"guardCutoffValue,omitzero"`
	GuardRandomization        *float64 `json:"guardRandomization,omitzero"`
	GuardMultiplier           *float64 `json:"guardMultiplier,omitzero"`
	GuardWeeklyIncrement      *float64 `json:"guardWeeklyIncrement,omitzero"`
	GuardReactionDistribution []int    `json:"guardReactionDistribution,omitzero"`
	DiplomacyModifier         *float64 `json:"diplomacyModifier,omitzero"`

	EncounterHolesSettings *template_entity.EncounterHolesSettings `json:"encounterHolesSettings,omitzero"`

	RandomHireEnableWeeklyUnitIncrement []bool `json:"randomHireEnableWeeklyUnitIncrement,omitzero"`
	RandomHireInitialUnitIncrement      []int  `json:"randomHireInitialUnitIncrement,omitzero"`

	GuardedContentPool   []string `json:"guardedContentPool,omitzero"`
	UnguardedContentPool []string `json:"unguardedContentPool,omitzero"`
	ResourcesContentPool []string `json:"resourcesContentPool,omitzero"`

	MandatoryContent   template_entity.StringList `json:"mandatoryContent,omitzero"`
	ContentCountLimits template_entity.StringList `json:"contentCountLimits,omitzero"`

	GuardedContentValue          *int `json:"guardedContentValue,omitzero"`
	GuardedContentValuePerArea   *int `json:"guardedContentValuePerArea,omitzero"`
	UnguardedContentValue        *int `json:"unguardedContentValue,omitzero"`
	UnguardedContentValuePerArea *int `json:"unguardedContentValuePerArea,omitzero"`
	ResourcesValue               *int `json:"resourcesValue,omitzero"`
	ResourcesValuePerArea        *int `json:"resourcesValuePerArea,omitzero"`

	MainObjects []template_entity.MainObject `json:"mainObjects,omitzero"`

	ZoneBiome        *template_entity.TypedRef `json:"zoneBiome,omitzero"`
	ContentBiome     *template_entity.TypedRef `json:"contentBiome,omitzero"`
	MetaObjectsBiome *template_entity.TypedRef `json:"metaObjectsBiome,omitzero"`

	CrossroadsPosition *int                   `json:"crossroadsPosition,omitzero"`
	Roads              []template_entity.Road `json:"roads,omitzero"`
}
