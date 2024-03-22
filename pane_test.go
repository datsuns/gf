package main

import (
	"os"
	"testing"
)

func TestPane(t *testing.T) {
	pwd, _ := os.Getwd()
	NewPane(Path(pwd))
}
