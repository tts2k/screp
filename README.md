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
$ ./screp go run . https://plato.stanford.edu/entries/exploitation/
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

### TODO
- [x] TOC print out as a tree
- [x] Scrape main text
- [x] CLI
- [x] JSON output
- [x] YAML output
- [ ] Use the result to do some silly things with [templ](https://templ.guide)
- [ ] Pandoc integration
