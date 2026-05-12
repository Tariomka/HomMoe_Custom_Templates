package gamerules

import "encoding/json"

// UnmarshalJSON allows GameRules to also decode templates where the WinConditions
// fields appear flat inside `gameRules` (alongside `bonuses`, `encounterHoles`,
// etc.) instead of being nested under `winConditions`. Notable example: the
// "Zookeeper" template. The nested `winConditions` block, when present, still
// takes precedence; flat sibling fields only fill values that the nested form
// omitted.
func (g *GameRules) UnmarshalJSON(data []byte) error {
	// Decode the standard GameRules layout (including any nested winConditions).
	type alias GameRules
	if err := json.Unmarshal(data, (*alias)(g)); err != nil {
		return err
	}

	// Also try to populate WinConditions from the same JSON object treated as a
	// flat WinConditions blob (Zookeeper-style flat templates).
	var flat WinConditions
	if err := json.Unmarshal(data, &flat); err != nil {
		// flat decode failure is non-fatal - the regular decode already succeeded.
		return nil
	}
	mergeWinConditionsZero(&g.WinConditions, flat)
	return nil
}

// mergeWinConditionsZero copies fields from src into dst, but only for fields
// where dst currently holds its zero value (so a nested `winConditions` block
// always wins over flat sibling keys).
func mergeWinConditionsZero(dst *WinConditions, src WinConditions) {
	dstBytes, err := json.Marshal(*dst)
	if err != nil {
		return
	}
	srcBytes, err := json.Marshal(src)
	if err != nil {
		return
	}
	var dstMap, srcMap map[string]json.RawMessage
	if err := json.Unmarshal(dstBytes, &dstMap); err != nil {
		return
	}
	if err := json.Unmarshal(srcBytes, &srcMap); err != nil {
		return
	}
	// Because every WinConditions field uses `omitempty`, any field present in
	// srcMap but absent from dstMap is a zero-valued dst field that should be
	// filled in from src.
	changed := false
	for k, v := range srcMap {
		if _, ok := dstMap[k]; !ok {
			dstMap[k] = v
			changed = true
		}
	}
	if !changed {
		return
	}
	merged, err := json.Marshal(dstMap)
	if err != nil {
		return
	}
	_ = json.Unmarshal(merged, dst)
}
