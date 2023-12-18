package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zoroqi/snippet/store"
	"strings"
)

func findOnlySnippet(snippets []store.Snippet) (store.Snippet, error) {
	if len(snippets) == 0 {
		return store.Snippet{}, ERR_FIND_NO_SCRIPT
	}
	if len(snippets) == 1 {
		return snippets[0], nil
	}
	snippet, err := filterSnippetsWithTui(snippets)
	if err != nil {
		return store.Snippet{}, err
	}
	if len(snippet) != 1 {
		if len(snippet) == 0 {
			return store.Snippet{}, ERR_FIND_NO_SCRIPT
		}
		return store.Snippet{}, ERR_FIND_MULTI_SCRIPT
	}
	return snippet[0], nil
}

func filterSnippetsWithTui(snippets []store.Snippet) ([]store.Snippet, error) {
	app := tview.NewApplication()

	list := tview.NewList()
	list.SetTitle("List").SetBorder(true)

	res := []store.Snippet{}
	addItem := func(t store.Snippet) {
		list.AddItem(t.Name, t.Description, ' ', func() {
			res = append(res, t)
			app.Stop()
		})
	}
	nameIndex := map[string][]store.Snippet{}
	for _, r := range snippets {
		nameIndex[r.Name] = append(nameIndex[r.Name], r)
	}

	for _, r := range snippets {
		addItem(r)
	}

	input := tview.NewInputField()
	input.SetLabel("key: ").
		SetBorder(true).
		SetTitle("Search")

	input.SetChangedFunc(func(t string) {
		list.Clear()
		for _, r := range snippets {
			if strings.Contains(strings.ToLower(r.Name+" "+r.Description), t) {
				addItem(r)
			}
		}
	})

	// jump list
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlJ:
			app.SetFocus(list)
			return nil
		default:
			return event
		}
	})
	// jumpu input and select list
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlK:
			app.SetFocus(input)
			return nil
		default:
			switch event.Rune() {
			case 'h':
				return tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone)
			case 'j':
				return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
			case 'k':
				return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
			case 'l':
				return tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone)
			default:
				return event
			}
		}
	})

	input.SetDoneFunc(func(key tcell.Key) {
		if tcell.KeyEnter == key {
			for i := 0; i < list.GetItemCount(); i++ {
				name, _ := list.GetItemText(i)
				res = append(res, nameIndex[name]...)
			}
			app.Stop()
		}
	})

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(input, 5, 2, true).
			AddItem(list, 0, 3, false),
			0, 2, false)

	if err := app.SetRoot(flex, true).SetFocus(input).Run(); err != nil {
		return nil, err
	}
	return res, nil
}
