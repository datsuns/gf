package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func changePane(app *tview.Application, appCtx *App, side PaneSide) *Pane {
	appCtx.Current = side
	app.SetFocus(appCtx.Panes[appCtx.Current].W)
	return appCtx.Panes[appCtx.Current]
}

func saveConfig(appCtx *App, cfg *Config) {
	cfg.Body.LeftPath = appCtx.Panes[LeftPane].CurPath()
	cfg.Body.RightPath = appCtx.Panes[RightPane].CurPath()
	cfg.Save()
}

func enterIncSearch(_ *tview.Application, appCtx *App) {
	appCtx.Mode = IncSearch
	pane := appCtx.Panes[appCtx.Current]
	pane.W.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorBlue))
}

func exitIncSearch(_ *tview.Application, appCtx *App) {
	appCtx.Mode = Normal
	pane := appCtx.Panes[appCtx.Current]
	pane.W.SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite))
	pane.T.Clear()
}

func mainHandlerNormal(app *tview.Application, appCtx *App, cfg *Config, event *tcell.EventKey) *tcell.EventKey {
	pane := appCtx.CurPane()
	switch event.Key() {
	case tcell.KeyEnter:
		if err := pane.Down(); err == nil {
			pane.Reload()
		}
	case tcell.KeyCtrlD:
		if pane.CurItem()+cfg.Body.ScrollLines < pane.ItemCount() {
			pane.SetItem(pane.CurItem() + cfg.Body.ScrollLines)
		} else {
			pane.SetItem(pane.W.GetItemCount() - 1)
		}
	case tcell.KeyCtrlU:
		pane.SetItem(pane.CurItem() - cfg.Body.ScrollLines)
	case tcell.KeyRune:
		switch event.Rune() {
		case 'f':
			enterIncSearch(app, appCtx)
		case 'h':
			pane = changePane(app, appCtx, LeftPane)
		case 'j':
			pane.SetItem(pane.CurItem() + 1)
		case 'k':
			if pane.CurItem() > 0 {
				pane.SetItem(pane.CurItem() - 1)
			}
		case 'l':
			pane = changePane(app, appCtx, RightPane)
		case 'u':
			if err := pane.Dir.Up(); err == nil {
				pane.Reload()
			}
		case 'q':
			saveConfig(appCtx, cfg)
			app.Stop()
		}
	}
	return event
}

func mainHandlerIncSearch(app *tview.Application, appCtx *App, _ *Config, event *tcell.EventKey) *tcell.EventKey {
	pane := appCtx.CurPane()
	switch event.Key() {
	case tcell.KeyEnter:
		exitIncSearch(app, appCtx)
	case tcell.KeyESC:
		exitIncSearch(app, appCtx)
	case tcell.KeyRune:
		pane.SetText(pane.GetText() + string(event.Rune()))
		candidate := pane.Find(pane.GetText())
		if len(candidate) > 0 {
			pane.SetItem(candidate[0])
		}
	}
	return event
}

func mainHandler(app *tview.Application, appCtx *App, cfg *Config, event *tcell.EventKey) *tcell.EventKey {
	if appCtx.Mode == IncSearch {
		return mainHandlerIncSearch(app, appCtx, cfg, event)
	}
	return mainHandlerNormal(app, appCtx, cfg, event)
}

func main() {
	var err error

	cfg, err := LoadConfig("gf.toml")
	if err != nil {
		panic(err)
	}
	appCtx, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}

	appCtx.Current = LeftPane
	appCtx.Panes[LeftPane].Reload()
	appCtx.Panes[RightPane].Reload()

	app := tview.NewApplication()

	flex := tview.NewFlex().
		AddItem(appCtx.Pane(LeftPane).W, 0, 1, true).
		AddItem(appCtx.Pane(RightPane).W, 0, 1, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return mainHandler(app, appCtx, cfg, event)
	})

	app.SetRoot(flex, true).SetFocus(flex)
	if err := app.Run(); err != nil {
		panic(err)
	}

}
