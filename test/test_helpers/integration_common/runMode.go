//go:build integration_test

package integration_common

import (
	"flag"
	"os"
	"testing"

	"gioui.org/app"
)

// updateSnapshots rewrites golden snapshots instead of validating against them.
// Registered at package init so testing's [flag.Parse] picks it up; pass it after
// -args: `go test -tags=integration_test ./test/integration/... -args -update`.
//
//nolint:gochecknoglobals // command-line flags are the idiomatic exception.
var updateSnapshots = flag.Bool("update", false, "rewrite golden snapshots instead of validating against them")

// IsUpdate reports whether the -update flag was passed to the test binary.
func IsUpdate() bool {
	return *updateSnapshots
}

// IsHeadless reports whether the "headed" argument was passed to the test
// binary (e.g. `go test ... -args headed`). [os.Args] is scanned directly rather
// than using a registered flag because TestMain must decide whether to start
// app.Main before the testing framework parses flags.
func IsHeadless() bool {
	for _, arg := range os.Args[1:] {
		if arg == "headed" || arg == "-headed" || arg == "--headed" {
			return false
		}
	}

	return true
}

// RunMain is the shared TestMain body for the gated suites. In windowed mode it
// runs the Gio event system on the main goroutine (required by app.Window) while
// the tests run on a separate goroutine. In headless mode app.Main is never
// started, so the package needs no display and is safe to run anywhere.
func RunMain(m *testing.M) {
	if IsHeadless() {
		os.Exit(m.Run())
	}
	go func() {
		os.Exit(m.Run())
	}()
	app.Main()
}
