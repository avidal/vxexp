package vxexp

import (
	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
)

// WidgetFunc is convenience adapter to turn an ordinary function into a [vxfw.Widget]
// WidgetFunc(fn) is a [vxfw.Widget] that calls fn during the Draw phase.
type WidgetFunc func(vxfw.DrawContext) (vxfw.Surface, error)

func (fn WidgetFunc) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	return fn(ctx)
}

// EventHandlerFunc is a convenience function to make stateless widgets.
// EventHandlerFunc(fn) is a [vxfw.EventHandler] that calls fn during the event handling phase.
type EventHandlerFunc func(vaxis.Event, vxfw.EventPhase) (vxfw.Command, error)

func (fn EventHandlerFunc) HandleEvent(event vaxis.Event, phase vxfw.EventPhase) (vxfw.Command, error) {
	return fn(event, phase)
}

// Draw implements [vxfw.Widget]
func (_ EventHandlerFunc) Draw(_ vxfw.DrawContext) (vxfw.Surface, error) { return vxfw.Surface{}, nil }

func ClampUint16(x, min, max uint16) uint16 {
	switch {
	case x == 0:
		return min
	case x < min:
		return min
	case x > max:
		return max
	default:
		return x
	}
}

// TightenContext returns a [vxfw.DrawContext] whose min and max along each dimension is as close
// to size as possible while still staying in-bounds. If either dimension of size is 0, it's
// ignored.
func TightenContext(ctx vxfw.DrawContext, size vxfw.Size) vxfw.DrawContext {
	out := vxfw.Size{
		Width:  ClampUint16(size.Width, ctx.Min.Width, ctx.Max.Width),
		Height: ClampUint16(size.Height, ctx.Min.Height, ctx.Max.Height),
	}
	return ctx.WithConstraints(out, out)
}

// LoosenContext takes a [vxfw.DrawContext] and returns a context with no minimum sizes.
func LoosenContext(ctx vxfw.DrawContext) vxfw.DrawContext {
	return ctx.WithMin(vxfw.Size{})
}
