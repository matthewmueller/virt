package virt

import (
	"io/fs"
	"time"
)

// FS is a virtual filesystem interface. It extends the fs.FS interface with
// methods for creating and removing files and directories.
type FS interface {
	fs.FS
	MkdirAll(path string, perm fs.FileMode) error
	WriteFile(name string, data []byte, perm fs.FileMode) error
	RemoveAll(path string) error
}

func Open(f *File) fs.File {
	if f.Mode.IsDir() {
		return &openDir{f, 0}
	}
	return &openFile{f, 0}
}

// Now may be overridden for testing purposes
var Now = func() time.Time {
	return time.Now()
}
