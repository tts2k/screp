package main

import (
	"fmt"
	"os"

	"screp/lib/scraper"

	flags "github.com/spf13/pflag"
)

func checkBoolFlagsConflict(flagList []string) error {
	hasFlagEnabled := false
	var enabledFlag string

	for _, flag := range flagList {
		value, _ := flags.CommandLine.GetBool(flag)
		if hasFlagEnabled && value {
			return fmt.Errorf("conflicting flags: %s, %s", enabledFlag, flag)
		}

		if value {
			hasFlagEnabled = true
			enabledFlag = flag
		}
	}

	return nil
}

func initFlags() error {
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
	flags.BoolP("html", "H", false, "Output to HTML")
	flags.BoolP("markdown", "m", false, "Output to Markdown")
	flags.BoolP("help", "h", false, "Print this help message")
	flags.CommandLine.SortFlags = false

	flags.Parse()

	err := checkBoolFlagsConflict([]string{"json", "yaml", "html", "markdown"})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := initFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		flags.Usage()
		return
	}

	// Process flags
	helpF, _ := flags.CommandLine.GetBool("help")
	if helpF {
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

	err = scr.Scrape()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// YAML
	yamlF, _ := flags.CommandLine.GetBool("yaml")
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
	jsonF, _ := flags.CommandLine.GetBool("json")
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
