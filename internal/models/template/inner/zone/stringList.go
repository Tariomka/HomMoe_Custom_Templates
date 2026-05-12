package zone

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
