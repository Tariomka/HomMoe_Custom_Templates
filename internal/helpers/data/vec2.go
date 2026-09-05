package data

import (
	"image"
	"math"
)

type Vec2[T INumeric] struct{ X, Y T }

func NewVec2[T INumeric](x, y T) Vec2[T] { return Vec2[T]{X: x, Y: y} }

func Vec2FromPoint[T INumeric](point image.Point) Vec2[T] {
	return Vec2[T]{X: T(point.X), Y: T(point.Y)}
}

func Transform[TOld, TNew INumeric](vector Vec2[TOld]) Vec2[TNew] {
	return NewVec2(TNew(vector.X), TNew(vector.Y))
}

// func (this Vec2[T]) Transform[TNew INumeric]() Vec2[TNew] {
// 	return NewVec2(TNew(this.X), TNew(this.Y))
// }

func (this Vec2[T]) Add(other Vec2[T]) Vec2[T] {
	this.X += other.X
	this.Y += other.Y
	return this
}

func (this Vec2[T]) Subtract(other Vec2[T]) Vec2[T] {
	this.X -= other.X
	this.Y -= other.Y
	return this
}

func (this Vec2[T]) MultiplyComponent(other Vec2[T]) Vec2[T] {
	this.X *= other.X
	this.Y *= other.Y
	return this
}

func (this Vec2[T]) MultiplyScalar(scalar T) Vec2[T] {
	this.X *= scalar
	this.Y *= scalar
	return this
}

func (this Vec2[T]) DivideComponent(other Vec2[T]) Vec2[T] {
	this.X /= other.X
	this.Y /= other.Y
	return this
}

func (this Vec2[T]) DivideScalar(scalar T) Vec2[T] {
	this.X /= scalar
	this.Y /= scalar
	return this
}

func (this Vec2[T]) SquaredLength() T { return this.X*this.X + this.Y*this.Y }

func (this Vec2[T]) DotProduct(other Vec2[T]) T { return this.X*other.X + this.Y*other.Y }

func (this Vec2[T]) CrossProduct(other Vec2[T]) T { return this.X*other.Y - this.Y*other.X }

func (this Vec2[T]) Distance() float64 { return math.Hypot(float64(this.X), float64(this.Y)) }

// RotateClockwise rotates the vector 90° clockwise: ( x, y ) → ( y, -x ).
func (this Vec2[T]) RotateClockwise() Vec2[T] { return NewVec2(this.Y, -this.X) }

func (this Vec2[T]) ToPoint() image.Point { return image.Pt(int(this.X), int(this.Y)) }

func (this Vec2[T]) ToPointRounded() image.Point {
	return NewVec2(math.Round(float64(this.X)), math.Round(float64(this.Y))).
		ToPoint()
}
