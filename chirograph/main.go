package main

import (
	"flag"
	"fmt"
	"github.com/tdewolff/canvas/pdf"
	"github.com/tdewolff/canvas/svg"
	"github.com/uncopied/tallystick"
	"log"
	"os"
	"strings"
)

func main() {
	outFile := flag.String("o", "out", "output file prefix, empty for stdout")
	outFormat := flag.String("f", "pdf", "output file format (svg or pdf)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `chirograph secured document in Go
https://github.com/uncopied/chirograph
Flags:
`)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Usage:
  	chirograph -o pdf "hello world"
`)
	}
	flag.Parse()

	content := "hello world!"
	if len(flag.Args()) == 0 {
		flag.Usage()
	} else {
		content = strings.Join(flag.Args(), " ")
	}
	t := chirograph.Chirograph{
		CertificateLabel:                content,
		PrimaryLinkURL:                  "PrimaryLinkURL",
		SecondaryLinkURL:                "SecondaryLinkURL",
		IssuerTokenURL:                  "IssuerTokenURL",
		OwnerTokenURL:                   "OwnerTokenURL",
		PrimaryAssetVerifierTokenURL:    "PrimaryAssetVerifierTokenURL",
		SecondaryAssetVerifierTokenURL:  "SecondaryAssetVerifierTokenURL",
		PrimaryOwnerVerifierTokenURL:    "PrimaryOwnerVerifierTokenURL",
		SecondaryOwnerVerifierTokenURL:  "SecondaryOwnerVerifierTokenURL",
		PrimaryIssuerVerifierTokenURL:   "PrimaryIssuerVerifierTokenURL",
		SecondaryIssuerVerifierTokenURL: "SecondaryIssuerVerifierTokenURL",
		MailToContentLeft:               "MailToContentLeft",
		MailToContentRight:              "MailToContentRight",
	}
	c := chirograph.Draw(&t)
	if *outFormat == "pdf" {
		c.WriteFile(*outFile+".pdf", pdf.Writer)
	} else if *outFormat == "svg" {
		c.WriteFile(*outFile+".svg", svg.Writer)
	} else {
		log.Fatal("unexpected format")
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

