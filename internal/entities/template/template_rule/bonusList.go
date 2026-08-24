package template_rule

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
)

// BonusList accepts either a single Bonus object or an array of Bonus objects.
// The "Wastelands" template declares `bonuses` as a single object rather than an
// array; this wrapper transparently normalizes that case.
type BonusList []Bonus

// UnmarshalJSONFrom allows decoding from either `{...}` or `[{...}, ...]`.
func (this *BonusList) UnmarshalJSONFrom(decoder *jsontext.Decoder) error {
	if decoder.PeekKind() != '{' {
		var list []Bonus
		if err := json.UnmarshalDecode(decoder, &list); err != nil {
			return err
		}
		*this = list

		return nil
	}

	var single Bonus
	if err := json.UnmarshalDecode(decoder, &single); err != nil {
		return err
	}
	*this = []Bonus{single}

	return nil
}

// UnmarshalJSON is the byte-slice entry point for decoders that do not use
// UnmarshalJSONFrom, which takes precedence wherever both are available.
func (this *BonusList) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*this = nil

		return nil
	}

	return this.UnmarshalJSONFrom(jsontext.NewDecoder(bytes.NewReader(data)))
}
