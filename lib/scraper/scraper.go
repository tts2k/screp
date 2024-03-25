package scraper

import (
	"fmt"
	"regexp"
	"strconv"

	"screp/lib/model"
	container "screp/lib/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Scraper struct {
	url       string
	article   model.Article
	collector *colly.Collector
}

func NewScarper(url string) *Scraper {
	return &Scraper{
		url:       url,
		collector: colly.NewCollector(),
		article:   model.Article{},
	}
}

func (scraper *Scraper) Article() model.Article {
	return scraper.article
}

func parseMainText(childs *goquery.Selection) []model.Section {
	result := []model.Section{}
	secStack := container.SectionStack{}

	headingRe := regexp.MustCompile("^[^A-Za-z]+")

	// Getting the correct section list
	var currSectList *[]model.Section
	headingLvl := 2

	updateCurrSectList := func() {
		if secStack.IsEmpty() {
			currSectList = &result
		} else {
			currSectList = &secStack.Peek().SubSections
		}
	}
	updateCurrSectList()

	childs.Each(func(_ int, s *goquery.Selection) {
		tagName := goquery.NodeName(s)

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
			newSection := model.Section{Title: headingRe.ReplaceAllString(s.Text(), "")}
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

func recurPrintTOC(sections []model.Section, prefix string, level int) {
	var pad string
	for i := 0; i < level; i++ {
		pad += "  "
	}

	for index, section := range sections {
		var sectionNumbering string

		if prefix == "" {
			sectionNumbering = strconv.Itoa(index+1) + "."
		} else {
			sectionNumbering = fmt.Sprintf("%s%d", prefix, index+1)
		}

		fmt.Printf("%s%s %s\n", pad, sectionNumbering, section.Title)

		// remove the dot on first level
		if prefix == "" {
			sectionNumbering = sectionNumbering[:1]
		}

		if len(section.SubSections) > 0 {
			recurPrintTOC(section.SubSections, sectionNumbering+".", level+1)
		}
	}
}

func (scraper *Scraper) PrintTOC() *Scraper {
	fmt.Println(scraper.article.Title)
	recurPrintTOC(scraper.article.Sections, "", 0)
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
