package main

import (
	"github.com/rivo/tview"
)

type PaneSide int

const (
	LeftPane = iota
	RightPane
)

type Pane struct {
	Dir *Dir
	W   *tview.List
}

func NewPane(path Path) (*Pane, error) {
	d, err := NewDir(path)
	if err != nil {
		return nil, err
	}
	w := tview.NewList().ShowSecondaryText(false)
	w.SetBorder(true)
	return &Pane{
		Dir: d,
		W:   w,
	}, nil
}

func (p *Pane) Cur() string {
	return p.Dir.Cur()
}

func (p *Pane) Selected() string {
	main, _ := p.W.GetItemText(p.W.GetCurrentItem())
	return main
}

func (p *Pane) DirSelected() bool {
	if p.W.GetItemCount() == 0 {
		return false
	}
	return p.Selected()[len(p.Selected())-1] == '/'
}

func (p *Pane) Reload() {
	p.W.Clear()
	p.W.SetTitle(p.Dir.Cur())

	for _, f := range p.Dir.Entries {
		if f.IsDir() {
			p.W.AddItem(f.Name()+"/", "", 0, nil)
		} else {
			p.W.AddItem(f.Name(), "", 0, nil)
		}
	}
}
