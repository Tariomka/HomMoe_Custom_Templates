# Batch M: coordinated persistence format (review §2.4 C3, §2.5, §2.6, §1.16)

Ship editor-state schema v3 in one coordinated change: typed ban/override lists (§2.5),
sectioned settings with identity/map/schema fields at the root (§2.6), and the C3 manual-edit
representation (structure + per-field overrides, everything else derived from the current
settings), which fixes §1.16. Plus the M-D5 to M-D8 extras. One migration chain keeps loading
every supported older schema.

Branch `AD/compact_state`, base owner HEAD `85752a8`. Design input:
[the Batch L report](../investigations/batch-l-compact-state.md) (§6, §7, §9–§13). Review items:
[§1.16, §2.4, §2.5, §2.6](../backlog/review-gpt-6-astra-09-07.md).

## For Future Agents
As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

Read [AGENTS.md](../../AGENTS.md) and [the handoff](../session-carry-forward.md) first.
Batches A–L are closed; never revisit them. Raw artifacts (baselines, benchmark output,
coverage profiles) go to gitignored `.agent/memories/batch-m/` (the owner's `air` watcher can
wipe `tmp/`). The owner reviews and commits **once, at the end** (M1): run every phase in
order, keep each phase green, and do not stop for commits between phases. Never stage, commit,
stash or switch branches.

## Owner decisions

### From the Batch L decision gate (2026-10-13)

| ID | Decision |
| --- | --- |
| M-D1 | Representation **C3**: store structure (zone name, tier, generator position/ring, manual position; `origin: added` + castles for editor-made zones; connection name, endpoints, type, `isUserAdded`). Everything else derived; stored only as `fieldOverrides` / `unset` when different. Roads are ordinary overrides. |
| M-D2 | Old full snapshots migrate as **derived when equal**, with a one-time notice. |
| M-D3 | Overrides by **value difference**; typing the default back makes the field derived. |
| M-D4 | §1.16 fixed inside M via C3; Batch Q folds into M. |
| M-D5 | Persist `generatorVersion` and show a "derived values changed" notice. |
| M-D6 | Omit default content rows behind a defaults version. |
| M-D7 | Fix the live per-frame cloning. |
| M-D8 | `tools/go.mod` golangci-lint v2.13.1 → v2.14.0. |
| M-D9 | Out of scope: compact (unindented) JSON; writing stays indented. |

### Planning answers (2026-10-13)

| ID | Decision |
| --- | --- |
| M1 | All phases run end to end; the owner reviews everything at the end (as Batch K Phases 2–4). |
| M2 | **Compact live model.** `ManualEditSettings` holds structure + typed overrides + `unset`. Apply captures; regeneration expands onto the fresh generation. No generation inside `file_service`. |
| M3 | Overrides are captured against a **fresh generation at Apply** (one `GenerateTemplate` per Apply). For non-deterministic layouts (Random topology or `RandomPortals`) zone `roads` are **always stored**. |
| M4 | A connection derives from the same-named connection of the fresh generation when its endpoints match; otherwise (user-added, unmatched) from the editor factory (`NewDefaultConnection` + `ChangeConnectionType`). |
| M5 | A zone derives from the same-named fresh zone when its profile is untouched. When re-profiled (tier **or** castle count differs from that zone) or `origin: added`, it derives from the neutral-zone factory with the stored tier and `castles` (stored only for these zones). A non-added name missing from the fresh generation is treated as added, with a warning. |
| M6 | Castle-setting reapply runs only on zones whose `mainObjects` are overridden or that are re-profiled; derived zones follow castle settings through the fresh generation. |
| M7, M26 | The derivation service lives in `internal/services/connection_editor` as a new file. That makes 5 interfaced implementations, so per AGENTS.md §4.2.2 all of the package's interfaces move to a new `connection_editor_interfaces` subpackage (behaviour-neutral; Wire regenerated). |
| M8 | Section keys: `players`, `neutralZones`, `castles`, `generation`, `gameRules`, `content`, `manualEdits`. Inner field keys keep their current names verbatim (including `isolateplayers`); only `bannedItems`, `bannedMagics`, `valueOverrides` change type. |
| M9 | Overrides are **typed structs** (one pointer field per overridable field) with a reflection parity test. |
| M10 | C3 keys as below; `generatorVersion` sits at the **root and is always written**. |
| M11 | `generatorVersion` is a dedicated integer constant, guarded by a golden-hash test that doubles as the determinism contract. |
| M12 | The M-D5 notice compares versions only (no fingerprint). |
| M13 | `content.contentDefaultsVersion` (current 1) with frozen per-version tables; default player rows omitted; unknown newer version refused; v0–v2 map to 1. **Cleared lists are written as `[]`** (fixes the pre-existing "cleared rows reload as defaults" quirk). |
| M14, M18 | M-D2 migration notice: an informational **modal with OK only**. M-D5: a status-line notice. |
| M15 | Malformed legacy ban/override lines are dropped **once at load** with line-numbered warnings shown in the load message. Generation parses nothing. |
| M16 | Old snapshots convert at the **first regeneration after load**, against the file's own settings. A save before any successful regeneration writes the lossless all-overridden v3 form. |
| M17 | That conversion **marks the document unsaved**. |
| M19 | The M-D5 notice fires only when the file has manual edits. |
| M20 | `contentDefaultsVersion` lives inside `content`. |
| M21 | M-D7: an explicit fix (the regeneration decision receives clones without the manual section) plus the compact model; a deterministic test and an ungated benchmark. |
| M22 | Commit a curated subset of the Batch L fixtures and one large hub v2 fixture under `test/test_helpers/testdata/`. |
| M23 | Coverage total ≥ the Phase 0 baseline; every new or changed function covered on every reachable branch (GUI-only code recorded in `test_observations.md`). Bonuses-panel and zone-editor goldens unchanged; the modal adds goldens; any other moved golden is explained and approved at the end. |
| M24 | Plan and implementation reviews by GPT-6.1 Sol. |
| M25 | Scope confirmed complete. |

