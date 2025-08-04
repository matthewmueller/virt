package virt

import (
	"io/fs"
	"os"
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
	return s.OpenFile(name, os.O_RDONLY, 0)
}

func (s *subFS) Stat(name string) (fs.FileInfo, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "Stat", Path: name, Err: fs.ErrInvalid}
	}
	return s.fs.Stat(path.Join(s.dir, name))
}

func (s *subFS) Lstat(name string) (fs.FileInfo, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "Lstat", Path: name, Err: fs.ErrInvalid}
	}
	return s.fs.Lstat(path.Join(s.dir, name))
}

func (s *subFS) OpenFile(name string, flag int, perm fs.FileMode) (RWFile, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "OpenFile", Path: name, Err: fs.ErrInvalid}
	}
	return s.fs.OpenFile(path.Join(s.dir, name), flag, perm)
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
