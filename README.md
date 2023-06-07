# `docsim`

[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![CI Status](https://github.com/hrs/docsim/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/hrs/docsim/actions/workflows/test.yml)

A simple, fast command-line tool for scoring the similarity of text documents.

``` console
$ docsim --query some-file.txt --show-scores ~/documents/notes
0.000  completely-dissimilar-file.txt
0.152  somewhat-similar-file.md
0.469  pretty-similar-file.md
0.872  very-similar-file.org
1.000  potentially-identical-file.txt
```

Given a query document and a collection of potential matches, `docsim` ranks
each document in the collection by its textual similarity to the query.

## Examples

Check the [`man` page] for the definitive documentation, but these should get
you started.

[`man` page]: ./man/docsim.1

Search for similar files in a given directory:

``` console
$ docsim --query some-file.txt ~/documents/notes
[...]
```

If no paths are provided `docsim` will search the current working directory.

``` console
$ docsim --query some-file.txt
[...]
```

Without a provided `--query` document `docsim` takes input from `STDIN`. This
means `docsim` can be used as an ad-hoc local search engine.

``` console
$ echo "Here's a query to search for." | docsim ~/documents/notes
[...]
```

Only show the top 3 matches, with the best at the top:

``` console
$ docsim --query some-file.txt --limit 3 --best-first ~/documents/notes
potentially-identical-file.txt
very-similar-file.org
pretty-similar-file.md
```

Find Go files similar to `main.go` in the current directory. Don't use stemming
or stoplists, since these aren't English documents.

``` console
$ docsim --query main.go --no-stemming --no-stoplist **/*.go
[...]
```

Notice that because `docsim` uses an English stoplist and an English stemming
algorithm, you'll almost certainly want to use the `--no-stoplist` and
`--no-stemming` flags if your documents are written in another language
(including source code).

Optionally, you can use the `--stoplist` flag to provide a custom stoplist. A
custom stoplist is just a text file of words to ignore, separated by whitespace.

**WARNING:** `docsim` doesn't respect `.ignore` files yet, so it'll try to
search through `.git` directories, `node_modules`, and so on. That should be
fixed in the near future.

## Installation

The easiest thing is probably to grab a [compiled binary][] appropriate to your
platform.

[compiled binary]: https://github.com/hrs/docsim/releases/latest

But if you've got a Go toolchain handy, you can also either:

``` console
$ git clone git@github.com:hrs/docsim.git
$ cd docsim
$ sudo make install
```

Or just:

``` console
$ go install github.com/hrs/docsim/docsim@latest
```

Note that using `go install` doesn't include the [`man` page][], which you can
optionally install manually by copying it into e.g. `/usr/local/share/man/man1`.

[`man` page]: ./man/docsim.1

## Running tests

Just use the supplied `make` task:

``` console
$ make test
```

## How it works

`docsim` uses [TF-IDF][] weighting and [cosine similarity][] to numerically
score the textual similarity between the query and every other document.

[TF-IDF]: https://en.wikipedia.org/wiki/Tf%E2%80%93idf
[cosine similarity]: https://en.wikipedia.org/wiki/Vector_space_model#Applications

"Textual similarity" roughly means "uses the same words." Each document is
parsed into a big bag of words, which are passed through a common English
[stoplist][], [stemmed][] (so "spins," "spinner," and "spinning" might all
reduce down to just "spin"), and assigned weights based on how often they appear
in the document and how rare they are in the corpus as a whole.

[stoplist]: https://en.wikipedia.org/wiki/Stop_word
[stemmed]: https://en.wikipedia.org/wiki/Stemming

We can think of each of these documents as a [vector in term space][], where
each word is a dimension with its weight as a magnitude. Two documents are
"similar," then, inasmuch as they point in the same direction, so we define
similarity by the size of the angle between them.

[vector in term space]: https://en.wikipedia.org/wiki/Vector_space_model

## Contributing

`docsim` is still in a nascent state, so I'm happy just writing the code myself
for now, but please feel free to [report any issues][] you encounter!

[report any issues]: https://github.com/hrs/docsim/issues
