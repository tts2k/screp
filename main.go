package main

import (
	"fmt"
	"os"

	"screp/lib"
)

func main() {
	args := os.Args[1:]

	scraper := lib.NewScarper(args[0])

	err := scraper.Scrape()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	scraper.PrintTOC()
}
