package variant_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

var layoutValues = registry.GetLayoutValues()

type ZoneBuilder struct {
	item entities.Zone
}

func NewZoneBuilder() *ZoneBuilder { return &ZoneBuilder{item: entities.Zone{}} }

func (this *ZoneBuilder) WithName(name string) *ZoneBuilder {
	this.item.Name = name
	return this
}

func (this *ZoneBuilder) WithSize(size float64) *ZoneBuilder {
	this.item.Size = size
	return this
}

func (this *ZoneBuilder) WithLayout(layout string) *ZoneBuilder {
	this.item.Layout = layout
	return this
}
func (this *ZoneBuilder) WithLayoutSpawns() *ZoneBuilder {
	return this.WithLayout(layoutValues.Spawns)
}
func (this *ZoneBuilder) WithLayoutCenter() *ZoneBuilder {
	return this.WithLayout(layoutValues.Center)
}

func (this *ZoneBuilder) WithGuardCutoffValue(value int) *ZoneBuilder {
	this.item.GuardCutoffValue = value
	return this
}

func (this *ZoneBuilder) WithGuardRandomization(randomization float64) *ZoneBuilder {
	this.item.GuardRandomization = randomization
	return this
}

func (this *ZoneBuilder) WithGuardMultiplier(multiplier float64) *ZoneBuilder {
	this.item.GuardMultiplier = multiplier
	return this
}

func (this *ZoneBuilder) WithGuardWeeklyIncrement(increment float64) *ZoneBuilder {
	this.item.GuardWeeklyIncrement = increment
	return this
}

func (this *ZoneBuilder) WithGuardReactionDistribution(distribution []int) *ZoneBuilder {
	this.item.GuardReactionDistribution = distribution
	return this
}

func (this *ZoneBuilder) WithDiplomacyModifier(modifier float64) *ZoneBuilder {
	this.item.DiplomacyModifier = modifier
	return this
}

func (this *ZoneBuilder) WithGuardedContentPool(pool []string) *ZoneBuilder {
	this.item.GuardedContentPool = pool
	return this
}

func (this *ZoneBuilder) WithUnguardedContentPool(pool []string) *ZoneBuilder {
	this.item.UnguardedContentPool = pool
	return this
}

func (this *ZoneBuilder) WithResourcesContentPool(pool []string) *ZoneBuilder {
	this.item.ResourcesContentPool = pool
	return this
}

func (this *ZoneBuilder) WithMandatoryContent(content ...string) *ZoneBuilder {
	this.item.MandatoryContent = content
	return this
}

func (this *ZoneBuilder) WithContentCountLimits(limits []string) *ZoneBuilder {
	this.item.ContentCountLimits = limits
	return this
}

func (this *ZoneBuilder) WithGuardedContentValue(value int) *ZoneBuilder {
	this.item.GuardedContentValue = value
	return this
}

func (this *ZoneBuilder) WithGuardedContentValuePerArea(value int) *ZoneBuilder {
	this.item.GuardedContentValuePerArea = value
	return this
}

func (this *ZoneBuilder) WithUnguardedContentValue(value int) *ZoneBuilder {
	this.item.UnguardedContentValue = value
	return this
}

func (this *ZoneBuilder) WithUnguardedContentValuePerArea(value int) *ZoneBuilder {
	this.item.UnguardedContentValuePerArea = value
	return this
}

func (this *ZoneBuilder) WithResourcesValue(value int) *ZoneBuilder {
	this.item.ResourcesValue = value
	return this
}

func (this *ZoneBuilder) WithResourcesValuePerArea(value int) *ZoneBuilder {
	this.item.ResourcesValuePerArea = value
	return this
}

func (this *ZoneBuilder) WithMainObjects(objects []entities.MainObject) *ZoneBuilder {
	this.item.MainObjects = objects
	return this
}

// WithBiome sets the same biome type for zone, content, and meta objects biomes. Effectively same as
//
//	builder.WithZoneBiome(biome).WithContentBiome(biome).WithMetaObjectsBiome(biome)
func (this *ZoneBuilder) WithBiome(biome entities.TypedRef) *ZoneBuilder {
	return this.WithZoneBiome(biome).WithContentBiome(biome).WithMetaObjectsBiome(biome)
}

// WithBiomeMatchMainObject sets "MatchMainObject" biome type for zone, content, and meta objects biomes. Effectively same as
//
//	builder.WithBiome(NewRefBuilder().BuildBiomeMatchMainObjectType(args...))
func (this *ZoneBuilder) WithBiomeMatchMainObject(args ...string) *ZoneBuilder {
	return this.WithBiome(NewRefBuilder().BuildBiomeMatchMainObjectType(args...))
}

// WithBiomeMatchZone sets "MatchZone" biome type for zone, content, and meta objects biomes. Effectively same as
//
//	builder.WithBiome(NewRefBuilder().BuildBiomeMatchZoneType(args...))
func (this *ZoneBuilder) WithBiomeMatchZone(args ...string) *ZoneBuilder {
	return this.WithBiome(NewRefBuilder().BuildBiomeMatchZoneType(args...))
}

func (this *ZoneBuilder) WithZoneBiome(biome entities.TypedRef) *ZoneBuilder {
	this.item.ZoneBiome = biome
	return this
}
func (this *ZoneBuilder) WithContentBiome(biome entities.TypedRef) *ZoneBuilder {
	this.item.ContentBiome = biome
	return this
}
func (this *ZoneBuilder) WithMetaObjectsBiome(biome entities.TypedRef) *ZoneBuilder {
	this.item.MetaObjectsBiome = biome
	return this
}
func (this *ZoneBuilder) WithCrossroadsPosition(position int) *ZoneBuilder {
	this.item.CrossroadsPosition = &position
	return this
}
func (this *ZoneBuilder) WithRoads(roads []entities.Road) *ZoneBuilder {
	this.item.Roads = roads
	return this
}

func (this *ZoneBuilder) WithEncounterHolesSettings(settings entities.EncounterHolesSettings) *ZoneBuilder {
	this.item.EncounterHolesSettings = &settings
	return this
}
func (this *ZoneBuilder) WithGeneratorPosition(position [2]float64) *ZoneBuilder {
	this.item.GeneratorPosition = &position
	return this
}
func (this *ZoneBuilder) WithGeneratorRing(ring int) *ZoneBuilder {
	this.item.GeneratorRing = &ring
	return this
}
func (this *ZoneBuilder) WithRandomHireEnableWeeklyUnitIncrement(values []bool) *ZoneBuilder {
	this.item.RandomHireEnableWeeklyUnitIncrement = values
	return this
}
func (this *ZoneBuilder) WithRandomHireInitialUnitIncrement(values []int) *ZoneBuilder {
	this.item.RandomHireInitialUnitIncrement = values
	return this
}
func (this *ZoneBuilder) Build() entities.Zone { return this.item }
