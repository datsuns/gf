package main

import (
	"os"
	"testing"
)

func TestApp(t *testing.T) {
	cur, _ := os.Getwd()
	c := &Config{LeftPath: Path(cur), RightPath: Path(cur)}
	a, e := NewApp(c)
	if e != nil {
		t.Fatalf("App constructor error %s. path[%v]", e.Error(), cur)
	}
	if a.Current != LeftPane {
		t.Fatalf("unexpected default pane %v != %v", LeftPane, a.Current)
	}
}
