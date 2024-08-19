package gui

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/homearchbishop/downgram/internal/gui/binding"
	"github.com/homearchbishop/downgram/internal/gui/router"
)

func Start() {
	go func() {
		window := &app.Window{}
		window.Option(app.Title("Downgram"))
		binding.UseWindow(window)
		if err := loop(window); err != nil {
			log.Fatalf("exiting due to error: %v", err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(window *app.Window) error {
	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, router.GetCurViewWidget()),
			)

			e.Frame(gtx.Ops)
		}
	}
}

func bindData(window *app.Window) {

}
