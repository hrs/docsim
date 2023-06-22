package corpus

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type termMap map[string]float64

type Document struct {
	path     string
	termFreq termMap
	tfIdf    termMap
	norm     float64
}

var explicitlyPermittedExtensions = map[string]bool{
	"markdown": true,
	"md":       true,
	"org":      true,
	"tex":      true,
	"txt":      true,
}

const apostropheRunes = "'’"

func ParseDocument(path string, config *Config) (*Document, error) {
	// Ensure that this is a text file
	if !isTextFile(path) {
		return nil, fmt.Errorf("not a text file, skipping: %s", path)
	}

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Parse the contents
	doc, err := NewDocument(file, config)
	if err != nil {
		return nil, err
	}

	doc.path = path
	return doc, nil
}

func parseTokens(rd io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(rd)
	scanner.Split(bufio.ScanWords)

	tokens := []string{}

	for scanner.Scan() {
		// Lower-case each word and replace curly apostrophes with single quotes.
		token := strings.Map(
			normalizeApostrophe,
			strings.ToLower(scanner.Text()),
		)

		// Split each token on non-alphanumeric characters (except single quotes, to
		// handle contractions)
		for _, word := range strings.FieldsFunc(token, splitToken) {
			// Since we didn't split on single quotes, we need to trim them off now.
			// We'd like "don't" to stay "don't", but "'hello" to become "hello".
			word = strings.Trim(word, apostropheRunes)

			// Similarly, we need to remove the common "'s" possessive case
			word = strings.TrimSuffix(word, "'s")

			if word != "" {
				tokens = append(tokens, word)
			}
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}

func NewDocument(rd io.Reader, config *Config) (*Document, error) {
	termCount := make(termMap)
	totalTermCount := 0.0

	tokens, err := parseTokens(rd)
	if err != nil {
		return nil, err
	}

	// Loop over the tokens, stem them if configured, pass them through the
	// stoplist if configured, and, for each that "should" "count", increment it
	// in the term map.
	for _, token := range tokens {
		if config.NoStoplist || !config.Stoplist.include(token) {
			if config.NoStemming {
				termCount[token]++
			} else {
				termCount[stem(token)]++
			}

			totalTermCount++
		}
	}

	// Scale the term frequency map according to the total number of terms in the document.
	termFreq := make(termMap)
	for term, count := range termCount {
		termFreq[term] = count / totalTermCount
	}

	return &Document{termFreq: termFreq}, nil
}

func splitToken(r rune) bool {
	return !(unicode.IsLetter(r) || unicode.IsNumber(r) || r == '\'')
}

func normalizeApostrophe(r rune) rune {
	for _, apostrophe := range apostropheRunes {
		if r == apostrophe {
			return '\''
		}
	}

	return r
}

func (doc *Document) normalizeTfIdf(invDocFreq termMap) {
	// Set the TF-IDF weights
	doc.tfIdf = make(termMap)
	for term, weight := range doc.termFreq {
		doc.tfIdf[term] = weight * invDocFreq[term]
	}

	// Now that we've set TF-IDF weights, we can save memory by removing the
	// original weights
	doc.termFreq = nil

	// Calculate and store the document's norm
	doc.norm = doc.calcNorm()
}

func (doc *Document) calcNorm() float64 {
	norm := 0.0
	for _, weight := range doc.tfIdf {
		norm += weight * weight
	}

	return math.Sqrt(norm)
}

func isTextFile(path string) bool {
	// If the file's extension is explicitly permitted, just use that.
	if hasPermittedExtension(path) {
		return true
	}

	// Try to get the file's MIME type from its extension, if that's
	// available.
	mimeType := mime.TypeByExtension(filepath.Ext(path))

	// If the MIME type can't be determined from the extension (maybe the file
	// lacks an extension, or maybe the local `mime.types` file doesn't exist)
	// instead try reading the first 4KB of the file and try to detect the content
	// type from that.
	if mimeType == "" {
		// Open the file for reading.
		file, err := os.Open(path)
		if err != nil {
			log.Fatalf("couldn't open file '%s': %s\n", path, err)
		}
		defer file.Close()

		// Read a byte slice of the file. Return false if the file's empty; there's
		// no reason to process it.
		data := make([]byte, 4096)
		_, err = file.Read(data)
		if err == io.EOF {
			return false
		} else if err != nil {
			log.Fatalf("couldn't read file '%s': %s\n", path, err)
		}
		mimeType = http.DetectContentType(data)
	}

	return strings.HasPrefix(mimeType, "text/")
}

func hasPermittedExtension(path string) bool {
	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	permitted, ok := explicitlyPermittedExtensions[ext]
	return ok && permitted
}
