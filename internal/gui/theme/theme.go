package theme

import (
	"image/color"

	"gioui.org/widget/material"
)

var Th *material.Theme = material.NewTheme()
var ThSet ThemeSet

type ThemeSet struct {
	LoginViewBg        color.NRGBA
	LoginViewHiBg      color.NRGBA
	LoginViewFg        color.NRGBA
	LoginViewHiFg      color.NRGBA
	LoginViewSubFg     color.NRGBA
	LoginViewHiBgBtnFg color.NRGBA

	HomeViewBg        color.NRGBA
	HomeViewSubBg     color.NRGBA
	HomeViewHiBg      color.NRGBA
	HomeViewSubHiBg   color.NRGBA
	HomeViewFg        color.NRGBA
	HomeViewHiFg      color.NRGBA
	HomeViewSubFg     color.NRGBA
	HomeViewHiBgBtnFg color.NRGBA
}

func UseTheme(themeSet ThemeSet) {
	ThSet = themeSet
}

func init() {
	UseTheme(Tg)
}
