# C# to Go Migration Summary

This document describes how the C# WPF project (`Olden Era - Template Editor`) was converted to a Go project (`Template Generator`).

## Project Scope

Both projects serve the same purpose: generating `.rmg.json` random map template files for Heroes of Might and Magic: Olden Era.

### C# Project (Original)
- **Type**: WPF Desktop Application
- **UI**: XAML-based GUI for template configuration
- **Main Class**: `TemplateGenerator.cs` (Services/TemplateGenerator.cs)
- **Models**: Unfrozen/ (JSON-serializable structures)
- **Framework**: .NET/C#

### Go Project (New)
- **Type**: Command-line Tool + Library
- **UI**: Command-line interface (flags-based)
- **Main Module**: internal/generator/template_generator.go
- **Models**: internal/models/ (struct-based)
- **Framework**: Go (golang)

## Architecture Comparison

### C# Structure → Go Structure

| Aspect | C# | Go |
|--------|----|----|
| Models | `Models/Unfrozen/` classes | `internal/models/template.go` structs |
| Settings | `Models/Generator/GeneratorSettings` | `internal/models/settings.go` structs |
| Generation | `Services/TemplateGenerator.cs` static class | `internal/generator/template_generator.go` functions |
| Content Management | `Services/ContentManagement/` | `internal/services/` |
| UI | WPF XAML (MainWindow.xaml) | CLI flags (`cmd/main.go`) |
| Tests | xUnit (`*.Tests.csproj`) | Go testing (`*_test.go`) |
| Data Files | `GameData/` | `data/GameData/` |

## Key Conversions

### 1. Data Models: C# Classes → Go Structs

**C# Example:**
```csharp
public class RmgTemplate
{
    public string Name { get; set; }
    public string Description { get; set; }
    public string Size { get; set; }
    public GameRules GameRules { get; set; }
    public List<Variant> Variants { get; set; }
}
```

**Go Equivalent:**
```go
type RmgTemplate struct {
    Name             string         `json:"name"`
    Description      string         `json:"description"`
    Size             string         `json:"size"`
    GameRules        GameRules      `json:"gameRules"`
    Variants         []Variant      `json:"variants"`
    MandatoryContent []ContentGroup `json:"mandatoryContent,omitempty"`
}
```

### 2. Generator Logic: Static Class → Package Functions

**C# Example:**
```csharp
public static class TemplateGenerator
{
    public static RmgTemplate Generate(GeneratorSettings settings)
    {
        // ...generation logic
    }

    private static void BuildVariant(GeneratorSettings settings)
    {
        // ...
    }
}
```

**Go Equivalent:**
```go
package generator

func Generate(settings *models.GeneratorSettings) (*models.RmgTemplate, error) {
    // ...generation logic
}

func buildVariant(settings *models.GeneratorSettings, ...) *models.Variant {
    // ...
}
```

### 3. Enumerations → Constants + Types

**C# Example:**
```csharp
public enum MapTopology
{
    Default,
    Chain,
    HubAndSpoke,
    SharedWeb,
    Random
}
```

**Go Equivalent:**
```go
type MapTopology string

const (
    TopologyDefault     MapTopology = "Default"
    TopologyChain       MapTopology = "Chain"
    TopologyHubAndSpoke MapTopology = "HubAndSpoke"
    TopologySharedWeb   MapTopology = "SharedWeb"
    TopologyRandom      MapTopology = "Random"
)
```

### 4. Builder Pattern

**C# Example:**
```csharp
var item = new ContentItemBuilder()
    .WithSid("my_sid")
    .Guarded()
    .Build();
```

**Go Equivalent:**
```go
item := NewContentItem("my_sid").
    Guarded().
    Build()
```

### 5. UI → CLI

**C#**: Multiple XAML windows with UI controls for template configuration

**Go**: Command-line flags
```bash
./template-generator \
  -name="My Template" \
  -players=4 \
  -size=L \
  -topology=Default \
  -output=./templates
```

## File Structure Comparison

