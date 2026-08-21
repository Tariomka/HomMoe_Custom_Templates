# Session Carry-Forward — 2026-08-21

## 1. Session goal

Plan batch **I** (backlog §2.1, the `EditorStateDto` rework, plus §1.5 folded in)
and produce its plan file. **No production code was to be written, and none was.**

## 2. Fixes applied

No production or test code was touched — this was a planning session.

Two **documentation** defects were fixed in
[todo/backlog-opus5.md](../todo/backlog-opus5.md):

- **Seven dangling links.** Every `../plans/` reference in the backlog pointed at
  a deleted file — the batch F, G and H plans are all gone, so §2.3, §5.1, §5.2,
  §5.3, §5.4's blockquote and both the G and H rows of the §8 batch table linked
  to nothing. Each entry now carries its own record instead of delegating to a
  document that no longer exists, which is what the doc-lifecycle rule asks for.
- **Two batch H notes with no durable home** were folded in before this file was
  overwritten (see §3).

The only other corrections were to the batch I plan itself, after an independent
review falsified six of its claims (see §7 and §8).

## 3. Features added / changed

None in the codebase.

**Two batch H notes were rescued into the backlog.** Both existed only in the
deleted batch H plan and in the previous version of this file, so the next
overwrite would have destroyed them:

- The **arm64 float-pin risk** is now recorded in §5.1, together with the reason
  rounding the pin is not an acceptable fix and what to do instead (an `InDelta`
  far tighter than a pixel).
- The **golden-footprint question** is now an open item in §5.2, with the numbers
  and the lever to pull if the size becomes a problem.

**Design settled for batch I** (16 owner questions across four rounds — all
answered, all forks closed; treat these as decided and **do not relitigate**):

- **Anonymous embedding all the way.** Nine behaviour-free entity group structs
  are embedded anonymously into a new `EditorStateModel`. Go's field promotion
  keeps `state.MapSize` compiling at every one of the ~1,000 access sites, and
  `encoding/json` flattens anonymously embedded structs, so the on-disk shape
  stays flat for free. This is what makes a 72-field regroup tractable.
- **The model becomes the runtime type.** `EditorStateDto` survives only at the
  load/save boundary, shrinking to
  `{SchemaVersion int; EditorState EditorStateModel}` with a **named** field plus
  custom `MarshalJSON`/`UnmarshalJSON` that re-flattens the wire format and adds
  `schemaVersion` as a sibling key (always written as `1`; a `0` read from an
  existing file is normalised through an explicit migration hook).
- **json tags live on the entity leaf fields.** Entities are strictly
  behaviour-free, so per §4.6 they need no tests of their own; all behaviour
  moves to the model.
- **Packages:** `internal/entities/editor_state/`,
  `internal/models/editor_state_model/` (named this way to dodge a package-name
  collision with the entity package), the three `internal/dtos/editorState*Dto.go`
  files move into `internal/dtos/editor_state_dto/`, content-row defaults move to
  `internal/common/common_zone_contents/`, and §1.5's per-panel view structs land
  in `app/gui/models/`.
- **Groups (72 fields, 9 groups):** `templateIdentity` (2), `mapSettings` (2),
  `playerSettings` (4), `neutralZoneSettings` (11), `castleSettings` (10),
  `generationSettings` (15), `gameRuleSettings` (16), `contentSettings` (10),
  `manualEditSettings` (2). There is deliberately **no** `hubZoneSettings` group —
  `HubZoneSize` went to `generationSettings` with the other two zone sizes, and
  `HubZoneCastles` to `castleSettings`.

## 4. File modifications

| File | Change |
| --- | --- |
| [plans/batch-i-editor-state-rework.md](../plans/batch-i-editor-state-rework.md) | **Created.** The batch I plan: design-decision table, field→group map, phase-ordering rationale, nine numbered hazards, seven phases each with checkbox items / Verification Plan / empty Phase Summary, and Final Recap + Deployment Plan placeholders. |
| [.agent/session-carry-forward.md](session-carry-forward.md) | **Rewritten** (this file) — replaced the batch H handoff. |
| [todo/backlog-opus5.md](../todo/backlog-opus5.md) | Repaired seven dangling `../plans/` links; folded the arm64 float-pin caveat into §5.1 and the golden-footprint open question into §5.2. |

Repository memory (`/memories/repo/conventions.md`, outside the repo tree) also
gained a "Batch I PLANNED" block recording the census numbers, the owner
decisions and the six review traps, because the plan file will be deleted at the
end of batch I per the doc-lifecycle rule.

**No file under `app/`, `internal/`, `test/`, `cmd/` or `data/` was touched.**

## 5. Tests added or updated

None. No test run was performed this session, because nothing was changed — the
last known-good state is the one batch H left behind (all suites green,
coverage 72.8 %).

The plan specifies the test work batch I must do, notably:

- A parsed-value golden round-trip test guarding the `.gen.json` wire shape,
  backed by a new all-fields fixture in `test/test_helpers/testdata/`. **No
  checked-in `.gen.json` fixture exists today** — the round trip is currently
  only exercised through temp files, which is why the golden has to be generated
  from *current* code in Phase 1, before anything moves.
