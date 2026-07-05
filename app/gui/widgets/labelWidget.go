package widgets

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

type LabelBuilder struct {
	item     material.LabelStyle
	textSize unit.Sp
}

func NewLabelBuilder(theme *material.Theme) *LabelBuilder {
	return &LabelBuilder{
		item:     material.Body1(theme, ""),
		textSize: theme.TextSize,
	}
}

func (this *LabelBuilder) WithText(text string) *LabelBuilder {
	this.item.Text = text
	return this
}

func (this *LabelBuilder) WithSize(size uint) *LabelBuilder {
	this.item.TextSize = this.textSize * unit.Sp(size) / 16
	return this
}
func (this *LabelBuilder) WithSizeDefault() *LabelBuilder { return this.WithSize(12) }
func (this *LabelBuilder) WithSizeSmall() *LabelBuilder   { return this.WithSize(10) }
func (this *LabelBuilder) WithSizeBig() *LabelBuilder     { return this.WithSize(14) }

func (this *LabelBuilder) WithColor(color color.NRGBA) *LabelBuilder {
	this.item.Color = color
	return this
}
func (this *LabelBuilder) WithColorDefault() *LabelBuilder { return this.WithColor(themes.ColorText) }
func (this *LabelBuilder) WithColorDim() *LabelBuilder     { return this.WithColor(themes.ColorTextDim) }
func (this *LabelBuilder) WithColorError() *LabelBuilder   { return this.WithColor(themes.ColorError) }

func (this *LabelBuilder) WithFont(font font.Font) *LabelBuilder {
	this.item.Font = font
	return this
}

func (this *LabelBuilder) WithAlignment(alignment text.Alignment) *LabelBuilder {
	this.item.Alignment = alignment
	return this
}

func (this *LabelBuilder) WithMaxLines(lineCount int) *LabelBuilder {
	this.item.MaxLines = lineCount
	this.item.Truncator = "…"
	return this
}

func (this *LabelBuilder) Build(gtx layout.Context) layout.Dimensions { return this.item.Layout(gtx) }

// NewLabelWidget returns a Widget that renders a small text line.
func NewLabelWidget(theme *material.Theme, text string, color color.NRGBA) layout.Widget {
	return NewLabelBuilder(theme).WithSizeDefault().WithText(text).WithColor(color).Build
}

// NewLabelBigWidget returns a Widget like [NewLabelWidget] but one size bigger.
func NewLabelBigWidget(theme *material.Theme, text string, color color.NRGBA) layout.Widget {
	return NewLabelBuilder(theme).WithSizeBig().WithText(text).WithColor(color).Build
}

// NewStyledLabelWidget returns a Widget that renders a small text line with a specific style.
func NewStyledLabelWidget(theme *material.Theme, text string, color color.NRGBA, style font.Font) layout.Widget {
	return NewLabelBuilder(theme).WithSizeDefault().WithText(text).WithColor(color).WithFont(style).Build
}

// NewDimmedLabelWidget returns a Widget that renders a small dimmed text line.
func NewDimmedLabelWidget(theme *material.Theme, text string) layout.Widget {
	return NewLabelBuilder(theme).WithSizeDefault().WithText(text).WithColorDim().Build
}
