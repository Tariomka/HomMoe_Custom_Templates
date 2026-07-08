# Unit Test Observations — code unreachable or untestable by unit tests

This file tracks implementation code that cannot be (fully) unit-tested
through public entry points, per the rule: never add helpers/seams to
implementation code just to make it testable. Each entry needs manual
investigation by the maintainer.

Format: `path` — reason — suggested action.

## Untestable / unreachable code

(populated during the unit test refactoring, phases 1-8)

- `internal/entities/template/template_rule/winConditions.go` `MergeWinConditionsIfDoesNotExist` — the four
  error-return branches (marshal of destination/source, unmarshal into maps, marshal of merged map) are
  unreachable: `WinConditions` contains only JSON-safe field types, so `json.Marshal`/`json.Unmarshal` on it
  can never fail. Function capped at ~77% coverage. Suggested action: none (defensive guards).
- `internal/models/config/config_inner/bonusEntry.go` `GetHash` — the fallback branch when `json.Marshal`
  fails is unreachable: `BonusEntry` holds only string/int fields, so marshalling never errors. Capped at 80%.
  Suggested action: none (defensive guard).

## Gio-UI-heavy files (covered by integration suite, not unit-testable)

(populated in phase 8)

## Dead code found while testing

(populated during the unit test refactoring)

- `internal/models/config/config_inner/bonusPresetType.go` `parseBonusPresetType` — private function with
  ZERO callers anywhere in the repo (the old `ParseBonusesJSON`/`SerializeBonuses` string round-trip it served
  was removed when `EditorStateDto.BonusesJSON` became `[]config_inner.BonusEntry` serialized via std json).
  Unreachable from any public entry point, stays 0%. Suggested action: delete the function (and the stale
  "see ParseBonusesJSON" comment in `internal/dtos/editorStateDto.go`).