- A rewrite of the equality drift guard in
  [test/unit/internal/dtos/editorStateDto/equalsIgnoringManualEdits_test.go](../test/unit/internal/dtos/editorStateDto/equalsIgnoringManualEdits_test.go#L211-L230),
  which walks `NumField()` at the top level only and so would silently stop
  covering every field once they move into embedded groups.

## 6. Git status snapshot

- **Branch:** `AD/fixing_some_stuff_08-12`
- **HEAD:** `94bc9a4 agents update` (in sync with `origin/`)
- **`git status --short`:** `?? plans/` — plus this file will now show as
  ` M .agent/session-carry-forward.md`, since it is tracked (committed in
  `4586f88`).

**Nothing was staged and nothing was committed.** The staging area was left
untouched, as required.

## 7. Rejections / things the user declined

Nothing the owner declined outright, but three of my proposals were **overruled
during the question rounds**, and the plan reflects their answer, not mine:

- I proposed the DTO embed the model anonymously (symmetry with everything
  else). The owner required a **named** `EditorState` field plus a
  `schemaVersion` sibling — which I had to point out would nest the JSON and
  break the flat-shape invariant. The owner's resolution was custom marshalling
  to keep the file flat.
- I proposed pushing the `editorState*Dto` internal structs into a subpackage of
  their own; the owner ruled there are too few of them to justify it, so they sit
  alongside the other DTOs.
- The independent review argued `internal/entities/editor_state/` contradicts
  AGENTS.md §4.4's placement table. **The owner's earlier decision stands** and
  the plan carries an explicit "this is deliberate, do not 'fix' it" note. The
  review's other half of that finding — that "commit the fixture" violates the
  no-commit rule — was accepted and reworded.

## 8. Open questions

1. **Owner sign-off on the plan itself.** Blocking; no phase should start
   without it.
2. **Phase 2 owner gate — the field→group table.** The 72-field split is the one
   decision the whole refactor is built on and the one that is expensive to
   revisit later, so Phase 2 deliberately ends waiting for an explicit ack.
3. **`internal/entities/editor_state/` vs AGENTS.md §4.4.** Settled by owner
   decision, but §4.4's table still says otherwise, so a future reviewer will
   likely raise it again. Worth either amending §4.4 or leaving the note in
   place permanently.
4. **Is a golden per handler action the right granularity?** Inherited from
   batch H and now recorded in backlog §5.2 rather than left to be lost. Not
   blocking anything, but it compounds with every new driving handler.

## 9. Next recommended actions

1. **Read and sign off** [plans/batch-i-editor-state-rework.md](../plans/batch-i-editor-state-rework.md).
2. **Then start Phase 1** of the plan: record before-coverage; build an
   all-fields DTO and marshal it with *current* code into
   `test/test_helpers/testdata/editorState_v0_flat.gen.json`; add the
   parsed-value golden round-trip test; move the three `editorState*Dto.go`
   files into `internal/dtos/editor_state_dto/`; strip the self-qualifiers; add
   the import to the four sibling DTOs that carry the type; update the ~12
   consumer packages file-by-file (**never** a bulk rewrite); move the mirrored
   unit-test folders.
3. Work strictly **Phases 1 → 7 in order** — the ordering is what keeps the
   build green at every boundary.
4. When batch I's plan is deleted at the end of §7, **fold its record into the
   backlog in the same pass** — the seven dead links fixed this session all came
   from skipping that step.

## 10. Carry-forward prompt

> Read `AGENTS.md` first, then `plans/batch-i-editor-state-rework.md` — the batch
> I plan (backlog §2.1, the `EditorStateDto` rework) is written, independently
> reviewed and corrected, but **no code has been written yet**. It is awaiting
> owner sign-off; ask before starting Phase 1.
>
> The hard rules, one line each: never modify `data/`,
> `internal/entities/template/` or `internal/registry/` without explicit
> approval; everything must build and run on both Windows and Linux (use
> `path/filepath`, chain PowerShell with `;`, never `&&`); every change ships
> with tests and unit coverage must not drop below 72.5 % (currently 72.8 %);
> durable multi-session work gets a plan file under `plans/`; never stage and
> never commit — the owner reviews, stages and commits, so leave the staging area
> alone entirely, and delete files with `Remove-Item`, never `git rm`; never
> change where `.rmg.json` is written and never persist the output directory;
> never run a bulk in-place rewrite over the repository; never run CI and never
> generate snapshot goldens in CI — generate them locally on the real GPU, always
> `-run`-scoped and scoped tightly enough not to match neighbouring tests.
>
> Where work left off: the plan is complete through all seven phases and the
> design is fully settled — nine behaviour-free entity groups anonymously
> embedded into `EditorStateModel`, which becomes the runtime type, with
> `EditorStateDto` shrinking to a versioned persistence shell whose custom
> marshalling keeps `.gen.json` flat. Do not relitigate the design decisions in
> §0 of the plan; they are the owner's answers to 16 questions. Three review
> findings are load-bearing and easy to trip over again: Phase 3 **must** keep
> DTO-signature shim methods or the build breaks on promoted methods;
> `omitempty` already conflates nil and empty on disk, so that distinction is
> in-memory only and must not be "fixed"; and regrouping reorders JSON keys, so
> every round-trip gate compares parsed objects, never bytes.
>
> Full handoff in `./.agent/session-carry-forward.md`. Note §9.4: when a plan
> file is deleted, its record must be folded into the backlog **in the same
> pass** — skipping that is what left seven dead links behind this session.
