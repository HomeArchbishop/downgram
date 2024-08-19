package binding

import (
	"gioui.org/app"
)

var window *app.Window

func UseWindow(w *app.Window) {
	window = w
}
