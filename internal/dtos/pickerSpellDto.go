package dtos

// PickerSpellDto is one spell catalogue row handed to the picker service.
// SchoolDisplayName is the resolved, human-readable school name; it is empty
// when the GUI catalogue has no display name for School.
type PickerSpellDto struct {
	Sid               string
	Name              string
	School            string
	SchoolDisplayName string
	Tier              int
}
