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

	flags.BoolP("help", "h", false, "Print this help message")
	flags.BoolP("json", "j", false, "Output to JSON")
	flags.BoolP("yaml", "y", false, "Output to YAML")
	flags.BoolP("html", "H", false, "Output to HTML")
	flags.BoolP("md", "m", false, "Output to Markdown")
	flags.BoolP("troff", "t", false, "Output to Markdown")
	flags.BoolP("verbose", "v", false, "Enable verbose output")
	flags.CommandLine.SortFlags = false

	flags.Parse()

	err := checkBoolFlagsConflict([]string{"json", "yaml", "html", "markdown", "troff"})
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
	article := scr.Article()

	// YAML
	yamlF, _ := flags.CommandLine.GetBool("yaml")
	if yamlF {
		yaml, err := article.YAML()
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
		json, err := article.JSON()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		fmt.Println(json)
		return
	}

	// Print plain text
	text := article.Troff()
	fmt.Println(text)
}
