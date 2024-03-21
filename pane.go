package main

import (
	"github.com/gizak/termui/v3/widgets"
)

type PaneSide int

const (
	LeftPane = iota
	RightPane
)

type Pane struct {
	Dir    *Dir
	Widget *widgets.List
}

func NewPane(path Path) (*Pane, error) {
	d, err := NewDir(path)
	if err != nil {
		return nil, err
	}
	return &Pane{Dir: d}, nil
}

func (p *Pane) Cur() string {
	return p.Dir.Cur()
}

func (p *Pane) Selected() string {
	return p.Widget.Rows[p.Widget.SelectedRow]
}

func (p *Pane) DirSelected() bool {
	if len(p.Widget.Rows) == 0 {
		return false
	}
	return p.Selected()[len(p.Selected())-1] == '/'
}

func (p *Pane) Reload() {
	p.Widget.Title = p.Dir.Cur()
	p.Widget.Rows = []string{}
	for _, f := range p.Dir.Entries {
		if f.IsDir() {
			p.Widget.Rows = append(p.Widget.Rows, f.Name()+"/")
		} else {
			p.Widget.Rows = append(p.Widget.Rows, f.Name())
		}
	}
}
