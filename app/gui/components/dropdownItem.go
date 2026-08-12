package components

import "gioui.org/widget"

type dropdownItem struct {
	label string
	row   widget.Clickable
}

func newDropdownItem(label string) *dropdownItem {
	return &dropdownItem{label: label}
}
