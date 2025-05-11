package vxlayout

import (
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"github.com/avidal/vxexp"
)

type flexible interface {
	vxfw.Widget

	// FlexFactor determines how much of the available space is proportioned to this widget
	FlexFactor() uint16
	// A loosely flexible widget may take less than its proportioned flex space.
	// [Flexible] is a widget that provides a loose fit to its child.
	FlexLoose() bool
}

// Expanded returns a [vxfw.Widget] that will expand in a flexible layout based on flex.
// Note that a flex of 0 means the widget will not be flexible at all and is typically a sign
// of a bug in your layout.
func Expanded(widget vxfw.Widget, flex uint16) vxfw.Widget {
	return flexbox{Widget: widget, flex: flex}
}

// Flex returns a [vxfw.Widget] that will loosely flex in a flexible layout based on flex.
// Note that a flex of 0 means the widget will not be flexible at all and is typically a sign
// of a bug in your layout.
func Flexible(widget vxfw.Widget, flex uint16) vxfw.Widget {
	return flexbox{Widget: widget, flex: flex, loose: true}
}

// Space returns a [vxfw.Widget] that will fill all available space in a flexible layout.
// Note that a Space(0) has no flex factor and will take ALL space that comes after it in the
// layout. If flex is > 0 (ie, Space(1)), Space will share remaining space proportionally with
// other flexible widgets.
func Space(flex uint16) vxfw.Widget {
	fn := vxexp.WidgetFunc(func(ctx vxfw.DrawContext) (vxfw.Surface, error) {
		return vxfw.NewSurface(ctx.Max.Width, ctx.Max.Height, nil), nil
	})

	return flexbox{Widget: fn, flex: flex}
}

type flexbox struct {
	vxfw.Widget
	flex  uint16
	loose bool
}

var (
	_ vxfw.Widget = flexbox{}
	_ flexible    = flexbox{}
)

func (b flexbox) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	return b.Widget.Draw(ctx)
}

func (b flexbox) FlexLoose() bool    { return b.loose }
func (b flexbox) FlexFactor() uint16 { return b.flex }
