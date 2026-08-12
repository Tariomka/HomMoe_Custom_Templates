package zone_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

// IZoneLabelProvider names the zones a variant is built from, letting callers
// substitute a deterministic labelling in tests.
type IZoneLabelProvider interface {
	CreatePlayerLabels(playerCount int) []string
	CreateNeutralZonePlans(configuration config.GeneratorConfig) neutral_zone.Plans
	GetHoldCityLabel(
		configuration config.GeneratorConfig,
		playerLabels []string,
		neutralZones neutral_zone.Plans) string
	CreateZoneName(label string, playerLabels []string) string
	CreateOrderedZoneLabels(
		configuration config.GeneratorConfig,
		playerLabels []string,
		neutralZones neutral_zone.Plans,
		isRing bool) []string
	CreateBalancedRingZoneLabels(playerLabels []string, neutralZones neutral_zone.Plans) []string
	CreateBalancedChainZoneLabels(playerLabels []string, neutralZones neutral_zone.Plans) []string
	CreateBalancedNeutralRingZoneLabels(neutralZones neutral_zone.Plans, playerCount int) []string
}
