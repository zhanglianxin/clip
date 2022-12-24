package main

import (
	"encoding/base64"
	"github.com/urfave/cli/v2"
	_ "howett.net/plist"
	"log"
	"os"
)

func toBase64(bs []byte) string {
	return base64.StdEncoding.EncodeToString(bs)
}

func main() {
	app := &cli.App{
		Name:  "clip",
		Usage: "make an excellent clip",
		Flags: flags,
		UsageText: "clip [global options] command [command options] [arguments...]\n" +
			"\t arguments: just place your args [label url icon]... order by order",
		Action: biz,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
