package data

type Vec3[T Numeric] struct{ X, Y, Z T }

func NewVec3[T Numeric](x, y, z T) Vec3[T] { return Vec3[T]{X: x, Y: y, Z: z} }
