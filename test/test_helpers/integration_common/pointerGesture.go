//go:build integration_test

package integration_common

import "gioui.org/f32"

// PointerGesture is one press-and-release a test queues with other gestures
// before a single frame. Secondary gestures are right mouse clicks, the rest
// are primary touch taps.
type PointerGesture struct {
	Position  f32.Point
	Secondary bool
}
