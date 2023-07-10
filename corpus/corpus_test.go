package corpus

import (
	"io/fs"
	"os"
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
