package variant_content

import "github.com/Tariomka/hommoe_custom_templates/internal/entities"

type OrientationBuilder struct {
	item entities.Orientation
}

func NewOrientationBuilder() *OrientationBuilder {
	return &OrientationBuilder{item: entities.Orientation{}}
}

func (this *OrientationBuilder) WithMode(mode string) *OrientationBuilder {
	this.item.Mode = mode
	return this
}
func (this *OrientationBuilder) WithZeroAngleZone(zone string) *OrientationBuilder {
	this.item.ZeroAngleZone = zone
	return this
}
func (this *OrientationBuilder) WithBaseAngleMin(angle int) *OrientationBuilder {
	this.item.BaseAngleMin = angle
	return this
}
func (this *OrientationBuilder) WithBaseAngleMax(angle int) *OrientationBuilder {
	this.item.BaseAngleMax = angle
	return this
}
func (this *OrientationBuilder) WithRandomAngleAmplitude(amplitude int) *OrientationBuilder {
	this.item.RandomAngleAmplitude = amplitude
	return this
}
func (this *OrientationBuilder) WithRandomAngleStep(step int) *OrientationBuilder {
	this.item.RandomAngleStep = step
	return this
}
func (this *OrientationBuilder) Build() entities.Orientation { return this.item }
