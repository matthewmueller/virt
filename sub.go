package virt

import (
	"io/fs"
	"path"
)

// Sub returns a new filesystem rooted at dir.
func Sub(fsys FS, dir string) (FS, error) {
	if !fs.ValidPath(dir) {
		return nil, &fs.PathError{Op: "Sub", Path: dir, Err: fs.ErrInvalid}
	}
	return &subFS{dir, fsys}, nil
}

type subFS struct {
	dir string
	fs  FS
}

func (s *subFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "Open", Path: name, Err: fs.ErrInvalid}
	}
	return s.fs.Open(path.Join(s.dir, name))
}

func (s *subFS) MkdirAll(name string, perm fs.FileMode) error {
	if !fs.ValidPath(name) {
		return &fs.PathError{Op: "MkdirAll", Path: name, Err: fs.ErrInvalid}
	}
	return s.fs.MkdirAll(path.Join(s.dir, name), perm)
}

func (s *subFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	if !fs.ValidPath(name) {
		return &fs.PathError{Op: "WriteFile", Path: name, Err: fs.ErrInvalid}
	}
	return s.fs.WriteFile(path.Join(s.dir, name), data, perm)
}

func (s *subFS) RemoveAll(name string) error {
	if !fs.ValidPath(name) {
		return &fs.PathError{Op: "RemoveAll", Path: name, Err: fs.ErrInvalid}
	}
	return s.fs.RemoveAll(path.Join(s.dir, name))
}
