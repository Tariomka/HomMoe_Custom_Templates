package dialogHost_test

import (
	"testing"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTheTopModalIsClosed_TheOneBeneathItStaysOpen(t *testing.T) {
	t.Parallel()
	// Arrange - a dialog may open another on top of itself, so closing the child
	// has to resume the parent rather than clear the stack.
	host := &drivers.DialogHost{}
	host.Open(stubDialog{})
	host.Open(stubDialog{})
	require.True(t, host.IsOpen())

	// Act
	host.Close()

	// Assert
	assert.True(t, host.IsOpen())
}

func TestWhenTheLastModalIsClosed_NothingIsOpen(t *testing.T) {
	t.Parallel()
	// Arrange
	host := &drivers.DialogHost{}
	host.Open(stubDialog{})

	// Act
	host.Close()

	// Assert
	assert.False(t, host.IsOpen())
}

func TestWhenNothingIsOpen_ClosingIsANoOp(t *testing.T) {
	t.Parallel()
	// Arrange
	host := &drivers.DialogHost{}

	// Act
	host.Close()

	// Assert
	assert.False(t, host.IsOpen())
}

// stubDialog is a stack occupant; DialogHost.Close never draws it.
type stubDialog struct{}

func (this stubDialog) Title() string { return "stub" }

func (this stubDialog) Body(gtx layout.Context, _ *material.Theme) (layout.Dimensions, bool) {
	return layout.Dimensions{Size: gtx.Constraints.Min}, false
}

func (this stubDialog) PreferredSize() (width, height unit.Dp) { return 0, 0 }
