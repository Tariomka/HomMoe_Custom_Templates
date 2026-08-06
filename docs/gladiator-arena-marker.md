# Crossed Swords marker — the `GladiatorArena` object

This note documents the "crossed swords" icon seen on three of the official
Olden Era template previews, what configuration produces it, and the preview
assets extracted for it.

## What the crossed swords mean

The crossed-swords icon marks a **Gladiator Arena** (`GladiatorArena`) object
placed on the map. It is the physical arena where heroes register and fight
for the *Gladiator Arena* victory condition. It is **not** a neutral zone of
its own — it is an object that is either dropped into a zone or attached to a
connection between two zones.

The icon appears in exactly three official templates:

| Template  | Where it sits in the preview        | How it is defined                                   |
| --------- | ----------------------------------- | --------------------------------------------------- |
| Blitz     | Inside a **high** (gold) neutral zone | Zone **main object** `"type": "GladiatorArena"`     |
| Symmetry  | Inside a **low** neutral zone        | **Connection** `"connectionType": "GladiatorArena"` |
| Helltide  | On the **central connector** (no zone) | **Connection** `"connectionType": "GladiatorArena"` |

## The exact configuration that determines the icon

There are **two** wire forms, both using the value `GladiatorArena`:

### 1. As a zone main object

A zone's `mainObjects` array contains an entry whose `type` is `GladiatorArena`.
The arena replaces the usual city/treasure object for that zone, so the
crossed swords render on top of that zone's quality bubble.

```jsonc
// data/ExampleTemplates/Blitz.rmg.json  -> zone "SuperTreasure-3"
{
  "name": "SuperTreasure-3",
  "layout": "zone_layout_supertreasure_zone",   // high-quality (gold) neutral zone
  "mainObjects": [
    {
      "type": "GladiatorArena",
      "placement": "Uniform",
      "placementArgs": [ "true", "0", "0" ]
    }
  ]
}
```

### 2. As a connection type

A connection between two zones declares `connectionType: "GladiatorArena"`.
The arena is drawn at the connection mid-point. Depending on the layout that
mid-point can fall on a bare connector (Helltide) or inside a neutral zone
that happens to sit there (Symmetry).

```jsonc
// data/ExampleTemplates/Helltide.rmg.json  -> "Win-Connection"
{
  "name": "Win-Connection",
  "from": "Center-Win-A",
  "to": "Center-Win-B",
  "connectionType": "GladiatorArena",
  "road": true,
  "guardValue": 5000000
}
```

```jsonc
// data/ExampleTemplates/Symmetry.rmg.json  -> "Arena-Connection"
{
  "name": "Arena-Connection",
  "from": "Side-A2",
  "to": "Side-B2",
  "connectionType": "GladiatorArena",
  "road": true,
  "guardValue": 45000
}
```

**Rule of thumb:** the crossed-swords marker is present iff some zone main
object OR some connection has the value `GladiatorArena`.

## Cross-reference against the other templates

Searching every template in `data/ExampleTemplates/` for the value
`GladiatorArena` (i.e. `: "GladiatorArena"`) returns **only** these three
matches:

- `Blitz.rmg.json` — `"type": "GladiatorArena"`
- `Helltide.rmg.json` — `"connectionType": "GladiatorArena"`
- `Symmetry.rmg.json` — `"connectionType": "GladiatorArena"`

No other template (Anarchy included) places a `GladiatorArena` object, so the
crossed-swords marker is unique to these three.

### Do not confuse it with the win-condition flag

The lower-case `gladiatorArena` flag under
`gameRules.winConditions` only *enables the arena victory rule*; it does not
place an arena object and therefore does not draw the icon. That flag is
`true` in **four** templates:

- Blitz, Helltide, Symmetry — which also place the object (icon shown), and
- **Jebus Outcast** — which enables the rule but places **no** `GladiatorArena`
  object, so it shows **no** crossed-swords icon.

This confirms the icon is driven by the placed object/connection, not by the
win-condition flag.

## Extracted preview assets

The crossed-swords glyph was extracted from the official `Helltide.png`
preview (the clean version sitting on the bare central connector, with the
connector line masked out) and matched to the existing 96x96 sprite format
(bubble center at `(48,48)`). The assets live under
[internal/services/asset_provider/assets/](../internal/services/asset_provider/assets):

| File                         | Contents                                                 |
| ---------------------------- | -------------------------------------------------------- |
| `gladiator_arena.png`        | Crossed swords only, transparent background              |
| `neutral_none_arena.png`     | Crossed swords over a **none** (open-ring) zone bubble   |
| `neutral_low_arena.png`      | Crossed swords over a **low** (bronze) zone bubble       |
| `neutral_medium_arena.png`   | Crossed swords over a **medium** (silver) zone bubble    |
| `neutral_high_arena.png`     | Crossed swords over a **high** (gold) zone bubble        |
| `neutral_highest_arena.png`  | Crossed swords over a **highest** (platinum) zone bubble |

Notes:

- `gladiator_arena.png` is the master glyph (solid, fully opaque swords). It is
  composited at ~0.85 scale (so the swords are about as big as the zone's outer
  ring, matching the Blitz preview) onto the `neutral_*.png` bubbles to produce
  the zone-background variants.
- The neutral-zone quality bubbles are: **none** = open ring (transparent
  center), **low** = light bronze, **medium** = silver, **high** = gold. The
  bronze **low** bubble is derived from the gold **high** bubble by recolouring
  the fill (keeping the ring and shading).
- The **medium** variant has no official source (no official template places an
  arena in a medium zone); it is synthesised from the master glyph for
  completeness of the sprite set.
- These files are embedded by the `//go:embed assets/*.png` directive in
  [internal/services/asset_provider/assetProvider.go](../internal/services/asset_provider/assetProvider.go).

## How this project places and draws the arena

The generator emits an arena whenever the template is in arena mode — either
the `gladiatorArena` win-condition rule is enabled, or the victory condition is
`win_condition_4` ("Guardian Arena") — as decided by
`config.GeneratorConfig.IsGladiatorArenaMode()`.

`providers.GladiatorArenaProvider.PlaceArena` then picks the wire form to use
based on what the chosen topology actually produced, in this order:

1. **Hub zone present** → a `GladiatorArena` main object is appended to the hub
   zone (`"placement": "Uniform"`, `"placementArgs": ["true","0","0"]`, matching
   Blitz).
2. **Otherwise, a neutral↔neutral connection exists** → the richest such
   connection (highest combined zone quality, ties broken by connection name)
   gets `"connectionType": "GladiatorArena"`.
3. **Otherwise** → the richest neutral zone (ties broken by zone name) receives
   the main object.
4. **No neutral zones at all** → nothing is placed.

Most topologies alternate player and neutral zones around a ring, so
neutral↔neutral connections often do not exist; step 3 is the common fallback.

Rendering mirrors the two forms:

- A zone carrying the main object renders the `neutral_<quality>_arena.png`
  bubble. There is no combined castle+arena artwork, so a zone with both a
  castle and an arena draws the **arena** sprite.
- An arena connection draws the plain connector line and then composites the
  `gladiator_arena.png` master glyph at the midpoint of the connection's
  quadratic Bézier curve, at 0.75 of the zone-sprite scale.

