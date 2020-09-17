package main

import (
	"log"
	"time"

	"fyne.io/fyne"
	//"fyne.io/fyne/app"
	//"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/theme"

        "golang.org/x/text/language"
        "golang.org/x/text/search"
)

func (t *Thing) find(term string) {
	if term == "" {
		t.list()
		return
	}

	matcher := search.New(language.English, search.Loose, search.IgnoreCase, search.IgnoreDiacritics)
	pattern := matcher.CompileString(term)

	rows := make([][]string, 0)
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

func (t *Thing) MainScreen() {
	bottom := widget.NewLabel("Bottom")

	var ti *time.Timer
	delay := 250 * time.Millisecond

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


	add := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		log.Println("add...")
	})

	top := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, add, s),
		add,
		s,
	)

	t.list()

	xtab := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, nil),
		t.table,
	)


	c := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(top, bottom, nil, nil),
		top,
		bottom,
		//widget.NewVScrollContainer(t.table),
		//t.table,
		xtab,
	)

	t.win.SetContent(c)
}

