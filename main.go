package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func changePane(app *tview.Application, ctx *App, side PaneSide) *Pane {
	ctx.Current = side
	app.SetFocus(ctx.Panes[ctx.Current].W)
	return ctx.Panes[ctx.Current]
}

func saveConfig(appCtx *App, cfg *Config) {
	cfg.Body.LeftPath = appCtx.Panes[LeftPane].CurPath()
	cfg.Body.RightPath = appCtx.Panes[RightPane].CurPath()
	cfg.Save()
}

func enterIncSearch(app *tview.Application, appCtx *App) {
	appCtx.Mode = IncSearch
}

func mainHandlerNormal(app *tview.Application, appCtx *App, cfg *Config, event *tcell.EventKey) *tcell.EventKey {
	pane := appCtx.Panes[appCtx.Current]
	switch event.Key() {
	case tcell.KeyEnter:
		if err := pane.Dir.Down(Path(pane.Selected())); err == nil {
			pane.Reload()
		}
	case tcell.KeyRune:
		switch event.Rune() {
		case 'f':
			enterIncSearch(app, appCtx)
		case 'h':
			pane = changePane(app, appCtx, LeftPane)
		case 'j':
			pane.W.SetCurrentItem(pane.W.GetCurrentItem() + 1)
		case 'k':
			if pane.W.GetCurrentItem() > 0 {
				pane.W.SetCurrentItem(pane.W.GetCurrentItem() - 1)
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

func mainHandlerIncSearch(app *tview.Application, appCtx *App, cfg *Config, event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEnter:
		appCtx.Mode = Normal
	case tcell.KeyESC:
		appCtx.Mode = Normal
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
		AddItem(appCtx.Panes[LeftPane].W, 0, 1, true).
		AddItem(appCtx.Panes[RightPane].W, 0, 1, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return mainHandler(app, appCtx, cfg, event)
	})

	app.SetRoot(flex, true).SetFocus(flex)
	if err := app.Run(); err != nil {
		panic(err)
	}

}
