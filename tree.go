package virt

import (
	"errors"
	"io/fs"
	"os"
	"sort"
	"strings"
	"time"
)

// File represents a file or directory in a virtual filesystem. Unlike virt.Map,
// the Tree filesystem can be traversed (similar to fstest.MapFS) and written to.
type Tree map[string]*File

var _ FS = (Tree)(nil)

func (fsys Tree) Open(path string) (fs.File, error) {
	return fsys.OpenFile(path, os.O_RDONLY, 0)
}

func (fsys Tree) OpenFile(path string, flag int, perm fs.FileMode) (RWFile, error) {
	if !fs.ValidPath(path) {
		return nil, &fs.PathError{Op: "Open", Path: path, Err: fs.ErrInvalid}
	} else if flag != os.O_RDONLY {
		return nil, &fs.PathError{
			Op:   "openfile",
			Path: path,
			Err:  errors.New("flag not currently supported"),
		}
	}
	file, ok := fsys[path]
	if ok {
		// Can be either a file or a empty directory
		file.Path = path
		if !file.IsDir() {
			return &openFile{file, flag, 0}, nil
		}
	}

	// The following logic is based on "testing/fstest".MapFS.Open
	// Directory, possibly synthesized.
	// Note that file can be nil here: the map need not contain explicit parent directories for all its files.
	// But file can also be non-nil, in case the user wants to set metadata for the directory explicitly.
	// Either way, we need to construct the list of children of this directory.
	var des []*DirEntry
	var need = make(map[string]bool)
	if path == "." {
		for fname, file := range fsys {
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
		for fname, file := range fsys {
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
func (t Tree) MkdirAll(path string, perm fs.FileMode) error {
	if !fs.ValidPath(path) {
		return &fs.PathError{Op: "MkdirAll", Path: path, Err: fs.ErrInvalid}
	} else if path == "." {
		return nil
	}
	// Don't create a directory unless we have to
	if _, err := fs.Stat(t, path); nil == err {
		return nil
	}
	t[path] = &File{path, nil, perm | fs.ModeDir, Now(), nil}
	return nil
}

// WriteFile writes a file
// TODO: WriteFile should fail if path.Dir(name) doesn't exist
func (t Tree) WriteFile(path string, data []byte, perm fs.FileMode) error {
	if !fs.ValidPath(path) {
		return &fs.PathError{Op: "WriteFile", Path: path, Err: fs.ErrInvalid}
	}
	t[path] = &File{path, data, perm, Now(), nil}
	return nil
}

// Remove removes a path
func (t Tree) RemoveAll(path string) error {
	if !fs.ValidPath(path) {
		return &fs.PathError{Op: "RemoveAll", Path: path, Err: fs.ErrInvalid}
	}
	stat, err := fs.Stat(t, path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	// Delete the path
	delete(t, path)
	// Only delete the file
	if !stat.IsDir() {
		return nil
	}
	// Need to delete the rest of the files
	dirpath := path + "/"
	for fpath := range t {
		if strings.HasPrefix(fpath, dirpath) {
			delete(t, fpath)
		}
	}
	return nil
}
