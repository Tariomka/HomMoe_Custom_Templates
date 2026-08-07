//go:build integration_test

package integration_common

import (
	"image"

	"gioui.org/f32"
)

type baseHandler struct {
	runner *AppRunner

	isRandomTopology bool
}

func NewHandler(runner *AppRunner) *baseHandler {
	runner.tb.Helper()
	handler := baseHandler{runner: runner}
	handler.setRandomTopology()
	return &handler
}

func (this *baseHandler) WithSnapshots() *baseHandler {
	this.runner.tb.Helper()
	this.runner.EnableSnapshots()
	return this
}

func (this *baseHandler) ClickGeneralTab() *baseHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(672, 60))
	this.runner.VerifySnapshot()
	return this
}

func (this *baseHandler) ClickLayoutAndZonesTab() *baseHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(789, 60))
	this.runner.VerifySnapshot()
	return this
}

func (this *baseHandler) ClickBonusesAndBansTab() *baseHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(933, 60))
	this.runner.VerifySnapshot()
	return this
}

func (this *baseHandler) setRandomTopology() *baseHandler {
	this.runner.tb.Helper()
	this.isRandomTopology = true
	previewMask := image.Rect(
		WindowWidth-470, 0,
		WindowWidth, WindowHeight)
	this.runner.MaskRect(previewMask)
	return this
}

// Need to add handlers for the content of each tab and dialog.
// This will be effectively most volatile test code but it is required to ensure the GUI functionality
// and to work equivalently to something like Playwright or Selenium, which are not available for Gio.
// The tests will be brittle but they will be the only way to ensure the GUI works as intended.
//
// Each handler will(for example switching between tabs still makes the tab buttons and toolbar buttons functional)
// or will not embed(for example dialogs disable all of the base tabs and other buttons so there is no point in checking the background differences)
// the base the baseHandler and each Click(action) method will return either self, or the next transition handler,
// depending on the action. This will allow for a fluent interface to be used in the tests.
//
// Each handler will require separate masks most likely (dialog mask will require the entire background to be masked, leaving only the dialog itself).
// Currently, the only mask is the mask that entirely covers the preview panel, because on startup,
// random topology is selected and generated, making it always differ, as well as blocking
// the output directory field because it differs between environments
// (local vs CI, local machines with or without installed Heroes Olden Era, operating systems, etc.)
// and some parts of the message text (like the timestamp and/or the file path of a successfully generated file)
// which will differ between environments and runs. This mask should be updated so that only "unstable"/"unpredictable"
// parts are covered, making the snapshot comparison more exact.
//
// The idea of multiple handlers for different tabs is also volatile - each handler probably will need to track the state,
// for example in General tab there are sliders, checkboxes and dropdown menus: for sliders its fine to calculate the position in place,
// checkboxes in certain scenarios change how the ui is layed out, and dropdowns aren't floating, they are rendered on the
// canvas like everything else, shifting the layout if not selected/reset. Something like clicking on
// "Allow non official larger map sizes" expands the Map Size dropdown.
//
// I envision that the NewHandler() will either provide baseHandler (and you will need to click a tab to get the GeneralTabHandler)
// or the GeneralTabHandler directly, a test that would change the map size,
// change the topology and then save the editor config would look something like this:
//
//	NewHandler(runner). // returns either baseHandler or GeneralTabHandler, depending on the implementation
//		WithSnapshots(). // enables snapshot validation on each action, returns self
//		ClickGeneralTab(). // returns GeneralTabHandler (either baseHandler produces GeneralTabHandler or returns self if its already GeneralTabHandler)
//		SelectAllowNonOfficialMapSizes(). // clicks the checkbox, returns self
//		SetMapSize(12). // if indexing from 1, there are only 11 base map sizes, so it requires the AllowNonOfficialMapSizes checkbox to be selected, otherwise it acts as passing 0 which would be selecting the dropdown widget again, returns self
//		ClickLayoutAndZonesTab(). // returns LayoutAndZonesTabHandler
//		SetTopology(2). // sets topology from Random to Ring, once again 0 is the base widget, everything after is the dropdown list value, if its out of bounds, then it defaults to 0, returns self
//		ClickSave() // clicks save button, and either clicks save on the save widget dialog if no file was saved to complete the save action and returns self, or just returns self if no dialog appears
//
// This looks like an entire framework, so I don't know if it's better to calculate positions in place or hardcode them, or a combination of both.
// There will also be an issue with scrolling, if a button is not visible.
//
// Also for some reason there is a difference of some of the rendered text between local (both this pc and steam produce identical results)
// and CI (looks like some of the text is grayed out like not finishing a rerender), so need to figure out how to fix that as well.
