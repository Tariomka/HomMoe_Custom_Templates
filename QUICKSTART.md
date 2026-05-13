# Quick Start — Olden Era Custom Templates GUI

This is a desktop GUI app, not a CLI. Launching it opens a window where
you configure a template across four tabs and then click **Generate** to
write a `.rmg.json` file.

## 1. Run the App

### Use the bundled binary

```powershell
.\bin\template-gui.exe
```

### Build from source

Requires Go **1.25.8+**.

```powershell
# Run directly
go run .

# Or build a binary first
go build -o bin/template-gui.exe .
.\bin\template-gui.exe
```

## 2. The Window

The window has three regions:

- **Toolbar** (top): `New`, `Open`, `Save`, `Save As`, `Templates folder`,
  `Discord`, `GitHub`, `Patch notes`.
- **Tabs** (centre): the four configuration tabs listed below.
- **Footer** (bottom): output folder picker, `Reveal output`,
  `Generate`, `Save template`.

### Tabs

#### 1. Map Setup
- **Template name** — used as the output filename (`<name>.rmg.json`).
- **Game mode** — `Classic` or `SingleHero` (only `Classic` is currently emitted).
- **Players** — 2..8.
- **Map size** — slider in tile units (64..240). Tick
  *"Allow experimental large map sizes (>240)"* to extend to 256..512.
- **Topology** — `Random` (default), `Ring`, `Hub`, `Chain`, `Shared Web`.

#### 2. Generation Options
- **Connectivity** — toggle roads, random portals (with max-portal slider),
  remote footholds, experimental balanced placement, player isolation,
  faction-matched player castles, minimum neutrals between players.
- **Advanced neutral zones** — when enabled, set per-quality counts split
  between *with castle* and *without castle* (Low / Medium / High).
  When disabled, just set the total `Neutral zone count`.
- **Zone sizes** — player / neutral / hub multipliers, hub castle count,
  guard randomization.
- **Content density** — resource and structure density percents
  (independent or via the unified `Content density` knob).

#### 3. Game Rules
- **Victory condition** — Standard / Lost Starting City / Hold City / Tournament.
  When Hold City is selected the generator also picks a neutral zone to
  host the city-hold target.
- **Heroes** — min, max, increment.
- **Faction laws XP %** and **Astrology XP %** — 25..200 (100 = baseline).
- **Lost-city / lost-hero / city-hold** day toggles.
- **Gladiator arena** — start delay and counter day.
- **Tournament** — first day, interval, points-to-win, save-army.

#### 4. Zone Content (EXP)
Add extra mandatory content items to seed into player zones in addition
to the defaults from `services.ZoneContentManager`.

## 3. Save / Load Settings

Settings are persisted as `.oetgs` JSON files (the
`models.SettingsFile` struct).

- **Save / Save As** — write the current widget state to disk.
- **Open** — load a `.oetgs` file back into the widgets.
- **New** — reset to defaults.

## 4. Generate a Template

1. Pick an output folder in the footer (`...` button next to the path).
2. Click **Generate**.
3. The app writes `<TemplateName>.rmg.json` into that folder.
4. **Reveal output** opens the folder in Explorer.

Drop the resulting `.rmg.json` into the game's templates directory and
pick it from the in-game Random Map Generator screen.

## 5. Programmatic Use

If you want to call the generator from Go without the GUI:

```go
package main

import (
    "encoding/json"
    "os"

    "github.com/Tariomka/hommoe_custom_templates/internal/models/generator"
    "github.com/Tariomka/hommoe_custom_templates/internal/services"
)

func main() {
    s := generator.NewGeneratorSettings() // sensible defaults
    s.TemplateName = "Programmatic Map"
    s.PlayerCount = 4
    s.MapSize = 160
    s.Topology = generator.TopologyDefault // Ring

    tmpl, err := services.Generate(s)
    if err != nil {
        panic(err)
    }

    data, _ := json.MarshalIndent(tmpl, "", "  ")
    _ = os.WriteFile("Programmatic Map.rmg.json", data, 0o644)
}
```

`generator.NewGeneratorSettings()` matches the GUI's defaults
(2 players, size 160, topology Random, Classic mode, etc.).

## 6. Map Sizes

Map size is the raw tile count (`sizeX = sizeZ = MapSize`). Common
presets the slider snaps to:

| Tiles | Notes                                  |
|-------|----------------------------------------|
|  64   | Tiny                                   |
|  96   | Small                                  |
| 128   | Medium                                 |
| 160   | Large (default)                        |
| 192   | Extra Large                            |
| 240   | Maximum non-experimental               |
| 256+  | Experimental (enable the checkbox)     |
| 512   | Largest experimental size              |

## 7. Topologies at a Glance

- **Random** — default; positions and connections are randomised.
- **Ring** (`TopologyDefault`) — players in a circle, neighbours connected.
- **Hub** (`TopologyHubAndSpoke`) — all players connect through a central hub zone.
- **Chain** — players strung in a line.
- **Shared Web** — players linked through shared neutral zones.

## 8. Troubleshooting

- **Window doesn't open** — Gio needs a working graphics stack. If
  running over Remote Desktop or in a headless container, run on a real
  desktop session instead.
- **`go build` complains about Go version** — install Go 1.25.8 or later
  (`go version`).
- **`Generate` does nothing** — check that the output folder is set and
  writable; the footer status line shows the last error.
- **Template name is empty** — the generator returns an error; pick a name.
- **Template won't load in the game** — verify the JSON is valid
  (`Get-Content map.rmg.json | python -m json.tool`) and compare against
  files in `data/ExampleTemplates/`.

## 9. Reference

- [README.md](README.md) — architecture, project layout, build/test commands.
- [MIGRATION.md](MIGRATION.md) — history of the C# WPF → Go GUI port.
- [CONVERSION_SUMMARY.md](CONVERSION_SUMMARY.md) — checkpoint summary of the port.
- `data/ExampleTemplates/` — 57 reference templates shipped with the game.
