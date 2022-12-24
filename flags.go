package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"time"
)

const megabyte = 1024 * 1024 // in bytes

var durationUntilRemoval uint  // mandatory
var removalDate time.Time      // mandatory
var payloadDisplayName string  // mandatory
var payloadDescription string  // mandatory
var payloadOrganization string // mandatory
var consent string             // mandatory

var needSign bool      // for sign
var certFile string    // for sign
var chainFile string   // for sign
var privkeyFile string // for sign

var flags = []cli.Flag{
	&cli.BoolFlag{
		Name:  "sign",
		Value: false,
		Usage: "Need signing",
		Action: func(ctx *cli.Context, b bool) error {
			needSign = b
			return nil
		},
	},
	&cli.StringFlag{
		Name:  "cert",
		Value: "certs/cert.pem",
		Usage: "ssl_certificate `FILE` for sign (in `PEM` format)",
		Action: func(ctx *cli.Context, v string) (err error) {
			if needSign && "" == v {
				log.Println("cert FILE specified, will use default")
			}
			return nil
		},
		Destination: &certFile,
	},
	&cli.StringFlag{
		Name:  "privkey",
		Value: "certs/privkey.pem",
		Usage: "ssl_certificate_key `FILE` for sign (in `PEM` format)",
		Action: func(ctx *cli.Context, v string) (err error) {
			if needSign && "" == v {
				log.Println("privkey FILE not specified, will use default")
			}
			return nil
		},
		Destination: &privkeyFile,
	},
	&cli.StringFlag{
		Name:  "chain",
		Value: "certs/chain.pem",
		Usage: "ssl_trusted_certificate `FILE` for sign (in `PEM` format)",
		Action: func(ctx *cli.Context, v string) (err error) {
			if needSign && "" == v {
				log.Println("chain FILE not specified, will use default")
			}
			return nil
		},
		Destination: &chainFile,
	},
	&cli.Float64Flag{
		Name:  "duration",
		Value: 0,
		Usage: "DurationUntilRemoval (in hours)", // mandatory
		Action: func(ctx *cli.Context, f float64) error {
			d, err := time.ParseDuration(fmt.Sprintf("%fh", f))
			if nil != err {
				return fmt.Errorf("parse duration value [%v] error", f)
			}
			durationUntilRemoval = uint(d.Seconds())
			return nil
		},
	},
	&cli.StringFlag{
		Name:  "date",
		Value: "",
		Usage: "RemovalDate (in format: 2006-01-02 15:04)", // mandatory
		Action: func(ctx *cli.Context, v string) error {
			var err error
			if removalDate, err = time.Parse("2006-01-02 15:04", v); nil != err {
				return fmt.Errorf("flag date value %v must be in format `2006-01-02 15:04`", v)
			}
			return nil
		},
	},
	/*&cli.TimestampFlag{
		Name:   "date",
		Value:  nil,
		Usage:  "RemovalDate (in format: 2006-01-02T15:04:05Z)", // mandatory
		Layout: isoFormat,
		Action: func(ctx *cli.Context, t *time.Time) error {
			fmt.Println(t)
			return nil
		},
	},*/
	&cli.StringFlag{
		Name:        "name",
		Value:       "Untitled",
		Usage:       "PayloadDisplayName", // mandatory
		Destination: &payloadDisplayName,
	},
	&cli.StringFlag{
		Name:        "desc",
		Value:       "Web Clip Installation",
		Usage:       "PayloadDescription", // mandatory
		Destination: &payloadDescription,
	},
	&cli.StringFlag{
		Name:        "org",
		Value:       "Excellent Inc.",
		Usage:       "PayloadOrganization", // mandatory
		Destination: &payloadOrganization,
	},
	&cli.StringFlag{
		Name:        "consent",
		Value:       "It will appear on the main screen.",
		Usage:       "ConsentText", // mandatory
		Destination: &consent,
	},
	/*&cli.StringFlag{
		Name:        "label",
		Value:       "",
		Usage:       "Web Clip Label (*required)",
		Destination: &clipLabel,
		Required:    true,
	},
	&cli.StringFlag{
		Name:  "url",
		Value: "",
		Usage: "Web Clip URL (*required)",
		Action: func(ctx *cli.Context, v string) error {
			if _, err := url.ParseRequestURI(v); nil != err {
				return fmt.Errorf("flag url value %v is not valid\n", v)
			}
			return nil
		},
		Destination: &clipUrl,
		Required:    true,
	},
	&cli.StringFlag{
		Name:  "icon",
		Value: "",
		Usage: "image `FILE` for the clip icon (400 x 400 pixels, size < 1MB, \n" +
			"\tand in `PNG` format for best)",
		Action: func(ctx *cli.Context, v string) error {
			var err error
			if iconBytes, err = os.ReadFile(v); nil != err {
				return err
			}
			if len(iconBytes) > megabyte {
				return fmt.Errorf("flag icon file %v size must less than 1MB\n", v)
			}
			return nil
		},
	},*/
}
