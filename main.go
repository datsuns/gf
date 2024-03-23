package main

import (
	"log/slog"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	logger *slog.Logger
)

func changePane(app *tview.Application, appCtx *App, side PaneSide) *Pane {
	appCtx.Current = side
	app.SetFocus(appCtx.CurPane().W)
	return appCtx.CurPane()
}

func saveConfig(appCtx *App, cfg *Config) {
	cfg.Body.LeftPath = appCtx.Pane(LeftPane).CurPath()
	cfg.Body.RightPath = appCtx.Pane(RightPane).CurPath()
	cfg.Save()
}

func enterIncSearch(_ *tview.Application, appCtx *App) {
	appCtx.Mode = IncSearch
	pane := appCtx.CurPane()
	pane.W.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorBlue))
}

func exitIncSearch(_ *tview.Application, appCtx *App) {
	appCtx.Mode = Normal
	pane := appCtx.CurPane()
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
			pane.SetItem(pane.ItemCount() - 1)
		}
	case tcell.KeyCtrlJ:
		m := tview.NewList().ShowSecondaryText(true)
		m.SetBorder(true)
		m.SetTitle("jump list")
		for name, path := range cfg.Body.JumpList {
			m.AddItem(name, path, 0, nil)
		}
		appCtx.JumpList = m
		appCtx.JumpSearch = ""
		appCtx.Mode = SelectJump
		app.SetRoot(m, false)
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
			if err := pane.Up(); err == nil {
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

func mainHandlerSelectJump(app *tview.Application, appCtx *App, _ *Config, event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEnter:
		n := appCtx.JumpList.GetCurrentItem()
		_, path := appCtx.JumpList.GetItemText(n)
		appCtx.CurPane().Jump(Path(path))
		appCtx.Mode = Normal
		app.SetRoot(appCtx.Root, false)
	case tcell.KeyESC:
		appCtx.Mode = Normal
		app.SetRoot(appCtx.Root, false)
	case tcell.KeyBS:
		if len(appCtx.JumpSearch) > 0 {
			appCtx.JumpSearch = appCtx.JumpSearch[:len(appCtx.JumpSearch)-1]
		} else {
			appCtx.JumpSearch = ""
		}
	case tcell.KeyRune:
		//appCtx.JumpSearch += string(event.Rune())
		//candidate := appCtx.JumpList.FindItems(appCtx.JumpSearch, "", false, true)
		//if len(candidate) > 0 {
		//	appCtx.JumpList.SetCurrentItem(candidate[0])
		//}
		n := appCtx.JumpList.GetCurrentItem()
		switch event.Rune() {
		case 'j':
			appCtx.JumpList.SetCurrentItem(n + 1)
		case 'k':
			if n > 0 {
				appCtx.JumpList.SetCurrentItem(n - 1)
			}
		}
	}
	return event
}

func mainHandler(app *tview.Application, appCtx *App, cfg *Config, event *tcell.EventKey) *tcell.EventKey {
	switch appCtx.Mode {
	case IncSearch:
		return mainHandlerIncSearch(app, appCtx, cfg, event)
	case SelectJump:
		return mainHandlerSelectJump(app, appCtx, cfg, event)
	default:
		return mainHandlerNormal(app, appCtx, cfg, event)
	}
}

func main() {
	var err error
	runlog, _ := os.Create("debug.txt")
	logger = slog.New(
		slog.NewTextHandler(runlog, nil),
	)

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

	appCtx.Root = flex
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return mainHandler(app, appCtx, cfg, event)
	})

	app.SetRoot(flex, true).SetFocus(flex)
	if err := app.Run(); err != nil {
		panic(err)
	}

}
