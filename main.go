package main

import (
	"log/slog"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	CofigFileName = "gf.toml"
	LogFileName   = "debug.txt"
)

var (
	logger *slog.Logger
)

func changePane(app *tview.Application, appCtx *App, side PaneSide) *Pane {
	appCtx.Current = side
	app.SetFocus(appCtx.CurPaneWidget())
	return appCtx.CurPane()
}

func saveConfig(appCtx *App, cfg *Config) {
	cfg.Body.LeftPath = appCtx.CurPath(LeftPane)
	cfg.Body.RightPath = appCtx.CurPath(RightPane)
	cfg.Save()
}

func openEditor(_ *tview.Application, appCtx *App, cfg *Config) {
	//logger.Info("openEditor", slog.Any("editor", cfg.Body.Editor), slog.Any("path", appCtx.CurPane().SelectedFullPath()))
	path := appCtx.CurPane().SelectedFullPath()
	ExecuteCommand(cfg.Body.Editor, path)
}

func enterIncSearch(_ *tview.Application, appCtx *App) {
	appCtx.Mode = IncSearch
	w := appCtx.CurPaneWidget()
	w.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorBlue))
}

func exitIncSearch(_ *tview.Application, appCtx *App) {
	appCtx.Mode = Normal
	w := appCtx.CurPaneWidget()
	w.SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite))
	appCtx.CurPane().ClearText()
}

// TODO imple
func createNewFile(app *tview.Application, appCtx *App, cfg *Config) {
}

// TODO imple
func createNewDirectory(app *tview.Application, appCtx *App, cfg *Config) {
}

func enterJumpListSelection(app *tview.Application, appCtx *App, cfg *Config) {
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
}

func scrollDown(pane *Pane, cfg *Config) {
	if pane.CurItem()+cfg.Body.ScrollLines < pane.ItemCount() {
		pane.SetItem(pane.CurItem() + cfg.Body.ScrollLines)
	} else {
		pane.SetItem(pane.ItemCount() - 1)
	}
}

func mainHandlerNormal(app *tview.Application, appCtx *App, cfg *Config, event *tcell.EventKey) *tcell.EventKey {
	pane := appCtx.CurPane()
	switch event.Key() {
	case tcell.KeyEnter:
		pane.Down()
	case tcell.KeyCtrlD:
		scrollDown(pane, cfg)
	case tcell.KeyCtrlE:
		createNewFile(app, appCtx, cfg)
	case tcell.KeyCtrlJ:
		enterJumpListSelection(app, appCtx, cfg)
	case tcell.KeyCtrlK:
		createNewDirectory(app, appCtx, cfg)
	case tcell.KeyCtrlU:
		pane.SetItem(pane.CurItem() - cfg.Body.ScrollLines)
	case tcell.KeyRune:
		switch event.Rune() {
		case 'e':
			openEditor(app, appCtx, cfg)
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
			pane.Up()
		case 'q':
			saveConfig(appCtx, cfg)
			app.Stop()
		}
	}
	return event
}

func updateIncSearch(pane *Pane) {
	candidate := pane.Find(pane.GetText())
	if len(candidate) > 0 {
		pane.SetItem(candidate[0])
	}
}

func mainHandlerIncSearch(app *tview.Application, appCtx *App, _ *Config, event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEnter:
		exitIncSearch(app, appCtx)
	case tcell.KeyESC:
		exitIncSearch(app, appCtx)
	case tcell.KeyBS:
		pane := appCtx.CurPane()
		pane.SetText(TrimLastOne(pane.GetText()))
		updateIncSearch(pane)
	case tcell.KeyRune:
		pane := appCtx.CurPane()
		pane.SetText(pane.GetText() + string(event.Rune()))
		updateIncSearch(pane)
	}
	return event
}

func backToNomal(app *tview.Application, appCtx *App) {
	appCtx.Mode = Normal
	app.SetRoot(appCtx.Root, false)
}

func selectJumpTarget(app *tview.Application, appCtx *App) {
	n := appCtx.JumpList.GetCurrentItem()
	_, path := appCtx.JumpList.GetItemText(n)
	appCtx.CurPane().Jump(Path(path))
	backToNomal(app, appCtx)
}

func mainHandlerSelectJump(app *tview.Application, appCtx *App, _ *Config, event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEnter:
		selectJumpTarget(app, appCtx)
	case tcell.KeyESC:
		backToNomal(app, appCtx)
	case tcell.KeyBS:
		appCtx.JumpSearch = TrimLastOne(appCtx.JumpSearch)
	case tcell.KeyRune:
		n := appCtx.JumpList.GetCurrentItem()
		switch event.Rune() {
		case 'j':
			appCtx.JumpList.SetCurrentItem(n + 1)
		case 'k':
			if n > 0 {
				appCtx.JumpList.SetCurrentItem(n - 1)
			}
		case 'q':
			backToNomal(app, appCtx)
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
	runlog, _ := os.Create(LogFileName)
	logger = slog.New(
		slog.NewTextHandler(runlog, nil),
	)

	cfg, err := LoadConfig(CofigFileName)
	if err != nil {
		panic(err)
	}
	appCtx, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}

	appCtx.Current = LeftPane
	appCtx.Reload(LeftPane)
	appCtx.Reload(RightPane)

	app := tview.NewApplication()

	flex := tview.NewFlex().
		AddItem(appCtx.PaneWidget(LeftPane), 0, 1, true).
		AddItem(appCtx.PaneWidget(RightPane), 0, 1, false)

	appCtx.Root = flex
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return mainHandler(app, appCtx, cfg, event)
	})

	app.SetRoot(flex, true).SetFocus(flex)
	if err := app.Run(); err != nil {
		panic(err)
	}

}
