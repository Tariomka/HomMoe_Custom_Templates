package handler_interfaces

type IGuiHandler interface {
	ITemplateHandler
	IStateHandler
	IPreviewHandler
	IZoneContentHandler
	IZoneEditorHandler
	IBonusHandler
	IPickerHandler
}
