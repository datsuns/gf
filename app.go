package main

import (
	"github.com/cockroachdb/errors"
	"github.com/rivo/tview"
)

type AppMode int

const (
	Normal = iota
	IncSearch
	SelectJump
	CreateNewFile
	CreateNewDirectory
	Rename
)

type App struct {
	Current        PaneSide
	Mode           AppMode
	Panes          []*Pane
	Root           *tview.Flex
	JumpSearch     string
	JumpList       *tview.List
	CreateCandiate *tview.InputField
	ErrorInfo      *tview.TextView
}

func NewApp(c *Config) (*App, error) {
	var e error
	left, e := NewPane(c.LeftPath())
	if e != nil {
		return nil, errors.Wrap(e, "NewPane(left)")
	}
	right, e := NewPane(c.RightPath())
	if e != nil {
		return nil, errors.Wrap(e, "NewPane(right)")
	}
	panes := make([]*Pane, 2, 2)
	panes[LeftPane] = left
	panes[RightPane] = right
	return &App{Current: LeftPane, Mode: Normal, Panes: panes}, nil
}

func (a *App) Pane(side PaneSide) *Pane {
	return a.Panes[side]
}

func (a *App) PaneWidget(side PaneSide) *tview.List {
	return a.Pane(side).W
}

func (a *App) Reload(side PaneSide) {
	a.Pane(side).Reload()
}

func (a *App) CurPane() *Pane {
	return a.Pane(a.Current)
}

func (a *App) CurPath(side PaneSide) Path {
	return a.Pane(side).CurPath()
}

func (a *App) CurPaneWidget() *tview.List {
	return a.Pane(a.Current).W
}
