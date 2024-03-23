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
	T   *tview.TextView
}

func NewPane(path Path) (*Pane, error) {
	d, err := NewDir(path)
	if err != nil {
		return nil, err
	}
	w := tview.NewList().ShowSecondaryText(false)
	w.SetBorder(true)
	t := tview.NewTextView()
	t.SetBorder(true)
	return &Pane{
		Dir: d,
		W:   w,
		T:   t,
	}, nil
}

func (p *Pane) Cur() string {
	return p.Dir.Cur()
}

func (p *Pane) CurPath() Path {
	return Path(p.Dir.Cur())
}

func (p *Pane) CurItem() int {
	return p.W.GetCurrentItem()
}

func (p *Pane) ItemCount() int {
	return p.W.GetItemCount()
}

func (p *Pane) SetItem(n int) *tview.List {
	return p.W.SetCurrentItem(n)
}

func (p *Pane) Find(s string) []int {
	return p.W.FindItems(s, "", false, true)
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

func (p *Pane) Up() error {
	return p.Dir.Up()
}

func (p *Pane) Down() error {
	return p.Dir.Down(Path(p.Selected()))
}

func (p *Pane) Jump(path Path) error {
	d, err := NewDir(path)
	if err != nil {
		return err
	}
	p.Dir = d
	p.Reload()
	return nil
}

func (p *Pane) GetText() string {
	return p.T.GetText(false)
}

func (p *Pane) SetText(s string) *tview.TextView {
	return p.T.SetText(s)
}
