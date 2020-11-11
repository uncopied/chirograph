package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	outFile := flag.String("o", "out", "output file prefix, empty for stdout")
	outFormat := flag.String("f", "svg", "output file format (svg or pdf)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `tallystick secured document in Go
https://github.com/uncopied/tallystick

Flags:
`)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Usage:
  	tallystick -o pdf "hello world"

`)
	}
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		checkError(fmt.Errorf("Error: no content given"))
	}

	content := strings.Join(flag.Args(), " ")


	if *outFile == "" {
		os.Stdout.WriteString(content)
	} else {
		fh, err := os.Create(*outFile + "."+ *outFormat)
		checkError(err)
		defer fh.Close()
		fh.WriteString(content)
		fh.Close()
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

