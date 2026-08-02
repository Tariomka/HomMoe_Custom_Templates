package variant_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

var castleQualities = registry.GetBuildingsConstructionSidValues()

type MainObjectBuilder struct {
	item entities.MainObject
}

func NewObjectBuilder() *MainObjectBuilder { return &MainObjectBuilder{item: entities.MainObject{}} }

func (this *MainObjectBuilder) WithTypeSpawn() *MainObjectBuilder {
	return this.withType(registry.GetMainObjectTypeValues().Spawn)
}
func (this *MainObjectBuilder) WithTypeCity() *MainObjectBuilder {
	return this.withType(registry.GetMainObjectTypeValues().City)
}
func (this *MainObjectBuilder) WithTypeAbandonedOutpost() *MainObjectBuilder {
	return this.withType(registry.GetMainObjectTypeValues().AbandonedOutpost)
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
	this.item.GuardChance = helpers.Clamp(chance, 0, 1)
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

func (this *MainObjectBuilder) WithFaction(factionMatchType string, arguments ...string) *MainObjectBuilder {
	this.item.Faction = new(NewRefBuilder().WithType(factionMatchType).WithArgs(arguments...).Build())
	return this
}
func (this *MainObjectBuilder) WithFactionFromList() *MainObjectBuilder {
	return this.WithFaction(registry.GetFactionTypeValues().FromList)
}
func (this *MainObjectBuilder) WithFactionMatch() *MainObjectBuilder {
	return this.WithFaction(registry.GetFactionTypeValues().Match, "0")
}

func (this *MainObjectBuilder) WithFactions(factions ...string) *MainObjectBuilder {
	this.item.Factions = append(this.item.Factions, factions...)
	return this
}

func (this *MainObjectBuilder) WithPlacementCenter() *MainObjectBuilder {
	return this.withPlacement(registry.GetPlacementValues().Center)
}
func (this *MainObjectBuilder) WithPlacementConnection() *MainObjectBuilder {
	return this.withPlacement(registry.GetPlacementValues().Connection)
}
func (this *MainObjectBuilder) WithPlacementNearZone() *MainObjectBuilder {
	return this.withPlacement(registry.GetPlacementValues().NearZone)
}
func (this *MainObjectBuilder) WithPlacementUniform() *MainObjectBuilder {
	return this.withPlacement(registry.GetPlacementValues().Uniform)
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
