package lib

type Section struct {
	Title       string
	Content     []string
	SubSections []Section
}

type Article struct {
	Title    string
	Preamble string
	Sections []Section
}
