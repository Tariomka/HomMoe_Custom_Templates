package data

type Vec3[T INumeric] struct{ X, Y, Z T }

func NewVec3[T INumeric](x, y, z T) Vec3[T] { return Vec3[T]{X: x, Y: y, Z: z} }
