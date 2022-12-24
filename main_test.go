package main

import (
	"crypto/md5"
	"fmt"
	"github.com/google/uuid"
	"howett.net/plist"
	_ "howett.net/plist"
	"os"
	"testing"
	"time"
)

func TestPlist(t *testing.T) {
	var err error
	var bs []byte
	bs, err = os.ReadFile("open_graph_logo.png")
	if nil != err {
		t.Fatal(err)
	}
	content := []*clip{NewClip("", bs, "https://www.google.com")}

	m := &mandatory{
		PayloadDisplayName: "Untitled",
		PayloadIdentifier:  "clip." + uuid.NewString(),
		PayloadType:        "Configuration",
		PayloadUUID:        uuid.New(),
		PayloadVersion:     1,
		PayloadContent:     content,
	}
	p, err := plist.MarshalIndent(m, plist.XMLFormat, "\t")
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("%s", p)
}

func TestToBase64(t *testing.T) {
	bs, err := os.ReadFile("open_graph_logo.png")
	if nil != err {
		t.Fatal(err)
	}
	str := toBase64(bs)
	sum := md5.Sum(bs)
	expected := "293abc9b5a3d79ed81c8f75f20ee7337"

	if expected != fmt.Sprintf("%x", sum) {
		t.Errorf("EXPECTED: %v, ACTUAL: %v", expected, sum)
	}
	t.Logf("BASE64: %s", str)

	var ti time.Time
	ti, err = time.Parse("2006-01-02 15:04:05", "2023-03-14 03:45:47")
	if nil != err {
		t.Fatal(err)
	}
	t.Log(ti.Format("2006-01-02T15:04:05Z"))
}
