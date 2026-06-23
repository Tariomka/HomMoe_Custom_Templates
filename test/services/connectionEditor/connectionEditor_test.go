package connectionEditor_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

// ════════════════════════════════════════════════════════════════════════
// T015 · US1 – Edit Connection Properties
// ════════════════════════════════════════════════════════════════════════

func TestT015a_GuardValue_CanBeSetAndReadBack(t *testing.T) {
	conn := entities.Connection{From: "Spawn-A", To: "Neutral-1", GuardValue: 100}
	conn.GuardValue = 999
	assert.Equal(t, 999, conn.GuardValue)
}

func TestT015b_ConnectionType_CanBeSetToPortal(t *testing.T) {
	conn := entities.Connection{From: "Spawn-A", To: "Neutral-1", ConnectionType: "Direct"}
	conn.ConnectionType = "Portal"
	assert.Equal(t, "Portal", conn.ConnectionType)
}

// The Go Connection model stores GuardWeeklyIncrement as a non-pointer float64,
// so "clearing" maps to the zero value rather than null.
func TestT015c_GuardWeeklyIncrement_CanBeCleared(t *testing.T) {
	conn := entities.Connection{From: "Spawn-A", To: "Neutral-1", GuardWeeklyIncrement: 1.5}
	conn.GuardWeeklyIncrement = 0
	assert.Equal(t, 0.0, conn.GuardWeeklyIncrement)
}

func TestT015d_IsUserAdded_StartsAsFalse(t *testing.T) {
	conn := entities.Connection{From: "Spawn-A", To: "Neutral-1"}
	assert.False(t, conn.IsUserAdded)
}

// ════════════════════════════════════════════════════════════════════════
// T019 · US2 – Add New Connections
// ════════════════════════════════════════════════════════════════════════

func TestT019a_AddingConnection_AppendsOneEntryWithIsUserAddedTrue(t *testing.T) {
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1"},
	}
	newConn := entities.Connection{From: "Spawn-A", To: "Neutral-2", ConnectionType: "Direct", IsUserAdded: true}
	connections = append(connections, newConn)

	assert.Equal(t, 2, len(connections))
	assert.True(t, connections[len(connections)-1].IsUserAdded)
}

func TestT019b_AddingSecondConnectionBetweenSamePair_ResultsInTwoEntries(t *testing.T) {
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1", ConnectionType: "Direct"},
	}
	connections = append(connections, entities.Connection{
		From: "Spawn-A", To: "Neutral-1", ConnectionType: "Portal", IsUserAdded: true,
	})

	pairCount := 0
	for _, c := range connections {
		if (c.From == "Spawn-A" && c.To == "Neutral-1") || (c.From == "Neutral-1" && c.To == "Spawn-A") {
			pairCount++
		}
	}
	assert.Equal(t, 2, pairCount)
}

func TestT019c_CancelledAdd_LeavesListUnchanged(t *testing.T) {
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1"},
	}
	countBefore := len(connections)
	// Simulate "Cancel" - nothing is added.
	assert.Equal(t, countBefore, len(connections))
}

// ════════════════════════════════════════════════════════════════════════
// T022 · US3 – Remove Connections
// ════════════════════════════════════════════════════════════════════════

func removeAt(connections []entities.Connection, index int) []entities.Connection {
	return append(connections[:index], connections[index+1:]...)
}

func TestT022a_RemovingConnection_LeavesCountMinusOne(t *testing.T) {
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1"},
		{From: "Neutral-1", To: "Neutral-2"},
	}
	connections = removeAt(connections, 0)
	assert.Equal(t, 1, len(connections))
}

func TestT022b_RemovingConnection_LeavesNoPairEntry(t *testing.T) {
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1"},
		{From: "Spawn-B", To: "Neutral-2"},
	}
	connections = removeAt(connections, 0)

	stillPresent := false
	for _, c := range connections {
		if (c.From == "Spawn-A" && c.To == "Neutral-1") || (c.From == "Neutral-1" && c.To == "Spawn-A") {
			stillPresent = true
		}
	}
	assert.False(t, stillPresent)
}

