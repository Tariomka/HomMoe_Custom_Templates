package template_variant

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
)

// StringList accepts either a single JSON string or an array of strings.
// Several `mandatoryContent` / `contentCountLimits` zone fields are declared as
// a single string in some templates (e.g. "Christmas Tree", "Crossroads") and as
// an array in others. Marshal always re-emits as an array for consistency.
type StringList []string

// UnmarshalJSONFrom allows decoding from either `"x"` or `["x", "y"]`.
func (this *StringList) UnmarshalJSONFrom(decoder *jsontext.Decoder) error {
	if decoder.PeekKind() != '"' {
		var list []string
		if err := json.UnmarshalDecode(decoder, &list); err != nil {
			return err
		}
		*this = list

		return nil
	}

	var single string
	if err := json.UnmarshalDecode(decoder, &single); err != nil {
		return err
	}
	*this = []string{single}

	return nil
}

// UnmarshalJSON is the byte-slice entry point for decoders that do not use
// UnmarshalJSONFrom, which takes precedence wherever both are available.
func (this *StringList) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*this = nil

		return nil
	}

	return this.UnmarshalJSONFrom(jsontext.NewDecoder(bytes.NewReader(data)))
}
