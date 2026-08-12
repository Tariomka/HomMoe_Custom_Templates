package models

type TopologyCapabilities struct {
	LayoutKind            TopologyLayoutKind
	UsesHub               bool
	UsesGeneratorPosition bool
	UsesGeneratorRing     bool
}
