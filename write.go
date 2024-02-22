package virt

import (
	"io/fs"
	"path"
)

// Write a fileystem to a directory. Unlike sync, it does not attempt to remove
// files that are not in the source filesystem. This is sugar on top of WriteFS,
// commonly used for testing.
func Write(fsys fs.FS, toDir string, subpaths ...string) error {
	return WriteFS(fsys, OS(toDir), subpaths...)
}

// WriteFS writes files from one filesystem to another at subpath. Unlike sync,
// it does not attempt to remove files that are not in the source filesystem.
func WriteFS(from fs.FS, to FS, subpaths ...string) error {
	target := path.Join(subpaths...)
	if target == "" {
		target = "."
	}
	return fs.WalkDir(from, target, func(fpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if fpath == "." {
			return nil
		}
		if d.IsDir() {
			mode := d.Type()
			// Many of the virtual filesystems don't set a mode. Writing these to an
			// actual filesystem will cause permission errors, so we'll use common
			// permissions when not explicitly set.
			if mode == 0 || mode == fs.ModeDir {
				mode = 0755 | fs.ModeDir
			}
			return to.MkdirAll(fpath, mode)
		}
		data, err := fs.ReadFile(from, fpath)
		if err != nil {
			return err
		}
		// Many of the virtual filesystems don't set a mode. Writing these to an
		// actual filesystem will cause permission errors, so we'll use common
		// permissions when not explicitly set.
		mode := d.Type()
		if mode == 0 {
			mode = 0644
		}
		return to.WriteFile(fpath, data, mode)
	})
}
