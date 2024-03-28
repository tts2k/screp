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

type Article struct {
	Title        string    `yaml:"title" json:"title"`
	Author       []string  `yaml:"author" json:"author"`
	Issued       string    `yaml:"issued" json:"issued"`
	Modified     string    `yaml:"modified" json:"modified"`
	Preamble     []string  `yaml:"preamble" json:"preamble"`
	Sections     []Section `yaml:"sections,omitempty" json:"sections,omitempty"`
	Bibliography []Section `yaml:"bibliography" json:"bibliography"`
}

func wrapTroffMacro(macro string, text string) string {
	return fmt.Sprintf(".%s\n%s\n", macro, text)
}

func recurTroff(sections []Section, level int, bib bool) string {
	var result string
	var header string
	if bib {
		header = "SH "
	} else {
		header = "NH "
	}
	header += strconv.Itoa(level)

	// Print body text
	for _, section := range sections {

		result += wrapTroffMacro(header, section.Title)
		for _, text := range section.Content {
			result += wrapTroffMacro("PP", text)
		}

		if len(section.SubSections) > 0 {
			result += recurTroff(section.SubSections, level+1, bib)
		}
	}
	return result
}

func (a *Article) Troff() string {
	// Title
	result := ".nr PO 0.6i\n"
	result += ".nr LL 7.05i\n"
	result += wrapTroffMacro("TL", a.Title)

	// Author
	for _, author := range a.Author {
		result += wrapTroffMacro("AU", author)
	}

	// Preamble
	for _, preamble := range a.Preamble {
		result += wrapTroffMacro("LP", preamble)
	}

	// Content
	result += recurTroff(a.Sections, 1, false)

	// Bibliography
	result += recurTroff(a.Bibliography, 1, true)

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