### Out of scope
Compact JSON (M-D9); Batches N, O, P; direct `GeneratorConfig` rejection; unrelated DTO
cleanup; the RMG output schema (`internal/entities/template_entity/`); `data/`; the registry;
the output path; every settled alternative in handoff §7. Handoff §8 stays verbatim.

### Must stay preserved
Everything in handoff §3 "Settled behavior": road tri-state and explicit Portal flags,
independent internal roads, nil-state preservation and source cloning, the PNG mask, the shared
curve builder, live-list hit tests, `markGraphDirty` on every list replacement and the
full-state status key, `zone_content` naming no DTO, `tools/go.mod` as the single lint-version
source, and all of Batch K (7 topologies, blocking rejection of retired/unknown IDs with the
document unchanged, `""` → Random with a warning, no-fallback `Resolve`, balanced tournaments).

## Plan review log

**Round 1 (GPT-6.1 Sol, 2026-10-13): REJECT, 14 findings, all applied.**
1. Added-zone provenance → D3a: `origin: added` read from a missing `GeneratorPosition` (every
   topology stamps one, the factory does not); legacy-unstamped snapshots by name. No dialog
   tracking needed; Undo/Revert carry the marker in the zone value.
2. Castle comparison → D1/D3a: `castles` present is the re-profile marker; absent means follow
   the fresh zone; only neutral zones are re-profiled.
3. Pending conversion across save/reload → D1/D2: `generatorVersion: 0` with edits is the
   persisted discriminator (`NeedsDerivation()`); no source-version metadata.
4. Conversion against the file's settings → D4: driver derivation baseline;
   `ConvertManualEdits` generates from it.
5. Capture failure → D4: all-overridden fallback, `generatorVersion` 0, baseline, unsaved.
6. Phase boundaries → §2.5 merged into the v3 phase; typed overrides in Phase 3; modal in
   Phase 4; facade/DTO/mock changes enumerated; Wire at each constructor change.
7. Interim text conversion → removed (merged phase); warnings delivered in the same phase.
8. Two-pass contract and tuning zone count → D3.
9. Nil/empty at the wire boundary → D1 tags and rules; JSON round-trip tests.
10. Expansion warnings → appended to the final generation status (D4).
11. Guard hash blind to editor geometry → model hash added; RandomPortals and tournament matrix.
12. Random export equality → scoped comparator in Phases 3 and 4.
13. Size/performance comparability → Phase 2 protocol and baseline; compact-vs-compact size gate.
14. Coverage per phase → every phase with Go changes compares coverage (Phase 1 has none).

**Round 2 (GPT-6.1 Sol): REJECT; 12 of 14 resolved, 5 findings, all applied.**
1. Pending baseline lost across save/reload → D4: documented conservative rule (values that
   differ stay overridden; lossless); combined acceptance test. Surfaced to the owner.
2. First-generation castle deltas → D4: `DiffCastleSettings(baseline, current)` when pending.
3. Conversion failure fell through into expansion → D4: explicit failure branch.
4. Connection provenance → D3b: `isUserAdded` derives from the factory first; the dialog's Undo
   baseline keeps the flag (two lines removed). Surfaced to the owner.
5. Phase 1 coverage claim → documented exception (no Go changes).

**Round 3 (GPT-6.1 Sol): APPROVE**, all round-2 findings resolved, no new findings; subject to the
owner confirming the conservative pending-baseline rule and the dialog Undo flag change.

## Design

### D1. v3 wire shape

```json
{
  "schemaVersion": 3,
  "generatorVersion": 1,
  "templateName": "Colosseum v4",
  "gameMode": "Classic",
  "mapSize": 368,
  "experimentalMapSizes": false,
  "players": { "playerCount": 6, "heroMin": 4, "...": "..." },
  "neutralZones": { "neutralZoneCount": 12, "...": "..." },
  "castles": { "advancedMode": false, "...": "..." },
  "generation": { "topology": "GeometricHub", "isolateplayers": false, "...": "..." },
  "gameRules": { "victoryCondition": "...", "...": "..." },
  "content": {
    "contentDefaultsVersion": 1,
    "bannedItems": [ { "sid": "item_a" } ],
    "bannedMagics": [],
    "valueOverrides": [ { "sid": "mine_gold", "guardValue": 12000 } ],
    "bonuses": [],
    "lowestNeutralContentRows": []
  },
  "manualEdits": {
    "zones": [
      { "name": "Spawn-A", "generatorPosition": { "X": 0.12, "Y": 0.5 }, "generatorRing": 0,
        "manualPosition": { "X": 0.12, "Y": 0.5 } },
      { "name": "Neutral-G", "quality": 2, "generatorPosition": { "X": 0.41, "Y": 0.63 },
        "generatorRing": 1, "manualPosition": { "X": 0.46, "Y": 0.63 },
        "fieldOverrides": { "size": 1.5 } },
      { "name": "Neutral-Z", "origin": "added", "quality": 2, "castles": 1,
        "manualPosition": { "X": 0.8, "Y": 0.2 }, "unset": ["mandatoryContent"] }
    ],
    "connections": [
      { "name": "Manual-G-Z", "from": "Neutral-G", "to": "Neutral-Z", "type": "Direct",
        "isUserAdded": true, "fieldOverrides": { "guardValue": 80000 }, "unset": ["road"] }
    ]
  }
}
```

- Root order: `schemaVersion`, `generatorVersion`, then `TemplateIdentity` and `MapSettings`
  fields flat, then the seven sections.
- `manualEdits` is omitted (`omitzero`) when there are no edits. `generatorVersion` is always
  written.
