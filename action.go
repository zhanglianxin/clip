package main

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/urfave/cli/v2"
	"net/url"
	"os"
)

func biz(cCtx *cli.Context) (err error) {
	var cert *x509.Certificate
	var chain []*x509.Certificate
	var privkey crypto.PrivateKey
	if cCtx.NArg() > 0 {
		if needSign {
			{
				if _, err = os.Stat(certFile); nil != err {
					return err
				}
				var certBs []byte
				certBs, err = os.ReadFile(certFile)
				if nil != err {
					return err
				}
				var certBlock *pem.Block
				certBlock, _ = pem.Decode(certBs)
				cert, err = x509.ParseCertificate(certBlock.Bytes)
				if nil != err {
					return err
				}
				if nil == cert {
					return fmt.Errorf("cert is required when need siging")
				}
			}
			{
				if _, err = os.Stat(chainFile); nil != err {
					return err
				}
				var chainBs []byte
				chainBs, err = os.ReadFile(chainFile)
				if nil != err {
					return err
				}
				var chainBlock *pem.Block
				chainBlock, _ = pem.Decode(chainBs)
				chain, err = x509.ParseCertificates(chainBlock.Bytes)
				if nil != err {
					return err
				}
				if nil == chain {
					return fmt.Errorf("chain is required when need siging")
				}
			}
			{
				if _, err = os.Stat(privkeyFile); nil != err {
					return err
				}
				var privBs []byte
				privBs, err = os.ReadFile(privkeyFile)
				if nil != err {
					return err
				}
				var privkBlock *pem.Block
				privkBlock, _ = pem.Decode(privBs)
				privkey, err = x509.ParsePKCS8PrivateKey(privkBlock.Bytes)
				if nil != err {
					return err
				}
				if nil == privkey {
					return fmt.Errorf("privkey is required when need siging")
				}
			}
		}
		argsLen := cCtx.Args().Len()
		// check args group by three
		grps := argsLen / 3
		clips := make([]*clip, grps)
		for g := 0; g < grps; g++ {
			c := NewClip("", nil, "")
			for i := 0; i < argsLen; i++ {
				arg := cCtx.Args().Get(i)
				switch i % 3 {
				case 0:
					if "" == arg {
						return fmt.Errorf("arg label is required")
					}
					c.Label = arg
				case 1:
					if "" == arg {
						return fmt.Errorf("arg url is required")
					}
					if _, err := url.ParseRequestURI(arg); nil != err {
						fmt.Println(err)
						return fmt.Errorf("url value [%v] is not valid\n", arg)
					}
					c.URL = arg
				case 2:
					var iconBytes []byte
					if "" != arg {
						var err error
						if iconBytes, err = os.ReadFile(arg); nil != err {
							return fmt.Errorf("read %s error: %v", arg, err)
						}
						if len(iconBytes) > megabyte {
							return fmt.Errorf("flag icon file [%v] size must less than 1MB\n", arg)
						}
					}
					c.Icon = iconBytes
				}
				clips[g] = c
			}
		}
		m := NewMandatory(payloadDisplayName, "", payloadOrganization,
			payloadDescription, consent, &removalDate, durationUntilRemoval)
		m.PayloadContent = clips

		unsignedBs, err := m.Plist("\t")
		if nil != err {
			return fmt.Errorf("marshal error [%v]", err)
		}

		fname := fmt.Sprintf("%s_unsigned.mobileconfig", payloadDisplayName)
		if err := os.WriteFile(fname, unsignedBs, os.ModePerm); nil != err {
			return fmt.Errorf("write to %s error: %v", fname, err)
		}

		if needSign {
			signed, err := m.Sign(cert, chain, privkey)
			if nil != err {
				return err
			}
			fname := fmt.Sprintf("%s_signed.mobileconfig", payloadDisplayName)
			if err := os.WriteFile(fname, signed, os.ModePerm); nil != err {
				return fmt.Errorf("write to %s error: %v", fname, err)
			}
		}
	}

	return nil
}
