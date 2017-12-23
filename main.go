package main

import (
	"fmt"
	"os"

	"github.com/Noah-Huppert/crime-map/parsers"
)

const file = "data/2017-10-12.pdf"

func main() {
	// Parse crimes
	r := parsers.NewReport(file)

	crimes, err := r.Parse()
	if err != nil {
		fmt.Printf("error parsing report: %s", err.Error())
		os.Exit(1)
	}

	// Print crimes
	for i, crime := range crimes {
		fmt.Printf("\n%d\n====\n%s\n", i+1, crime)
	}
}