- `playerZoneContentRows` is omitted when equal to the defaults of `contentDefaultsVersion`;
  any cleared row list is written as `[]`; a never-set neutral list (nil) is omitted as today.
- `fieldOverrides` keys are the `.rmg.json` (`template_entity`) JSON names; values use the same
  entity types, so they read like the RMG file.
- `unset` lists field keys that are nil in the edited item while the derivation has a value.
  A nil value is **never** an override value: nil → `unset`, `[]` / `false` / `0` → override.
  For non-deterministic layouts a nil `roads` is always written as `unset: ["roads"]`,
  independent of the derivation (M3).
- Tags: v3 content-row slices use `omitzero` (nil omitted, `[]` written, M13); override struct
  fields are pointers with `omitzero` (a non-nil pointer to `false`, `0` or `[]` is written).
- `origin` is written only as `"added"`; `castles` only for added or re-profiled neutral zones,
  and its presence on a non-added zone **is** the re-profile marker; `quality` whenever the zone
  has a recorded tier.
- `generatorVersion` is the generator the derived values were established against. `0` means
  "never established": the manual edits are a lossless all-overridden snapshot waiting for
  conversion (legacy migration, a save before the first successful regeneration, or a failed
  capture). It is the reload-safe pending-conversion discriminator; no extra key is needed.

### D2. Entities, models, config

- **Frozen `internal/entities/editor_state/editor_state_v2/`**: a byte-for-byte copy of today's
  flat v2 entity files (same JSON tags), read-only like `editor_state_v1`, with a doc comment.
- **v3 entity (`internal/entities/editor_state/`)**: `EditorState` embeds `SchemaOptions`
  (now `SchemaVersion` + `GeneratorVersion`), `TemplateIdentity`, `MapSettings` flat, and holds
  seven named section fields with the M8 keys. `CurrentEditorStateSchemaVersion = 3`.
  New one-type-per-file entities: `BanEntry{Sid}`, `ValueOverrideEntry{Sid, GuardValue}`,
  `ManualZoneEdit`, `ManualConnectionEdit`, `ZoneFieldOverrides`, `ConnectionFieldOverrides`,
  `ZoneOrigin`. `ContentSettings` gains `ContentDefaultsVersion` and typed lists.
  `ManualZoneSave` / `ManualConnectionSave` leave the v3 package (they live on in v2/v1).
- **Models (`internal/models/editor_state_model/`)**: the group models keep embedding their
  entity groups, so promoted field access (`state.PlayerCount`) is unchanged. New wrappers
  `BanEntry`, `ValueOverrideEntry` (like `BonusEntry`), `ManualZoneEdit`,
  `ManualConnectionEdit`, `ZoneFieldOverrides`, `ConnectionFieldOverrides` (model types over
  `template_model` field types), and `ManualEditSettings{Zones, Connections}`. Business logic on
  the models: `Equals`, `Clone`, `HasManualEdits`, `EditorState.NeedsDerivation()`
  (`HasManualEdits() && GeneratorVersion == 0`), the all-overridden encoding
  (`ManualEditSettingsFromSnapshot`) and its derivation-free inverse, override apply/diff, and
  the default-rows table lookup (`GetDefaultPlayerZoneContentRowsForVersion`).
- **Config (`internal/models/config/`)**: `BanEntry`, `ValueOverrideEntry` aliases;
  `GeneratorConfig.BannedItems/BannedMagics []BanEntry`, `ValueOverrides []ValueOverrideEntry`
  (replaces `ValueOverridesText`).
- **Generator version**: `GeneratorOutputVersion = 1` in a new
  `internal/common/common_generator/generatorOutputVersion.go`. The model keeps the loaded value
  (`ToModel` no longer overwrites it). On save the mapper writes `0` when `NeedsDerivation()`,
  otherwise `GeneratorOutputVersion`. The driver sets the in-memory value to current after a
  successful capture or conversion, and after showing the M-D5 notice.

### D3. Capture and expand (C3 semantics)

New `ManualEditDeltaService` in `internal/services/connection_editor/manualEditDeltaService.go`
(interface `IManualEditDeltaService` in `connection_editor_interfaces`). Models only; it never
calls the generator (the handler generates and passes the fresh template, which also avoids a
`connection_editor` ↔ `template_generator` cycle).

- `Capture(request ManualEditCaptureRequest) ManualEditSettings`: inputs are the fresh variant
  (derivation source), the edited zones/connections (variant 0 of the template **after**
  `UpdateTemplate`), the `isUserAdded` flags, the tuning/road options, whether the layout is
  non-deterministic, and whether the snapshot is legacy-unstamped (no zone carries a
  `GeneratorPosition`, i.e. v0/v1). Pass 1, zones in order: classify (D3a), build the
  derivation, store the structure, store each non-structure field whose value differs
  (`reflect.DeepEqual` on the typed field) as an override, list nil fields whose derivation is
  non-nil in `unset`; for non-deterministic layouts always store `roads` (or `unset`). Pass 2,
  connections in order (D3b), derived against the **edited** zones. Order is the edited order.
- `Expand(request ManualEditExpandRequest) ManualEditExpansion`: pass 1 rebuilds **every** zone
  (derivation + structure + overrides + `unset`); pass 2 rebuilds every connection against the
  rebuilt zones and their player names; pass 3 runs castle reapply (M6) when castle options
  changed, on neutral zones carrying `castles` and on zones with an overridden `mainObjects`,
  then re-diffs those zones so their stored `castles`/overrides follow. Tuning everywhere uses
  `len(all rebuilt zones)`; the M6 filter never shrinks that count. Returns zones, connections,
  the possibly-updated `ManualEditSettings`, a changed flag and warnings.
