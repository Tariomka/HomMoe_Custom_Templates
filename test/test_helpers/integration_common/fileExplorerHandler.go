//go:build integration_test

package integration_common

// FileExplorerHandler drives the load/save file dialog. It holds the originating
// BaseHandler but deliberately does not embed it: an open dialog's scrim absorbs
// every background click, so promoting the tab clicks would be a lie. It also
// takes no snapshot, because the dialog lists the per-machine templates
// directory (AGENTS.md 2.7).
type FileExplorerHandler struct {
	base *BaseHandler
}

func (this *FileExplorerHandler) IsOpen() bool {
	this.base.runner.tb.Helper()
	return this.base.runner.DialogsOpen()
}

func (this *FileExplorerHandler) Close() *BaseHandler {
	this.base.runner.tb.Helper()
	this.base.runner.CloseTopDialog()
	return this.base
}
