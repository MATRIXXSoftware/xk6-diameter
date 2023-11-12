package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fiorix/go-diameter/v4/diam/dict"
	log "github.com/sirupsen/logrus"
)

func main() {
	var dictionary string
	var output string

	flag.StringVar(&dictionary, "dictionary", "", "Dictionary")
	flag.StringVar(&output, "output", "dict.js", "Output file")
	flag.Parse()

	if dictionary == "" {
		log.Fatalf("Dictionary not found\n")
	}

	file, err := os.Open(dictionary)
	if err != nil {
		log.Fatalf("Error opening dictioanry: %s\n", err)
	}
	defer file.Close()

	parser, err := dict.NewParser("dict/base.xml", dictionary)
	if err != nil {
		log.Fatalf("Error parsing dictioanry: %s\n", err)
	}

	w, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %s\n", err)
	}
	defer w.Close()

	PrintFlags(w)
	PrintAvpCode(w, parser)
	PrintVendorId(w, parser)
}

func PrintFlags(w io.Writer) {
	fmt.Fprintf(w, "export const flags = {\n")
	fmt.Fprintf(w, "    Vbit: 0x80,\n")
	fmt.Fprintf(w, "    Mbit: 0x40,\n")
	fmt.Fprintf(w, "    Pbit: 0x20,\n")
	fmt.Fprintf(w, "}\n")
	fmt.Fprintf(w, "\n")
}

func PrintAvpCode(w io.Writer, parser *dict.Parser) {
	fmt.Fprintf(w, "export const avpCode = {\n")
	for _, app := range parser.Apps() {
		fmt.Fprintf(w, "    // %s\n", app.Name)
		for _, avp := range app.AVP {
			name := strings.ReplaceAll(avp.Name, "-", "")
			fmt.Fprintf(w, "    %s: %d,\n", name, avp.Code)
		}
	}
	fmt.Fprintf(w, "}\n")
	fmt.Fprintf(w, "\n")

}

func PrintVendorId(w io.Writer, parser *dict.Parser) {
	fmt.Fprintf(w, "export const vendorId = {\n")
	for _, app := range parser.Apps() {
		for _, vendor := range app.Vendor {
			fmt.Fprintf(w, "    %s: %d,\n", vendor.Name, vendor.ID)
		}
	}
	fmt.Fprintf(w, "}\n")
	fmt.Fprintf(w, "\n")
}
