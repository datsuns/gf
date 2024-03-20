package main

type PaneSide int

const (
	LeftPane = iota
	RightPane
)

type Pane struct {
	Dir *Dir
}

func NewPane(path Path) (*Pane, error) {
	d, err := NewDir(path)
	if err != nil {
		return nil, err
	}
	return &Pane{Dir: d}, nil
}
