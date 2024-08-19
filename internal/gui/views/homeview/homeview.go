package homeview

import (
	"bytes"
	"image"
	"strconv"

	"golang.org/x/image/draw"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/homearchbishop/downgram/internal/cache"
	"github.com/homearchbishop/downgram/internal/gui/binding"
	"github.com/homearchbishop/downgram/internal/gui/theme"
	"github.com/homearchbishop/downgram/internal/tg"
	"github.com/ncruces/zenity"
)

const (
	idMode = iota
	usernameMode
)

var (
	usernameAndIDEditor = &widget.Editor{}

	searchClickable = &widget.Clickable{}

	idModeClickable       = &widget.Clickable{}
	usernameModeClickable = &widget.Clickable{}

	dirPathSelectorClickable = &widget.Clickable{}

	getMoreClickable  = &widget.Clickable{}
	downloadClickable = &widget.Clickable{}

	mediaItemClickableList = []*widget.Clickable{}

	currentSearchMode = binding.BindVar(idMode)

	downloadDirPath = binding.BindVar(cache.Get("downloadDir", ""))

	searchTipStr = binding.BindVar("")
	detailTipStr = binding.BindVar("")

	list = &layout.List{Axis: layout.Vertical}

	imgByteList   = binding.BindSlice(&[][]byte{})
	mediaInfoList = binding.BindSlice(&[]tg.MediaInfo{})
	selectedList  = binding.BindSlice(&[]bool{})
)

func createHomeViewInputModeTabBtn(mode, txt string, clickable *widget.Clickable, clickAction func()) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			btnColor := theme.ThSet.HomeViewHiBg
			if (currentSearchMode.Get() != idMode && mode == "id") || (currentSearchMode.Get() != usernameMode && mode == "username") {
				btnColor = theme.ThSet.HomeViewSubBg
			}
			paint.FillShape(gtx.Ops, btnColor, clip.Rect{
				Min: image.Point{X: 0, Y: 0},
				Max: image.Point{X: 140, Y: 40},
			}.Op())
			gtx.Constraints.Min.X = 140
			label := material.Label(theme.Th, unit.Sp(16), txt)
			label.Color = theme.ThSet.HomeViewHiBgBtnFg
			label.Alignment = text.Middle
			label.Layout(gtx)
			return clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				for {
					ev, _ := gtx.Event(pointer.Filter{
						Target: clickable,
						Kinds:  pointer.Enter | pointer.Leave | pointer.Press | pointer.Release,
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
					} else if x.Kind == pointer.Leave {
						pointer.CursorDefault.Add(gtx.Ops)
					} else if x.Kind == pointer.Press {
						pointer.CursorPointer.Add(gtx.Ops)
					} else if x.Kind == pointer.Release {
						pointer.CursorPointer.Add(gtx.Ops)
						clickAction()
					}
				}

				return layout.Dimensions{Size: image.Point{X: 140, Y: 40}}
			})
		})
	})
}

func createHomeViewInputModeTab() layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		inset := layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(16)}
		return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				createHomeViewInputModeTabBtn("id", "ID", idModeClickable, func() {
					currentSearchMode.Set(idMode)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Spacer{Width: unit.Dp(4)}.Layout(gtx)
				}),
				createHomeViewInputModeTabBtn("username", "username", usernameModeClickable, func() {
					currentSearchMode.Set(usernameMode)
				}),
			)
		})
	})
}

