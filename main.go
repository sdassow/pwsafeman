package main

import (
	"io"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"

	// https://godoc.org/github.com/lucasepe/pwsafe
	// https://github.com/lucasepe/pwsafe/blob/master/dbFile.go
	"github.com/lucasepe/pwsafe"
)

type Thing struct {
	app        fyne.App
	win        fyne.Window
	input      io.ReadCloser
	db         pwsafe.V3
	fileInput  *enterSelectEntry
	table      *widget.Box
	configRoot string
}

func NewThing() *Thing {
	t := &Thing{}
	t.app = app.NewWithID("pwsafeman")
	t.win = t.app.NewWindow("pwsafeman")

	t.win.Resize(fyne.Size{400, 300})

	//t.configRoot = app.Storage().RootURI().String()

	return t
}

//func (t *Thing) OpenDatabase(file, secret string) {
//_, err := t.db.Decrypt(file, secret)
//}

func main() {
	thing := NewThing()
	thing.LoginScreen()
	thing.win.ShowAndRun()
}
