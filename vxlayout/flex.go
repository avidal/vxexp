package vxlayout

import (
	"math"

	"git.sr.ht/~rockorager/vaxis/vxfw"
)

// Determines how children in a [Row] or [Column] are laid out on the cross axis.
// The default is [CrossAxisCenter], which centers children in the available cross axis space.
// Use [CrossAxisStart] to align children to the top of a [Row] or left of a [Column], and vice
// versa for [CrossAxisEnd].
// [CrossAxisStretch] will force all children to fill the maximum space on the cross axis.
type CrossAxisAlignment int

const (
	CrossAxisCenter CrossAxisAlignment = iota
	CrossAxisStart
	CrossAxisEnd
	CrossAxisStretch
)

// Determines how children in a [Row] or [Column] are laid out on the main axis.
// The default is [MainAxisStart], which places children at the left or top of a [Row] or [Column]
// respectively.
// Use [MainAxisCenter] to center children, or one of the Space variants to distribute the space
// elsewhere.
// Note that these options do nothing if one of the children has a flex factor > 0, as those
// children will take the available space.
type MainAxisAlignment int

const (
	MainAxisStart MainAxisAlignment = iota
	MainAxisEnd
	MainAxisCenter
	MainAxisSpaceBetween
	MainAxisSpaceAround
	MainAxisSpaceEvenly
)

type Options struct {
	MainAxis  MainAxisAlignment
	CrossAxis CrossAxisAlignment

	// Gap controls how much space is placed between each child before the children are sized.
	Gap uint16
}

// Row returns a [vxfw.Widget] that lays out children horizontally.
func Row(children []vxfw.Widget, options Options) vxfw.Widget {
	return &flex{children: children, options: options, orientation: Horizontal}
}

// Column returns a [vxfw.Widget] that lays out children vertically.
func Column(children []vxfw.Widget, options Options) vxfw.Widget {
	return &flex{children: children, options: options, orientation: Vertical}
}

type flex struct {
	children []vxfw.Widget
	options  Options

	orientation Orientation
}

var _ vxfw.Widget = flex{}

