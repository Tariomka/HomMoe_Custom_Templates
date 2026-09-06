# Batch T — `template` becomes `template_entity`, and the alias shim dies

Backlog §2.5. Rename the protected `.rmg.json` entity tree from the extremely generic
`template` to `template_entity` (six subpackages renamed to match the existing
`template_model` suffix convention), and delete the outer alias-of-alias shim
`internal/entities/types.go`, moving its 28 consumers onto `template_entity.X`.
Pure naming: no field, tag, type or behaviour changes anywhere, and no wire-format
movement is possible by construction.

## Owner approvals on record (this session)

1. **Protected-tree edit approved.** AGENTS.md §2.1 forbids touching
   `internal/entities/template/` without explicit approval; the owner scheduled this
   item as batch T, which is that approval. The edit is confined to `package` clauses,
   sibling import paths and the directory names themselves. **No JSON tag, field name,
   field order or type may move.**
2. **The frozen `editor_state_v1` snapshot gets an import-only edit.** Batch S declared it
   "data only, never edited"; the owner approved changing `entities.Zone`/`entities.Connection`
   to the new spelling in `manualZoneSave.go` and `manualConnectionSave.go` (plus the doc
   comment in `editorState.go` that names them). It is the same type under a new name, so
   the frozen bytes cannot move.
3. **Suffix naming**, mirroring `internal/models/template_model/template_variant_model`.
4. **`template_entity/types.go` stays** and keeps re-exporting all 40 subpackage types.
5. **One phase**, one gate run at the end.

## Corrections to §2.5 as written (found during scoping — do not follow the item verbatim)

- Its **step 1 is already done**: a canonical `types.go` already lives *inside* the
  protected tree at `internal/entities/template/types.go`, aliasing the six subpackages.
  The outer `internal/entities/types.go` is an alias *of* those aliases and its own doc
  comment says it should be removed. Nothing moves *into* the tree; the outer file is
  deleted. The item's "open question for the owner" is therefore moot.
- Its **step 1 says `git mv`** — forbidden. AGENTS.md §2.5 says never stage; `git mv`
  stages. Use `Move-Item`.
- Its **step 4 says `gopls rename`** — not available to this agent. Every import line and
  package clause is edited individually with the edit tools. No bulk in-place rewrite
  (AGENTS.md §2.6), and **no `Get-Content`/`Set-Content` round-trip on a `.go` file**.
- `goimports` is banned. Imports are fixed by hand against the compiler.

## Naming map

| Old | New |
| --- | --- |
| `internal/entities/template` (`package template`) | `internal/entities/template_entity` (`package template_entity`) |
| `template/template_common` | `template_entity/template_common_entity` |
| `template/template_content` | `template_entity/template_content_entity` |
| `template/template_layout` | `template_entity/template_layout_entity` |
| `template/template_override` | `template_entity/template_override_entity` |
| `template/template_rule` | `template_entity/template_rule_entity` |
| `template/template_variant` | `template_entity/template_variant_entity` |
| `internal/entities` alias shim (`entities.Zone`) | **deleted** — consumers say `template_entity.Zone` |

Directory name equals package name throughout, before and after, so no import alias is
needed anywhere.

## Measured blast radius (taken before any edit)

- 31 `.go` files under the protected tree, all needing a `package` clause edit:
  root 2, `template_common` 1, `template_content` 7, `template_layout` 4,
  `template_override` 1, `template_rule` 5, `template_variant` 11.
- 45 files / 52 lines reference the `internal/entities/template` import path.
- 28 files import the `internal/entities` shim (5 production, 23 test).
- Docs naming the old path: `AGENTS.md` (§2.1, §4.4 table, TL;DR), `README.md`,
  `.agent/backlog/backlog-opus5.md`, `.agent/promt_templates/review-prompt.md`,
  `.agent/memories/*`.

## Phase 1: Rename and de-shim

Status: Complete

- [x] Capture the pre-change coverage figure and confirm the gate baseline is green
      (build, vet, lint, testlayoutcheck, all three suites).
- [x] Record a hash/listing of the protected tree's file contents so the "no content
      change beyond the package clause" claim can be proved at the end.
- [x] `Move-Item` the six subdirectories to their `_entity` names, then `Move-Item`
      `internal/entities/template` → `internal/entities/template_entity`.
- [x] Edit the `package` clause in all 31 moved files.
- [x] Fix the intra-tree sibling imports.
- [x] Update the import lines across the consumer files.
- [x] Delete `internal/entities/types.go` (`Remove-Item`).
- [x] Repoint the 28 shim consumers, including the two frozen `editor_state_v1`
      files, per approval 2.