func createHomeViewSearchBtn(txt string, clickable *widget.Clickable, clickAction func()) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			paint.FillShape(gtx.Ops, theme.ThSet.HomeViewHiBg, clip.Rect{
				Min: image.Point{X: 0, Y: 0},
				Max: image.Point{X: 288, Y: 40},
			}.Op())
			gtx.Constraints.Min.X = 288
			label := material.Label(theme.Th, unit.Sp(16), txt)
			label.Color = theme.ThSet.HomeViewHiBgBtnFg
			label.Alignment = text.Middle
			label.Layout(gtx)
			return clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				for {
					ev, _ := gtx.Event(pointer.Filter{
						Target: clickable,
						Kinds:  pointer.Enter | pointer.Leave | pointer.Press | pointer.Release,
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
					} else if x.Kind == pointer.Leave {
						pointer.CursorDefault.Add(gtx.Ops)
					} else if x.Kind == pointer.Press {
						pointer.CursorPointer.Add(gtx.Ops)
					} else if x.Kind == pointer.Release {
						pointer.CursorPointer.Add(gtx.Ops)
						clickAction()
					}
				}

				return layout.Dimensions{Size: image.Point{X: 288, Y: 40}}
			})
		})
	})
}

func createHomeViewInput(placeholder string, widget *widget.Editor) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		inset := layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(16)}
		return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			editor := material.Editor(theme.Th, widget, placeholder)
			editor.Color = theme.ThSet.HomeViewSubFg
			editor.TextSize = unit.Sp(16)

			gtx.Constraints.Min.X = int(unit.Dp(300 - 16))
			gtx.Constraints.Max.X = int(unit.Dp(300 - 16))
			gtx.Constraints.Min.Y = int(unit.Dp(40))
			editor.Editor.SingleLine = true

			dims := editor.Layout(gtx)
			paint.FillShape(gtx.Ops, theme.ThSet.HomeViewHiFg, clip.Rect{
				Min: image.Point{X: 0, Y: dims.Size.Y},
				Max: image.Point{X: dims.Size.X, Y: dims.Size.Y + 3},
			}.Op())

			return dims
		})
	})
}

func createHomeViewSearchTip(tipStr string) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(32)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = 288
			label := material.Label(theme.Th, unit.Sp(16), tipStr)
			label.Color = theme.ThSet.HomeViewSubFg
			label.Alignment = text.Start
			label.Layout(gtx)
			return layout.Dimensions{Size: image.Point{X: 288, Y: 200}}
		})
	})
}

func createHomeViewDirPathSelector(clickable *widget.Clickable) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(32)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			paint.FillShape(gtx.Ops, theme.ThSet.HomeViewHiBg, clip.Rect{
				Min: image.Point{X: 0, Y: 0},
				Max: image.Point{X: 288, Y: 40},
			}.Op())
			gtx.Constraints.Min.X = 288
			clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Label(theme.Th, unit.Sp(16), "Select download path")
				label.Color = theme.ThSet.HomeViewHiBgBtnFg
				label.Alignment = text.Middle
				label.Layout(gtx)
				for {
					ev, _ := gtx.Event(pointer.Filter{
						Target: clickable,
						Kinds:  pointer.Enter | pointer.Leave | pointer.Press | pointer.Release,
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
					} else if x.Kind == pointer.Leave {
						pointer.CursorDefault.Add(gtx.Ops)
					} else if x.Kind == pointer.Press {
						pointer.CursorPointer.Add(gtx.Ops)
					} else if x.Kind == pointer.Release {
						pointer.CursorPointer.Add(gtx.Ops)
						folder, err := zenity.SelectFile(
							zenity.Title("Select download directory"),
							zenity.Directory(),
						)
						if err == nil {
							downloadDirPath.Set(folder)
							cache.Update("downloadDir", folder)
						}
					}
				}
				return layout.Dimensions{Size: image.Point{X: 288, Y: 40}}
			})
			return layout.Dimensions{Size: image.Point{X: 288, Y: 40}}
		})
	})
}