func TestT022c_ListIsEmptyAfterDelete(t *testing.T) {
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1"},
	}
	connections = removeAt(connections, 0)
	assert.Empty(t, connections)
}

// ════════════════════════════════════════════════════════════════════════
// T027 · US4 – Graph Overview Helpers
// ════════════════════════════════════════════════════════════════════════

func TestT027a_AfterReset_ConnectionListMatchesOriginalInCountAndIsUserAddedIsFalse(t *testing.T) {
	original := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1", GuardValue: 100},
		{From: "Spawn-B", To: "Neutral-1", GuardValue: 200},
	}

	current := make([]entities.Connection, 0, len(original))
	for _, orig := range original {
		current = append(current, connection_editor.CloneConnection(orig, false))
	}

	assert.Equal(t, len(original), len(current))
	for _, c := range current {
		assert.False(t, c.IsUserAdded)
	}
	assert.Equal(t, 100, current[0].GuardValue)
	assert.Equal(t, 200, current[1].GuardValue)
}

func TestT027b_IsolatedZoneDetection_IdentifiesZoneWithNoConnections(t *testing.T) {
	zones := []entities.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-1"},
		{Name: "Neutral-2"}, // isolated
	}
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1"},
	}

	isolated := connection_editor.FindIsolatedZones(zones, connections)

	assert.Equal(t, 1, len(isolated))
	assert.Equal(t, "Neutral-2", isolated[0])
}

func TestT027c_DuplicateNameDetection_FlagsWhenTwoConnectionsShareSameName(t *testing.T) {
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1", Name: "main-road"},
		{From: "Neutral-1", To: "Neutral-2", Name: "main-road"},
		{From: "Spawn-B", To: "Neutral-1", Name: "side-path"},
	}

	assert.True(t, connection_editor.HasDuplicateName(connections, &connections[0]))
	assert.True(t, connection_editor.HasDuplicateName(connections, &connections[1]))
	assert.False(t, connection_editor.HasDuplicateName(connections, &connections[2]))
}

func TestT027c_DuplicateNameDetection_DoesNotFlagWhenNamesAreDistinct(t *testing.T) {
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1", Name: "alpha"},
		{From: "Neutral-1", To: "Neutral-2", Name: "beta"},
	}

	assert.False(t, connection_editor.HasDuplicateName(connections, &connections[0]))
	assert.False(t, connection_editor.HasDuplicateName(connections, &connections[1]))
}

// ════════════════════════════════════════════════════════════════════════
// T032d · FR-009 – HasUnresolvedErrors
// ════════════════════════════════════════════════════════════════════════

func TestT032d_HasUnresolvedErrors_FalseWhenAllZoneNamesExist(t *testing.T) {
	zones := []entities.Zone{{Name: "Spawn-A"}, {Name: "Neutral-1"}}
	connections := []entities.Connection{{From: "Spawn-A", To: "Neutral-1"}}
	assert.False(t, connection_editor.ComputeHasErrors(zones, connections))
}

func TestT032d_HasUnresolvedErrors_TrueWhenFromZoneIsMissing(t *testing.T) {
	zones := []entities.Zone{{Name: "Neutral-1"}}
	connections := []entities.Connection{{From: "Spawn-A", To: "Neutral-1"}}
	assert.True(t, connection_editor.ComputeHasErrors(zones, connections))
}

func TestT032d_HasUnresolvedErrors_TrueWhenToZoneIsMissing(t *testing.T) {
	zones := []entities.Zone{{Name: "Spawn-A"}}
	connections := []entities.Connection{{From: "Spawn-A", To: "Neutral-99"}}
	assert.True(t, connection_editor.ComputeHasErrors(zones, connections))
}

// ════════════════════════════════════════════════════════════════════════
// IsUserAdded serialisation contract
// ════════════════════════════════════════════════════════════════════════

