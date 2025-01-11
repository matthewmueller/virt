package virt

import (
	"io"
	"io/fs"
	"time"
)

// FS is a virtual filesystem interface. It extends the fs.FS interface with
// methods for creating and removing files and directories.
type FS interface {
	fs.FS
	OpenFile(name string, flag int, perm fs.FileMode) (VFile, error)
	MkdirAll(path string, perm fs.FileMode) error
	WriteFile(name string, data []byte, perm fs.FileMode) error
	RemoveAll(path string) error
}

// VFile is a virtual file interface. It extends the fs.File interface with
// methods for writing to files.
type VFile interface {
	fs.File
	io.ReadCloser
}

// Now may be overridden for testing purposes
var Now = func() time.Time {
	return time.Now()
}