func createImgItem(gtx layout.Context, imgBytes []byte, mediaInfo tg.MediaInfo, clickable *widget.Clickable, i int, wh int) layout.Dimensions {
	if imgBytes == nil {
		defer clip.Rect{Max: image.Pt(100, 100)}.Push(gtx.Ops).Pop()
		paint.ColorOp{Color: theme.ThSet.HomeViewSubBg}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return layout.Dimensions{Size: image.Point{X: wh, Y: wh}}
	}

	imgWH := wh
	imgStartXY := 0

	if selectedList.Get(i) {
		paint.FillShape(gtx.Ops, theme.ThSet.HomeViewHiBg, clip.Stroke{
			Path:  clip.RRect{Rect: image.Rect(6, 6, wh-6, wh-6)}.Path(gtx.Ops),
			Width: 8,
		}.Op())
		imgWH = wh - 12
		imgStartXY = 6
	}

	imgReader := bytes.NewReader(imgBytes)
	img, _, _ := image.Decode(imgReader)
	minLen := img.Bounds().Size().X
	if img.Bounds().Size().Y < minLen {
		minLen = img.Bounds().Size().Y
	}
	x0 := (img.Bounds().Size().X - minLen) / 2
	y0 := (img.Bounds().Size().Y - minLen) / 2
	x1 := x0 + minLen
	y1 := y0 + minLen
	dst := image.NewRGBA(image.Rect(0, 0, imgWH, imgWH))
	draw.NearestNeighbor.Scale(dst, dst.Bounds(), img, image.Rect(x0, y0, x1, y1), draw.Over, nil)
	op.Offset(image.Point{X: imgStartXY, Y: imgStartXY}).Add(gtx.Ops)
	paint.NewImageOp(dst).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	op.Offset(image.Point{X: -imgStartXY, Y: -imgStartXY}).Add(gtx.Ops)

	if mediaInfo.MediaType == "document" {
		maxSizeMb := mediaInfo.VideoSize / 1024 / 1024
		label := material.Label(theme.Th, unit.Sp(10), strconv.Itoa(int(mediaInfo.Duration))+"s|"+
			strconv.FormatInt(maxSizeMb, 10)+"MB\n"+strconv.FormatInt(mediaInfo.MediaID, 10))
		label.Color = theme.ThSet.HomeViewHiBgBtnFg
		label.Alignment = text.Start
		label.Layout(gtx)
	} else {
		label := material.Label(theme.Th, unit.Sp(10), strconv.FormatInt(mediaInfo.MediaID, 10))
		label.Color = theme.ThSet.HomeViewHiBgBtnFg
		label.Alignment = text.Start
		label.Layout(gtx)
	}

	clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		for {
			ev, _ := gtx.Event(pointer.Filter{
				Target: clickable,
				Kinds:  pointer.Enter | pointer.Leave | pointer.Press | pointer.Release,
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
			} else if x.Kind == pointer.Leave {
				pointer.CursorDefault.Add(gtx.Ops)
			} else if x.Kind == pointer.Press {
				pointer.CursorPointer.Add(gtx.Ops)
			} else if x.Kind == pointer.Release {
				pointer.CursorPointer.Add(gtx.Ops)
				selectedList.Set(i, !selectedList.Get(i))
			}
		}
		return layout.Dimensions{Size: image.Point{X: wh, Y: wh}}
	})

	return layout.Dimensions{Size: image.Point{X: wh, Y: wh}}
}

func createWrapLayout(gtx layout.Context, maxWidth int) layout.Dimensions {
	const itemWidth = 100
	const itemHeight = 100

	imgBytesList := *imgByteList.Val()
	mediaInfoList := *mediaInfoList.Val()

	var flexChildren []layout.FlexChild
	var rowChildren []layout.FlexChild
	currentWidth := 0

	for i, imgBytes := range imgBytesList {
		i := i
		imgBytes := imgBytes
		mediaInfo := mediaInfoList[i]
		mediaItemClickable := mediaItemClickableList[i]
		if currentWidth+itemWidth > maxWidth {
			completedRowChildren := rowChildren
			flexChildren = append(flexChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, completedRowChildren...)
			}))
			rowChildren = nil
			currentWidth = 0
		}
		rowChildren = append(rowChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return createImgItem(gtx, imgBytes, mediaInfo, mediaItemClickable, i, itemWidth)
		}))
		currentWidth += itemWidth
	}

	if len(rowChildren) > 0 {
		completedRowChildren := rowChildren
		flexChildren = append(flexChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, completedRowChildren...)
		}))
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, flexChildren...)
}