- **D3a zone classification** (only neutral zones are ever re-profiled or added):
  - **added**: no `GeneratorPosition` (every topology stamps one; the editor factory does not),
    or, for a legacy-unstamped snapshot, a name absent from the fresh variant. Derived from
    `ZoneEditorService.NewDefaultNeutralZone(label, tier, castles, tuning, roads)`. This also
    covers delete-then-re-add of a reused label, and Undo/Revert need no dialog tracking because
    the marker travels with the zone value.
  - **re-profiled**: generated, same-named fresh zone exists, and its tier or
    `CountZoneCastles` differs. Stores `castles`; derived from the factory like added zones.
  - **generated**: everything else; derived from the same-named fresh zone. At expand, absent
    `castles` means "follow the fresh zone", so castle-setting changes reach it unchanged.
  - At expand, a generated name missing from the fresh variant is rebuilt by the factory
    (stored tier, or Medium; 1 castle) with the warning `zone "X" is no longer generated; kept
    as an editor-made zone`, and becomes `origin: added` in the returned settings.
- **D3b connection derivation**: a connection with `isUserAdded` always derives from the factory
  (M4), before any name matching. Otherwise the same-named fresh connection with exactly the
  stored `from` and `to`; otherwise the factory. Factory = `NewDefaultConnection` over the
  (edited or rebuilt) zone list and player names (+ `ChangeConnectionType` when the type
  differs), then the stored name.
- **`isUserAdded` provenance in the dialog**: `ZoneEditorDialog.setEditingSet` and
  `undoSessionEdits` currently clear `IsUserAdded` on the Undo baseline (`zoneEditorDialog.go`,
  `clone.IsUserAdded = false` twice), so Undo then Apply silently drops the persisted flag of
  connections added in earlier sessions. Both lines are removed so the flag survives opening and
  Undo (owner-approved with the plan; see the approval note).
- **Determinism flag**: `topology == Random || RandomPortals`, computed in the handler from the
  configuration.
- Roads stay subject to `UpdateTemplate`'s `Reconcile` after expansion, exactly as today.

### D4. Handler and driver flows

`ITemplateHandler` changes, mirrored in the `IGuiHandler` facade forwarding methods
(`guiHandler.go`), the handler interfaces, `test/test_helpers/templateHandlerMock.go` and the
handler/facade unit tests; new DTOs in `internal/dtos/` one per file:
- **Add** `CaptureManualEdits(dtos.ManualEditCaptureRequestDto) (dtos.ManualEditCaptureDto, error)`:
  validates, generates fresh from the given settings (manual edits are not generator input),
  calls `Capture`. Used by Apply.
- **Add** `ExpandManualEdits(dtos.ManualEditExpandRequestDto) dtos.ManualEditExpansionDto`:
  takes the already-generated fresh template, the state and the castle changes; calls `Expand`.
- **Add** `ConvertManualEdits(dtos.ManualEditConvertRequestDto) (dtos.ManualEditCaptureDto, error)`:
  inverts the all-overridden form (no derivation needed), generates fresh from the given
  **baseline settings**, and re-captures (M16).
- **Remove** `ReapplyCastleSettings` and `dtos.CastleSettingsReapplyRequestDto` (folded into
  expand), with their facade method, mock and tests.

Driver (`app/gui/drivers/`):
- **Derivation baseline.** `State.derivationBaseline *editor_state_model.EditorState` (settings
  clone without the manual section) and `State.announceConversion bool`. Set in
  `handleLoadState` when the loaded state `NeedsDerivation()` (baseline = the loaded settings,
  announce = true), and on a capture failure (baseline = the settings at Apply, announce
  unchanged). Cleared by `Reset`, a new load, `ClearManualEdits`, a mode transition that drops
  edits, and a successful conversion. Never persisted: after a save and reload while conversion
  is still pending (`generatorVersion: 0` survives the save) the reloaded settings become the
  baseline. If the user changed settings before that save, conversion then compares the old
  snapshot against the changed settings: values that differ are kept as overrides. That is
  lossless and conservative (those fields stay pinned, as today) rather than wrongly derived;
  persisting a separate baseline was rejected as format weight for a rare double-failure path
  (owner-approved with the plan). The same conservative rule applies when the baseline itself
  cannot be generated: conversion then uses the current settings.
- `ApplyEditedZones`: unchanged revert-to-base detection; `handleUpdateTemplate`; then
  `CaptureManualEdits` on the updated template's variant 0. Success: `SetManualEdits(compact)`,
  in-memory `generatorVersion` = current, `flagAsUnsaved` on change. **Failure:** store the
  lossless all-overridden form with `generatorVersion` 0 (so `NeedsDerivation()` holds and a
  save writes it), set the baseline, `flagAsUnsaved`, and report `Applied … ‼ Error: …`.
- `handleGenerateTemplate` → `reapplyManualEdits`:
  - **Pending (`NeedsDerivation()`)**: castle changes for this generation are
    `DiffCastleSettings(baseline, current)` (the decision service returns none after a load
    because `previous` is nil). Call `ConvertManualEdits` with the baseline (or the current
    settings if the baseline cannot generate). **Success:** store the compact form, set the
    in-memory `generatorVersion` to current, clear the baseline, `flagAsUnsaved` (M17), open the
    M-D2 modal once if `announceConversion`, then continue with expansion below. **Failure:**
    a separate branch: reapply the derivation-free snapshot with `handleUpdateTemplate`, skip
    expansion and any compact-state replacement, keep `generatorVersion` 0 and the baseline, and
    report the error in the status.
  - **Expansion**: `ExpandManualEdits` on the fresh template with the castle changes;
    `handleUpdateTemplate(zones, connections)`; store the returned compact form and
    `flagAsUnsaved` when it changed. Expansion and conversion warnings are collected and
  appended to the final generation status (after "(Manual zone edits reapplied.)"), so the
  closing `SetStatus` no longer hides them. If the loaded `generatorVersion` is non-zero and
  differs from current and edits exist, the M-D5 suffix is appended once and the in-memory
  value set to current.
