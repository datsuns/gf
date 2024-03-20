package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApp(t *testing.T) {
	cur, _ := os.Getwd()
	p := Path(filepath.Join(cur, "gf.toml"))
	c := &Config{Path: p, Body: &ConfigEntry{LeftPath: Path(cur), RightPath: Path(cur)}}
	a, e := NewApp(c)
	if e != nil {
		t.Fatalf("App constructor error %s. path[%v]", e, cur)
	}
	if a.Current != LeftPane {
		t.Fatalf("unexpected default pane %v != %v", LeftPane, a.Current)
	}
}