func createHomeViewGetMoreAndDownloadBtn(txt string, clickable *widget.Clickable, clickAction func()) layout.FlexChild {
	if len(*mediaInfoList.Val()) == 0 {
		return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: image.Point{X: 0, Y: 0}}
		})
	}
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			paint.FillShape(gtx.Ops, theme.ThSet.HomeViewHiBg, clip.Rect{
				Min: image.Point{X: 0, Y: 0},
				Max: image.Point{X: 320, Y: 40},
			}.Op())
			gtx.Constraints.Min.X = 320
			label := material.Label(theme.Th, unit.Sp(16), txt)
			label.Color = theme.ThSet.HomeViewHiBgBtnFg
			label.Alignment = text.Middle
			label.Layout(gtx)
			return clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				for {
					ev, _ := gtx.Event(pointer.Filter{
						Target: clickable,
						Kinds:  pointer.Enter | pointer.Leave | pointer.Press | pointer.Release,
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
					} else if x.Kind == pointer.Leave {
						pointer.CursorDefault.Add(gtx.Ops)
					} else if x.Kind == pointer.Press {
						pointer.CursorPointer.Add(gtx.Ops)
					} else if x.Kind == pointer.Release {
						pointer.CursorPointer.Add(gtx.Ops)
						clickAction()
					}
				}

				return layout.Dimensions{Size: image.Point{X: 320, Y: 40}}
			})
		})
	})
}

func HomeViewWidget(gtx layout.Context) layout.Dimensions {
	paint.Fill(gtx.Ops, theme.ThSet.HomeViewBg)

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			width := 300

			gtx.Constraints.Max.X = width
			gtx.Constraints.Min.X = width

			paint.FillShape(gtx.Ops, theme.ThSet.HomeViewSubHiBg, clip.Rect{
				Min: image.Point{X: 0, Y: 0},
				Max: image.Point{X: width, Y: gtx.Constraints.Max.Y},
			}.Op())

			inset := layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(4), Right: unit.Dp(4)}
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				var currentInput layout.FlexChild
				if currentSearchMode.Get() == idMode {
					usernameAndIDEditor.Filter = "1234567890-"
					currentInput = createHomeViewInput("e.g. 3321450088", usernameAndIDEditor)
				} else {
					usernameAndIDEditor.Filter = ""
					currentInput = createHomeViewInput("e.g. examplegp", usernameAndIDEditor)
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					createHomeViewInputModeTab(),
					currentInput,
					createHomeViewSearchBtn("Search", searchClickable, func() {
						search(currentSearchMode.Get(), usernameAndIDEditor.Text())
					}),
					createHomeViewSearchTip(searchTipStr.Get()),
					createHomeViewDirPathSelector(dirPathSelectorClickable),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.Label(theme.Th, unit.Sp(16), "download to:\n"+downloadDirPath.Get())
						label.Color = theme.ThSet.HomeViewFg
						label.Alignment = text.Start
						label.Layout(gtx)
						return layout.Dimensions{Size: image.Point{X: 300, Y: 40}}
					}),
				)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					inset := layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(8), Right: unit.Dp(8)}
					return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.Label(theme.Th, unit.Sp(16), detailTipStr.Get()).Layout(gtx)
					})
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return list.Layout(gtx, 1, func(gtx layout.Context, i int) layout.Dimensions {
						inset := layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(8), Right: unit.Dp(8)}
						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return createWrapLayout(gtx, gtx.Constraints.Max.X)
						})
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						createHomeViewGetMoreAndDownloadBtn("Get more", getMoreClickable, func() {
							updateNextMediaListByBtn()
						}),
						createHomeViewGetMoreAndDownloadBtn("Download", downloadClickable, func() {
							downloadMediaToDir(downloadDirPath.Get(), *mediaInfoList.Val(), *selectedList.Val())
							allFalseSelectedList := make([]bool, len(*selectedList.Val()))
							selectedList.Overwrite(&allFalseSelectedList)
						}),
					)
				}),
			)
		}),
	)
}