- `app/gui/models.EditorState`: `SetManualEdits(ManualEditSettings) bool`,
  `GetManualEdits() ManualEditSettings`, `SetGeneratorVersion(int)`;
  `GetManualZones/GetManualConnections` removed; `ClearManualEdits` and `ApplyModeTransition`
  clear the compact form (and the driver clears the baseline).
- `ZoneEditorDialog`, `markGraphDirty`, hit tests, revert-to-base: unchanged (the editor still
  opens on `lastTemplate`; Revert-to-Base still clears the manual edits on Apply).

### D5. Migration chain

- `EditorStateMigrator.Load(path, target) ([]string, error)` dispatches on the probed
  `schemaVersion`: `> 3` refused (unchanged error); `3` → v3 repository, then
  `contentDefaultsVersion > 1` refused with a new `UnsupportedContentDefaultsVersionError`
  (document unchanged, same route as an unsupported schema); `2` → frozen v2 repository →
  `MigrateToV3`; `< 2` → v1 repository → `MigrateToV2` (now returns `editor_state_v2`) →
  `MigrateToV3`.
- `MigrateToV3(v2) (v3, []string)`: settings into sections; ban text → `BanEntry` per trimmed
  non-empty line, order and duplicates kept; override text → `ValueOverrideEntry` per valid line
  using today's rules (first `=`, trimmed SID non-empty, `strconv.Atoi` on the trimmed value),
  malformed lines dropped with today's messages (`line N: '…' is not sid=value`, `… has a
  non-numeric value`, N counting raw lines); `contentDefaultsVersion = 1`; manual snapshot →
  all-overridden C3 form (structure fields lifted, every other non-nil field an override, nil
  fields in `unset`); `generatorVersion = 0`.
- `FileService.LoadSettingsFile(path) (*model, []string, error)` passes the migration warnings
  through (no source-version metadata is needed: `generatorVersion: 0` marks pending
  conversion).
  `stateHandler.LoadState` returns the warnings in a new `EditorStateValidationDto.LoadWarnings`;
  `handleLoadState` prints them before the validation warnings in the existing
  `Loaded … (adjusted: …)` message.
- Seeding: the v3 load still decodes over the default entity (absent keys keep defaults).
  Player rows absent → the frozen table for `contentDefaultsVersion` (filled by the mapper).
- New read-only `LegacyV2EditorStateRepository` (`IFileRepository[editor_state_v2.EditorState]`)
  beside the v1 one; `newV1Seed`/`newV2Seed` derive seeds from the v3 default entity.

### D6. GUI notices

- New `app/gui/dialogs/messageDialog.go` (`IDialog`): title, body, single OK button, closes on
  OK/Escape. Opened by the driver through `DialogHost` on conversion. Text:
  "Older file format converted" / "This file was saved in an older format. Values you did not
  edit in the zone editor now follow the current settings and generator. Save to keep the new
  format." Exact wording confirmed at the final review.
- M-D5 status suffix: "Saved with generator version N; values you did not edit now follow
  version M and may differ."

## Phase 0: Baseline
Status: Not started

- [ ] Confirm branch `AD/compact_state`, HEAD `85752a8`, `git status --short` (only the handoff
  and this plan expected).
- [ ] `golangci-lint cache clean`.
- [ ] `go build ./...`; `go vet ./...` with no tag, `-tags=integration_test`,
  `-tags='integration_test,gui'`.
- [ ] Coverage: `go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=.agent/memories/batch-m/baseline/coverage.txt' ./test/unit/...`
  then `go tool cover '-func=…'` (recapture; the last figure 73.0% predates `85752a8`).
- [ ] `go test ./test/... -count=1`; `go test -tags=integration_test ./test/integration/... -count=1`;
  `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` (record `.failure`
  count, must be 0).
- [ ] GPU-free benchmarks: `go test '-bench=.' '-run=xxx' ./test/performance/... '-benchtime=20x' '-timeout=120s'`.
- [ ] `go run ./cmd/testlayoutcheck .`; `gofmt -l .` (excluding `data/`); `wire diff ./internal/composition/...`.
- [ ] `golangci-lint run ./... --issues-exit-code=0`, check `$LASTEXITCODE`.
- [ ] Save all outputs under `.agent/memories/batch-m/baseline/`.

### Verification Plan
- Every command exits 0; lint 0 issues; 0 `.failure` files; `wire diff` empty. Record the
  coverage total as M23's floor. Any red result stops the batch and is reported.

### Phase Summary
_(write when phase completes)_

## Phase 1: Linter bump (M-D8)
Status: Not started

- [ ] In `tools/` only: `go get github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0`
  then `go mod tidy`; accept MVS-forced indirect changes (H3). `wire`, `gcov2lcov`, `go` line and
  the `tool` block unchanged.
- [ ] `git ls-files --eol tools/go.mod tools/go.sum` stays `i/lf w/lf`.
- [ ] `go list -m` in `tools/` reports v2.14.0 (the CI step reads this).

### Verification Plan
- `go mod tidy -diff` exits 0 in root and `tools/`; lint 0 issues; `go build ./...` PASS.
- No Go source changes in this phase, so no coverage run (the only phase without one).

### Phase Summary
_(write when phase completes)_

## Phase 2: Per-frame cloning (M-D7) and persistence baselines
Status: Not started

- [ ] `editor_state_model.EditorState.CloneWithoutManualEdits()` (settings deep copy; manual
  section zeroed).
- [ ] `app/gui/models.EditorState`: `GetCurrentSettings()`, `GetPreviousSettings()`,
  `GetNextSettings()` built on it; `AutoRegenerate` feeds only these to `DecideRegeneration`.
  `NextStateSetFromCurrent` keeps the full current state.
