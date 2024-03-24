package lib

import (
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Scraper struct {
	url       string
	article   Article
	collector *colly.Collector
}

func NewScarper(url string) *Scraper {
	return &Scraper{
		url:       url,
		collector: colly.NewCollector(),
		article:   Article{},
	}
}

func (scraper *Scraper) Article() Article {
	return scraper.article
}

func parseMainText(childs *goquery.Selection) []Section {
	result := []Section{}
	secStack := SectionStack{}

	headingRe := regexp.MustCompile("^[^A-Za-z]+")

	headingLvl := 2

	childs.Each(func(_ int, s *goquery.Selection) {
		tagName := goquery.NodeName(s)

		// Getting the correct section list
		var currSectList *[]Section

		updateCurrSectList := func() {
			if secStack.IsEmpty() {
				currSectList = &result
			} else {
				currSectList = &secStack.Peek().SubSections
			}
		}
		updateCurrSectList()

		// Heading tag
		if tagName[0] == 'h' {
			currTagHeadingLvl := int(tagName[1] - '0')

			// When heading is deeper than current heading level
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
			newSection := Section{Title: headingRe.ReplaceAllString(s.Text(), "")}
			*currSectList = append(*currSectList, newSection)

			return
		}

		// Unsupported tags
		if tagName != "p" {
			fmt.Println("Skipped unsupported tag:", tagName)
			return
		}

		// p tag
		currSect := &(*currSectList)[len(*currSectList)-1]
		currSect.Content = append(currSect.Content, s.Text())
	})

	return result
}

func recurPrintTOC(sections []Section, level int) {
	var pad string
	for i := 0; i < level; i++ {
		pad += "  "
	}

	for _, section := range sections {
		fmt.Printf("%s%s\n", pad, section.Title)
		if len(section.SubSections) > 0 {
			recurPrintTOC(section.SubSections, level+1)
		}
	}
}

func (scraper *Scraper) PrintTOC() *Scraper {
	println(scraper.article.Title)
	recurPrintTOC(scraper.article.Sections, 0)
	return scraper
}

func (scraper *Scraper) Scrape() error {
	scraper.collector.OnHTML("div[id=aueditable] h1", func(e *colly.HTMLElement) {
		scraper.article.Title = e.Text
	})

	scraper.collector.OnHTML("div[id=preamble]", func(e *colly.HTMLElement) {
		scraper.article.Preamble = e.Text
	})

	scraper.collector.OnHTML("div[id=main-text]", func(e *colly.HTMLElement) {
		scraper.article.Sections = parseMainText(e.DOM.Children())
	})

	err := scraper.collector.Visit(scraper.url)
	if err != nil {
		return err
	}

	scraper.collector.Wait()

	return nil
}
