package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fiorix/go-diameter/v4/diam/dict"
	log "github.com/sirupsen/logrus"
)

func main() {
	var dictionary string
	flag.StringVar(&dictionary, "dictionary", "", "Dictionary")

	flag.Parse()

	if dictionary == "" {
		log.Fatalf("Dictionary not found\n")
	}

	file, err := os.Open(dictionary)
	if err != nil {
		log.Fatalf("Error opening dictioanry: %s\n", err)
	}
	defer file.Close()

	log.Infof("Dictionary '%s' exist\n", dictionary)

	parser, err := dict.NewParser("dict/base.xml", dictionary)
	if err != nil {
		log.Fatalf("Error parsing dictioanry: %s\n", err)
	}

	for _, app := range parser.Apps() {
		log.Infof("App: %v\n", app.AVP)
		for _, avp := range app.AVP {
			name := strings.ReplaceAll(avp.Name, "-", "")
			fmt.Printf("    %s: %d,\n", name, avp.Code)
			//log.Infof("    name: %s code: %d\n", name, avp.Code)
		}
	}

}
