package main

import (
	"log"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"golang.org/x/text/language"
	"golang.org/x/text/search"

	"github.com/lucasepe/pwsafe"
)

func (t *Thing) find(term string) {
	if term == "" {
		t.list2()
		return
	}

	matcher := search.New(language.English, search.Loose, search.IgnoreCase, search.IgnoreDiacritics)
	pattern := matcher.CompileString(term)

	recs := make([]pwsafe.Record, 0)
	for _, title := range t.db.List() {
		rec, found := t.db.GetRecord(title)
		if !found {
			continue
		}

		start, _ := pattern.IndexString(rec.Title + " " + rec.Username + " " + rec.Group)
		if start == -1 {
			continue
		}
		log.Printf("found: %s", title)

		recs = append(recs, rec)
	}

	t.updateList(recs)
}

func (t *Thing) list() {
	rows := make([][]string, 0)
	for _, title := range t.db.List() {
		log.Printf("entry: %+v", title)

		rec, found := t.db.GetRecord(title)
		if !found {
			continue
		}

		row := []string{
			rec.Title,
			rec.Title,
			rec.Username,
			rec.Group,
			rec.Title,
		}

		rows = append(rows, row)
	}

	t.updateTable(rows)
}

func (t *Thing) updateTable(rows [][]string) {
	table := makeTable(
		[]string{"", "Title", "Username", "Group", ""},
		rows,
	)

	if t.table == nil {
		t.table = table
		return
	}

	t.table.Children[0] = table.Children[0]
	t.table.Children[1] = table.Children[1]
	t.table.Children[2] = table.Children[2]
	t.table.Children[3] = table.Children[3]
	t.table.Children[4] = table.Children[4]
	t.table.Refresh()
}

func makeTable(headings []string, rows [][]string) *widget.Box {

	columns := rowsToColumns(headings, rows)

	objects := make([]fyne.CanvasObject, len(columns))
	for k, col := range columns {
		box := widget.NewVBox(widget.NewLabelWithStyle(
			headings[k],
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		))

		if k == 0 {
			for _, val := range col {
				box.Append(widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
					log.Printf("tap: %s", val)
				}))
			}
		} else if k == len(headings)-1 {
			for _, val := range col {
				box.Append(widget.NewButtonWithIcon("", theme.FileIcon(), func() {
					log.Printf("tap: %s", val)
				}))
			}
		} else {
			for _, val := range col {
				box.Append(widget.NewLabel(val))
			}
		}

		objects[k] = box
	}

	return widget.NewHBox(objects...)
}

func rowsToColumns(headings []string, rows [][]string) [][]string {
	columns := make([][]string, len(headings))
	for _, row := range rows {
		for colK := range row {
			columns[colK] = append(columns[colK], row[colK])
		}
	}
	return columns
}

func (t *Thing) list2() {
	entries := t.db.List()
	recs := make([]pwsafe.Record, 0)
	for _, title := range entries {
		rec, found := t.db.GetRecord(title)
		if !found {
			continue
		}
		recs = append(recs, rec)
	}
	t.updateList(recs)
}

func (t *Thing) clearClipboardSoon() {
	// stop previous timer
	if t.cliptimer != nil {
		t.cliptimer.Stop()
	}

	// after N seconds clear clipboard again
	go func() {
		t.cliptimer = time.AfterFunc(10 * time.Second, func() {
			t.win.Clipboard().SetContent("")
			t.bottom.Text = "Clipboard cleared..."
			t.bottom.Refresh()
		})
	}()
}

func (t *Thing) updateList(recs []pwsafe.Record) {
	objects := make([]fyne.CanvasObject, 0, len(recs))
	for _, rec := range recs {
		box := widget.NewHBox(
			widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
				t.ShowViewScreen(rec)
			}),
			widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
				// put password into clipboard
				t.win.Clipboard().SetContent(rec.Password)
				t.bottom.Text = "Copied password to clipboard..."
				t.bottom.Refresh()
				t.clearClipboardSoon()
			}),
			widget.NewLabelWithStyle(
				rec.Title,
				fyne.TextAlignLeading,
				fyne.TextStyle{Bold: true},
			),
			widget.NewLabel(rec.Username),
		)

		objects = append(objects, box)
	}

	if t.table == nil {
		t.table = widget.NewVBox(objects...)
	} else {
		t.table.Children = objects
	}
	t.table.Refresh()
}

func (t *Thing) ShowMainScreen() {
	t.list2()
	t.win.SetContent(t.MainScreen)
	t.win.Show()
}

func (t *Thing) MakeMainScreen() *fyne.Container {
	t.bottom = widget.NewLabel("Bottom")

	var ti *time.Timer
	delay := 250 * time.Millisecond

	// search field
	s := widget.NewEntry()
	s.ActionItem = widget.NewIcon(theme.SearchIcon())
	s.OnChanged = func(text string) {
		log.Printf("search: %s", text)

		// clear timer regardless of state
		if ti != nil && !ti.Stop() {
			<-ti.C
		}
		ti = time.AfterFunc(delay, func() {
			log.Println("timer kicking...")

			t.find(text)

			ti = nil
		})

	}

	// button to add entry
	add := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		log.Println("add...")
		t.ShowAddScreen()
	})

	// top line with add button left and search on right
	top := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, add, s),
		add,
		s,
	)

	// show all entries
	t.list2()

	// inner table
	xtab := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, nil),
		t.table,
	)

	c := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(top, t.bottom, nil, nil),
		top,
		t.bottom,
		//widget.NewVScrollContainer(t.table),
		//t.table,
		xtab,
	)

	return c
}
