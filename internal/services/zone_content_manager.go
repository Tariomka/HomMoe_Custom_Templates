package services

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// ZoneContentManager builds mandatory content for zones
type ZoneContentManager struct{}

// BuildPlayerZoneMandatoryContent builds content for player zones
func (m *ZoneContentManager) BuildPlayerZoneMandatoryContent(scale float64) []models.ContentItem {
	return []models.ContentItem{
		{
			SID:     "hero_starting",
			Count:   1,
			Guarded: false,
		},
		{
			SID:     "town_starting",
			Count:   1,
			Guarded: false,
		},
	}
}

// BuildLowNeutralMandatoryContent builds content for low-tier neutral zones
func (m *ZoneContentManager) BuildLowNeutralMandatoryContent(scale float64) []models.ContentItem {
	return []models.ContentItem{
		{
			SID:     "encounter_easy",
			Count:   2,
			Guarded: true,
		},
		{
			SID:     "mine_basic",
			Count:   1,
			Guarded: true,
		},
	}
}

// BuildMediumNeutralMandatoryContent builds content for medium-tier neutral zones
func (m *ZoneContentManager) BuildMediumNeutralMandatoryContent(scale float64) []models.ContentItem {
	return []models.ContentItem{
		{
			SID:     "encounter_medium",
			Count:   3,
			Guarded: true,
		},
		{
			SID:     "mine_advanced",
			Count:   2,
			Guarded: true,
		},
		{
			SID:     "town_neutral",
			Count:   1,
			Guarded: true,
		},
	}
}

// BuildHighNeutralMandatoryContent builds content for high-tier neutral zones
func (m *ZoneContentManager) BuildHighNeutralMandatoryContent(scale float64) []models.ContentItem {
	return []models.ContentItem{
		{
			SID:     "encounter_hard",
			Count:   4,
			Guarded: true,
		},
		{
			SID:     "mine_epic",
			Count:   3,
			Guarded: true,
		},
		{
			SID:     "dragon_utopia",
			Count:   1,
			Guarded: true,
		},
	}
}

// BuildAllContentCountLimits builds count limits for all content
func (m *ZoneContentManager) BuildAllContentCountLimits() map[string]int {
	return map[string]int{
		"hero_starting":    999,
		"town_starting":    999,
		"town_neutral":     10,
		"encounter_easy":   50,
		"encounter_medium": 50,
		"encounter_hard":   30,
		"mine_basic":       20,
		"mine_advanced":    20,
		"mine_epic":        10,
		"dragon_utopia":    5,
	}
}

// ContentItemBuilder provides a fluent builder for content items
type ContentItemBuilder struct {
	item models.ContentItem
}

// NewContentItemBuilder creates a new builder
func NewContentItemBuilder() *ContentItemBuilder {
	return &ContentItemBuilder{
		item: models.ContentItem{},
	}
}

// WithSID sets the SID
func (b *ContentItemBuilder) WithSID(sid string) *ContentItemBuilder {
	b.item.SID = sid
	return b
}

// Guarded marks the item as guarded
func (b *ContentItemBuilder) Guarded() *ContentItemBuilder {
	b.item.Guarded = true
	return b
}

// Mine marks the item as a mine
func (b *ContentItemBuilder) Mine() *ContentItemBuilder {
	b.item.Mine = true
	return b
}

// SoloEncounter marks the item as a solo encounter
func (b *ContentItemBuilder) SoloEncounter() *ContentItemBuilder {
	b.item.SoloEncounter = true
	return b
}

// WithCount sets the count
func (b *ContentItemBuilder) WithCount(count int) *ContentItemBuilder {
	b.item.Count = count
	return b
}

// AddRule adds a placement rule
func (b *ContentItemBuilder) AddRule(rule models.PlacementRule) *ContentItemBuilder {
	b.item.PlacementRules = append(b.item.PlacementRules, rule)
	return b
}

// RoadDistance sets road distance variation
func (b *ContentItemBuilder) RoadDistance(distance int) *ContentItemBuilder {
	b.item.RoadDistanceVariation = map[string]int{
		"min": distance - 5,
		"max": distance + 5,
	}
	return b
}

// Build returns the built content item
func (b *ContentItemBuilder) Build() models.ContentItem {
	return b.item
}

// NewContentItem creates a new content item with fluent builder
func NewContentItem(sid string) *ContentItemBuilder {
	return NewContentItemBuilder().WithSID(sid)
}
