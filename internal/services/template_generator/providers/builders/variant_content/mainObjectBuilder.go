package variant_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

var (
	castleQualities = registry.GetBuildingsConstructionSidValues()
	objectTypes     = registry.GetMainObjectTypeValues()
	placements      = registry.GetPlacementValues()
)

type MainObjectBuilder struct {
	item entities.MainObject
}

func NewObjectBuilder() *MainObjectBuilder { return &MainObjectBuilder{item: entities.MainObject{}} }

func (this *MainObjectBuilder) WithTypeSpawn() *MainObjectBuilder {
	return this.withType(objectTypes.Spawn)
}
func (this *MainObjectBuilder) WithTypeCity() *MainObjectBuilder {
	return this.withType(objectTypes.City)
}
func (this *MainObjectBuilder) WithSpawn(spawn string) *MainObjectBuilder {
	this.item.Spawn = spawn
	return this
}
func (this *MainObjectBuilder) WithOwner(owner string) *MainObjectBuilder {
	this.item.Owner = owner
	return this
}
func (this *MainObjectBuilder) WithNoGuardWhenOwned() *MainObjectBuilder {
	this.item.RemoveGuardIfHasOwner = true
	return this
}
func (this *MainObjectBuilder) WithGuardChance(chance float64) *MainObjectBuilder {
	this.item.GuardChance = helpers.ClampFloat(chance, 0, 1)
	return this
}
func (this *MainObjectBuilder) WithGuardValue(value int) *MainObjectBuilder {
	this.item.GuardValue = value
	return this
}
func (this *MainObjectBuilder) WithGuardRandomization(randomization float64) *MainObjectBuilder {
	this.item.GuardRandomization = randomization
	return this
}
func (this *MainObjectBuilder) WithGuardWeeklyIncrement(increment float64) *MainObjectBuilder {
	this.item.GuardWeeklyIncrement = increment
	return this
}
func (this *MainObjectBuilder) WithCastleQuality(sid string) *MainObjectBuilder {
	this.item.BuildingsConstructionSid = sid
	return this
}
func (this *MainObjectBuilder) WithCastleQualityDefault() *MainObjectBuilder {
	return this.WithCastleQuality(castleQualities.Default)
}
func (this *MainObjectBuilder) WithCastleQualityPoor() *MainObjectBuilder {
	return this.WithCastleQuality(castleQualities.Poor)
}
func (this *MainObjectBuilder) WithCastleQualityMedium() *MainObjectBuilder {
	return this.WithCastleQuality(castleQualities.Medium)
}
func (this *MainObjectBuilder) WithCastleQualityRich() *MainObjectBuilder {
	return this.WithCastleQuality(castleQualities.Rich)
}
func (this *MainObjectBuilder) WithCastleQualityExtraRich() *MainObjectBuilder {
	return this.WithCastleQuality(castleQualities.ExtraRich)
}
func (this *MainObjectBuilder) WithCastleQualityUltraRich() *MainObjectBuilder {
	return this.WithCastleQuality(castleQualities.UltraRich)
}
func (this *MainObjectBuilder) WithFaction(matchOrFromListType string, arguments ...string) *MainObjectBuilder {
	this.item.Faction = &entities.TypedRef{Type: matchOrFromListType, Args: arguments}
	return this
}
func (this *MainObjectBuilder) WithFactions(factions ...string) *MainObjectBuilder {
	this.item.Factions = append(this.item.Factions, factions...)
	return this
}
func (this *MainObjectBuilder) WithPlacementCenter() *MainObjectBuilder {
	return this.withPlacement(placements.Center)
}
func (this *MainObjectBuilder) WithPlacementConnection() *MainObjectBuilder {
	return this.withPlacement(placements.Connection)
}
func (this *MainObjectBuilder) WithPlacementNearZone() *MainObjectBuilder {
	return this.withPlacement(placements.NearZone)
}
func (this *MainObjectBuilder) WithPlacementUniform() *MainObjectBuilder {
	return this.withPlacement(placements.Uniform)
}
func (this *MainObjectBuilder) WithPlacementArgs(arguments ...string) *MainObjectBuilder {
	this.item.PlacementArgs = append(this.item.PlacementArgs, arguments...)
	return this
}
func (this *MainObjectBuilder) WithHoldCityWinCon() *MainObjectBuilder {
	this.item.HoldCityWinCon = true
	return this
}
func (this *MainObjectBuilder) WithKeyObject() *MainObjectBuilder {
	this.item.IsKeyObject = true
	return this
}
func (this *MainObjectBuilder) WithWeeklyUnitIncrement() *MainObjectBuilder {
	this.item.EnableWeeklyUnitIncrement = true
	return this
}
func (this *MainObjectBuilder) WithInitialUnitIncrement(initialIncrement int) *MainObjectBuilder {
	this.item.InitialUnitIncrement = initialIncrement
	return this
}
func (this *MainObjectBuilder) Build() entities.MainObject { return this.item }

func (this *MainObjectBuilder) withType(objectType string) *MainObjectBuilder {
	this.item.Type = objectType
	return this
}
func (this *MainObjectBuilder) withPlacement(placement string) *MainObjectBuilder {
	this.item.Placement = placement
	return this
}
