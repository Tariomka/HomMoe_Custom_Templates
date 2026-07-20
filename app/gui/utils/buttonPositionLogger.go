package utils

import (
	"context"
	"image"
	"log/slog"
	"math"

	"gioui.org/io/input"
	"gioui.org/io/semantic"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
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
// center point of every labeled button widget. The raw semantics bounds are in
// physical pixels (already scaled by the OS display scale), so they are converted
// back to device-independent pixels using the frame metric. The integration test
// harness lays the editor out at 1600x900 with PxPerDp=1, so the logged dp
// coordinates are directly usable as synthetic click points whenever the logged
// frame size is also 1600x900 dp (the default window size).
// It is a no-op when debug logging is disabled.
func (this *ButtonPositionLogger) LogButtonPositions(operations *op.Ops, metric unit.Metric, frameSize image.Point) {
	if !this.logger.Enabled(context.Background(), slog.LevelDebug) {
		return
	}

	this.logger.Debug("====== New Frame ======",
		slog.String("frameSizeDp", pixelsToDp(metric, frameSize).String()))
	this.router.Frame(operations)
	for _, node := range this.router.AppendSemantics(nil) {
		if node.Desc.Class != semantic.Button || node.Desc.Label == "" {
			continue
		}

		bounds := image.Rectangle{
			Min: pixelsToDp(metric, node.Desc.Bounds.Min),
			Max: pixelsToDp(metric, node.Desc.Bounds.Max),
		}
		center := bounds.Min.Add(bounds.Max).Div(2)
		this.logger.Debug("Button position",
			slog.String("label", node.Desc.Label),
			slog.String("center", center.String()),
			slog.String("bounds", bounds.String()))
	}
}

// pixelsToDp converts a point from physical pixels to device-independent pixels,
// rounding to the nearest whole dp.
func pixelsToDp(metric unit.Metric, point image.Point) image.Point {
	return image.Pt(
		int(math.Round(float64(point.X)/float64(metric.PxPerDp))),
		int(math.Round(float64(point.Y)/float64(metric.PxPerDp))))
}

// AddButtonSemantics records the button class and label in a nested, handler-free
// clip area sized to the button. The input router keeps semantics of such areas
// intact (areas with an input handler lose them unless a gesture filter is
// registered), which lets utils.ButtonPositionLogger resolve every button's
// absolute window bounds from a replayed frame.
func AddButtonSemantics(operations *op.Ops, label string, size image.Point) {
	area := clip.Rect(image.Rectangle{Max: size}).Push(operations)
	semantic.Button.Add(operations)
	semantic.LabelOp(label).Add(operations)
	area.Pop()
}
