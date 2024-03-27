package scraper

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"screp/lib/model"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"gopkg.in/yaml.v3"
)

// Global
var headingRegex = regexp.MustCompile("^[^A-Za-z]+")

type Config struct {
	Verbose bool
	URL     string
}

type Scraper struct {
	article   model.Article
	collector *colly.Collector
	config    Config
}

func NewScraper(config Config) *Scraper {
	return &Scraper{
		collector: colly.NewCollector(),
		article:   model.Article{},
		config:    config,
	}
}

func (scraper *Scraper) Article() model.Article {
	return scraper.article
}

func processPText(text string) string {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")
	return text
}

func parseMainText(config *Config, childs *goquery.Selection) []model.Section {
	result := []model.Section{}
	secStack := SectionStack{}

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
			newSection := model.Section{Title: headingRegex.ReplaceAllString(s.Text(), "")}
			*currSectList = append(*currSectList, newSection)

			return
		}

		// Unsupported tags
		if tagName != "p" {
			if config.Verbose {
				fmt.Println("Skipped unsupported tag:", tagName)
			}
			return
		}

		// p tag
		currSect := &(*currSectList)[len(*currSectList)-1]
		currSect.Content = append(currSect.Content, processPText(s.Text()))
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
	scraper.collector.OnHTML(`meta[name="DC.title"]`, func(e *colly.HTMLElement) {
		scraper.article.Title = e.Attr("content")
	})

	scraper.collector.OnHTML(`meta[name="DC.creator"]`, func(e *colly.HTMLElement) {
		scraper.article.Author = append(scraper.article.Author, e.Attr("content"))
	})

	scraper.collector.OnHTML(`meta[name="DCTERMS.issued"]`, func(e *colly.HTMLElement) {
		scraper.article.Issued = e.Attr("content")
	})

	scraper.collector.OnHTML(`meta[name="DCTERMS.modified"]`, func(e *colly.HTMLElement) {
		scraper.article.Modified = e.Attr("content")
	})

	scraper.collector.OnHTML("div[id=preamble] > p", func(e *colly.HTMLElement) {
		scraper.article.Preamble = append(scraper.article.Preamble, processPText(e.Text))
	})

	scraper.collector.OnHTML("div[id=main-text]", func(e *colly.HTMLElement) {
		scraper.article.Sections = parseMainText(&scraper.config, e.DOM.Children())
	})

	scraper.collector.OnHTML("div[id=bibliography] > ul", func(e *colly.HTMLElement) {
		// Primary
		if e.Index == 0 {
			scraper.article.Bibliography.Primary = e.DOM.
				Find("li").
				Map(func(_ int, s *goquery.Selection) string {
					return processPText(s.Text())
				})
			return
		}

		// Secondary
		scraper.article.Bibliography.Secondary = e.DOM.
			Find("li").
			Map(func(_ int, s *goquery.Selection) string {
				return processPText(s.Text())
			})
	})

	err := scraper.collector.Visit(scraper.config.URL)
	if err != nil {
		return err
	}

	scraper.collector.Wait()

	return nil
}

func (scraper *Scraper) YAML() (string, error) {
	out, err := yaml.Marshal(scraper.article)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (scraper *Scraper) JSON() (string, error) {
	out, err := json.Marshal(scraper.article)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
