package main

import (
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func genPane(app *App, side PaneSide) *widgets.List {
	w := widgets.NewList()
	p := app.Panes[side]
	p.Widget = w
	p.Reload()

	if side == LeftPane {
		w.SetRect(0, 00, 60, 30)
	} else {
		w.SetRect(60, 00, 120, 30)
	}
	return w
}

func rollbackSelected(w *widgets.List) {
	w.SelectedRowStyle = ui.NewStyle(ui.ColorWhite)
	ui.Render(w)
}

func changeTo(app *App, side PaneSide) *Pane {
	app.Current = side
	ret := app.Panes[app.Current]
	ret.Widget.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)
	ui.Render(ret.Widget)
	return ret
}

func mainHandler(app *App) {
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		pane := app.Panes[app.Current]
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			pane.Widget.ScrollDown()
		case "k", "<Up>":
			pane.Widget.ScrollUp()
		case "l":
			rollbackSelected(pane.Widget)
			pane = changeTo(app, RightPane)
		case "h":
			rollbackSelected(pane.Widget)
			pane = changeTo(app, LeftPane)
		case "u":
			if err := pane.Dir.Up(); err == nil {
				pane.Widget.SelectedRow = 0
				pane.Reload()
			}
			ui.Render(pane.Widget)
		case "<Enter>":
			if err := pane.Dir.Down(Path(pane.Selected())); err == nil {
				pane.Widget.SelectedRow = 0
				pane.Reload()
			}
			ui.Render(pane.Widget)
		}

		if pane.DirSelected() {
			pane.Widget.SelectedRowStyle = ui.NewStyle(ui.ColorCyan)
		} else {
			pane.Widget.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)
		}
		ui.Render(pane.Widget)
	}
}

func main() {
	var err error
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	cfg, err := LoadConfig("gf.toml")
	if err != nil {
		panic(err)
	}
	app, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}

	app.Current = LeftPane
	genPane(app, LeftPane)
	genPane(app, RightPane)
	pane := app.Panes[app.Current]
	if pane.DirSelected() {
		pane.Widget.SelectedRowStyle = ui.NewStyle(ui.ColorCyan)
	} else {
		pane.Widget.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)
	}

	ui.Render(app.Panes[LeftPane].Widget)
	ui.Render(app.Panes[RightPane].Widget)
	mainHandler(app)
}
