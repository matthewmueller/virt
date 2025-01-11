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
	OpenFile(name string, flag int, perm fs.FileMode) (RWFile, error)
	MkdirAll(path string, perm fs.FileMode) error
	WriteFile(name string, data []byte, perm fs.FileMode) error
	RemoveAll(path string) error
}

// RWFile is a virtual file interface. It extends fs.FS to support reading and
// writing files.
type RWFile interface {
	fs.File
	io.WriteCloser
}

// Now may be overridden for testing purposes
var Now = func() time.Time {
	return time.Now()
}
