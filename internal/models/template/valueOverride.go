package template

// ValueOverride overrides the default guard value of a specific object SID, optionally per variant index.
type ValueOverride struct {
	SID        string `json:"sid"`
	Variant    int    `json:"variant"`
	GuardValue int    `json:"guardValue"`
}
