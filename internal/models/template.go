package models

// RmgTemplate represents the top-level template structure for .rmg.json files
type RmgTemplate struct {
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	TemplateFilename string         `json:"templateFilename"`
	Size             string         `json:"size"` // e.g., "L", "XL", "2XL"
	GameRules        GameRules      `json:"gameRules"`
	Variants         []Variant      `json:"variants"`
	MandatoryContent []ContentGroup `json:"mandatoryContent,omitempty"`
}

// Variant represents a single map variant
type Variant struct {
	Name        string       `json:"name"`
	Zones       []Zone       `json:"zones"`
	Connections []Connection `json:"connections"`
}

// Zone represents a single zone in the map
type Zone struct {
	Name              string          `json:"name"`
	Type              string          `json:"type"` // "player" or "neutral"
	Letter            string          `json:"letter"`
	Owner             int             `json:"owner,omitempty"`
	DefenseValue      int             `json:"defenseValue,omitempty"`
	Layout            ZoneLayout      `json:"layout"`
	GuardSettings     GuardSettings   `json:"guardSettings"`
	BiomeSelectors    []BiomeSelector `json:"biomeSelectors,omitempty"`
	ContentPools      ContentPools    `json:"contentPools"`
	MandatoryContents []string        `json:"mandatoryContents,omitempty"`
	MainObjects       []MainObject    `json:"mainObjects,omitempty"`
	Roads             []string        `json:"roads,omitempty"`
	WaterSlots        []WaterSlot     `json:"waterSlots,omitempty"`
	ConnectionRules   []string        `json:"connectionRules,omitempty"`
}

