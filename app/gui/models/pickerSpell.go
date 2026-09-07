package models

// PickerSpell is one spell catalogue row handed to the picker view model.
// SchoolDisplayName is the resolved, human-readable school name; it is empty
// when the catalogue has no display name for School.
type PickerSpell struct {
	Sid               string
	Name              string
	School            string
	SchoolDisplayName string
	Tier              int
}
