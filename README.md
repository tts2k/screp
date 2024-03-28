# Screp
A prototype article scraper for the [Stanford Encyclopedia of Philosophy](https://plato.stanford.edu). To be used for something else (like a web app).

### Usage
```
Screp is a scraper for the Stanford Encyclopedia of Philosophy (SEP).

Usage:
  screp [flags] <url>

Flags:
  -v, --verbose   Enable verbose output
  -j, --json      Output to JSON
  -y, --yaml      Output to YAML
  -h, --help      Print this help message

```
### Example
#### Run

``` sh
$ ./screp https://plato.stanford.edu/entries/exploitation/
```

#### JSON output
```sh
$ ./screp https://plato.stanford.edu/entries/exploitation/ --json | jq .preamble[0]

"To exploit someone is to take unfair advantage of them. It is to use another person’s vulnerability
for one’s own benefit. Of course, benefitting from another’s vulnerability is not always morally
wrong—we do not condemn a chess player for exploiting a weakness in his opponent’s defence, for
instance. But some forms of advantage-taking do seem to be clearly wrong, and it is this normative
sense of exploitation that is of primary interest to moral and political philosophers."
```

#### Troff output
If you prefer to read from your terminal, screp can output Troff with ms macro. Using it should be the same, more or less
![simplicity](docs/simplicity.gif)

Leveraging the power of [Groff](https://www.gnu.org/software/groff/), you can also turn the output into a pdf file
```sh
./screp https://plato.stanford.edu/entries/simplicity | groff -k -ms -E -Tpdf > output.pdf
```

Groff can turn the output into html too! Although the upcoming Templ render feature will output much better formatted html
```sh
./screp https://plato.stanford.edu/entries/simplicity | groff -k -ms -E -Thtml > output.html
```

#### Maybe grill some [fishes](https://fishshell.com/) while chewing [gums](https://github.com/charmbracelet/gum), recorded in [vhs](https://github.com/charmbracelet/vhs)
![derrida](docs/derrida.gif)

### TODO
- [x] TOC print out as a tree
- [x] Scrape main text
- [x] CLI
- [x] JSON output
- [x] YAML output
- [ ] Use the result to do some silly things with [templ](https://templ.guide)
- [ ] Markdown output
- [ ] Blockquote parsing
- [ ] Table parsing

### Advanced nice-to-have features
- [ ] Scrape images and math fields
