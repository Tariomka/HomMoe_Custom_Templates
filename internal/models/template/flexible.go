package template

import "encoding/json"

// StringList accepts either a single JSON string or an array of strings.
// Several `mandatoryContent` / `contentCountLimits` zone fields are declared as
// a single string in some templates (e.g. "Christmas Tree", "Crossroads") and as
// an array in others. Marshal always re-emits as an array for consistency.
type StringList []string

// UnmarshalJSON allows decoding from either `"x"` or `["x", "y"]`.
func (s *StringList) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*s = nil
		return nil
	}
	switch data[0] {
	case '"':
		var single string
		if err := json.Unmarshal(data, &single); err != nil {
			return err
		}
		*s = []string{single}
	case 'n':
		*s = nil
	default:
		var arr []string
		if err := json.Unmarshal(data, &arr); err != nil {
			return err
		}
		*s = arr
	}
	return nil
}

// BonusList accepts either a single Bonus object or an array of Bonus objects.
// The "Wastelands" template declares `bonuses` as a single object rather than an
// array; this wrapper transparently normalises that case.
type BonusList []Bonus

// UnmarshalJSON allows decoding from either `{...}` or `[{...}, ...]`.
func (b *BonusList) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*b = nil
		return nil
	}
	switch data[0] {
	case '{':
		var single Bonus
		if err := json.Unmarshal(data, &single); err != nil {
			return err
		}
		*b = []Bonus{single}
	case 'n':
		*b = nil
	default:
		var arr []Bonus
		if err := json.Unmarshal(data, &arr); err != nil {
			return err
		}
		*b = arr
	}
	return nil
}