// ZoneLayout defines the structure of a zone
type ZoneLayout struct {
	Type      string `json:"type"`
	Template  string `json:"template,omitempty"`
	TownType  string `json:"townType,omitempty"`
	OffsetX   int    `json:"offsetX,omitempty"`
	OffsetY   int    `json:"offsetY,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	BlockSize int    `json:"blockSize,omitempty"`
}

// GuardSettings defines guard parameters for a zone
type GuardSettings struct {
	Randomization                    float64 `json:"randomization,omitempty"`
	ReactionDistribution             string  `json:"reactionDistribution,omitempty"`
	Multiplier                       float64 `json:"multiplier,omitempty"`
	ReactionAggression               float64 `json:"reactionAggression,omitempty"`
	ReactionAggravating              float64 `json:"reactionAggravating,omitempty"`
	ReactionDefensive                float64 `json:"reactionDefensive,omitempty"`
	ReactionWary                     float64 `json:"reactionWary,omitempty"`
	CreereMutabilityFactor           float64 `json:"creereMutabilityFactor,omitempty"`
	LevelForMultiplierWithoutMonster float64 `json:"levelForMultiplierWithoutMonster,omitempty"`
}

// BiomeSelector defines biome selection for a zone
type BiomeSelector struct {
	Name    string             `json:"name"`
	Weights map[string]float64 `json:"weights,omitempty"`
}

// ContentPools groups content items by type
type ContentPools struct {
	Guarded   []ContentItem `json:"guarded,omitempty"`
	Unguarded []ContentItem `json:"unguarded,omitempty"`
	Resources []ContentItem `json:"resources,omitempty"`
	Plants    []ContentItem `json:"plants,omitempty"`
}

// ContentItem represents a single content item with placement rules
type ContentItem struct {
	SID                   string                 `json:"sid"`
	Count                 int                    `json:"count,omitempty"`
	GuardType             string                 `json:"guardType,omitempty"`
	Guarded               bool                   `json:"guarded,omitempty"`
	Mine                  bool                   `json:"mine,omitempty"`
	SoloEncounter         bool                   `json:"soloEncounter,omitempty"`
	PlacementRules        []PlacementRule        `json:"placementRules,omitempty"`
	RoadDistanceVariation interface{}            `json:"roadDistanceVariation,omitempty"`
	OffsetX               int                    `json:"offsetX,omitempty"`
	OffsetY               int                    `json:"offsetY,omitempty"`
	AdditionalData        map[string]interface{} `json:"additionalData,omitempty"`
}

// PlacementRule defines where content can be placed
type PlacementRule struct {
	RuleType  string `json:"ruleType"`
	Value     string `json:"value,omitempty"`
	Distance  int    `json:"distance,omitempty"`
	Direction string `json:"direction,omitempty"`
}

// Connection links two zones together
type Connection struct {
	FromZone             string                 `json:"from"`
	ToZone               string                 `json:"to"`
	GuardZone            string                 `json:"guardZone,omitempty"`
	PortalPlacement      PortalPlacement        `json:"portalPlacement"`
	GuardEscapeRouting   string                 `json:"guardEscapeRouting,omitempty"`
	RoadRoute            string                 `json:"roadRoute,omitempty"`
	HasRoad              bool                   `json:"hasRoad,omitempty"`
	GuardIncrement       int                    `json:"guardIncrement,omitempty"`
	GuardIncrementFactor float64                `json:"guardIncrementFactor,omitempty"`
	AdditionalData       map[string]interface{} `json:"additionalData,omitempty"`
}

// PortalPlacement defines portal placement rules
type PortalPlacement struct {
	From PortalEndpoint `json:"from"`
	To   PortalEndpoint `json:"to"`
}

// PortalEndpoint defines where a portal appears
type PortalEndpoint struct {
	ZoneName string   `json:"zoneName"`
	Rules    []string `json:"rules,omitempty"`
}

// MainObject represents a main structure (castle, foothold, etc.)
type MainObject struct {
	Type     string `json:"type"`
	Location string `json:"location,omitempty"`
	Owner    int    `json:"owner,omitempty"`
	OffsetX  int    `json:"offsetX,omitempty"`
	OffsetY  int    `json:"offsetY,omitempty"`
}

// WaterSlot represents water placement
type WaterSlot struct {
	X      int `json:"x,omitempty"`
	Y      int `json:"y,omitempty"`
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}

// GameRules defines game rules and win conditions
type GameRules struct {
	HeroCount        int                    `json:"heroCount,omitempty"`
	AllowedHeroes    []string               `json:"allowedHeroes,omitempty"`
	AllowedTowns     []string               `json:"allowedTowns,omitempty"`
	WinConditions    []WinCondition         `json:"winConditions,omitempty"`
	LossConditions   []LossCondition        `json:"lossConditions,omitempty"`
	Bonuses          []Bonus                `json:"bonuses,omitempty"`
	ExpenseModifiers ExpenseModifiers       `json:"expenseModifiers,omitempty"`
	AdditionalData   map[string]interface{} `json:"additionalData,omitempty"`
}

// WinCondition defines a win condition
type WinCondition struct {
	Type       string                 `json:"type"`
	Condition  string                 `json:"condition,omitempty"`
	Target     interface{}            `json:"target,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// LossCondition defines a loss condition
type LossCondition struct {
	Type       string                 `json:"type"`
	Condition  string                 `json:"condition,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// Bonus defines a game bonus
type Bonus struct {
	Type      string      `json:"type"`
	Value     interface{} `json:"value,omitempty"`
	Parameter string      `json:"parameter,omitempty"`
	For       []string    `json:"for,omitempty"`
}

// ExpenseModifiers defines faction expense modifiers
type ExpenseModifiers struct {
	Castle     float64 `json:"Castle,omitempty"`
	Rampart    float64 `json:"Rampart,omitempty"`
	Tower      float64 `json:"Tower,omitempty"`
	Inferno    float64 `json:"Inferno,omitempty"`
	Necropolis float64 `json:"Necropolis,omitempty"`
	Dungeon    float64 `json:"Dungeon,omitempty"`
	Stronghold float64 `json:"Stronghold,omitempty"`
	Fortress   float64 `json:"Fortress,omitempty"`
	Conflux    float64 `json:"Conflux,omitempty"`
}

// ContentGroup represents mandatory content items grouped by theme
type ContentGroup struct {
	Name  string        `json:"name"`
	Items []ContentItem `json:"items,omitempty"`
}
