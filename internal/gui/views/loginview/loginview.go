package loginview

import (
	"image"
	"log"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/homearchbishop/downgram/internal/cache"
	"github.com/homearchbishop/downgram/internal/gui/binding"
	"github.com/homearchbishop/downgram/internal/gui/theme"
	"github.com/homearchbishop/downgram/internal/gui/views"
)

const (
	phoneSection = iota
	codeSection
)

var (
	appIDEditor       widget.Editor
	appHashEditor     widget.Editor
	phoneNumberEditor widget.Editor
	codeEditor        widget.Editor
	proxyEditor       widget.Editor

	loginBtnClickable = &widget.Clickable{}

	section = binding.BindVar(phoneSection)

	codeChan chan string
)

func createLoginViewInputLabel(label string) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		inset := layout.Inset{}
		return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			label := material.Label(theme.Th, unit.Sp(12), label)
			label.Color = theme.ThSet.LoginViewHiFg
			label.Font.Weight = 300
			return label.Layout(gtx)
		})
	})
}

func createLoginViewInput(placeholder string, widget *widget.Editor) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		inset := layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(16)}
		return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			editor := material.Editor(theme.Th, widget, placeholder)
			editor.Color = theme.ThSet.LoginViewSubFg
			editor.TextSize = unit.Sp(16)

			gtx.Constraints.Min.X = int(unit.Dp(500))
			gtx.Constraints.Max.X = int(unit.Dp(500))
			gtx.Constraints.Min.Y = int(unit.Dp(40))
			editor.Editor.SingleLine = true

			dims := editor.Layout(gtx)
			paint.FillShape(gtx.Ops, theme.ThSet.LoginViewHiFg, clip.Rect{
				Min: image.Point{X: 0, Y: dims.Size.Y},
				Max: image.Point{X: dims.Size.X, Y: dims.Size.Y + 3},
			}.Op())

			return dims
		})
	})
}

func createLoginViewButton(label string, action func()) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		inset := layout.Inset{Top: unit.Dp(22)}
		return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			button := material.Button(theme.Th, loginBtnClickable, label)
			button.Background = theme.ThSet.LoginViewHiBg
			button.Color = theme.ThSet.LoginViewHiBgBtnFg
			button.TextSize = unit.Sp(16)
			button.CornerRadius = unit.Dp(5)

			gtx.Constraints.Min.X = int(unit.Dp(500))
			gtx.Constraints.Max.X = int(unit.Dp(500))
			gtx.Constraints.Min.Y = int(unit.Dp(40))

			for {
				ev, _ := gtx.Event(pointer.Filter{
					Target: loginBtnClickable,
					Kinds:  pointer.Enter | pointer.Leave,
				})
				if ev == nil {
					break
				}
				x, ok := ev.(pointer.Event)
				if !ok {
					continue
				}
				if x.Kind == pointer.Enter {
					pointer.CursorPointer.Add(gtx.Ops)
				} else {
					pointer.CursorDefault.Add(gtx.Ops)
				}
			}

			for loginBtnClickable.Clicked(gtx) {
				action()
			}

			return button.Layout(gtx)
		})
	})
}

func LoginViewWidget(gtx layout.Context) layout.Dimensions {
	paint.Fill(gtx.Ops, theme.ThSet.LoginViewBg)

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if section.Get() == phoneSection {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				createLoginViewInputLabel("App ID:"),
				createLoginViewInput("", &appIDEditor),
				createLoginViewInputLabel("App Hash:"),
				createLoginViewInput("", &appHashEditor),
				createLoginViewInputLabel("Phone Number:"),
				createLoginViewInput("", &phoneNumberEditor),
				createLoginViewInputLabel("Proxy:"),
				createLoginViewInput("", &proxyEditor),
				createLoginViewButton("Login", func() {
					login(loginInfo{
						appID:    appIDEditor.Text(),
						apphash:  appHashEditor.Text(),
						phone:    phoneNumberEditor.Text(),
						proxyStr: proxyEditor.Text(),
						newLogin: phoneNumberEditor.Text() != cache.Get("phone", ""),
					}, func() string {
						codeChan = make(chan string)
						section.Set(codeSection)
						return <-codeChan
					}, func() {
						views.Router["routeToHome"]()
						log.Println("logined")
					})
				}),
			)
		} else {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				createLoginViewInputLabel("Code:"),
				createLoginViewInput("", &codeEditor),
				createLoginViewButton("Login", func() {
					codeChan <- codeEditor.Text()
					close(codeChan)
				}),
			)
		}
	})
}

func init() {
	appIDEditor.SetText(cache.Get("appid", ""))
	appHashEditor.SetText(cache.Get("apphash", ""))
	phoneNumberEditor.SetText(cache.Get("phone", ""))
	proxyEditor.SetText(cache.Get("proxy", ""))
	codeEditor.SetText("")

	appIDEditor.Filter = "0123456789"
	appHashEditor.Filter = "0123456789abcdefghijklmnopqrstuvwxyz"
	phoneNumberEditor.Filter = "+0123456789()"
	codeEditor.Filter = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
}
