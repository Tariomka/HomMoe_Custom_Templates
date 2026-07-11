package models

// SidMapping pairs a string ID (sid) with a human-readable name.
type SidMapping struct {
	Sid  string `json:"sid"`
	Name string `json:"name"`
}
