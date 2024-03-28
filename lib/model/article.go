package model

import (
	"encoding/json"
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Section struct {
	Title       string    `yaml:"title" json:"title"`
	Content     []string  `yaml:"content" json:"content"`
	SubSections []Section `yaml:"subSections,omitempty" json:"subSections,omitempty"`
}

type Bibliography struct {
	Primary   []string `yaml:"primary" json:"primary"`
	Secondary []string `yaml:"secondary,omitempty" json:"secondary,omitempty"`
}

type Article struct {
	Title        string       `yaml:"title" json:"title"`
	Author       []string     `yaml:"author" json:"author"`
	Issued       string       `yaml:"issued" json:"issued"`
	Modified     string       `yaml:"modified" json:"modified"`
	Preamble     []string     `yaml:"preamble" json:"preamble"`
	Sections     []Section    `yaml:"sections,omitempty" json:"sections,omitempty"`
	Bibliography Bibliography `yaml:"bibliography" json:"bibliography"`
}

func wrapTroffMacro(macro string, text string) string {
	return fmt.Sprintf(".%s\n%s\n", macro, text)
}

func recurTroff(sections []Section, level int) string {
	var result string
	header := "NH " + strconv.Itoa(level)

	for _, section := range sections {

		// Print body text
		result += wrapTroffMacro(header, section.Title)
		for _, text := range section.Content {
			result += wrapTroffMacro("PP", text)
		}

		if len(section.SubSections) > 0 {
			result += recurTroff(section.SubSections, level+1)
		}
	}

	return result
}

func (a *Article) Troff() string {
	// Title
	result := wrapTroffMacro("TL", a.Title)

	// Author
	for _, author := range a.Author {
		result += wrapTroffMacro("AU", author)
	}

	// Preamble
	for _, preamble := range a.Preamble {
		result += wrapTroffMacro("LP", preamble)
	}

	// Content
	result += result + recurTroff(a.Sections, 1)

	return result
}

func (a *Article) YAML() (string, error) {
	out, err := yaml.Marshal(*a)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (a *Article) JSON() (string, error) {
	out, err := json.Marshal(*a)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
