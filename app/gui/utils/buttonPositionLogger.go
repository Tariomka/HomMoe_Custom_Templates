package utils

import (
	"context"
	"log/slog"

	"gioui.org/io/input"
	"gioui.org/io/semantic"
	"gioui.org/op"
)

// ButtonPositionLogger replays frame operations into a debug-only input router
// and logs the absolute window coordinates of every labeled button.
type ButtonPositionLogger struct {
	router input.Router
	logger *slog.Logger
}

// NewButtonPositionLogger returns a ButtonPositionLogger that writes to the given logger.
func NewButtonPositionLogger(logger *slog.Logger) *ButtonPositionLogger {
	return &ButtonPositionLogger{logger: logger}
}

// LogButtonPositions replays the frame operations and logs the label, bounds and
// center point (in absolute window coordinates) of every labeled button widget.
// It is a no-op when debug logging is disabled.
func (this *ButtonPositionLogger) LogButtonPositions(operations *op.Ops) {
	if !this.logger.Enabled(context.Background(), slog.LevelDebug) {
		return
	}

	this.logger.Debug("====== New Frame ======")
	this.router.Frame(operations)
	for _, node := range this.router.AppendSemantics(nil) {
		if node.Desc.Class != semantic.Button || node.Desc.Label == "" {
			continue
		}

		bounds := node.Desc.Bounds
		center := bounds.Min.Add(bounds.Max).Div(2)
		this.logger.Debug("Button position",
			slog.String("label", node.Desc.Label),
			slog.String("center", center.String()),
			slog.String("bounds", bounds.String()))
	}
}
