package virt

import (
	"errors"
	"io/fs"
	"os"
	"sort"
	"strings"
	"time"
)

// List is meant to be a simple list of files. It's not a tree of files
// and you can't walk it. Use Tree if you need a more capable filesystem.
// This filesytem is not safe for concurrent use.
type List []*File

var _ FS = (*List)(nil)

func (fsys List) Open(path string) (fs.File, error) {
	return fsys.OpenFile(path, os.O_RDONLY, 0)
}

func (fsys List) OpenFile(path string, flag int, perm fs.FileMode) (VFile, error) {
	if !fs.ValidPath(path) {
		return nil, &fs.PathError{Op: "openfile", Path: path, Err: fs.ErrInvalid}
	} else if flag != os.O_RDONLY {
		return nil, &fs.PathError{
			Op:   "openfile",
			Path: path,
			Err:  errors.New("flag not currently supported"),
		}
	}
	file, ok := fsys.find(path)
	if !ok {
		return nil, fs.ErrNotExist
	}
	// Found a file or directory
	file.Path = path
	if !file.IsDir() {
		return &openFile{file, flag, 0}, nil
	}
	// The following logic is based on "testing/fstest".MapFS.Open
	// Directory, possibly synthesized.
	// Note that file can be nil here: the map need not contain explicit parent directories for all its files.
	// But file can also be non-nil, in case the user wants to set metadata for the directory explicitly.
	// Either way, we need to construct the list of children of this directory.
	var des []*DirEntry
	var need = make(map[string]bool)
	if path == "." {
		for _, file := range fsys {
			fname := file.Path
			i := strings.Index(fname, "/")
			if i < 0 {
				if fname != "." {
					file.Path = fname
					des = append(des, file.Entry())
				}
			} else {
				need[fname[:i]] = true
			}
		}
	} else {
		prefix := path + "/"
		for _, file := range fsys {
			fname := file.Path
			if strings.HasPrefix(fname, prefix) {
				felem := fname[len(prefix):]
				i := strings.Index(felem, "/")
				if i < 0 {
					file.Path = felem
					des = append(des, file.Entry())
				} else {
					need[fname[len(prefix):len(prefix)+i]] = true
				}
			}
		}
		// If the directory name is not in the map,
		// and there are no children of the name in the map,
		// then the directory is treated as not existing.
		if file == nil && des == nil && len(need) == 0 {
			return nil, &fs.PathError{Op: "open", Path: path, Err: fs.ErrNotExist}
		}
	}
	for _, fi := range des {
		delete(need, fi.Name())
	}
	for path := range need {
		dir := &File{path, nil, fs.ModeDir, time.Time{}, nil}
		des = append(des, dir.Entry())
	}
	sort.Slice(des, func(i, j int) bool {
		return des[i].Name() < des[j].Name()
	})
	// Create a new directory if it wasn't found previously.
	if file == nil {
		file = &File{path, nil, fs.ModeDir, time.Time{}, nil}
	}
	// Return the synthesized entries as a directory.
	file.Entries = des
	return &openDir{file, flag, 0}, nil
}

// Mkdir create a directory.
func (fsys *List) MkdirAll(path string, perm fs.FileMode) error {
	if path == "." {
		return nil
	}
	file, ok := fsys.find(path)
	if ok {
		if file.IsDir() {
			return nil
		}
		return &fs.PathError{
			Op:   "MkdirAll",
			Path: path,
			Err:  fs.ErrExist,
		}
	}
	*fsys = append(*fsys, &File{path, nil, perm | fs.ModeDir, time.Time{}, nil})
	return nil
}

// WriteFile writes a file
func (fsys *List) WriteFile(path string, data []byte, perm fs.FileMode) error {
	file, ok := fsys.find(path)
	if ok {
		if file.IsDir() {
			return &fs.PathError{
				Op:   "WriteFile",
				Path: path,
				Err:  fs.ErrExist,
			}
		}
		file.Data = data
		file.Mode = perm
		return nil
	}
	*fsys = append(*fsys, &File{path, data, perm, time.Time{}, nil})
	return nil
}

// Remove removes a path
func (fsys *List) RemoveAll(path string) error {
	idx := fsys.indexOf(path)
	if idx < 0 {
		return nil
	}
	*fsys = append((*fsys)[:idx], (*fsys)[idx+1:]...)
	return nil
}

func (l List) find(path string) (f *File, ok bool) {
	for _, file := range l {
		if file.Path == path {
			return file, true
		}
	}
	return nil, false
}

func (l List) indexOf(path string) (i int) {
	for i, file := range l {
		if file.Path == path {
			return i
		}
	}
	return -1
}
