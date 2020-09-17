package main

import (
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type enterSelectEntry struct {
	widget.SelectEntry
}

func (e *enterSelectEntry) onEnter() {
	log.Printf("ENTER")
}

func newEnterSelectEntry() *enterSelectEntry {
	entry := &enterSelectEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *enterSelectEntry) KeyDown(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyReturn {
		e.onEnter()
	} else {
		e.SelectEntry.KeyDown(key)
	}
}