func (f flex) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	var totalFlex, maxCross uint16
	surfaces := make([]vxfw.Surface, len(f.children))
	used := f.options.Gap * uint16(len(f.children)-1)
	maxMain := f.orientation.MainAxis(ctx.Max)

	// First pass, lay out all non flexible children.
	// Determine how much space they've used on the main axis, and the largest size on the
	// cross axis.
	for i, child := range f.children {
		if c, ok := child.(flexible); ok {
			// If the flex factor is 0, this is the same as being intrinsically sized
			factor := c.FlexFactor()
			if factor > 0 {
				totalFlex += factor
				continue
			}
		}

		surface, err := child.Draw(instrinsicConstraint(ctx, f.orientation, f.options.CrossAxis))
		if err != nil {
			return vxfw.Surface{}, err
		}

		surfaces[i] = surface
		used += f.orientation.MainAxis(surface.Size)
		maxCross = f.orientation.CrossMax(surface.Size, maxCross)
	}

	// The remaining space is divided among flexible children
	remaining := maxMain - used

	for i, child := range f.children {
		c, ok := child.(flexible)
		// Non-flexible children, or children with a flex factor of 0, were laid out in the
		// first pass above.
		if !ok || c.FlexFactor() == 0 {
			continue
		}

		size := uint16(0)

		if i == len(f.children)-1 {
			// last child gets all of the remaining space
			// NOTE: This is calculated from the space used in this layout pass
			// Whereas flex unit distribution is calculated based on what was remaining
			// after the first pass was completed
			size = maxMain - used
		} else {
			// otherwise, size is based on the flex factor
			size = (remaining * c.FlexFactor()) / totalFlex
		}

		// If c is FlexLoose, we loosen the minimum constraint to 0
		// Otherwise (the default), the child must take a tight constraint
		cons := flexibleConstraint(ctx, f.orientation, f.options.CrossAxis, size)
		if c.FlexLoose() {
			cons = f.orientation.Loosen(cons)
		} else {
			cons = f.orientation.Tighten(cons)
		}

		surface, err := c.Draw(cons)
		if err != nil {
			return vxfw.Surface{}, err
		}

		surfaces[i] = surface
		used += f.orientation.MainAxis(surface.Size)
		maxCross = f.orientation.CrossMax(surface.Size, maxCross)
	}

	// We have all of our surfaces, we know our constraints, it's time to finalize the layout.
	// Each child is placed within the parent surface based on the layout options.
	// TODO: Implement MainAxisSize option which allows the main axis to take min (size of all
	// children) or max. For now, we'll always use the max.
	// Note that even if we implement min main axis size, it becomes irrelevant if any of the
	// children are flexible (and not loose)
	size := f.orientation.Size(maxMain, maxCross)
	surface := vxfw.Surface{
		Size:     size,
		Children: make([]vxfw.SubSurface, len(surfaces)),
	}

	/*
		Distribution.

		First, apply f.options.Gap *between* each child
		If main axis start, no other adjustments
		if main axis end, offset starts at remaining
		If main axis center, offset starts at remaining / 2
		If main axis space between, gap increases by remaining / (len(children)-1)
		If main axis space around, gap increases by remaining / len(children), offset starts at gap / 2
		If main axis space evenly, gap increases by remaining / len(children)+1, offset starts at gap

		Cross offset is simpler:

		if cross axis start, offset is 0
		if cross axis end, offset is max cross - child cross
		if cross axis center, offset is (max cross - child cross) / 2
		if cross axis stretch, offset is 0 (child is already tight in the cross axis)
	*/

	remaining = maxMain - used
	var offset, gap uint16
	var nchildren uint16 = uint16(len(f.children))

	gap = f.options.Gap

	switch f.options.MainAxis {
	case MainAxisEnd:
		offset = remaining
	case MainAxisCenter:
		offset = remaining / 2
	case MainAxisSpaceBetween:
		// Place all remaining space between the children
		gap += remaining / (nchildren - 1)
	case MainAxisSpaceAround:
		// Place all remaining space between children, with half that space on each end
		chunk := remaining / nchildren
		gap += chunk
		offset = chunk / 2
	case MainAxisSpaceEvenly:
		// Place all remaining space between, before, and after children equally
		chunk := remaining / (nchildren + 1)
		gap += chunk
		offset = chunk
	}

	// Iterate over children, applying spacing and distribution options
	var cross uint16
	for i, child := range surfaces {
		// If this is not the first child, add gap to offset
		if i > 0 {
			offset += gap
		}

		cross = 0
		switch f.options.CrossAxis {
		case CrossAxisEnd:
			cross = maxCross - f.orientation.CrossAxis(child.Size)
		case CrossAxisCenter:
			cross = (maxCross - f.orientation.CrossAxis(child.Size)) / 2
		}

		origin := f.orientation.Origin(int(offset), int(cross))
		surface.Children[i] = vxfw.SubSurface{
			Origin:  origin,
			Surface: child,
		}

		offset += f.orientation.MainAxis(child.Size)
	}

	return surface, nil
}

// instrinsicConstraint takes a [vxfw.DrawContext] and returns a new one with the main axis
// unbound, and the cross axis adjusted based on the crossalign.
// This constraint is used to compute instrinsic sizes of non-[Flexible] children in the first
// layout pass.
func instrinsicConstraint(ctx vxfw.DrawContext, o Orientation, crossalign CrossAxisAlignment) vxfw.DrawContext {
	out := vxfw.DrawContext(ctx)
	switch o {
	case Horizontal:
		out.Max.Width = math.MaxUint16
		if crossalign == CrossAxisStretch {
			out.Min.Height = out.Max.Height
		}
	case Vertical:
		out.Max.Height = math.MaxUint16
		if crossalign == CrossAxisStretch {
			out.Min.Width = out.Max.Width
		}
	}
	return out
}

// flexibleConstraint is like [instrinsicConstraint] but sets the main axis to the specified size.
// This is used to layout a flexible child after determining its portion of the available space.
func flexibleConstraint(ctx vxfw.DrawContext, o Orientation, crossalign CrossAxisAlignment, size uint16) vxfw.DrawContext {
	out := vxfw.DrawContext(ctx)
	switch o {
	case Horizontal:
		out.Max.Width = size
		if crossalign == CrossAxisStretch {
			out.Min.Height = out.Max.Height
		}
	case Vertical:
		out.Max.Height = size
		if crossalign == CrossAxisStretch {
			out.Min.Width = out.Max.Width
		}
	}
	return out
}
