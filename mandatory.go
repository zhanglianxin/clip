package main

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"fmt"
	"github.com/google/uuid"
	"go.mozilla.org/pkcs7"
	"howett.net/plist"
	"log"
	"time"
)

type consentText struct {
	Default string `plist:"default,omitempty" json:"default,omitempty"`
}

type mandatory struct {
	PayloadDisplayName   string       // required
	PayloadIdentifier    string       // required
	PayloadOrganization  string       `plist:",omitempty" json:",omitempty"`
	PayloadDescription   string       `plist:",omitempty" json:",omitempty"`
	ConsentText          *consentText `plist:",omitempty" json:",omitempty"`
	PayloadType          string
	PayloadUUID          uuid.UUID
	PayloadVersion       uint       // 1
	RemovalDate          *time.Time `plist:",omitempty" json:",omitempty" time:"2006-01-02T15:04:05Z"`
	DurationUntilRemoval uint       `plist:",omitempty" json:",omitempty"`
	PayloadContent       []*clip
}

func NewMandatory(name string, identifier string, org string, desc string,
	consent string, removalDate *time.Time, duration uint) *mandatory {
	if "" == name {
		name = "Untitled"
	}
	if "" == identifier {
		identifier = "clip." + uuid.NewString()
	}
	if nil == removalDate && 0 == duration {
		log.Fatalf("Either [RemovalDate] or [DurationUntilRemoval] is required\n")
	}

	m := &mandatory{
		PayloadDisplayName:  name,
		PayloadIdentifier:   identifier,
		PayloadOrganization: org,
		PayloadDescription:  desc,
		ConsentText:         nil,
		PayloadType:         "Configuration",
		PayloadUUID:         uuid.New(),
		PayloadVersion:      1,
	}
	if "" != consent {
		m.ConsentText = &consentText{Default: consent}
	}
	if nil != removalDate && !removalDate.IsZero() {
		m.RemovalDate = removalDate
	} else if 0 != duration {
		m.DurationUntilRemoval = duration
	} else {
		log.Fatalf("Neither [RemovalDate] nor [DurationUntilRemoval] is truthy\n")
	}

	return m
}

// Plist return unsigned plist in byte slice
func (m *mandatory) Plist(indent string) ([]byte, error) {
	if "" != indent {
		return plist.MarshalIndent(m, plist.XMLFormat, indent)
	} else {
		return plist.Marshal(m, plist.XMLFormat)
	}
}

func (m *mandatory) UnmarshalXml(content []byte) error {
	format, err := plist.Unmarshal(content, m)
	if nil != err {
		return err
	}
	if plist.XMLFormat != format {
		return fmt.Errorf("unmarshal in XML format error")
	}
	return nil
}

// SignAndDetach
//
//	ref: https://pkg.go.dev/go.mozilla.org/pkcs7#section-readme
func (m *mandatory) SignAndDetach(cert *x509.Certificate, parents []*x509.Certificate,
	privkey crypto.PrivateKey) (signed []byte, err error) {
	var content []byte
	content, err = m.Plist("\t")
	toBeSigned, err := pkcs7.NewSignedData(content)
	if err != nil {
		err = fmt.Errorf("Cannot initialize signed data: %s", err)
		return
	}
	if err = toBeSigned.AddSignerChain(cert, privkey, parents,
		pkcs7.SignerInfoConfig{}); err != nil {
		err = fmt.Errorf("Cannot add signer: %s", err)
		return
	}

	// Detach signature, omit if you want an embedded signature
	toBeSigned.Detach()

	signed, err = toBeSigned.Finish()
	if err != nil {
		err = fmt.Errorf("Cannot finish signing data: %s", err)
		return
	}

	// Verify the signature
	// pem.Encode(os.Stdout, &pem.Block{Type: "PKCS7", Bytes: signed})
	p7, err := pkcs7.Parse(signed)
	if err != nil {
		err = fmt.Errorf("Cannot parse our signed data: %s", err)
		return
	}

	// since the signature was detached, reattach the content here
	p7.Content = content

	if bytes.Compare(content, p7.Content) != 0 {
		err = fmt.Errorf("Our content was not in the parsed data:\n"+
			"\tExpected: %s\n\tActual: %s", content, p7.Content)
		return
	}
	if err = p7.Verify(); err != nil {
		err = fmt.Errorf("Cannot verify our signed data: %s", err)
		return
	}

	return signed, nil
}

// Sign no detach
func (m *mandatory) Sign(cert *x509.Certificate, parents []*x509.Certificate,
	privkey crypto.PrivateKey) (signed []byte, err error) {
	var content []byte
	content, err = m.Plist("\t")
	toBeSigned, err := pkcs7.NewSignedData(content)
	if nil != err {
		err = fmt.Errorf("cannot initialize signed data: %s", err)
		return
	}
	if err = toBeSigned.AddSignerChain(cert, privkey, parents,
		pkcs7.SignerInfoConfig{}); nil != err {
		err = fmt.Errorf("cannot add signer chain: %s", err)
		return
	}

	signed, err = toBeSigned.Finish()
	if nil != err {
		err = fmt.Errorf("cannot finish signing data: %s", err)
		return
	}

	// Verify the signature
	// pem.Encode(os.Stdout, &pem.Block{Type: "PKCS7", Bytes: signed})
	p7, err := pkcs7.Parse(signed)
	if nil != err {
		err = fmt.Errorf("cannot parse our signed data: %s", err)
		return
	}

	if 0 != bytes.Compare(content, p7.Content) {
		err = fmt.Errorf("our content was not in the parsed data:\n"+
			"\tExpected: %s\n\tActual: %s", content, p7.Content)
		return
	}
	truststore := x509.NewCertPool()
	for _, cert := range parents {
		truststore.AddCert(cert)
	}
	if err = p7.VerifyWithChain(truststore); nil != err {
		err = fmt.Errorf("cannot verify our signed data: %s", err)
		return
	}

	return signed, nil
}
