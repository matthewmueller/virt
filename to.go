package virt

import "io/fs"

// To converts a *virt.File to an fs.File.
func To(f *File) fs.File {
	if f.Mode.IsDir() {
		return &openDir{f, 0}
	}
	return &openFile{f, 0}
}
