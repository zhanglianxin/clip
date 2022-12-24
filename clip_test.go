package main

import (
	"howett.net/plist"
	"os"
	"testing"
)

func TestNewClip(t *testing.T) {
	bs, err := os.ReadFile("test.jpg")
	if nil != err {
		t.Fatal(err)
	}
	c := NewClip("test", bs, "https://www.google.com")
	p, err := plist.MarshalIndent(c, plist.XMLFormat, "\t")
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("%s", p)
}
