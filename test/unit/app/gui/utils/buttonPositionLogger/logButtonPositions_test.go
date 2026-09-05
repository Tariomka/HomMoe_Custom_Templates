package buttonPositionLogger_test

import (
	"context"
	"fmt"
	"image"
	"log/slog"
	"math"
	"testing"

	"gioui.org/io/event"
	"gioui.org/io/semantic"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	levelDebug = slog.LevelDebug
	levelInfo  = slog.LevelInfo
)

func TestWhenOpsContainOffsetButton_LogsAbsoluteCenter(t *testing.T) {
	t.Parallel()
	// Arrange
	offset := image.Pt(gofakeit.Number(1, 500), gofakeit.Number(1, 500))
	size := image.Pt(gofakeit.Number(10, 200), gofakeit.Number(10, 200))
	expectedCenter := offset.Add(offset.Add(size)).Div(2)
	handler := newRecordingHandler(levelDebug)
	logger := utils.NewButtonPositionLogger(newSlogLogger(handler))
	operations := new(op.Ops)
	appendButtonOps(operations, new(int), offset, size, gofakeit.Word())

	// Act
	logger.LogButtonPositions(operations, identityMetric(), testFrameSize())

	// Assert
	buttons := buttonRecords(handler)
	require.Len(t, buttons, 1)
	assert.Equal(t, expectedCenter.String(), attrValue(buttons[0], "center"))
}

func TestWhenOpsContainLabeledButton_LogsLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	label := gofakeit.Word()
	handler := newRecordingHandler(levelDebug)
	logger := utils.NewButtonPositionLogger(newSlogLogger(handler))
	operations := new(op.Ops)
	appendButtonOps(operations, new(int), image.Pt(10, 10), image.Pt(50, 20), label)

	// Act
	logger.LogButtonPositions(operations, identityMetric(), testFrameSize())

	// Assert
	buttons := buttonRecords(handler)
	require.Len(t, buttons, 1)
	assert.Equal(t, label, attrValue(buttons[0], "label"))
}

func TestWhenOpsContainMultipleButtons_LogsEachButton(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newRecordingHandler(levelDebug)
	logger := utils.NewButtonPositionLogger(newSlogLogger(handler))
	operations := new(op.Ops)
	appendButtonOps(operations, new(int), image.Pt(10, 10), image.Pt(50, 20), gofakeit.Word())
	appendButtonOps(operations, new(int), image.Pt(200, 300), image.Pt(80, 30), gofakeit.Word())

	// Act
	logger.LogButtonPositions(operations, identityMetric(), testFrameSize())

	// Assert
	assert.Len(t, buttonRecords(handler), 2)
}

func TestWhenDebugLoggingIsDisabled_LogsNothing(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newRecordingHandler(levelInfo)
	logger := utils.NewButtonPositionLogger(newSlogLogger(handler))
	operations := new(op.Ops)
	appendButtonOps(operations, new(int), image.Pt(10, 10), image.Pt(50, 20), gofakeit.Word())

	// Act
	logger.LogButtonPositions(operations, identityMetric(), testFrameSize())

	// Assert
	assert.Empty(t, handler.records)
}

func TestWhenOpsContainNoButtons_LogsNothing(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newRecordingHandler(levelDebug)
	logger := utils.NewButtonPositionLogger(newSlogLogger(handler))
	operations := new(op.Ops)
	area := clip.Rect(image.Rectangle{Max: image.Pt(50, 20)}).Push(operations)
	area.Pop()

	// Act
	logger.LogButtonPositions(operations, identityMetric(), testFrameSize())

	// Assert
	assert.Empty(t, buttonRecords(handler))
}

func TestWhenLoggingFrame_EmitsSingleNewFrameRecord(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newRecordingHandler(levelDebug)
	logger := utils.NewButtonPositionLogger(newSlogLogger(handler))
	operations := new(op.Ops)

	// Act
	logger.LogButtonPositions(operations, identityMetric(), testFrameSize())

	// Assert
	frameMarkers := make([]slog.Record, 0, len(handler.records))
	for _, record := range handler.records {
		if record.Message == "====== New Frame ======" {
			frameMarkers = append(frameMarkers, record)
		}
	}
	assert.Len(t, frameMarkers, 1)
}

func TestWhenButtonHasNoLabel_SkipsButton(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newRecordingHandler(levelDebug)
	logger := utils.NewButtonPositionLogger(newSlogLogger(handler))
	operations := new(op.Ops)
	appendButtonOps(operations, new(int), image.Pt(10, 10), image.Pt(50, 20), "")

	// Act
	logger.LogButtonPositions(operations, identityMetric(), testFrameSize())

	// Assert
	assert.Empty(t, buttonRecords(handler))
}

