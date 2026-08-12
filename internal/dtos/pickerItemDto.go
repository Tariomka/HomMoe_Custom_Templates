package dtos

// PickerItemDto is one bannable-item catalogue row handed to the picker
// service. The catalogue itself lives in the GUI layer, which may not be
// imported from internal, so the rows travel as data.
type PickerItemDto struct {
	Sid      string
	Name     string
	Category string
}
