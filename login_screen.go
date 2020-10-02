package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/widget"
)

var fileExtensions []string = []string{".pwsafe3", ".pwsafe", "psafe3"}

func (t *Thing) getFiles() ([]string, error) {
	filesStr := t.app.Preferences().String("files")
	files := make([]string, 0)

	// list to array
	if filesStr != "" {
		if err := json.Unmarshal([]byte(filesStr), &files); err != nil {
			dialog.ShowError(err, t.win)
			return nil, err
		}
	}

	return files, nil
}

func (t *Thing) addFile(file string) {
	files, err := t.getFiles()
	if err != nil {
		dialog.ShowError(err, t.win)
		return
	}

	// put file to top, rest follows
	dedups := make(map[string]bool)
	nfiles := make([]string, 0, len(files))
	if file != "" {
		nfiles = append(nfiles, file)
	}

	for _, f := range files {
		if f == file || f == "" {
			continue
		}
		if _, exists := dedups[f]; exists {
			continue
		}
		dedups[f] = true
		nfiles = append(nfiles, f)
	}
	files = nfiles

	// save back
	filesStr, err := json.Marshal(files)
	if err != nil {
		dialog.ShowError(err, t.win)
		return
	}
	t.app.Preferences().SetString("files", string(filesStr))

	t.fileInput.SetOptions(files)
	t.fileInput.Text = file
	t.fileInput.Refresh()
}

func (t *Thing) ShowLoginScreen() {
	t.win.SetContent(t.LoginScreen)
	t.win.Show()
}

func (t *Thing) MakeLoginScreen() *widget.Box {

	hello := widget.NewLabel("Hello Fyne!")

	files, err := t.getFiles()
	if err != nil {
		dialog.ShowError(err, t.win)
	}
	log.Printf("files: %+v", files)

	//t.fileInput = widget.NewSelectEntry(files)
	t.fileInput = newEnterSelectEntry()
	if len(files) > 0 {
		t.addFile(files[0])
	} else {
		t.addFile("")
	}

	//t.fileInput.KeyDown = func(text string) {
	//	log.Printf("got input file: %+v", text)
	//}

	browse := widget.NewButton("Browse", func() {
		log.Println("browse...")
		fd := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
			if err == nil && r == nil {
				return
			}
			if err != nil {
				dialog.ShowError(err, t.win)
				return
			}
			log.Printf("what now: %s", r.URI())
			t.input = r
			t.addFile(strings.TrimPrefix(r.URI().String(), "file://"))
		}, t.win)
		fd.SetFilter(storage.NewExtensionFileFilter(fileExtensions))
		fd.Show()
	})
	password := widget.NewPasswordEntry()

	c := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, browse),
		t.fileInput,
		//widget.NewHScrollContainer(file),
		browse,
	)

	f := widget.NewForm(
		widget.NewFormItem("File", c),
		widget.NewFormItem("Password", password),
	)

	return widget.NewVBox(
		hello,
		f,
		widget.NewHBox(
			widget.NewButton("New", func() {
				log.Println("browse...")
				fd := dialog.NewFileSave(func(w fyne.URIWriteCloser, err error) {
					if err == nil && w == nil {
						return
					}
					if err != nil {
						dialog.ShowError(err, t.win)
						return
					}
					log.Printf("what now: %s", w.URI())
					//t.input = w
					t.addFile(strings.TrimPrefix(w.URI().String(), "file://"))
				}, t.win)
				fd.SetFilter(storage.NewExtensionFileFilter(fileExtensions))
				fd.Show()

			}),
			widget.NewButton("Open", func() {
				pwfile := t.fileInput.Text
				log.Printf("try this: %+v", password.Text)
				if t.input == nil {
					var err error
					t.input, err = os.Open(pwfile)
					if err != nil {
						dialog.ShowError(err, t.win)
						log.Printf("failed to open file: %s", pwfile)
						t.input = nil
						return
					}
				}
				if t.input != nil {
					_, err := t.db.Decrypt(t.input, password.Text)
					if err != nil {
						dialog.ShowError(err, t.win)
						log.Printf("failed to open file: %v", err)
						return
					}
					// after login put file to front
					t.addFile(pwfile)
					t.ShowMainScreen()
				}
			}),
		),
	)
}
