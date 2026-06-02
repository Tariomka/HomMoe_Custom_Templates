package variant_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

const (
	typeSpawn = "Spawn"
	typeCity  = "City"

	castleQualityDefault   = "default_buildings_construction"
	castleQualityPoor      = "poor_buildings_construction"
	castleQualityMedium    = "medium_buildings_construction"
	castleQualityRich      = "rich_buildings_construction"
	castleQualityExtraRich = "extra_rich_buildings_construction"
	castleQualityUltraRich = "ultra_rich_buildings_construction"
	castleQualityArcade    = "arcade_buildings_construction"
	castleQualityArmy      = "army_buildings_construction"
	castleQualityChosenOne = "chosen_one_buildings_construction"
	castleQualityFull      = "full_buildings_construction"
	castleQualityMassacre  = "massacre_buildings_construction"
	castleQualitySiege     = "siege_buildings_construction"

	placementCenter     = "Center"
	placementConnection = "Connection"
	placementNearZone   = "NearZone"
	placementUniform    = "Uniform"
)

type MainObjectBuilder struct {
	item template.MainObject
}

func NewObjectBuilder() *MainObjectBuilder { return &MainObjectBuilder{item: template.MainObject{}} }

func (this *MainObjectBuilder) WithTypeSpawn() *MainObjectBuilder { return this.withType(typeSpawn) }
func (this *MainObjectBuilder) WithTypeCity() *MainObjectBuilder  { return this.withType(typeCity) }
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
func (this *MainObjectBuilder) WithCastleQualityDefault() *MainObjectBuilder {
	return this.withBuildingsQuality(castleQualityDefault)
}
func (this *MainObjectBuilder) WithCastleQualityPoor() *MainObjectBuilder {
	return this.withBuildingsQuality(castleQualityPoor)
}
func (this *MainObjectBuilder) WithCastleQualityMedium() *MainObjectBuilder {
	return this.withBuildingsQuality(castleQualityMedium)
}
func (this *MainObjectBuilder) WithCastleQualityRich() *MainObjectBuilder {
	return this.withBuildingsQuality(castleQualityRich)
}
func (this *MainObjectBuilder) WithCastleQualityExtraRich() *MainObjectBuilder {
	return this.withBuildingsQuality(castleQualityExtraRich)
}
func (this *MainObjectBuilder) WithCastleQualityUltraRich() *MainObjectBuilder {
	return this.withBuildingsQuality(castleQualityUltraRich)
}
func (this *MainObjectBuilder) WithCastleQualityArcade() *MainObjectBuilder {
	return this.withBuildingsQuality(castleQualityArcade)
}
func (this *MainObjectBuilder) WithCastleQualityArmy() *MainObjectBuilder {
	return this.withBuildingsQuality(castleQualityArmy)
}
func (this *MainObjectBuilder) WithCastleQualityChosenOne() *MainObjectBuilder {
	return this.withBuildingsQuality(castleQualityChosenOne)
}
func (this *MainObjectBuilder) WithCastleQualityFull() *MainObjectBuilder {
	return this.withBuildingsQuality(castleQualityFull)
}
func (this *MainObjectBuilder) WithCastleQualityMassacre() *MainObjectBuilder {
	return this.withBuildingsQuality(castleQualityMassacre)
}
func (this *MainObjectBuilder) WithCastleQualitySiege() *MainObjectBuilder {
	return this.withBuildingsQuality(castleQualitySiege)
}
func (this *MainObjectBuilder) WithFaction(matchOrFromListType string, arguments ...string) *MainObjectBuilder {
	this.item.Faction = &template.TypedRef{Type: matchOrFromListType, Args: arguments}
	return this
}
func (this *MainObjectBuilder) WithFactions(factions ...string) *MainObjectBuilder {
	this.item.Factions = append(this.item.Factions, factions...)
	return this
}
func (this *MainObjectBuilder) WithPlacementCenter() *MainObjectBuilder {
	return this.withPlacement(placementCenter)
}
func (this *MainObjectBuilder) WithPlacementConnection() *MainObjectBuilder {
	return this.withPlacement(placementConnection)
}
func (this *MainObjectBuilder) WithPlacementNearZone() *MainObjectBuilder {
	return this.withPlacement(placementNearZone)
}
func (this *MainObjectBuilder) WithPlacementUniform() *MainObjectBuilder {
	return this.withPlacement(placementUniform)
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
func (this *MainObjectBuilder) Build() template.MainObject { return this.item }

func (this *MainObjectBuilder) withType(objectType string) *MainObjectBuilder {
	this.item.Type = objectType
	return this
}
func (this *MainObjectBuilder) withBuildingsQuality(sid string) *MainObjectBuilder {
	this.item.BuildingsConstructionSid = sid
	return this
}
func (this *MainObjectBuilder) withPlacement(placement string) *MainObjectBuilder {
	this.item.Placement = placement
	return this
}
