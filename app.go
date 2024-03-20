package main

type App struct {
	Current PaneSide
	Panes   []*Pane
}

func NewApp(c *Config) (*App, error) {
	var e error
	left, e := NewPane(c.LeftPath)
	if e != nil {
		return nil, e
	}
	right, e := NewPane(c.RightPath)
	if e != nil {
		return nil, e
	}
	panes := make([]*Pane, 2, 2)
	panes[LeftPane] = left
	panes[RightPane] = right
	return &App{Current: LeftPane, Panes: panes}, nil
}
