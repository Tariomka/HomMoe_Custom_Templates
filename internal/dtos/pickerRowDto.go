package dtos

// PickerRowDto is one line of a picker's flattened, filtered row model: either
// a group header (with the number of matching entries below it) or a leaf entry.
type PickerRowDto struct {
	IsGroupHeader   bool
	Group           string
	GroupMatchCount int
	Entry           PickerEntryDto
}
