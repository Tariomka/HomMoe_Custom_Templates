package variant_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

var waterTypes = registry.GetWaterTypeValues()

type BorderBuilder struct {
	item entities.Border
}

func NewBorderBuilder() *BorderBuilder { return &BorderBuilder{item: entities.Border{}} }

func (this *BorderBuilder) WithCornerRadius(radius float64) *BorderBuilder {
	this.item.CornerRadius = radius
	return this
}
func (this *BorderBuilder) WithObstaclesWidth(width int) *BorderBuilder {
	this.item.ObstaclesWidth = width
	return this
}
func (this *BorderBuilder) WithObstaclesNoise(amplitude float64, frequency int) *BorderBuilder {
	this.item.ObstaclesNoise = []entities.Noise{{Amplitude: amplitude, Frequency: frequency}}
	return this
}
func (this *BorderBuilder) WithWaterWidth(width int) *BorderBuilder {
	this.item.WaterWidth = width
	return this
}
func (this *BorderBuilder) WithWaterNoise(amplitude float64, frequency int) *BorderBuilder {
	this.item.WaterNoise = []entities.Noise{{Amplitude: amplitude, Frequency: frequency}}
	return this
}
func (this *BorderBuilder) WithWaterTypeWaterGrass() *BorderBuilder {
	return this.withWaterType(waterTypes.WaterGrass)
}
func (this *BorderBuilder) Build() entities.Border { return this.item }

func (this *BorderBuilder) withWaterType(waterType string) *BorderBuilder {
	this.item.WaterType = waterType
	return this
}
