package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	cur, _ := os.Getwd()
	p := Path(filepath.Join(cur, "gf.toml"))
	c, e := LoadConfig(p)
	if e != nil {
		t.Fatalf("config load error %v", e)
	}
	if c.Path != p {
		t.Fatalf("invalid config path [%v]-[%v]", p, c.Path)
	}
}
