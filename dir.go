package main

import (
	"errors"
	"os"
	"path/filepath"
)

type Path string

type Dir struct {
	Pwd     Path
	Entries []os.FileInfo
}

func NewDir(path Path) (*Dir, error) {
	var err error
	abs, err := filepath.Abs(string(path))
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(abs); err != nil {
		return nil, err
	}
	ret := &Dir{Pwd: path, Entries: []os.FileInfo{}}
	if err := ret.Glob(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (d *Dir) Cur() string {
	return string(d.Pwd)
}

func (d *Dir) Up() error {
	p := filepath.Dir(d.Cur())
	if _, err := os.Stat(p); err != nil {
		return err
	}
	d.Pwd = Path(p)
	return d.Glob()
}

func (d *Dir) Down(child Path) error {
	p := filepath.Join(d.Cur(), string(child))
	if f, err := os.Stat(p); err != nil {
		return err
	} else {
		if f.IsDir() == false {
			return errors.New("ChangeDirError not dir : " + string(child))
		}
	}
	d.Pwd = Path(p)
	return d.Glob()
}

func (d *Dir) Glob() error {
	latched := []os.FileInfo{}
	found, err := filepath.Glob(filepath.Join(d.Cur(), "*"))
	if err != nil {
		return err
	}
	for _, f := range found {
		s, err := os.Stat(f)
		if err != nil {
			return err
		}
		latched = append(latched, s)
	}
	d.Entries = latched
	return nil
}
