package zoneEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneHasThreeCastles_RebuildsTheCastleRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := template_model.Zone{
		Name:        "Neutral-Z",
		MainObjects: []template_model.MainObject{{Type: "City"}, {Type: "City"}, {Type: "City"}},
	}

	// Act
	test_helpers.NewZoneEditorService().RebuildCastleRoads(&zone)

	// Assert
	assert.Equal(t, []string{"1", "2"}, roadTargetArgs(zone, "MainObject"))
}

func TestWhenZoneHasStaleCastleRoads_DropsTheStaleOnes(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := template_model.Zone{
		Name:        "Neutral-Z",
		MainObjects: []template_model.MainObject{{Type: "City"}},
		Roads: []template_model.Road{
			{
				From: mainObjectZeroRef(),
				To:   template_model.TypedRef{Type: "MainObject", Args: []string{"1"}},
			},
			{
				From: mainObjectZeroRef(),
				To:   template_model.TypedRef{Type: "MainObject", Args: []string{"2"}},
			},
		},
	}

	// Act
	test_helpers.NewZoneEditorService().RebuildCastleRoads(&zone)

	// Assert
	assert.Empty(t, roadTargetArgs(zone, "MainObject"))
}

func TestWhenZoneHasAConnectionRoad_PreservesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := template_model.Zone{
		Name:        "Neutral-Z",
		MainObjects: []template_model.MainObject{{Type: "City"}, {Type: "City"}},
		Roads: []template_model.Road{
			{From: mainObjectZeroRef(), To: template_model.TypedRef{Type: "Connection", Args: []string{"Rnd-A-Z"}}},
		},
	}

	// Act
	test_helpers.NewZoneEditorService().RebuildCastleRoads(&zone)

	// Assert
	assert.True(t, roadTargets(zone, "Connection")["Rnd-A-Z"])
}

// roadTargetArgs returns the first argument of every road targeting the given
// reference type, in road order.
func roadTargetArgs(zone template_model.Zone, referenceType string) []string {
	var args []string
	for _, road := range zone.Roads {
		if road.To.Type == referenceType && len(road.To.Args) > 0 {
			args = append(args, road.To.Args[0])
		}
	}
	return args
}
