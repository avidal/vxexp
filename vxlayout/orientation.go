package vxlayout

import (
	"git.sr.ht/~rockorager/vaxis/vxfw"
)

type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

// MainAxis returns the MainAxis of size
func (o Orientation) MainAxis(size vxfw.Size) uint16 {
	main, _ := o.Axes(size)
	return main
}

// CrossAxis returns the cross axis of size
func (o Orientation) CrossAxis(size vxfw.Size) uint16 {
	_, cross := o.Axes(size)
	return cross
}

// Axes returns the main and cross axes of size.
func (o Orientation) Axes(size vxfw.Size) (uint16, uint16) {
	var main, cross uint16
	switch o {
	case Horizontal:
		main = size.Width
		cross = size.Height
	case Vertical:
		main = size.Height
		cross = size.Width
	}
	return main, cross
}

// CrossMax takes a size and a constraint and returns the max of size or constraint on the cross
// axis.
func (o Orientation) CrossMax(size vxfw.Size, constraint uint16) uint16 {
	if o == Horizontal && size.Height > constraint {
		return size.Height
	} else if o == Vertical && size.Width > constraint {
		return size.Width
	}

	return constraint
}

// Size takes main and cross axis and returns a [vxfw.Size] based on orientation
func (o Orientation) Size(main, cross uint16) vxfw.Size {
	var out vxfw.Size
	switch o {
	case Horizontal:
		out.Width, out.Height = main, cross
	case Vertical:
		out.Height, out.Width = main, cross
	}
	return out
}

// Origin returns a [vxfw.RelativePoint] with main and cross offsets.
func (o Orientation) Origin(main, cross int) vxfw.RelativePoint {
	var p vxfw.RelativePoint
	switch o {
	case Horizontal:
		p.Col = main
		p.Row = cross
	case Vertical:
		p.Row = main
		p.Col = cross
	}
	return p
}

// OffsetOrigin applies main and cross axis offsets to origin and returns a new
// [vxfw.RelativePoint]
func (o Orientation) OffsetOrigin(origin vxfw.RelativePoint, main, cross int) vxfw.RelativePoint {
	var p vxfw.RelativePoint
	p.Col, p.Row = origin.Col, origin.Row
	switch o {
	case Horizontal:
		p.Col += main
		p.Row += cross
	case Vertical:
		p.Row += main
		p.Col += cross
	}
	return p
}

// Loosen returns a [vxfw.DrawContext] with the main axis set to 0
func (o Orientation) Loosen(ctx vxfw.DrawContext) vxfw.DrawContext {
	out := vxfw.DrawContext(ctx)
	switch o {
	case Horizontal:
		out.Min.Width = 0
	case Vertical:
		out.Min.Height = 0
	}
	return out
}

// Tighten returns a [vxfw.DrawContext] with the main axis minimum increased to equal the maximum.
func (o Orientation) Tighten(ctx vxfw.DrawContext) vxfw.DrawContext {
	out := vxfw.DrawContext(ctx)
	switch o {
	case Horizontal:
		out.Min.Width = out.Max.Width
	case Vertical:
		out.Min.Height = out.Max.Height
	}
	return out
}
