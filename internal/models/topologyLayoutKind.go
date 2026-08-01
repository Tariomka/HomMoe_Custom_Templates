package models

type TopologyLayoutKind uint8

const (
	TopologyLayoutRingHub TopologyLayoutKind = iota
	TopologyLayoutScatter
	TopologyLayoutFixedGeometry
)
