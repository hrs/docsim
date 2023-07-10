package corpus

import (
	"io/fs"
	"math"
	"os"
	"strings"
	"testing"
)

func TestIsParsableFileRegular(t *testing.T) {
	f, err := os.CreateTemp("", "docsim-")
	if err != nil {
		t.Errorf("got unexpected error %v", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	fileInfo, err := f.Stat()
	if err != nil {
		t.Errorf("got unexpected error %v", err)
	}

	if !isParsableFile(fs.FileInfoToDirEntry(fileInfo), &Config{Stoplist: DefaultStoplist}) {
		t.Errorf("expected regular file %s to be parsable", f.Name())
	}
}

func TestIsParsableFileDirectory(t *testing.T) {
	name, err := os.MkdirTemp("", "docsim-")
	if err != nil {
		t.Errorf("got unexpected error %v", err)
	}
	defer os.RemoveAll(name)

	fileInfo, err := os.Lstat(name)
	if err != nil {
		t.Errorf("got unexpected error %v", err)
	}

	if isParsableFile(fs.FileInfoToDirEntry(fileInfo), &Config{Stoplist: DefaultStoplist}) {
		t.Errorf("expected directory %s not to be parsable", name)
	}
}

func TestIsParsableFileSymlink(t *testing.T) {
	// Create a target file to symlink to
	f, err := os.CreateTemp("", "docsim-")
	if err != nil {
		t.Errorf("got unexpected error %v", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	// Create and delete a tempfile to get a temporary name for the symlink
	symlinkFile, err := os.CreateTemp("", "docsim-")
	if err != nil {
		t.Errorf("got unexpected error %v", err)
	}
	symlinkName := symlinkFile.Name()
	symlinkFile.Close()
	os.Remove(symlinkName)

	// Create a symlink with that name targeting the original file
	err = os.Symlink(f.Name(), symlinkName)
	if err != nil {
		t.Errorf("got unexpected error %v", err)
	}
	defer os.Remove(symlinkName)

	fileInfo, err := os.Lstat(symlinkName)
	if err != nil {
		t.Errorf("got unexpected error %v", err)
	}

	// if !FollowSymlinks, the symlink SHOULDN'T be parsable
	if isParsableFile(fs.FileInfoToDirEntry(fileInfo), &Config{Stoplist: DefaultStoplist, FollowSymlinks: false}) {
		t.Errorf("with FollowSymlinks false, expected symlink %s not to be parsable", symlinkName)
	}

	// if FollowSymlinks, the symlink SHOULD be parsable
	if !isParsableFile(fs.FileInfoToDirEntry(fileInfo), &Config{Stoplist: DefaultStoplist, FollowSymlinks: true}) {
		t.Errorf("with FollowSymlinks true, expected symlink %s to be parsable", symlinkName)
	}
}

func makeDoc(s string, t *testing.T) *Document {
	config := Config{
		NoStemming: true,
		NoStoplist: true,
	}

	doc, err := NewDocument(strings.NewReader(s), &config)
	if err != nil {
		t.Errorf("got unexpected error creating document: %v", err)
	}

	return doc
}

func TestNewCorpus(t *testing.T) {
	docs := []*Document{
		makeDoc("a b c", t),
		makeDoc("a b", t),
		makeDoc("a", t),
	}

	corpus := NewCorpus(docs)

	tests := []struct {
		term string
		freq float64
	}{
		{"a", math.Log(3.0 / 3.0)},
		{"b", math.Log(3.0 / 2.0)},
		{"c", math.Log(3.0 / 1.0)},
	}

	for _, tc := range tests {
		got := corpus.invDocFreq[termID(tc.term)]
		if !approxEq(got, tc.freq) {
			t.Errorf("expected IDF(%s) to be %0.4f, got %0.4f", tc.term, tc.freq, got)
		}
	}
}
