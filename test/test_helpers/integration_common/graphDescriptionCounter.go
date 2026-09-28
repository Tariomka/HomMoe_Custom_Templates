//go:build integration_test

package integration_common

import (
	"sync"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// GraphDescriptionCounter is the real GUI handler with its zone-editor graph
// diagnostics counted, so a test can tell which frames rebuilt them. The
// counts are locked because a headed run lays the window out on its render
// goroutine too.
type GraphDescriptionCounter struct {
	handler_interfaces.IGuiHandler

	mu    sync.Mutex
	calls int
	last  dtos.ZoneEditorGraphDto
}

func NewGraphDescriptionCounter(handler handler_interfaces.IGuiHandler) *GraphDescriptionCounter {
	return &GraphDescriptionCounter{IGuiHandler: handler}
}

func (this *GraphDescriptionCounter) DescribeZoneEditorGraph(
	zones []template_model.Zone,
	connections []template_model.Connection) dtos.ZoneEditorGraphDto {
	graph := this.IGuiHandler.DescribeZoneEditorGraph(zones, connections)
	this.mu.Lock()
	defer this.mu.Unlock()
	this.calls++
	this.last = graph

	return graph
}

// Calls reports how many times the graph has been described so far.
func (this *GraphDescriptionCounter) Calls() int {
	this.mu.Lock()
	defer this.mu.Unlock()

	return this.calls
}

// Last reports the most recent description, which is what the status line shows.
func (this *GraphDescriptionCounter) Last() dtos.ZoneEditorGraphDto {
	this.mu.Lock()
	defer this.mu.Unlock()

	return this.last
}
