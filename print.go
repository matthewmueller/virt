package virt

import (
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/xlab/treeprint"
)

// Print out a virtual filesystem.
func Print(fsys fs.FS, subpaths ...string) (string, error) {
	dir := path.Join(subpaths...)
	if dir == "" {
		dir = "."
	}
	tree := treeprint.New()
	tree.SetValue(dir)
	err := fs.WalkDir(fsys, dir, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if path == "." {
			return nil
		}
		parent := tree
		for _, element := range strings.Split(filepath.ToSlash(path), "/") {
			existing := parent.FindByValue(element)
			if existing != nil {
				parent = existing
			} else {
				parent = parent.AddBranch(element)
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return tree.String(), nil
}