func TestWhenMetricScalesPixels_LogsCenterInDp(t *testing.T) {
	t.Parallel()
	// Arrange
	const pxPerDp = 1.75
	offsetDp := image.Pt(gofakeit.Number(1, 400), gofakeit.Number(1, 400))
	sizeDp := image.Pt(gofakeit.Number(8, 200), gofakeit.Number(8, 200))
	offsetPx := image.Pt(int(float32(offsetDp.X)*pxPerDp), int(float32(offsetDp.Y)*pxPerDp))
	sizePx := image.Pt(int(float32(sizeDp.X)*pxPerDp), int(float32(sizeDp.Y)*pxPerDp))
	expectedCenter := offsetDp.Add(offsetDp.Add(sizeDp)).Div(2)
	handler := newRecordingHandler(levelDebug)
	logger := utils.NewButtonPositionLogger(newSlogLogger(handler))
	operations := new(op.Ops)
	appendButtonOps(operations, new(int), offsetPx, sizePx, gofakeit.Word())

	// Act
	logger.LogButtonPositions(operations, unit.Metric{PxPerDp: pxPerDp, PxPerSp: pxPerDp}, testFrameSize())

	// Assert
	buttons := buttonRecords(handler)
	require.Len(t, buttons, 1)
	assert.InDelta(t, 0, centerDistance(t, attrValue(buttons[0], "center"), expectedCenter), 1)
}

func TestWhenMetricScalesPixels_LogsFrameSizeInDp(t *testing.T) {
	t.Parallel()
	// Arrange
	const pxPerDp = 2
	frameSizeDp := image.Pt(gofakeit.Number(100, 2000), gofakeit.Number(100, 2000))
	frameSizePx := frameSizeDp.Mul(int(pxPerDp))
	handler := newRecordingHandler(levelDebug)
	logger := utils.NewButtonPositionLogger(newSlogLogger(handler))
	operations := new(op.Ops)

	// Act
	logger.LogButtonPositions(operations, unit.Metric{PxPerDp: pxPerDp, PxPerSp: pxPerDp}, frameSizePx)

	// Assert
	require.Len(t, handler.records, 1)
	assert.Equal(t, frameSizeDp.String(), attrValue(handler.records[0], "frameSizeDp"))
}

// recordingHandler is a [slog.Handler] that captures every record it receives.
type recordingHandler struct {
	level   slog.Level
	records []slog.Record
}

func newRecordingHandler(level slog.Level) *recordingHandler {
	return &recordingHandler{level: level}
}

func (this *recordingHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= this.level
}

func (this *recordingHandler) Handle(_ context.Context, record slog.Record) error {
	this.records = append(this.records, record)
	return nil
}

func (this *recordingHandler) WithAttrs(_ []slog.Attr) slog.Handler { return this }

func (this *recordingHandler) WithGroup(_ string) slog.Handler { return this }

func newSlogLogger(handler slog.Handler) *slog.Logger {
	return slog.New(handler)
}

// identityMetric maps 1 dp to 1 px, so logged dp coordinates equal raw pixels.
func identityMetric() unit.Metric {
	return unit.Metric{PxPerDp: 1, PxPerSp: 1}
}

// testFrameSize is an arbitrary frame size for tests that do not assert on it.
func testFrameSize() image.Point {
	return image.Pt(1600, 900)
}

// appendButtonOps records the ops a labeled clickable button produces:
// an offset transform, a clipped input area, and a nested handler-free
// area carrying the button semantics (mirrors widgets.addButtonSemantics).
func appendButtonOps(operations *op.Ops, tag event.Tag, offset image.Point, size image.Point, label string) {
	transform := op.Offset(offset).Push(operations)
	area := clip.Rect(image.Rectangle{Max: size}).Push(operations)
	event.Op(operations, tag)
	semanticArea := clip.Rect(image.Rectangle{Max: size}).Push(operations)
	semantic.Button.Add(operations)
	semantic.LabelOp(label).Add(operations)
	semanticArea.Pop()
	area.Pop()
	transform.Pop()
}

func attrValue(record slog.Record, key string) string {
	value := ""
	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key == key {
			value = attr.Value.String()
			return false
		}
		return true
	})
	return value
}

// buttonRecords filters the captured records down to per-button log entries,
// dropping the intended once-per-call "====== New Frame ======" marker.
func buttonRecords(handler *recordingHandler) []slog.Record {
	filtered := make([]slog.Record, 0, len(handler.records))
	for _, record := range handler.records {
		if record.Message == "Button position" {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

// centerDistance parses a logged "(x,y)" center attribute and returns its
// Chebyshev distance from the expected point, so scaled coordinates can be
// asserted with a rounding tolerance.
func centerDistance(t *testing.T, logged string, expected image.Point) float64 {
	t.Helper()
	var x, y int
	_, err := fmt.Sscanf(logged, "(%d,%d)", &x, &y)
	require.NoError(t, err)
	return math.Max(math.Abs(float64(x-expected.X)), math.Abs(float64(y-expected.Y)))
}
