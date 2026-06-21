package template_rule

import "encoding/json"

// BonusList accepts either a single Bonus object or an array of Bonus objects.
// The "Wastelands" template declares `bonuses` as a single object rather than an
// array; this wrapper transparently normalizes that case.
type BonusList []Bonus

// UnmarshalJSON allows decoding from either `{...}` or `[{...}, ...]`.
func (this *BonusList) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*this = nil
		return nil
	}
	switch data[0] {
	case '{':
		var single Bonus
		if err := json.Unmarshal(data, &single); err != nil {
			return err
		}
		*this = []Bonus{single}
	case 'n':
		*this = nil
	default:
		var arr []Bonus
		if err := json.Unmarshal(data, &arr); err != nil {
			return err
		}
		*this = arr
	}
	return nil
}
