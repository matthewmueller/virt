package virt

import (
	"io/fs"
	"os"
	"path/filepath"
)

// OS creates a new OS filesystem rooted at the given directory.
// TODO: create an os_windows for opening on multiple drives
// with the same API:
// https://github.com/golang/go/issues/44279#issuecomment-955766528
type OS string

var _ FS = (OS)("")

func (dir OS) Open(name string) (fs.File, error) {
	return os.DirFS(string(dir)).Open(name)
}

func (dir OS) Stat(name string) (fs.FileInfo, error) {
	return fs.Stat(os.DirFS(string(dir)), name)
}

func (dir OS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(os.DirFS(string(dir)), name)
}

func (dir OS) Lstat(name string) (fs.FileInfo, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "Lstat", Path: name, Err: fs.ErrInvalid}
	}
	return os.Lstat(filepath.Join(string(dir), name))
}

func (dir OS) Readlink(name string) (string, error) {
	if !fs.ValidPath(name) {
		return "", &fs.PathError{Op: "Readlink", Path: name, Err: fs.ErrInvalid}
	}
	return os.Readlink(filepath.Join(string(dir), name))
}

func (dir OS) OpenFile(name string, flag int, perm fs.FileMode) (RWFile, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "OpenFile", Path: name, Err: fs.ErrInvalid}
	}
	return os.OpenFile(filepath.Join(string(dir), name), flag, perm)
}

func (dir OS) MkdirAll(path string, perm fs.FileMode) error {
	if !fs.ValidPath(path) {
		return &fs.PathError{Op: "mkdirall", Path: path, Err: fs.ErrInvalid}
	}
	return os.MkdirAll(filepath.Join(string(dir), path), perm)
}

func (dir OS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	if !fs.ValidPath(name) {
		return &fs.PathError{Op: "WriteFile", Path: name, Err: fs.ErrInvalid}
	}
	if perm&fs.ModeSymlink != 0 {
		return os.Symlink(string(data), filepath.Join(string(dir), name))
	}
	return os.WriteFile(filepath.Join(string(dir), name), data, perm)
}

func (dir OS) RemoveAll(path string) error {
	if !fs.ValidPath(path) {
		return &fs.PathError{Op: "RemoveAll", Path: path, Err: fs.ErrInvalid}
	}
	return os.RemoveAll(filepath.Join(string(dir), path))
}
