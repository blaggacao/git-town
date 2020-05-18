package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// fileData describes a memoized file
type fileData struct {
	relPath string      // the relative path of the file
	content []byte      // the file content
	perms   os.FileMode // the permissions
}

// DirectoryCloner creates many copies of a given directory efficiently.
type DirectoryCloner struct {
	files []fileData // files to create
	dirs  []fileData // directories to create
}

// NewDirectoryCloner provides a new directoryCloner instance ready to create clones of the given directory.
func NewDirectoryCloner(src string) (result *DirectoryCloner, err error) {
	result = &DirectoryCloner{}
	err = filepath.Walk(src, func(absPath string, fi os.FileInfo, e error) error {
		relPath, err := filepath.Rel(src, absPath)
		if err != nil {
			return fmt.Errorf("cannot make path %q relative to %q", absPath, src)
		}
		if relPath == "." {
			return nil
		}
		if fi.IsDir() {
			result.dirs = append(result.dirs, fileData{relPath: relPath, perms: fi.Mode()})
			return nil
		}
		content, err := ioutil.ReadFile(absPath)
		if err != nil {
			return fmt.Errorf("cannot read source file %q: %w", absPath, err)
		}
		result.files = append(result.files, fileData{relPath: relPath, content: content, perms: fi.Mode()})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create a new directoryCloner for %q: %w", src, err)
	}
	return result, nil
}

// CreateCopy creates a copy of the memoized files in the given directory.
func (dc *DirectoryCloner) CreateCopy(target string) error {
	err := os.Mkdir(target, 0744)
	if err != nil {
		return fmt.Errorf("cannot create directory %q", err)
	}
	for d := range dc.dirs {
		fullPath := filepath.Join(target, dc.dirs[d].relPath)
		err := os.Mkdir(fullPath, dc.dirs[d].perms)
		if err != nil {
			return fmt.Errorf("cannot create directory %q: %w", fullPath, err)
		}
	}
	for f := range dc.files {
		absPath := filepath.Join(target, dc.files[f].relPath)
		err := ioutil.WriteFile(absPath, dc.files[f].content, dc.files[f].perms)
		if err != nil {
			return fmt.Errorf("cannot create file %q: %w", absPath, err)
		}
	}
	return nil
}
