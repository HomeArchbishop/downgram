package router

import (
	"gioui.org/layout"
	"github.com/homearchbishop/downgram/internal/gui/views/homeview"
	"github.com/homearchbishop/downgram/internal/gui/views/loginview"
)

const (
	HomeView int = iota
	LoginView
)

var ViewWidgetMap = map[int]layout.Widget{
	HomeView:  homeview.HomeViewWidget,
	LoginView: loginview.LoginViewWidget,
}
