package main

import (
	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/log"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/text"
	"github.com/avidal/vxexp/vxlayout"
)

type App struct {
	infobar vxfw.Widget
	screens []vxfw.Widget
	index   int
}

func (a *App) CaptureEvent(ev vaxis.Event) (vxfw.Command, error) {
	switch ev := ev.(type) {
	case vaxis.Key:
		if ev.Matches('c', vaxis.ModCtrl) {
			return vxfw.QuitCmd{}, nil
		}
		if ev.Matches('l') || ev.Matches('L') {
			a.changeScreen()
			return vxfw.ConsumeAndRedraw(), nil
		}
	}
	return nil, nil
}

func (a *App) changeScreen() {
	if a.index == len(a.screens)-1 {
		a.index = 0
	} else {
		a.index += 1
	}
}

func makeScreen1() vxfw.Widget {
	filler := vaxis.Cell{
		Character: vaxis.Character{
			Grapheme: "·", Width: 1,
		},
		Style: vaxis.Style{
			Background: vaxis.ColorNavy,
		},
	}
	defaultopts := vxlayout.Options{}

	header := func(s string) vxfw.Widget {
		t := text.New(s)
		t.Style.Background = vaxis.ColorGray
		return vxlayout.Row([]vxfw.Widget{t}, defaultopts)
	}

	return vxlayout.Column([]vxfw.Widget{
		// row 1
		header("Three widgets with space distributed evenly."),
		vxlayout.Constrained(vxlayout.Row([]vxfw.Widget{
			text.New("ONE"),
			text.New("TWO"),
			text.New("THREE"),
		}, vxlayout.Options{
			MainAxis: vxlayout.MainAxisSpaceEvenly,
		}), &vxfw.Size{Height: 1}, &vxfw.Size{Height: 1}),

		// row 2
		header("Three widgets center aligned with a 2 col gap."),
		vxlayout.Constrained(vxlayout.Row([]vxfw.Widget{
			text.New("ONE"),
			text.New("TWO"),
			text.New("THREE"),
		}, vxlayout.Options{
			MainAxis: vxlayout.MainAxisCenter,
			Gap:      2,
		}), &vxfw.Size{Height: 1}, &vxfw.Size{Height: 1}),

		// row 3
		header("Three widgets right aligned with a 2 col gap and a 20 column constrained filler."),
		vxlayout.Constrained(vxlayout.Row([]vxfw.Widget{
			text.New("ONE"),
			text.New("TWO"),
			text.New("THREE"),
			vxlayout.Flexible(vxlayout.Constrained(vxlayout.Fill(filler), &vxfw.Size{Width: 20}, &vxfw.Size{Width: 20}), 1),
		}, vxlayout.Options{
			MainAxis: vxlayout.MainAxisEnd,
			Gap:      2,
		}), &vxfw.Size{Height: 1}, &vxfw.Size{Height: 1}),

		// row 4
		header("Three widgets with a 2 col gap and filler in between."),
		vxlayout.Constrained(vxlayout.Row([]vxfw.Widget{
			text.New("ONE"),
			vxlayout.Expanded(vxlayout.Fill(filler), 1),
			text.New("TWO"),
			vxlayout.Expanded(vxlayout.Fill(filler), 1),
			text.New("THREE"),
		}, vxlayout.Options{
			Gap: 2,
		}), &vxfw.Size{Height: 1}, &vxfw.Size{Height: 1}),
	}, vxlayout.Options{
		Gap: 1,
	})
}

func (a *App) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	root := vxfw.NewSurface(ctx.Max.Width, ctx.Max.Height, a)

	infobar, err := a.infobar.Draw(vxfw.DrawContext{
		Max:        vxfw.Size{Width: ctx.Max.Width, Height: 1},
		Characters: ctx.Characters,
	})
	if err != nil {
		return vxfw.Surface{}, err
	}

	screen, err := a.screens[a.index].Draw(ctx.WithMax(vxfw.Size{
		Width: ctx.Max.Width, Height: ctx.Max.Height - 1,
	}))
	if err != nil {
		return vxfw.Surface{}, err
	}

	root.AddChild(0, 0, infobar)
	root.AddChild(0, 1, screen)
	return root, nil
}

func main() {
	vxapp, err := vxfw.NewApp(vaxis.Options{})
	if err != nil {
		log.Error("Couldn't create a new app: %v", err)
	}

	app := &App{
		infobar: text.New("Ctrl+C to quit, L (or l) to switch layouts."),
		screens: []vxfw.Widget{makeScreen1()},
	}

	vxapp.Run(app)
}