- [ ] Unit tests: the new model method (manual section empty, settings equal, mutation
  isolation for every cloned slice); the three GUI-model getters (nil previous/next, isolation).
- [ ] Driver test in `test/unit/app/gui/drivers/stateGeneration/autoRegenerate_test.go` (with a
  recording regeneration handler): with manual edits stored, every DTO handed to
  `DecideRegeneration` has no manual section; `HasManualEdits` still true afterwards.
- [ ] Commit the geohub6p move fixture now (M22) as
  `test/test_helpers/testdata/editorState_v2_geohub6p_move.gen.json`.
- [ ] Ungated GPU-free benchmarks in `test/performance/`:
  - `BenchmarkAutoRegenerate_IdleFrame` on that fixture and on a no-edit state;
  - `BenchmarkEditorStateFile_SaveLoad`: setup loads the fixture into a model (and, from
    Phase 4 on, lets the first regeneration convert it); the measured op is
    `FileService.SaveSettings` + `LoadSettingsFile` in `b.TempDir()`.
- [ ] **Measurement protocol** (used again in Phase 4 and Phase 6): `go test '-bench=…'
  '-run=xxx' '-benchmem' '-benchtime=1s' '-count=5' ./test/performance/...`, median of the five
  ns/op, B/op, allocs/op, saved to `.agent/memories/batch-m/bench/`. Add the fixture and both
  benchmarks **first**, record them on the unchanged code (before the clone fix and before any
  format change: the persistence baseline), then implement the clone fix and record again.

### Verification Plan
- Build, unit suite PASS; the deterministic driver test is the gate; benchmarks recorded.
- Coverage run and compare to the Phase 0 floor (every touched function at full reachable
  coverage).

### Phase Summary
_(write when phase completes)_

## Phase 3: Schema v3 with typed lists, sections, migration and versions (§2.5, §2.6, M-D5 storage, M-D6)
Status: Not started

§2.5 and §2.6 ship together so that no interim text conversion exists: legacy text is parsed
only by `MigrateToV3`, and its warnings reach the GUI in the same phase. Manual edits are
persisted in the C3 **shape** but always in the all-overridden form with `generatorVersion` 0
(the live model is still expanded until Phase 4), so the format changes without changing
behaviour. Wire is regenerated (`wire gen`, then `wire diff`) at every constructor or provider
change, before building.

**§2.5 typed lists**
- [ ] Entity `BanEntry`, `ValueOverrideEntry`; model wrappers with wrap/unwrap/clone; config
  aliases. Model `ContentSettings` typed fields; `EqualsIgnoringManualEdits` compares them
  (`slices.Equal`, nil-vs-empty as for bonuses) and the fuzz-parity test covers them; `Clone`
  deep-copies them.
- [ ] `GeneratorConfig` typed fields; `GeneratorConfigMapper` forwards them.
- [ ] `GameRulesProvider.CreateValueOverrides(configuration) []template_model.ValueOverride`
  (`Variant: -1`, order and duplicates kept, no warnings) and `CreateGlobalBans` from typed lists
  (nil when both empty); its interface and every caller; the generator's warning list loses only
  the override parse warnings (which now come from the load).
- [ ] Bonuses panel and `ValueOverridePickerDialog` / item and spell pickers work on typed
  entries; `splitNonEmptyLines`, `overrideLabel` line parsing and `overrideSids` become typed
  accessors. No text parsing remains outside `MigrateToV3`.

**§2.6 sections, v3 entity, versions**
- [ ] Freeze `editor_state_v2` (copy of today's flat entity files + `ManualZoneSave`,
  `ManualConnectionSave`); add `LegacyV2EditorStateRepository`; Wire provider; migrator
  constructor.
- [ ] v3 entity per D1/D2 (tags per D1); `CurrentEditorStateSchemaVersion = 3`;
  `SchemaOptions.GeneratorVersion`; `common_generator.GeneratorOutputVersion = 1`; mapper rules
  for both version fields (D2).
- [ ] `contentDefaultsVersion`: `GetDefaultPlayerZoneContentRowsForVersion(1)` = today's table
  (frozen); mapper omits equal player rows, writes cleared lists as `[]`, fills absent player
  rows from the table on load.
- [ ] Manual section: typed override structs (entity + model) with apply/diff helpers, the
  reflection parity test (every `template_model.Zone`/`Connection` field is structure or
  overridable; the entity override structs mirror `template_entity` field-for-field),
  `ManualEditSettingsFromSnapshot` and its derivation-free inverse; the mapper uses them while
  the live model is still expanded.
- [ ] Migrator per D5 (`MigrateToV2` → `editor_state_v2`, new `MigrateToV3`, new
  `UnsupportedContentDefaultsVersionError`, warnings returned); `FileService.LoadSettingsFile`
  returns warnings; `EditorStateValidationDto.LoadWarnings`; `stateHandler.LoadState`;
  `handleLoadState` message; facade and mocks.
- [ ] Fixtures (M22) under `test/test_helpers/testdata/`: `editorState_v3_flat.gen.json`
  (all-fields), a v2 file with malformed override text, v2 with cleared player rows, a v3 file
  with `contentDefaultsVersion: 2` and one with `schemaVersion: 4`; the curated Batch L subset
  (square2p: move, add-zone, delete-zone, add-connection, delete-connection, edit-connection,
  quality, field-override, nil-empty; random8p: move, add-zone, nil-empty; geohub6p:
  delete-zone, plus the move fixture from Phase 2), `.gen.json` only, about 1.2 MB.

