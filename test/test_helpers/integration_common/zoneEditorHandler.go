//go:build integration_test

package integration_common

// ZoneEditorHandler drives the manual zone editor dialog. It holds the
// originating BaseHandler but deliberately does not embed it: an open dialog's
// scrim absorbs every background click, so promoting the tab clicks would be a
// lie. It also takes no snapshot, because the dialog draws the freshly
// randomised Random-topology layout, which no golden can hold.
type ZoneEditorHandler struct {
	base *BaseHandler
}

func (this *ZoneEditorHandler) IsOpen() bool {
	this.base.runner.tb.Helper()
	return this.base.runner.DialogsOpen()
}

func (this *ZoneEditorHandler) Close() *BaseHandler {
	this.base.runner.tb.Helper()
	this.base.runner.CloseTopDialog()
	return this.base
}
