package data

type Vec4[T INumeric] struct{ X, Y, Z, W T }

func NewVec4[T INumeric](x, y, z, w T) Vec4[T] { return Vec4[T]{X: x, Y: y, Z: z, W: w} }
