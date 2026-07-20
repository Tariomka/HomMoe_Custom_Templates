# Backlog

Small future-work items moved out of code comments (godox purge, review §5.5).

- **Preview sub-pixel precision**: `internal/models/preview/previewLayout.go` —
  `Layout.Positions` uses `map[string]image.Point`; switching to a `Vec2`
  (float64) type would give sub-pixel precision in preview rendering and zone
  editor geometry. Related: `Zone.GeneratorPosition *[2]float64` in the
  read-only schema would ideally share that Vec2 type.

- remove the [2]float from the template entities for example entities/template/template_variant/zone.go. It's easier to use Vec2 instead

- **`createTopologyAdjacency` dead Chain/Ring branches**:
  `internal/services/zones/zoneLabelProvider.go` — the `case TopologyChain` and
  `case TopologyRing, TopologyCircles` branches (plus the `isIsolated` guard
  they use) are unreachable: the only caller, `GetHoldCityLabel`, gates on
  `IsHubCityToHold()` (= `Topology == HubAndSpoke && IsCityHoldMode()`), so the
  switch always takes `default`. Verified 2026-07 (review §5.5): single private
  call site, single production caller (`templateGenerator.Generate`), branches
  never reachable in any commit (all three symbols born together in `bb50aab`),
  0% test coverage on them. A removal was implemented and verified
  (build/tests/lint green, coverage 64.1→64.2) but ROLLED BACK by the owner to
  keep the topology-aware adjacency as a starting point. Decide eventually:
  - either delete the branches (pure ratio win, zero behavior change), or
  - start using them: extend hold-city (or other adjacency-based features) to
    Chain/Ring/Circles topologies, which would also fix the `default` branch
    modelling Hub & Spoke as a sequential ring instead of its real star graph.

- need to add untracked zone tier property to Entities and/or method. This will be the source of truth for all tier related operations.
  As currently there is no reading of rmg.json files, deserialization doesn't matter so the entity can assume the tier from the content currently inside the zone.
  Or maybe Zone doesn't actually need a property, neutralZone.Quality can be used to track and infer this information as well as neutralZone.Profile can have this property.

- need to consolidate road distances - there are multiple implementations and UI uses services directly

- need to add that anything inside app should only use entities, models, handlers and commons, not services

- need to use common (either from commons or models) values for template generation

- either random portals or connections in general to hub have incorrect guard values