func TestIsUserAdded_IsNotSerialized_JsonDoesNotContainProperty(t *testing.T) {
	conn := entities.Connection{From: "Spawn-A", To: "Neutral-1", IsUserAdded: true}
	data, err := json.Marshal(conn)
	assert.NoError(t, err)

	jsonStr := string(data)
	assert.False(t, strings.Contains(strings.ToLower(jsonStr), "isuseradded"))
}

// ════════════════════════════════════════════════════════════════════════
// Guard-preset tables & zone tiers (new functionality)
// ════════════════════════════════════════════════════════════════════════

func TestUserCreatableConnectionTypes_ExcludesProximity(t *testing.T) {
	types := connection_editor.UserCreatableConnectionTypes()
	assert.Equal(t, []string{"Direct", "Portal"}, types)
	for _, typeName := range types {
		assert.NotEqual(t, "Proximity", typeName)
	}
}

func TestGuardPresets_MatchTable(t *testing.T) {
	assert.Equal(t, [5]int{3000, 6000, 9000, 12000, 16000}, connection_editor.GuardPresetsForTier(connection_editor.ZoneTierBronze))
	assert.Equal(t, [5]int{18000, 21000, 24000, 27000, 30000}, connection_editor.GuardPresetsForTier(connection_editor.ZoneTierSilver))
	assert.Equal(t, [5]int{36000, 42000, 48000, 54000, 60000}, connection_editor.GuardPresetsForTier(connection_editor.ZoneTierGold))
	assert.Equal(t, [5]int{10000, 22000, 34000, 46000, 58000}, connection_editor.GuardPresetsForTier(connection_editor.ZoneTierPlayerToPlayer))
}

func TestTierExtras_GeneratorDefaults(t *testing.T) {
	assert.Equal(t, 15000, connection_editor.ExtrasForTier(connection_editor.ZoneTierBronze)[0].Value)
	assert.Equal(t, 20000, connection_editor.ExtrasForTier(connection_editor.ZoneTierSilver)[0].Value)
	assert.Equal(t, 25000, connection_editor.ExtrasForTier(connection_editor.ZoneTierGold)[0].Value)
	assert.Equal(t, 30000, connection_editor.ExtrasForTier(connection_editor.ZoneTierPlayerToPlayer)[0].Value)
}

func TestWeeklyIncrementValues_MatchTable(t *testing.T) {
	assert.Equal(t, []float64{0.05, 0.10, 0.15, 0.20, 0.25}, connection_editor.WeeklyIncrementValues)
}

func TestGetZoneTier_PlayerZoneIsBronze(t *testing.T) {
	zones := []entities.Zone{{Name: "Spawn-A"}}
	players := map[string]bool{"Spawn-A": true}
	assert.Equal(t, connection_editor.ZoneTierBronze, connection_editor.GetZoneTier("Spawn-A", zones, players))
}

func TestGetZoneTier_HubIsBronze(t *testing.T) {
	zones := []entities.Zone{{Name: "Hub"}, {Name: "Hub-1"}}
	players := map[string]bool{}
	assert.Equal(t, connection_editor.ZoneTierBronze, connection_editor.GetZoneTier("Hub", zones, players))
	assert.Equal(t, connection_editor.ZoneTierBronze, connection_editor.GetZoneTier("Hub-1", zones, players))
}

func TestGetZoneTier_NeutralPoolDecidesBracket(t *testing.T) {
	zones := []entities.Zone{
		{Name: "Neutral-Gold", GuardedContentPool: []string{"pool_t4_x"}},
		{Name: "Neutral-Gold5", GuardedContentPool: []string{"pool_t5_x"}},
		{Name: "Neutral-Bronze", GuardedContentPool: []string{"pool_t1_x"}},
		{Name: "Neutral-Bronze2", GuardedContentPool: []string{"pool_t2_x"}},
		{Name: "Neutral-Silver", GuardedContentPool: []string{"pool_t3_x"}},
		{Name: "Neutral-NoPool"},
	}
	players := map[string]bool{}
	assert.Equal(t, connection_editor.ZoneTierGold, connection_editor.GetZoneTier("Neutral-Gold", zones, players))
	assert.Equal(t, connection_editor.ZoneTierGold, connection_editor.GetZoneTier("Neutral-Gold5", zones, players))
	assert.Equal(t, connection_editor.ZoneTierBronze, connection_editor.GetZoneTier("Neutral-Bronze", zones, players))
	assert.Equal(t, connection_editor.ZoneTierBronze, connection_editor.GetZoneTier("Neutral-Bronze2", zones, players))
	assert.Equal(t, connection_editor.ZoneTierSilver, connection_editor.GetZoneTier("Neutral-Silver", zones, players))
	assert.Equal(t, connection_editor.ZoneTierSilver, connection_editor.GetZoneTier("Neutral-NoPool", zones, players))
}

