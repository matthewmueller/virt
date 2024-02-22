package virt

import (
	"io/fs"
)

// Map is a simple in-memory filesystem.
// This filesytem is not safe for concurrent use.
type Map map[string]string

// Map only implements fs.FS because we can't make directories
// or store permission bits in a map.
var _ fs.FS = (Map)(nil)

func (m Map) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{
			Op:   "Open",
			Path: name,
			Err:  fs.ErrInvalid,
		}
	}
	return toTree(m).Open(name)
}

func toTree(m map[string]string) Tree {
	tree := Tree{}
	for path, data := range m {
		tree[path] = &File{Data: []byte(data)}
	}
	return tree
}
