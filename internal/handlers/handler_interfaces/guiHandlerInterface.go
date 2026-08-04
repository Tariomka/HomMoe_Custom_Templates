package handler_interfaces

type IGuiHandler interface {
	ITemplateHandler
	IStateHandler
	IPreviewHandler
	IContentRuleHandler
	IZoneEditorHandler
}