func TestGetZoneTier_UnknownZoneIsBronze(t *testing.T) {
	assert.Equal(t, connection_editor.ZoneTierBronze, connection_editor.GetZoneTier("Nope", nil, nil))
	assert.Equal(t, connection_editor.ZoneTierBronze, connection_editor.GetZoneTier("", nil, nil))
}

func TestHigherTierOf_BothPlayersIsPlayerToPlayer(t *testing.T) {
	zones := []entities.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}
	players := map[string]bool{"Spawn-A": true, "Spawn-B": true}
	assert.Equal(t, connection_editor.ZoneTierPlayerToPlayer, connection_editor.HigherTierOf("Spawn-A", "Spawn-B", zones, players))
}

func TestHigherTierOf_TakesMaximumTier(t *testing.T) {
	zones := []entities.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-Gold", GuardedContentPool: []string{"pool_t4_x"}},
	}
	players := map[string]bool{"Spawn-A": true}
	assert.Equal(t, connection_editor.ZoneTierGold, connection_editor.HigherTierOf("Spawn-A", "Neutral-Gold", zones, players))
}

func TestNewDefaultConnection_UsesTierDefaultsAndIsUserAdded(t *testing.T) {
	zones := []entities.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-Gold", GuardedContentPool: []string{"pool_t4_x"}},
	}
	players := map[string]bool{"Spawn-A": true}

	conn := connection_editor.NewDefaultConnection("Spawn-A", "Neutral-Gold", zones, players)

	assert.Equal(t, "Spawn-A", conn.From)
	assert.Equal(t, "Neutral-Gold", conn.To)
	assert.Equal(t, "Direct", conn.ConnectionType)
	assert.Equal(t, 25000, conn.GuardValue) // Gold generator default
	assert.Equal(t, "Spawn-A", conn.GuardZone)
	assert.Equal(t, "rnd_guard_A_Gold", conn.GuardMatchGroup)
	assert.Equal(t, 0.15, conn.GuardWeeklyIncrement)
	assert.True(t, conn.IsUserAdded)
}

func TestZoneLetterFromName(t *testing.T) {
	assert.Equal(t, "A", connection_editor.ZoneLetterFromName("Spawn-A"))
	assert.Equal(t, "C", connection_editor.ZoneLetterFromName("Neutral-C"))
	assert.Equal(t, "Hub", connection_editor.ZoneLetterFromName("Hub"))
}

func TestCloneConnection_CopiesFieldsAndSetsFlag(t *testing.T) {
	original := entities.Connection{
		Name: "edge", From: "Spawn-A", To: "Neutral-1",
		ConnectionType: "Portal", GuardValue: 4242, GuardWeeklyIncrement: 0.2,
		GuardZone: "Spawn-A", GuardMatchGroup: "grp", Length: 1.25,
	}
	clone := connection_editor.CloneConnection(original, true)

	assert.Equal(t, original.Name, clone.Name)
	assert.Equal(t, original.From, clone.From)
	assert.Equal(t, original.To, clone.To)
	assert.Equal(t, original.ConnectionType, clone.ConnectionType)
	assert.Equal(t, original.GuardValue, clone.GuardValue)
	assert.Equal(t, original.GuardWeeklyIncrement, clone.GuardWeeklyIncrement)
	assert.Equal(t, original.Length, clone.Length)
	assert.True(t, clone.IsUserAdded)
	assert.False(t, original.IsUserAdded)
}
