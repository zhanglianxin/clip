package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
	"testing"
	"time"
)

func TestNewMandatory(t *testing.T) {
	// m := NewMandatory("", "", time.Now().Add(24*time.Hour), 0)
	m := NewMandatory("", "", "Hello Clip", "Web Clip Installation",
		"", nil, uint(time.Now().Add(24*time.Hour).Unix()))
	p, err := m.Plist("\t")
	if nil != err {
		t.Fatal(err)
	}
	t.Log(string(p))
}

func TestMandatory_Plist(t *testing.T) {
	bs, err := os.ReadFile("test.jpg")
	if nil != err {
		t.Fatal(err)
	}
	c := NewClip("test", nil, "https://www.google.com")
	ti := time.Now().Add(24 * time.Hour)
	m := NewMandatory("", "", "Apple Inc.", "Web Clip Install",
		"It will appear on the main screen", &ti, 0)
	clips := []*clip{c}
	m.PayloadContent = clips

	bs, err = json.MarshalIndent(m, "", "\t")
	if nil != err {
		t.Fatal(err)
	}
	t.Log(string(bs))

	unsigned, err := m.Plist("\t")
	if nil != err {
		t.Fatal(err)
	}
	os.WriteFile("test.mobileconfig", unsigned, os.ModePerm)
}

func TestMandatory_UnmarshalXml(t *testing.T) {
	content, err := os.ReadFile("test.mobileconfig")
	if nil != err {
		t.Fatal(err)
	}
	var m mandatory
	if err := m.UnmarshalXml(content); nil != err {
		t.Fatal(err)
	}
	t.Logf("%#v", m)
}

func TestMandatory_Sign(t *testing.T) {
	m := NewMandatory("", "", "Hello Clip", "Web Clip Installation",
		"", nil, uint(time.Now().Add(24*time.Hour).Unix()))
	fullChainBs, _ := os.ReadFile("certs/fullchain.pem")
	chainBs, _ := os.ReadFile("certs/chain.pem")
	privBs, _ := os.ReadFile("certs/privkey.pem")
	certBs, _ := os.ReadFile("certs/cert.pem")

	fullChainBlock, _ := pem.Decode(fullChainBs)
	fullChain, err := x509.ParseCertificate(fullChainBlock.Bytes)
	if nil != err {
		t.Fatal(err)
	}
	chainBlock, _ := pem.Decode(chainBs)
	chain, err := x509.ParseCertificates(chainBlock.Bytes)
	if nil != err {
		t.Fatal(err)
	}
	privkBlock, _ := pem.Decode(privBs)
	privk, err := x509.ParsePKCS8PrivateKey(privkBlock.Bytes)
	if nil != err {
		t.Fatal(err)
	}

	certBlock, _ := pem.Decode(certBs)
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if nil != err {
		t.Fatal(err)
	}
	t.Log(fullChain.Issuer)
	t.Log(cert.Issuer)

	signed, err := m.Sign(cert, chain, privk)
	if nil != err {
		t.Fatal(err)
	}

	t.Log("succeed, file size:", len(signed))
}