- [x] Regenerate wire — **not needed**, `wire_gen.go` never named the old path.
- [x] Update the prose comments that spell `entities.Zone` / `entities.Connection`.
- [x] `gofmt -w` on the explicit list produced by `gofmt -l` — three files.
- [x] Update the docs so the protected path in AGENTS.md §2.1 is accurate.
- [x] Close backlog §2.5, open row U in §8, refresh the coverage note.

### Verification Plan

Run every gate from backlog §9 and compare against the pre-change capture:

- `go build ./...` → clean
- `go vet ./...` and `go vet -tags='integration_test,gui' ./...` → clean
- `gofmt -l ./app ./internal ./test ./cmd` → empty
- `go run ./cmd/testlayoutcheck .` → `test-layout check passed`
- `go test ./test/unit/... -count=1` → pass
- `go test -tags=integration_test ./test/integration/... -count=1` → pass
- `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` → pass,
  **without `-update`**; no golden may move — a rename cannot change a pixel
- `golangci-lint-v2 run ./...` → **0 issues**
- Coverage → **≥ 72.5 %**, expected to land within rounding of the 74.3 % baseline; a
  pure rename adds and removes no statements
- **Wire-format proof:** the `.rmg.json` and `.gen.json` fixture tests pass unmodified,
  and `data/`, `internal/registry/` are untouched
- **Content proof:** a diff of the moved tree shows only `package` clauses and import
  paths changed — zero JSON tags, field names or field orders moved

### Phase Summary

Done. Every gate was run on the finished tree.

| Gate | Result |
| --- | --- |
| `go build ./...`, `go vet ./...`, `go vet -tags='integration_test,gui' ./...` | 0 |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | passed |
| `wire diff ./internal/composition/...` | no diff |
| `go test ./test/unit/... -count=1` | pass |
| `go test -tags=integration_test ./test/integration/... -count=1` | pass |
| `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | pass, **no `-update`**, no golden moved |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Coverage | 74.2 % → **74.3 %** (floor 72.5 %) |
| Content proof | all **31** protected files byte-identical to `HEAD` once the rename is normalized away |

Three things worth carrying forward:

1. **The item's own recipe was wrong twice**, and the corrections now live in backlog
   §2.5: `git mv` stages, and its step 1 had already been done long before.
2. **A dead architecture gate was found and revived.**
   `test/unit/architecture/construction/construction_test.go` keys on a single import
   path constant that pointed at the shim. Since batches P/Q no production file
   imported it, so the gate scanned an empty set and passed silently. Repointing it
   made it fire on seven `ToXEntity` converters.
3. **"Use builders" is not reachable for entities.** Since batch Q every builder returns
   a *model*; none imports `internal/entities` and none may. Entities are constructed in
   production in exactly one place — the Model→Entity converter seam — so by owner
   decision `internal/models/template_model/` joins `internal/services/builders/` in
   `shouldSkipPath`, with the reasoning in a comment beside it.

One cosmetic fix rode along: a `golines` finding in
`test/unit/.../manualConnectionSave/clone_test.go`, where the longer type name pushed a
literal past the limit. Wrapped by hand in house style, not via `--fix`.

Also corrected while in the file: the `editor_state_v1` package doc pointed at
`internal/services/editor_state_migration`, which does not exist — the real package is
`internal/services/file_service/editor_state_migrator`. A stale pointer left by batch S
in the one file whose whole job is to explain the freeze.

## Final Recap

Batch T renamed the `.rmg.json` entity tree from the generic `template` to
`template_entity`, with the six subpackages taking the `_entity` suffix that mirrors the
existing `template_model` convention, and deleted `internal/entities/types.go` — the
alias-of-alias shim whose own comment had been asking to go. Its 28 consumers now name
the real package: `entities.Zone` is `template_entity.Zone`. Two owner-approved
protected edits: the 31 package clauses inside the tree, and an import-only edit to the
two frozen `editor_state_v1` files. Nothing else moved, proved file by file against
`HEAD`.

The unplanned find is the more valuable half. The construction gate had been dead since
batch P/Q and nobody could have noticed, because a gate keyed to an unused import path
fails open. It is live again, with the Model→Entity converter seam exempted for the same
reason builders are.

## Deployment Plan

1. Review the diff. The moves used `Move-Item`, so git sees deletes plus untracked files
   until they are staged — `git status` looks noisier than the change is. Staging
   resolves them into renames.
2. Nothing to migrate or regenerate. `.rmg.json` and `.gen.json` are untouched,
   `wire_gen.go` is unchanged, no golden or fixture moved. There is **no runtime
   behaviour change**, so this batch needs no in-app smoke test — unlike batch S, whose
   smoke test is still outstanding underneath it.
3. Commit, then delete this plan file. Backlog §2.5 and §8 row T stand alone.
4. Batch S remains committed-but-unpushed below this work; decide whether to push both
   together.
