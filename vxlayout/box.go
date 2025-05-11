package vxlayout

import (
	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"github.com/avidal/vxexp"
)

// Constrained is a [vxfw.Widget] that constrains a widget by min and max size.
// Min and max are pointers to [vxfw.Size] so that the caller can indicate a lack of constraint.
// If either axis of size is 0, that axis is ignored for constraint purposes.
// Constrained can be useful for laying out a Row where you want to ensure a maximum height.
func Constrained(widget vxfw.Widget, minSize, maxSize *vxfw.Size) vxfw.Widget {
	return vxexp.WidgetFunc(func(ctx vxfw.DrawContext) (vxfw.Surface, error) {
		if minSize != nil {
			if minSize.Width > ctx.Min.Width {
				ctx.Min.Width = minSize.Width
			}
			if minSize.Height > ctx.Min.Height {
				ctx.Min.Height = minSize.Height
			}
		}
		if maxSize != nil {
			if maxSize.Width != 0 && maxSize.Width < ctx.Max.Width {
				ctx.Max.Width = maxSize.Width
			}
			if maxSize.Height != 0 && maxSize.Height < ctx.Max.Height {
				ctx.Max.Height = maxSize.Height
			}
		}

		return widget.Draw(ctx)
	})
}

// Sized is a [vxfw.Widget] that passes a fixed size to its child widget as long as size fits the
// incoming constraints. If either axis of size is 0, that size is ignored.
// This can be used to place a widget that normally cannot be unconstrained (such as an infinite
// list) into a flexible layout. See also [Limited].
// This is a shortcut for Constrained(widget, size, size)
func Sized(widget vxfw.Widget, size vxfw.Size) vxfw.Widget {
	return Constrained(widget, &size, &size)
}

// Limited is a [vxfw.Widget] that limits its child by size only if the incoming constraint is
// unlimited. If either axis of size is 0, that axis is ignored.
func Limited(widget vxfw.Widget, size vxfw.Size) vxfw.Widget {
	return vxexp.WidgetFunc(func(ctx vxfw.DrawContext) (vxfw.Surface, error) {
		if ctx.Max.HasUnboundedWidth() && size.Width > 0 {
			ctx.Max.Width = size.Width
		}
		if ctx.Max.HasUnboundedHeight() && size.Height > 0 {
			ctx.Max.Height = size.Height
		}

		return widget.Draw(ctx)
	})
}

// Fill returns a [vxfw.Widget] that fills its space with the supplied cell.
// Note that Fill will take all available space. It's primarily useful to diagnose layouts, and
// will usually be contained in a [Flexible]
func Fill(cell vaxis.Cell) vxfw.Widget {
	return vxexp.WidgetFunc(func(ctx vxfw.DrawContext) (vxfw.Surface, error) {
		surface := vxfw.NewSurface(ctx.Max.Width, ctx.Max.Height, nil)
		surface.Fill(cell)
		return surface, nil
	})
}
