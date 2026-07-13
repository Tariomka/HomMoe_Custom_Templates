# Backlog

Small future-work items moved out of code comments (godox purge, review §5.5).

- **Preview sub-pixel precision**: `internal/models/preview/previewLayout.go` —
  `Layout.Positions` uses `map[string]image.Point`; switching to a `Vec2`
  (float64) type would give sub-pixel precision in preview rendering and zone
  editor geometry. Related: `Zone.GeneratorPosition *[2]float64` in the
  read-only schema would ideally share that Vec2 type.
