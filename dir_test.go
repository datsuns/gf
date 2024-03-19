package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDir_simple(t *testing.T) {
	var err error
	_, err = NewDir("abc")
	if err == nil {
		t.Fatalf("path [abc] should not exists")
	}

	dir, _ := os.Getwd()
	d, err := NewDir(Path(dir))
	if err != nil {
		t.Fatalf("path [.] should exists. [%v]", err)
	}
	if d.Cur() != dir {
		t.Fatalf("path should be [%v] but [%v]", dir, d.Cur())
	}
	if len(d.Entries) == 0 {
		t.Fatalf("entry should exists [%v]", d.Entries)
	}
}

func TestDir_Up(t *testing.T) {
	dir, _ := os.Getwd()
	d, _ := NewDir(Path(dir))

	up := filepath.Dir(dir)
	err := d.Up()
	if err != nil {
		t.Fatalf("path [%v] should exists. [%v]", up, err)
	}
}
