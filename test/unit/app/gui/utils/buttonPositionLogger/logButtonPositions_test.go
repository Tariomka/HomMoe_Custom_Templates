package buttonPositionLogger_test

import (
	"context"
	"image"
	"log/slog"
	"testing"

	"gioui.org/io/event"
	"gioui.org/io/semantic"
	"gioui.org/op"
	"gioui.org/op/clip"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	levelDebug = slog.LevelDebug
	levelInfo  = slog.LevelInfo
)

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
	logger.LogButtonPositions(operations)

	// Assert
	require.Len(t, handler.records, 1)
	assert.Equal(t, expectedCenter.String(), attrValue(handler.records[0], "center"))
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
	logger.LogButtonPositions(operations)

	// Assert
	require.Len(t, handler.records, 1)
	assert.Equal(t, label, attrValue(handler.records[0], "label"))
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
	logger.LogButtonPositions(operations)

	// Assert
	assert.Len(t, handler.records, 2)
}

func TestWhenDebugLoggingIsDisabled_LogsNothing(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newRecordingHandler(levelInfo)
	logger := utils.NewButtonPositionLogger(newSlogLogger(handler))
	operations := new(op.Ops)
	appendButtonOps(operations, new(int), image.Pt(10, 10), image.Pt(50, 20), gofakeit.Word())

	// Act
	logger.LogButtonPositions(operations)

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
	logger.LogButtonPositions(operations)

	// Assert
	assert.Empty(t, handler.records)
}

func TestWhenButtonHasNoLabel_SkipsButton(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newRecordingHandler(levelDebug)
	logger := utils.NewButtonPositionLogger(newSlogLogger(handler))
	operations := new(op.Ops)
	appendButtonOps(operations, new(int), image.Pt(10, 10), image.Pt(50, 20), "")

	// Act
	logger.LogButtonPositions(operations)

	// Assert
	assert.Empty(t, handler.records)
}
