package dtos

// PickerEntryDto is one selectable row of a multi-select picker, already mapped
// out of its source catalogue and ready to be searched and rendered.
type PickerEntryDto struct {
	ID       string
	Group    string // category / school; ignored when the picker is flat
	Label    string // primary display text
	Badge    string // optional leading badge (e.g. "[T3]")
	Trailing string // optional dim trailing text (e.g. the raw SID)
	Haystack string // lowercased search text
}
