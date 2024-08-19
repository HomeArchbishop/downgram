package router

import (
	"gioui.org/layout"
	"github.com/homearchbishop/downgram/internal/gui/binding"
	"github.com/homearchbishop/downgram/internal/gui/views"
)

var CurViewName = binding.BindVar(LoginView)

func GetCurViewWidget() layout.Widget {
	return ViewWidgetMap[CurViewName.Get()]
}

func init() {
	views.RegisterRouterFunc("routeToHome", func() {
		CurViewName.Set(HomeView)
	})
}
