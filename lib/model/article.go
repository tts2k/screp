package model

type Section struct {
	Title       string    `yaml:"title" json:"title"`
	Content     []string  `yaml:"content" json:"content"`
	SubSections []Section `yaml:"subSections,omitempty" json:"subSections,omitempty"`
}

type Article struct {
	Title        string    `yaml:"title" json:"title"`
	Created      string    `yaml:"created" json:"created"`
	LastModified string    `yaml:"lastModified,omitempty" json:"lastModified,omitempty"`
	Preamble     []string  `yaml:"preamble" json:"preamble"`
	Sections     []Section `yaml:"sections,omitempty" json:"sections,omitempty"`
}
