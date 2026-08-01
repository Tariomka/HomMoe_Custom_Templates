package models

type ZoneBiomeMatchPolicy int

const (
	ZoneBiomeMatchZone ZoneBiomeMatchPolicy = iota
	ZoneBiomeMatchPrimaryMainObjectWhenPresent
)
