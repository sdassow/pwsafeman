package main

import (
	//"log"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/lucasepe/pwsafe"
)

func (t *Thing) ShowViewScreen(rec pwsafe.Record) {
	//t.record = rec;
	t.win.SetContent(t.ViewScreen)
	t.win.Show()
}

func (t *Thing) MakeViewScreen() *fyne.Container {
	bottom := widget.NewLabel("Bottom")
	right := widget.NewLabel("Right")

	// button to add entry
	back := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		t.ShowMainScreen()
	})

	// top line with add button left and search on right
	top := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, back, right),
		back,
		right,
	)

	// view record
	data := widget.NewLabel("DATA")

	// inner table
	xtab := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, nil),
		data,
	)

	c := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(top, bottom, nil, nil),
		top,
		bottom,
		//widget.NewVScrollContainer(t.table),
		//t.table,
		xtab,
	)

	return c
}

