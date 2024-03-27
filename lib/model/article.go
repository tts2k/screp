package model

type Section struct {
	Title       string    `yaml:"title" json:"title"`
	Content     []string  `yaml:"content" json:"content"`
	SubSections []Section `yaml:"subSections,omitempty" json:"subSections,omitempty"`
}

type Article struct {
	Title    string    `yaml:"title" json:"title"`
	Author   []string  `yaml:"author" json:"author"`
	Issued   string    `yaml:"issued" json:"issued"`
	Modified string    `yaml:"modified" json:"modified"`
	Preamble []string  `yaml:"preamble" json:"preamble"`
	Sections []Section `yaml:"sections,omitempty" json:"sections,omitempty"`
}
