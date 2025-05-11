package vxlayout

import (
	"testing"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/text"
)

func TestFlexRow(t *testing.T) {
	layout := Row([]vxfw.Widget{
		text.New("abc"),
		Expanded(text.New("def"), 1),
		Expanded(text.New("ghi"), 1),
		Expanded(text.New("jkl\nnmo"), 1),
	}, Options{})

	ctx := vxfw.DrawContext{Max: vxfw.Size{Width: 16, Height: 16}, Characters: vaxis.Characters}
	surface, err := layout.Draw(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if surface.Size.Width != 16 {
		t.Logf("wrong flex width, got=%d, want=16", surface.Size.Width)
		t.Fail()
	}

	if surface.Size.Height != 2 {
		t.Logf("wrong flex height, got=%d, want=2", surface.Size.Height)
		t.Fail()
	}

	if len(surface.Children) != 4 {
		t.Logf("wrong number of flex children, got=%d, want=4", len(surface.Children))
		t.Fail()
	}

	// col moves forward by the width of each child, used to assert origins
	col := 0

	// expected widths of each child
	// first child should be 3 since it's not flexible
	// remaining 3 children have equal flex and share the remaining (16-3) columns
	// since 13 / 3 is 4 that leaves 1 extra which is assigned to the last child
	widths := []uint16{
		3,
		3 + 1,
		3 + 1,
		3 + 1 + 1,
	}

	for i, want := range widths {
		child := surface.Children[i]

		got := child.Surface.Size.Width
		if got != want {
			t.Logf("wrong width for child %d, got=%d, want=%d", i, got, want)
			t.Fail()
		}
		if child.Origin.Col != col {
			t.Logf("wrong origin for child %d, got=%d, want=%d", i, child.Origin.Col, col)
			t.Fail()
		}
		col += int(want)
	}
}

func TestFlexColumn(t *testing.T) {
	// Wrap Column in Constrained so Space doesn't take all 16 columns wide
	layout := Constrained(Column([]vxfw.Widget{
		text.New("abc"),
		Expanded(text.New("def"), 1),
		Space(1),
		Expanded(text.New("jkl\nnmo"), 2),
	}, Options{}), nil, &vxfw.Size{Width: 3})

	ctx := vxfw.DrawContext{Max: vxfw.Size{Width: 16, Height: 16}, Characters: vaxis.Characters}
	surface, err := layout.Draw(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if surface.Size.Width != 3 {
		t.Logf("wrong flex width, got=%d, want=3", surface.Size.Width)
		t.Fail()
	}

	if surface.Size.Height != 16 {
		t.Logf("wrong flex height, got=%d, want=16", surface.Size.Height)
		t.Fail()
	}

	if len(surface.Children) != 4 {
		t.Logf("wrong number of flex children, got=%d, want=4", len(surface.Children))
		t.Fail()
	}

	// row moves forward by the height of each child, used to assert origins
	row := 0

	// first child is not flexible, so it takes 1 row, leaving 15 to distribute
	// 2 and 3 are both flex 1, 4 is flex 2 so it takes as much available space
	// as the others.
	// 15 / 4 == 3 (+3 leftover), giving 3, 3, and 3+3 (flex 2) + leftover 3 to the last
	// child.
	heights := []uint16{
		1,
		3,
		3,
		9,
	}

	for i, want := range heights {
		child := surface.Children[i]

		got := child.Surface.Size.Height
		if got != want {
			t.Logf("wrong height for child %d, got=%d, want=%d", i, got, want)
			t.Fail()
		}
		if child.Origin.Row != row {
			t.Logf("wrong origin for child %d, got=%d, want=%d", i, child.Origin.Row, row)
			t.Fail()
		}
		row += int(want)
	}
}
