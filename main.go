package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"

	"screp/lib"
)

func parseMainText(childs *goquery.Selection) []lib.Section {
	result := []lib.Section{}
	secStack := lib.SectionStack{}

	headingRe := regexp.MustCompile("^[^A-Za-z]+")

	headingLvl := 2

	childs.Each(func(_ int, s *goquery.Selection) {
		tagName := goquery.NodeName(s)

		// Getting the correct section list
		var currSectList *[]lib.Section

		updateCurrSectList := func() {
			if secStack.IsEmpty() {
				currSectList = &result
			} else {
				currSectList = &secStack.Peek().SubSections
			}
		}
		updateCurrSectList()

		// Heading tag
		// When heading is deeper than current heading level
		if tagName[0] == 'h' {
			currTagHeadingLvl := int(tagName[1] - '0')

			if currTagHeadingLvl > headingLvl {
				headingLvl++
				currSect := &(*currSectList)[len(*currSectList)-1]
				secStack.Push(currSect)
				updateCurrSectList()
			}

			// When heading isn't deeper than current heading level
			for currTagHeadingLvl < headingLvl {
				headingLvl--
				secStack.Pop()
				updateCurrSectList()
			}

			// When heading levels are equal
			newSection := lib.Section{Title: headingRe.ReplaceAllString(s.Text(), "")}
			*currSectList = append(*currSectList, newSection)

			return
		}

		// Unsupported tags
		if tagName != "p" {
			fmt.Println("Skipped unsupported tag:", tagName)
			return
		}

		// TODO: p tag
	})

	return result
}

func printToc(sections []lib.Section, level int) {
	var pad string
	for i := 0; i < level; i++ {
		pad += "  "
	}

	for _, section := range sections {
		fmt.Printf("%s%s\n", pad, section.Title)
		if len(section.SubSections) > 0 {
			printToc(section.SubSections, level+1)
		}
	}
}

func main() {
	args := os.Args[1:]

	article := lib.Article{}
	collector := colly.NewCollector()

	collector.OnHTML("div[id=aueditable] h1", func(e *colly.HTMLElement) {
		article.Title = e.Text
	})

	collector.OnHTML("div[id=preamble]", func(e *colly.HTMLElement) {
		article.Preamble = e.Text
	})

	collector.OnHTML("div[id=main-text]", func(e *colly.HTMLElement) {
		article.Sections = parseMainText(e.DOM.Children())
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	err := collector.Visit(args[0])
	if err != nil {
		panic(err)
	}

	collector.Wait()

	// print TOC
	printToc(article.Sections, 0)
}
