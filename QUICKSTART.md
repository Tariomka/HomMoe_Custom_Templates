# Quick Start Guide - Template Generator (Go)

## Installation

### Option 1: Build from Source

Requirements:
- Go 1.25.8 or later

```bash
# Clone or navigate to the project
cd "Template Generator"

# Build the binary
go build -o template-generator ./cmd

# The binary is ready to use
./template-generator -help
```

### Option 2: Use the Pre-built Binary

If a pre-built binary is available, download it and run:
```bash
./template-generator -help
```

## Basic Usage

### Generate a Simple Template

```bash
./template-generator \
  -name="My First Map" \
  -players=4 \
  -size=L \
  -output=.
```

This creates: `My First Map.rmg.json`

### Generate with Custom Settings

```bash
./template-generator \
  -name="Epic Battle" \
  -players=6 \
  -size=2XL \
  -topology=HubAndSpoke \
  -game=Heroic \
  -roads=true \
  -portals=true \
  -footholds=true \
  -cityhold=true \
  -output=./my_templates
```

## Command-Line Options

```
-name string
    Template name (default "Generated Template")

-players int
    Number of players (default 4)

-size string
    Map size: S, M, L, XL, 2XL (default "L")

-topology string
    Topology type: Default, Chain, HubAndSpoke, SharedWeb, Random (default "Default")

-game string
    Game mode: Classic, Blitz, Heroic (default "Classic")

-roads boolean
    Allow roads between zones (default true)

-portals boolean
    Allow portals between zones (default true)

-footholds boolean
    Allow footholds (default false)

-cityhold boolean
    Enable city hold win condition (default false)

-output string
    Output directory (default ".")
```

## Examples

### Tournament Map (4 Players)
```bash
./template-generator \
  -name="Tournament" \
  -players=4 \
  -size=L \
  -topology=Default \
  -output=./templates
```

### Large Free-for-All (8 Players)
```bash
./template-generator \
  -name="Free for All" \
  -players=8 \
  -size=2XL \
  -topology=Random \
  -roads=false \
  -output=./templates
```

### Hub and Spoke (6 Players)
```bash
./template-generator \
  -name="Spoke" \
  -players=6 \
  -size=XL \
  -topology=HubAndSpoke \
  -output=./templates
```

### Blitz Mode (2v2)
```bash
./template-generator \
  -name="Blitz 2v2" \
  -players=4 \
  -size=M \
  -game=Blitz \
  -topology=Chain \
  -output=./templates
```

## Understanding Map Sizes

- **S (Small)**: Fast-paced, testing maps - 72x72 tiles
- **M (Medium)**: Quick games - 108x108 tiles
- **L (Large)**: Standard play - 144x144 tiles (default)
- **XL (Extra Large)**: Long games - 180x180 tiles
- **2XL (Double Extra Large)**: Epic battles - 216x216 tiles

## Understanding Topologies

### Default (Ring)
Players arranged in a circle, each connected to neighbors.
```
    1
  6   2
5       3
  4
```

### Chain
Players in a line, each connected linearly.
```
1 -- 2 -- 3 -- 4
```

### HubAndSpoke
All players connected through a central neutral zone.
```
    1
    |
5--HUB--2
    |
    4
    |
    3
```

### SharedWeb
Players connected through neutral zones.
```
1 -- N1 -- 2 -- N2 -- 3
```

### Random
Random connections with Delaunay triangulation logic.
```
(variable topology based on zone positions)
```

## File Output

The generator creates a `.rmg.json` file with:
- Template metadata (name, description, size)
- Game rules (win conditions, hero counts)
- Zones (player and neutral)
- Connections (roads, portals, guard zones)
- Layout information

## Troubleshooting

### "command not found: template-generator"
- Make sure you're in the correct directory
- Use `./template-generator` with the leading `./`
- Or add to PATH: `export PATH=$PATH:$(pwd)`

### Output directory doesn't exist
- Create the directory first: `mkdir -p ./my_templates`
- Or use an existing directory like `.` (current directory)

### File already exists
- The tool will overwrite existing files
- Rename or move existing `.rmg.json` files first

### Template won't load in game
- Verify the JSON is valid: `cat your-template.rmg.json | python -m json.tool`
- Check example templates for correct structure
- Ensure map size and zone count are appropriate

## Programmatic Usage (Go)

```go
package main

import (
    "encoding/json"
    "os"

    "github.com/Tariomka/hommoe_custom_templates/internal/generator"
    "github.com/Tariomka/hommoe_custom_templates/internal/models"
)

func main() {
    settings := &models.GeneratorSettings{
        TemplateName: "Programmatic Map",
        PlayerCount: 4,
        MapSize: "L",
        Topology: models.TopologyDefault,
    }

    template, err := generator.Generate(settings)
    if err != nil {
        panic(err)
    }

    // Marshal to JSON and save
    data, _ := json.MarshalIndent(template, "", "  ")
    os.WriteFile("output.rmg.json", data, 0644)
}
```

## Reference

- See `README.md` for technical details
- See `MIGRATION.md` for C# to Go conversion notes
- Example templates: `data/GameData/ExampleTemplates/`

## Support

For issues or questions:
1. Check the README.md
2. Review AGENTS.md for development guidelines
3. Examine example templates for format reference