**Tests**
- [ ] Entry JSON; model wrap/unwrap/clone/isolation; generator config mapping;
  `createValueOverrides_test.go` rewritten (multiple entries, duplicates, order, empty → nil,
  `Variant=-1`); `CreateGlobalBans` (items only, magics only, both empty → nil); exported
  `GlobalBans` / `ValueOverrides` identical to today's output for an equivalent legacy text.
- [ ] v3 entity JSON: root vs sections asserted key by key (exactly the six root keys, exactly the
  seven section keys); nil vs `[]` vs `false` vs `0` in overrides and rows through real JSON
  save/load; the frozen v2 decode.
- [ ] `MigrateToV3`: every group, typed lists (blanks, duplicates, order, each malformed kind with
  raw line numbers), manual snapshot → all-overridden with nil/empty, `generatorVersion` 0.
- [ ] Migrator dispatch for v0/v1/v2/v3/4 and defaults version 2 (document unchanged); mapper
  (defaults omission, cleared `[]`, version rules, absent rows → table); `FileService` /
  `stateHandler` / driver warning propagation.
- [ ] Rewrite `editorStateWireFormat_integration_test.go` for v3 with v0/v1/v2 as inputs; extend
  `editorStateRoundTrip_integration_test.go`: every curated fixture loads → saves → loads to an
  identical model and a byte-stable second save; export and preview identical with the
  **scoped comparator** (deterministic layouts: full export bytes and preview layout; Random:
  raw zones and connections plus preview layout, template-level `orientation` excluded, as
  report §5 records).
- [ ] Update `allFieldsEditorState.go` for v3.

### Verification Plan
- Build, vet × 3, unit, integration (both), GUI tagged suite PASS with 0 `.failure`;
  testlayoutcheck; lint 0; `wire diff` shows only the v2 repository provider and migrator
  constructor.
- Coverage run and compare to the Phase 0 floor; list each changed function's coverage.

### Phase Summary
_(write when phase completes)_

## Phase 4: C3 semantics and §1.16 (M-D1 to M-D5)
Status: Not started

- [ ] Behaviour-neutral move first: `connection_editor` interfaces → new
  `internal/services/connection_editor/connection_editor_interfaces/` (4 files), all imports,
  mocks and Wire updated; `wire gen` + `wire diff`, build and unit green before continuing.
- [ ] `ManualEditDeltaService` (D3) + interface + Wire; `templateHandler` gains the dependency
  and the three methods (D4), mirrored in the `IGuiHandler` facade, interfaces, mocks and tests;
  `ReapplyCastleSettings`, its DTO, facade method, mock and tests removed;
  `ManualReapplyService` used from `Expand` per M6.
- [ ] Live model switches to compact (D2/D4); mapper persists the compact form directly.
- [ ] `messageDialog.go` (D6) moved here so the conversion path is complete in this phase.
- [ ] Driver flows (D4): Apply capture with the all-overridden failure fallback; baseline and
  `announceConversion`; conversion + modal + M17; expansion; warnings appended to the final
  status; M-D5 (`DecideManualEditReapplication` returns `GeneratorVersionChanged`, compared in
  the regeneration decision service, which may import `internal/common`).
- [ ] **Generator version guard + determinism contract** (untagged integration test
  `test/integration/generatorOutput_integration_test.go`, production APIs only). Matrix: the
  7 topologies with `RandomPortals` off; Square and Geometric with `RandomPortals` on; a
  tournament 2p. Each generated twice.
  - Deterministic cases: two runs identical, and two hashes equal to
    `testdata/generatorOutputHashes.json` keyed by `GeneratorOutputVersion`: the RMG entity JSON
    **and** a canonical model JSON of variant 0 including `GeneratorPosition`, `GeneratorRing`
    and `Quality` (which the RMG mapper drops).
  - Random topology and `RandomPortals`: identical zone set and identical non-road zone content
    (model, including tier) across runs.
  - A hash change at the same version fails ("bump `GeneratorOutputVersion`"); a version with no
    recorded hashes fails ("regenerate with `-update`"); `-update` rewrites only that file.
- [ ] Unit tests (mirrored layout, one function per file): `Capture` (untouched → no overrides;
  move; field override; nil → `unset`; `[]`/`false`/`0` → override; re-tier → factory +
  `castles`; castles-only → factory; added zone; delete-then-re-add of a reused label → added;
  legacy-unstamped snapshot → added by name; deleted zone absent; non-deterministic → roads
  always stored including nil; connection matched / unmatched / endpoint mismatch / type change /
  user-added, derived against edited zones; order kept); `Expand` (inverse of each; two-pass
  order; tuning uses all zones; missing name → factory + warning + `added`; castle changes on
  player, hub and neutral zones: derived zones follow the fresh generation, re-profiled and
  `mainObjects`-overridden zones get the reapply and a re-diff, several selectively eligible
  zones; changed flag); `ConvertManualEdits` (baseline settings used, not current);
  handler methods (validation rejection, generation error, determinism flag); model
  `Equals`/`Clone`/`HasManualEdits`/`NeedsDerivation`/override helpers; regeneration decision
  `GeneratorVersionChanged`; driver flows (`test/unit/app/gui/drivers/`): capture failure →
  all-overridden, unsaved, Save writes `generatorVersion: 0`, Exit prompts; warnings survive in
  the final status; baseline cleared by Reset / load / mode transition.
- [ ] **Conversion acceptance** (integration): legacy load → save without any successful
  generation → reload → conversion still happens; legacy load → change density and castles
  before the first successful regeneration → conversion uses the load-time baseline, castle
  deltas baseline → current reach re-profiled and `mainObjects`-overridden zones, expansion
  uses the new settings; legacy load → change density → save while pending → reload → convert:
  no value lost, the density-dependent fields stay overridden (the documented conservative
  rule); conversion failure with non-empty castle changes → snapshot reapplied unchanged,
  `generatorVersion` 0 kept, error in the status.
