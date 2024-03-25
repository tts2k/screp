package main

import (
	"fmt"
	"os"

	"screp/lib/scraper"
)

func main() {
	args := os.Args[1:]

	scr := scraper.NewScarper(args[0])

	err := scr.Scrape()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	scr.PrintTOC()
}
