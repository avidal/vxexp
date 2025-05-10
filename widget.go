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