- [ ] Dialog provenance: existing user-added edge → reopen → Undo → Apply keeps `isUserAdded`;
  a user-added connection whose name and endpoints equal a fresh generated connection still
  derives from the factory.
- [ ] **Round-trip acceptance** (integration, curated fixtures, scoped comparator from Phase 3):
  load v2 → first regeneration → raw zones/connections equal the v2 snapshot, export and preview
  equal; then save → load → regenerate: equal again.
- [ ] **§1.16 regressions** (integration through the driver, `drivers.NewUIState`): after one move,
  changing each setting in report §6 (`MapSize`, resource density, structure density, player /
  neutral / hub zone size, neutral-stack strength, border-guard strength incl. connection
  `guardValue`, abandoned outposts on and count, City Hold on a hub topology, hub castles' biomes,
  advanced-mode guard randomization) makes the affected fields equal a fresh generation; an
  explicit override of the same field survives; castle control (player castles change still
  applied); no-op control (template rename changes no zone).
- [ ] Identity tests: delete then re-add the same label (`origin: added`), parallel edges,
  user-added naming, order preserved, Undo and Revert-to-Base then Apply.
- [ ] **Size gate** (integration test, deterministic): for the geohub6p move fixture,
  `json.Compact` of the v2 `manualZones` + `manualConnections` values vs `json.Compact` of the
  converted v3 `manualEdits` value; v3 ≤ 30% (report §13). The indented file sizes are recorded
  separately in the phase summary.
- [ ] **Performance**: rerun the Phase 2 protocol. `BenchmarkEditorStateFile_SaveLoad` median
  ns/op not above the Phase 2 baseline; `BenchmarkAutoRegenerate_IdleFrame` allocs/op and B/op
  not above the Phase 2 "after" figures.
- [ ] Update `manualCastleReapply_integration_test.go`, `manualEditDirtyTracking_integration_test.go`,
  `manualZoneTierPersistence_integration_test.go`, `zoneEditorRevertToBase_integration_test.go`
  to the new semantics (each change justified in the phase summary).

### Verification Plan
- Build, vet × 3, unit, integration (both), GUI tagged PASS with 0 unexplained `.failure`;
  benchmarks per protocol; testlayoutcheck; lint 0; `wire gen` then `wire diff` clean;
  architecture layering test passes with no allow-list growth.
- Coverage run and compare to the Phase 0 floor; every branch of `Capture`, `Expand`, the three
  handler methods and the driver flows covered or recorded as GUI-only.

### Phase Summary
_(write when phase completes)_

## Phase 5: GUI notices and regression flows
Status: Not started

- [ ] GUI-only statements of the modal and the driver's dialog opening recorded in
  `test_observations.md`.
- [ ] GUI integration (`integration_test && gui`): loading a v2 file with manual edits opens the
  modal once, OK closes it, the document is unsaved; a v2 file without edits opens no modal and
  stays saved; the M-D5 status suffix for a v3 file with an older non-zero `generatorVersion` and
  edits, absent without edits; malformed-legacy load message lists the line warnings.
- [ ] Manual-edit flows through the window harness: move, add zone, delete zone, quality,
  connection edit, apply, revert to base, save, load, settings change (border-guard strength) —
  each asserts the exported zone values and the persisted `manualEdits` shape.
- [ ] Bonuses panel: add/remove ban item, ban spell, override; save; load; values and order
  preserved; goldens unchanged.
- [ ] New modal goldens accepted only from `.failure` files, listed for the owner.

### Verification Plan
- GUI tagged suite PASS (twice, `-count=1` then cached), 0 unexplained `.failure`; all other
  gates as Phase 4, including a coverage comparison.

### Phase Summary
_(write when phase completes)_

## Phase 6: Final gate, review and close-out
Status: Not started

- [ ] Full gate: build, vet × 3, coverage (≥ Phase 0 floor; every new/changed function at full
  reachable coverage, gaps recorded), `./test/...`, both tagged suites, both benchmark suites
  (protocol from Phase 2), testlayoutcheck, `gofmt -l` on new files, lint 0 after cache clean,
  `wire diff` clean, `go mod tidy -diff` in both modules.
- [ ] Independent implementation review (GPT-6.1 Sol); apply or record each finding.
- [ ] Review document: §2.4, §2.5, §2.6 and §1.16 marked fixed in place (pending the owner's
  commit) with Progress paragraphs; §9 rows M and Q (Q folded into M); progress line and totals.
  §8 untouched. Report §15: "answered by M-D1 to M-D9".
- [ ] `test_observations.md`; README/QUICKSTART where they describe the `.gen.json` shape or the
  manual-edit snapshot; `.agent/memories/` notes (format, traps).
- [ ] Handoff refreshed per AGENTS.md §5.2 for the owner's review.

### Verification Plan
- All gates green; review verdict recorded; the owner's end review is the approval.

### Phase Summary
_(write when phase completes)_

## Risks and open points for the reviewer
- Random layouts: connection derivation will mostly fall back to the factory (names and
  endpoints differ per generation), so most Random connections carry overrides. Accepted by M3/M4.
- One generation per Apply (M3) and one extra per conversion: acceptable cost (a few ms),
  measured in Phase 4.
- The parity test is the guard against a new `Zone` field silently becoming non-persisted.
- The pending-conversion state is `generatorVersion == 0` with edits (persisted) plus the
  driver's derivation baseline (transient, rebuilt from the loaded settings after reload).
- Added-zone provenance relies on every topology stamping `GeneratorPosition`; the guard test's
  model hash and a capture test pin that invariant.
- Phase 3 is large and, to avoid any interim text conversion, the typed lists and the v3 entity
  switch together: the build may be red between its sub-groups and must be green at the phase
  end.

## Final Recap
_(write when all phases complete)_

## Deployment Plan
_(write when all phases complete)_
