package main

import (
	"fmt"
	"os"

	"screp/lib/scraper"

	flags "github.com/spf13/pflag"
)

func initFlags() {
	flags.Usage = func() {
		fmt.Fprintln(os.Stderr,
			"Screp is a scraper for the Stanford Encyclopedia of Philosophy (SEP).\n\n"+
				"Usage:\n"+
				"  screp [flags] <url>\n\n"+
				"Flags:",
		)
		flags.PrintDefaults()
	}

	flags.BoolP("verbose", "v", false, "Enable verbose output")
	flags.BoolP("json", "j", false, "Output to JSON")
	flags.BoolP("yaml", "y", false, "Output to YAML")
	flags.BoolP("help", "h", false, "Print this help message")
	flags.CommandLine.SortFlags = false

	flags.Parse()
}

func main() {
	initFlags()

	// Process flags
	helpF, _ := flags.CommandLine.GetBool("help")
	if helpF {
		flags.Usage()
		return
	}

	jsonF, _ := flags.CommandLine.GetBool("json")
	yamlF, _ := flags.CommandLine.GetBool("yaml")
	if jsonF && yamlF {
		fmt.Fprintln(os.Stderr, "Only json and yaml flags cannot be used together.")
		flags.Usage()
		return
	}

	verboseF, _ := flags.CommandLine.GetBool("verbose")
	url := flags.CommandLine.Arg(0)
	if url == "" {
		fmt.Fprintln(os.Stderr, "Please provide an URL.")
		flags.Usage()
		return
	}

	// Begin initializing scraper
	config := scraper.Config{
		URL:     url,
		Verbose: verboseF,
	}

	scr := scraper.NewScraper(config)

	err := scr.Scrape()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// YAML
	if yamlF {
		yaml, err := scr.YAML()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		fmt.Println(yaml)
		return
	}

	// JSON
	if jsonF {
		json, err := scr.JSON()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		fmt.Println(json)
		return
	}

	// Print table of content
	scr.PrintTOC()
}
