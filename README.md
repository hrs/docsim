# `docsim`

[![Release version](https://img.shields.io/github/v/release/hrs/docsim)](https://github.com/hrs/docsim/releases/latest)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![CI Status](https://github.com/hrs/docsim/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/hrs/docsim/actions/workflows/test.yml)

A local, in-memory search tool. Query and compare your text documents from the
terminal, with results ranked by textual similarity.

``` console
$ docsim --show-scores --limit 3 --best-first "search query" ~/documents/notes
0.472  very-relevant-file.txt
0.123  slightly-similar-file.org
0.000  completely-unrelated-file.md
```

`docsim` is an information retrieval tool, so it's different from other search
tools like `grep`, `ripgrep`, `ag`, and so on. Those tools are all great, but
they search for literal text matches, and sometimes we want to know, "what notes
are *most similar* to this query, or to this other note?"

If I search for "chunky bacon," I still want to see documents that talk about
"chunks of bacon." And, below those, I probably want to see notes that discuss
regular "bacon," even if it's not chunky. `docsim` uses a few different
[information retrieval algorithms][] to provide a ranked list of text documents.

[information retrieval algorithms]: #how-it-works

It's also slower and more memory-intensive than e.g. `grep`, of course, since it
does more work. But performance *is* a goal, and on a mid-range machine it'll
process a few thousand documents without notable lag.

This all sounds complicated, but `docsim` aspires to be easy to use and to
behave like a good UNIX citizen. It's a single binary that operates on plain
files and streams. No servers, no daemons, no dependencies on Docker containers
or `scikit-learn`, not even any persistent indexes or caches to get out of sync.
Searching local documents with information retrieval algorithms shouldn't be any
harder than using `grep`!

## Examples

Check the [`man` page][] for the definitive documentation, but these should get
you started.

[`man` page]: ./man/docsim.1

If no paths are provided `docsim` will search the current working directory.

``` console
$ docsim "here's a search query"
[...]
```

Use the `--stdin` flag to read the search query from `STDIN` instead of a string
argument.

``` console
$ echo "Here's another query to search for." | docsim --stdin ~/documents/notes
[...]
```

Search for similar files in a given directory:

``` console
$ docsim --file some-file.txt ~/documents/notes
[...]
```

Find Go files similar to `main.go` in the current directory. Don't use natural
language processing techniques like [stemming or stoplists][], since these
aren't English documents:

[stemming or stoplists]: #how-it-works

``` console
$ docsim --file main.go --no-stemming --no-stoplist **/*.go
[...]
```

Note that because `docsim` uses an English stoplist and an English stemming
algorithm, you'll almost certainly want to use the `--no-stoplist` and
`--no-stemming` flags if your documents are written in another language
(including source code).

Optionally, you can use the `--stoplist` flag to provide a custom stoplist. A
custom stoplist is just a text file of words to ignore, separated by whitespace.

**WARNING:** `docsim` doesn't respect `.ignore` or `.gitignore` files yet, so
it'll try to search through `.git` directories, `node_modules`, and so on. That
should be fixed in the near future.

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
reduce down to just "spin"). Those terms are assigned weights based on how often
they appear in the document and how rare they are in the corpus as a whole.

[stoplist]: https://en.wikipedia.org/wiki/Stop_word
[stemmed]: https://en.wikipedia.org/wiki/Stemming

We can think of each of these documents as a [vector in term space][], where
each term is a dimension with its weight as a magnitude. Two documents are
"similar," then, inasmuch as they point in the same direction, so we define
similarity by the size of the angle between them.

[vector in term space]: https://en.wikipedia.org/wiki/Vector_space_model

## Contributing

`docsim` is still in a nascent state, so I'm happy just writing the code myself
for now, but please feel free to [report any issues][] you encounter!

[report any issues]: https://github.com/hrs/docsim/issues
