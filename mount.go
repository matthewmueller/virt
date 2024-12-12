package virt

import (
	"io/fs"
	"strings"
)

func Mount(dir string, fsys fs.FS) fs.FS {
	dirfs := Tree{dir: &File{Path: dir, Mode: fs.ModeDir}}
	return &mountFS{dir, fsys, dirfs}
}

type mountFS struct {
	dir   string
	fs    fs.FS
	dirfs fs.FS
}

func (m *mountFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "Open", Path: name, Err: fs.ErrInvalid}
	}

	prefix := m.dir + "/"

	// Lookup within the mounted filesystem.
	if name == m.dir {
		return m.fs.Open(".")
	} else if strings.HasPrefix(name, prefix) {
		return m.fs.Open(strings.TrimPrefix(name, prefix))
	}

	// Lookup within the generated filesystem.
	return m.dirfs.Open(name)
}
