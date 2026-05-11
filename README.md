# Heroes of Might and Magic: Olden Era - Template Generator (Go)

This is a Go implementation of the random map template generator for Heroes of Might and Magic: Olden Era.

## Project Structure

```
.
├── cmd/                          # CLI entry point
│   └── main.go                  # Template generator command-line tool
├── internal/
│   ├── generator/               # Core generation logic
│   │   ├── template_generator.go
│   │   └── template_generator_test.go
│   ├── models/                  # Data structures
│   │   ├── template.go          # .rmg.json template models
│   │   └── settings.go          # Generator settings models
│   └── services/                # Content and utility services
│       └── zone_content_manager.go
├── data/
│   └── GameData/                # Example templates and generator data
│       ├── ExampleTemplates/    # Reference .rmg.json files (57 examples)
│       └── GeneratorData/       # Game configuration and content pools
├── go.mod                        # Go module definition
└── README.md                     # This file
```

## Features

- **Template Generation**: Create `.rmg.json` random map template files
- **Multiple Topologies**: Support for Default (Ring), Chain, HubAndSpoke, SharedWeb, and Random topologies
- **Zone Management**: Player and neutral zones with quality tiers (Low, Medium, High)
- **Content Scaling**: Automatic content scaling based on map size and zone count
- **Game Rules**: Support for various win conditions (Classic Victory, City Hold, Gladiator Arena, Tournaments)
- **CLI Interface**: Command-line tool for generating templates
- **Unit Tests**: Comprehensive test coverage

## Building the Project

```bash
# Build the CLI tool
go build -o template-generator ./cmd

# Run tests
go test ./...

# Run with verbose output
go test -v ./...
```

## Usage

### Command-Line Tool

```bash
# Generate a basic 4-player template
./template-generator -name="My Template" -players=4 -size=L -topology=Default -output=.

# Generate with custom settings
./template-generator \
  -name="Custom Map" \
  -players=6 \
  -size=XL \
  -topology=HubAndSpoke \
  -mode=Blitz \
  -roads=true \
  -portals=true \
  -footholds=true \
  -cityhold=false \
  -output=./templates
```

### Programmatic Usage

```go
package main

import (
    "github.com/Tariomka/hommoe_custom_templates/internal/generator"
    "github.com/Tariomka/hommoe_custom_templates/internal/models"
)

settings := &models.GeneratorSettings{
    TemplateName: "My Template",
    PlayerCount:  4,
    MapSize:      "L",
    Topology:     models.TopologyDefault,
}

template, err := generator.Generate(settings)
```

## Map Topologies

- **Default (Ring)**: Players arranged in a circle with neighbors connected
- **Chain**: Linear arrangement of zones
- **HubAndSpoke**: All player zones connect through a central hub
- **SharedWeb**: Players connected through neutral zones in a spoke pattern
- **Random**: Random zone positions computed with Delaunay triangulation

## Supported Map Sizes

- `S` - Small
- `M` - Medium
- `L` - Large
- `XL` - Extra Large
- `2XL` - Double Extra Large

## Game Modes

- `Classic` - Standard mode
- `Blitz` - Fast-paced mode
- `Heroic` - Hard mode

## Win Conditions

- **Classic Victory**: Capture all towns
- **City Hold**: Hold a specific town for a number of days
- **Gladiator Arena**: Gladiator arena matches
- **Tournaments**: Tournament challenges

## Example Templates

57 reference templates are included in `data/GameData/ExampleTemplates/`. These serve as design references and examples for understanding the `.rmg.json` format.

## Generator Data

The generator uses configuration files in `data/GameData/GeneratorData/`:

- `generator_config.json` - Meta objects and item definitions
- `generator_environment_assets.json` - Biome and environment settings
- `generator_stats_config.json` - Hero stat modifiers
- `content_lists/` - Named groupings of items
- `content_pools/` - Quality tier pools (t2, t3, t4, t5)
- `encounter_templates/` - Encounter definitions
- `zone_layouts/` - Zone structure templates

## Testing

Run the full test suite:

```bash
go test ./...
```

Run tests for a specific package:

```bash
go test ./internal/generator
go test ./internal/services
```

Run specific tests:

```bash
go test ./internal/generator -run TestGenerateUsesRequestedSettings
```

## Architecture

### Core Components

1. **Models** (`internal/models/`)
   - `RmgTemplate` - Top-level template structure
   - `Variant` - Map variant with zones and connections
   - `Zone` - Individual zone with content and guards
   - `GeneratorSettings` - User configuration

2. **Generator** (`internal/generator/`)
   - `Generate()` - Main generation entry point
   - `buildTopologyAdjacency()` - Zone connectivity logic
   - `buildGameRules()` - Win condition construction
   - `buildZones()` - Zone creation
   - `ComputeContentScale()` - Dynamic content scaling

3. **Services** (`internal/services/`)
   - `ZoneContentManager` - Mandatory content by zone type
   - `ContentItemBuilder` - Fluent builder for content items

### Generation Flow

```
GeneratorSettings
    ↓
Generate()
├─→ buildNeutralZonePlan()      # Plan neutral zones
├─→ buildTopologyAdjacency()    # Map zone connections
├─→ buildGameRules()            # Create win conditions
└─→ buildVariant()              # Create zones and connections
    ↓
RmgTemplate (JSON serializable)
```

## Notes

- Zones are lettered A-AF (supporting up to 32 zones)
- Guard randomization defaults to 0.05
- Content scaling uses a gentle sqrt curve, clamped to 0.5-2.5x
- Default connection count per zone is 2
- Player zones come first (A onwards), neutral zones follow

## Related Documentation

- See `AGENTS.md` for development guidelines
- See the C# project (Olden Era - Template Editor) for the original implementation

## License

See the main project repository for license information.
