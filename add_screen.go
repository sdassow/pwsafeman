package main

import (
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/lucasepe/pwsafe"
)

func (t *Thing) ShowAddScreen() {
	t.win.SetContent(t.AddScreen)
	t.win.Show()
}

func (t *Thing) MakeAddScreen() *fyne.Container {
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

	ftitle := widget.NewEntry()
	fgroup := widget.NewEntry()
	fusername := widget.NewEntry()
	fpassword := widget.NewPasswordEntry()
	furl := widget.NewEntry()
	fnotes := widget.NewMultiLineEntry()

	//data := widget.NewLabel("DATA")
	form := widget.NewForm()
	form.Append("Title", ftitle)
	form.Append("Group", fgroup)
	form.Append("Username", fusername)
	form.Append("Password", fpassword)
	form.Append("URL", furl)
	form.Append("Notes", fnotes)
	form.CancelText = "Cancel"
	form.SubmitText = "Ok"

	form.OnCancel = func() {
		t.ShowMainScreen()
	}
	form.OnSubmit = func() {
		rec := pwsafe.Record{
			Title: ftitle.Text,
			Group: fgroup.Text,
			Username: fusername.Text,
			Password: fpassword.Text,
			URL: furl.Text,
			Notes: fnotes.Text,
		}
		log.Printf("ADD RECORD: %+v", rec)

		t.db.SetRecord(rec)
		t.SaveDb()
	}

	form.Refresh()

	// inner table
	xtab := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, nil),
		form,
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

