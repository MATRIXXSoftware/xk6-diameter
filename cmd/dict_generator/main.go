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
	flag.StringVar(&output, "output", "const.js", "Output file")
	flag.Parse()

	parser := dict.Default

	if dictionary != "" {
		file, err := os.Open(dictionary)
		if err != nil {
			log.Fatalf("Error opening dictioanry: %s\n", err)
		}
		defer file.Close()

		parser.Load(file)
		if err != nil {
			log.Fatalf("Error parsing dictioanry: %s\n", err)
		}
	}

	w, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %s\n", err)
	}
	defer w.Close()

	PrintCmd(w)
	PrintAppId(w)
	PrintFlags(w)
	PrintAvpCode(w, parser)
	PrintVendorId(w, parser)
}

func PrintAppId(w io.Writer) {
	fmt.Fprintf(w, "export const appId = {\n")
	fmt.Fprintf(w, "	%-35s %d,\n", "Accounting:", 3)
	fmt.Fprintf(w, "	%-35s %d,\n", "ChargingControl:", 4)
	fmt.Fprintf(w, "	%-35s %d,\n", "Gx:", 16777238)
	fmt.Fprintf(w, "	%-35s %d,\n", "Sy:", 16777302)
	fmt.Fprintf(w, "}\n")
	fmt.Fprintf(w, "\n")
}

func PrintCmd(w io.Writer) {
	fmt.Fprintf(w, "export const cmd = {\n")
	fmt.Fprintf(w, "	%-35s %d,\n", "AA:", 265)
	fmt.Fprintf(w, "	%-35s %d,\n", "Accounting:", 271)
	fmt.Fprintf(w, "	%-35s %d,\n", "CreditControl:", 272)
	fmt.Fprintf(w, "	%-35s %d,\n", "ReAuth:", 258)
	fmt.Fprintf(w, "	%-35s %d,\n", "SessionTermination:", 275)
	fmt.Fprintf(w, "	%-35s %d,\n", "SpendingLimit:", 8388635)
	fmt.Fprintf(w, "}\n")
	fmt.Fprintf(w, "\n")
}

func PrintFlags(w io.Writer) {
	fmt.Fprintf(w, "export const flags = {\n")
	fmt.Fprintf(w, "	%-35s 0x%x,\n", "Vbit:", 0x80)
	fmt.Fprintf(w, "	%-35s 0x%x,\n", "Mbit:", 0x40)
	fmt.Fprintf(w, "	%-35s 0x%x,\n", "Pbit:", 0x20)
	fmt.Fprintf(w, "}\n")
	fmt.Fprintf(w, "\n")
}

func PrintAvpCode(w io.Writer, parser *dict.Parser) {
	fmt.Fprintf(w, "export const avpCode = {\n")
	for _, app := range parser.Apps() {
		fmt.Fprintf(w, "    // %s\n", app.Name)
		for _, avp := range app.AVP {
			name := strings.ReplaceAll(avp.Name, "-", "")
			fmt.Fprintf(w, "    %-35s %d,\n", name+":", avp.Code)
		}
	}
	fmt.Fprintf(w, "}\n")
	fmt.Fprintf(w, "\n")

}

func PrintVendorId(w io.Writer, parser *dict.Parser) {
	vendorIds := make(map[uint32]struct{})
	exists := struct{}{}

	fmt.Fprintf(w, "export const vendorId = {\n")
	for _, app := range parser.Apps() {
		for _, vendor := range app.Vendor {
			// Remove duplicate vendorId
			_, found := vendorIds[vendor.ID]
			if found {
				continue
			}
			vendorIds[vendor.ID] = exists

			fmt.Fprintf(w, "    %-35s %d,\n", vendor.Name+":", vendor.ID)
			// fmt.Fprintf(w, "    %s: %d,\n", vendor.Name, vendor.ID)
		}
	}
	fmt.Fprintf(w, "}\n")
	fmt.Fprintf(w, "\n")
}
