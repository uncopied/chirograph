package main

import (
	"flag"
	"fmt"
	"github.com/tdewolff/canvas/pdf"
	"github.com/tdewolff/canvas/svg"
	"github.com/uncopied/chirograph"
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
		TopHelper: "This is a chirograph by UNCOPIED with 5 parts (CENTER,TOP,BOTTOM,LEFT,RIGHT):\n-1 Print it on quality 100+ gram paper (ex. C by Clairfontaine recommended https://uncopied.art/blog/papers or other ISO-9706).\n-2 Cut it using scissors into 5 parts along the blue lines (slight imperfections are OK, make it unique and tamper-proof).\n-3 Optionally, add your own identical physical mark or signature to each of the 5 parts (without corrupting the QR Codes).\n-4 Glue its CENTER part to the back of your painting or limited edition print (ex. 3M recommended https://uncopied.art/blog/glues).",
		LeftHelper:  "Please, do not forget to send this LEFT part by post. We will activate the digital certificate within 3 days of receiving the physical copy in our PO Box. The LEFT part will be archived by UNCOPIED internally.",
		RightHelper: "Please, do not forget to send this RIGHT part by post. We will activate the digital certificate within 3 days of receiving the physical copy in our PO Box. The RIGHT part will be archived by a third party.",
		BottomHelper: "-5 Glue the other 4 parts to each 4 copies of the legal documentation (make sure labels correspond).\n-6 Retain one copy of the legal document (the one with the TOP chirograph part).\n-7 Send one copy of the legal document to the other contracting party (the one with the BOTTOM chirograph part).\n-8 Send the other two copies of the legal documents to UNCOPIED PO Box by post (the ones with LEFT and RIGHT chirograph parts).",
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

