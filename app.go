package main

import (
	"github.com/cockroachdb/errors"
)

type AppMode int

const (
	Normal = iota
	IncSearch
)

type App struct {
	Current PaneSide
	Mode    AppMode
	Panes   []*Pane
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
