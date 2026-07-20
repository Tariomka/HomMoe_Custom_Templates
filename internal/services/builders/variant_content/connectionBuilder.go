package variant_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

type ConnectionBuilder struct {
	item entities.Connection
}

func NewConnectionBuilder() *ConnectionBuilder {
	return &ConnectionBuilder{item: entities.Connection{}}
}

func (this *ConnectionBuilder) WithName(name string) *ConnectionBuilder {
	this.item.Name = name
	return this
}

func (this *ConnectionBuilder) WithFrom(from string) *ConnectionBuilder {
	this.item.From = from
	return this
}

func (this *ConnectionBuilder) WithTo(to string) *ConnectionBuilder {
	this.item.To = to
	return this
}

func (this *ConnectionBuilder) WithConnectionTypeDirect() *ConnectionBuilder {
	return this.withConnectionType(registry.GetConnectionTypeValues().Direct)
}
func (this *ConnectionBuilder) WithConnectionTypePortal() *ConnectionBuilder {
	return this.withConnectionType(registry.GetConnectionTypeValues().Portal)
}
func (this *ConnectionBuilder) WithConnectionTypeProximity() *ConnectionBuilder {
	return this.withConnectionType(registry.GetConnectionTypeValues().Proximity)
}

func (this *ConnectionBuilder) WithSimTurnSquad() *ConnectionBuilder {
	this.item.SimTurnSquad = true
	return this
}

func (this *ConnectionBuilder) WithRoad(road bool) *ConnectionBuilder {
	this.item.Road = &road
	return this
}

func (this *ConnectionBuilder) WithGuardZone(guardZone string) *ConnectionBuilder {
	this.item.GuardZone = guardZone
	return this
}

func (this *ConnectionBuilder) WithGuardEscape(guardEscape bool) *ConnectionBuilder {
	this.item.GuardEscape = guardEscape
	return this
}

func (this *ConnectionBuilder) WithGuardValue(guardValue int) *ConnectionBuilder {
	this.item.GuardValue = guardValue
	return this
}

func (this *ConnectionBuilder) WithGuardRandomization(guardRandomization float64) *ConnectionBuilder {
	this.item.GuardRandomization = guardRandomization
	return this
}

func (this *ConnectionBuilder) WithGuardWeeklyIncrement(guardWeeklyIncrement float64) *ConnectionBuilder {
	this.item.GuardWeeklyIncrement = guardWeeklyIncrement
	return this
}

func (this *ConnectionBuilder) WithGatePlacementCenter() *ConnectionBuilder {
	return this.withGatePlacement(registry.GetGatePlacementValues().Center)
}

func (this *ConnectionBuilder) WithLength(length float64) *ConnectionBuilder {
	this.item.Length = length
	return this
}

func (this *ConnectionBuilder) WithGuardMatchGroup(guardMatchGroup string) *ConnectionBuilder {
	this.item.GuardMatchGroup = guardMatchGroup
	return this
}

func (this *ConnectionBuilder) WithPortalPlacementRulesFrom(rules ...entities.PlacementRule) *ConnectionBuilder {
	this.item.PortalPlacementRulesFrom = append(this.item.PortalPlacementRulesFrom, rules...)
	return this
}
func (this *ConnectionBuilder) WithPortalPlacementRulesTo(rules ...entities.PlacementRule) *ConnectionBuilder {
	this.item.PortalPlacementRulesTo = append(this.item.PortalPlacementRulesTo, rules...)
	return this
}

func (this *ConnectionBuilder) Build() entities.Connection { return this.item }

func (this *ConnectionBuilder) withConnectionType(connectionType string) *ConnectionBuilder {
	this.item.ConnectionType = connectionType
	return this
}

func (this *ConnectionBuilder) withGatePlacement(gatePlacement string) *ConnectionBuilder {
	this.item.GatePlacement = gatePlacement
	return this
}