### C#
```
Olden Era - Template Editor/
├── Models/
│   ├── Unfrozen/
│   │   ├── RmgTemplate.cs
│   │   ├── Variant.cs
│   │   └── Zone.cs
│   └── Generator/
│       ├── GeneratorSettings.cs
│       └── MapTopology.cs
├── Services/
│   ├── TemplateGenerator.cs
│   └── ContentManagement/
│       ├── ZoneContentManager.cs
│       └── ContentItemBuilder.cs
├── MainWindow.xaml
├── MainWindow.xaml.cs
└── GameData/
    ├── ExampleTemplates/
    └── GeneratorData/
```

### Go
```
Template Generator/
├── internal/
│   ├── models/
│   │   ├── template.go
│   │   └── settings.go
│   ├── generator/
│   │   ├── template_generator.go
│   │   └── template_generator_test.go
│   └── services/
│       ├── zone_content_manager.go
│       └── zone_content_manager_test.go
├── cmd/
│   └── main.go
├── data/
│   └── GameData/
│       ├── ExampleTemplates/
│       └── GeneratorData/
├── go.mod
└── README.md
```

## Core Logic Preservation

The following logic was preserved exactly from C# to Go:

1. **Topology Adjacency Building**
   - Default (Ring): Players in circle, neighbors connected
   - Chain: Linear arrangement
   - HubAndSpoke: All connect to central hub
   - SharedWeb: Players connected through neutral zones
   - Random: Basic random connections

2. **Zone Planning**
   - Neutral zone quality tiers (Low, Medium, High)
   - Zone lettering (A-AF for up to 32 zones)
   - Player and neutral zone creation

3. **Game Rules**
   - Win conditions construction
   - Hero count initialization
   - Default victory condition

4. **Content Scaling**
   - Size-based multipliers
   - Zone count adjustment
   - Clamping to bounds (0.5-2.5x)

## Testing

### C# Tests → Go Tests

| C# Test (xUnit) | Go Test |
|-----------------|---------|
| `TestGenerateUsesRequestedSettings` | `TestGenerateUsesRequestedSettings` |
| `TestGenerateDefaultTopologyCreatesExpectedZones` | `TestGenerateDefaultTopologyCreatesRingConnections` |
| `TestComputeContentScale` | `TestComputeContentScale` |
| `TestGenerateWithRoadsDisabledLeavesZones` | `TestGenerateWithRoadsDisabled` |

**Test Execution:**
```bash
# Go
go test ./...

# C#
dotnet test "Olden Era - Template Editor.slnx"
```

## Building & Running

### C# (Original)
```powershell
# Build
dotnet build "Olden Era - Template Editor.slnx"

# Run (WPF GUI)
dotnet run --project "Olden Era - Template Editor/Olden Era - Template Editor.csproj"

# Test
dotnet test "Olden Era - Template Editor.slnx"
```

### Go (New)
```bash
# Build
go build -o template-generator ./cmd

# Run CLI
./template-generator -name="Map" -players=4 -size=L

# Test
go test ./...
```

## Dependencies

### C# Project
- .NET Framework
- xUnit (testing)
- System.Text.Json (JSON serialization)

### Go Project
- Standard library only (no external dependencies)
- Uses `encoding/json` for JSON marshaling

## Data Files Reused

Both projects use identical data files:
- `data/GameData/ExampleTemplates/` - 57 reference templates
- `data/GameData/GeneratorData/` - Configuration and content pools
- No changes needed to game data

## JSON Output

The `.rmg.json` output format is identical between projects. Generated templates can be used interchangeably.

**Sample Output:**
```json
{
  "name": "My Template",
  "description": "...",
  "size": "L",
  "gameRules": {...},
  "variants": [{...}]
}
```

## Performance Considerations

- **Go**: Compiled binary, faster startup, no runtime dependencies
- **C#**: Requires .NET runtime, heavier (WPF dependencies)
- **CLI**: No UI overhead, pure generation logic

## Future Enhancements

Potential additions to Go project:
1. REST API server for template generation
2. Configuration file support (YAML/TOML)
3. Batch template generation
4. Template validation against game rules
5. Template preview/visualization

## Notes

- Go project focuses on generation logic, removing UI complexity
- All core generation algorithms remain unchanged
- Test coverage maintained from original project
- Example templates and game data are identical references
- Error handling improved with Go's error return style
- JSON marshaling automatic through struct tags
