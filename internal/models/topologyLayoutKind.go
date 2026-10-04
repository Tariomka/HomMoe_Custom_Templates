package models

type TopologyLayoutKind uint8

const (
	TopologyLayoutGeneric TopologyLayoutKind = iota // draws zones on an outer ring around any explicitly named hub.
	TopologyLayoutScatter
	TopologyLayoutFixedGeometry
)
