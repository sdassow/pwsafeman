package main

import (
	"io"
	"bytes"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/dialog"
	//"fyne.io/fyne/storage"

	// https://godoc.org/github.com/lucasepe/pwsafe
	// https://github.com/lucasepe/pwsafe/blob/master/dbFile.go
	"github.com/lucasepe/pwsafe"

	"github.com/sbertrang/atomic"
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

/*
type myuri struct {
	path string
}

func (m myuri) String() string {
	return m.Scheme() + m.path
}

func (m myuri) Extension() string {
	return ""
}

func (m myuri) MimeType() string {
	return ""
}

func (m myuri) Scheme() string {
	return "file://"
}
*/

func (t *Thing) SaveDb() {
/*
	uri := myuri{t.fileInput.Text}
	w, err := storage.SaveFileToURI(uri)
	if err != nil {
		dialog.ShowError(err, t.win)
	}
	defer w.Close()
*/
	var w bytes.Buffer

	_, err := t.db.Encrypt(&w)
	if err != nil {
		dialog.ShowError(err, t.win)
	}

	if err := atomic.WriteFile(t.fileInput.Text, &w); err != nil {
		dialog.ShowError(err, t.win)
	}
	t.MainScreen()
}

func main() {
	thing := NewThing()
	thing.LoginScreen()
	thing.win.ShowAndRun()
}
