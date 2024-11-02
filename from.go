package virt

import (
	"io"
	"io/fs"
	"path"
	"sort"
)

// From a file to a virtual file
func From(path string, file fs.File) (entry *File, err error) {
	// Get the stats
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	// Copy the directory data over
	if stat.IsDir() {
		return fromDir(path, file, stat)
	}
	return fromFile(path, file, stat)
}

// FromEntry converts a fs.DirEntry to a virtual DirEntry
func FromEntry(path string, de fs.DirEntry) (*DirEntry, error) {
	info, err := de.Info()
	if err != nil {
		return nil, err
	}
	return &DirEntry{
		Path:    path,
		Mode:    info.Mode(),
		ModTime: info.ModTime(),
		Size:    info.Size(),
	}, nil
}

func fromDir(dirPath string, file fs.File, stat fs.FileInfo) (entry *File, err error) {
	vdir := &File{
		Path:    dirPath,
		ModTime: stat.ModTime(),
		Mode:    stat.Mode(),
	}
	if dir, ok := file.(fs.ReadDirFile); ok {
		des, err := dir.ReadDir(-1)
		if err != nil {
			return nil, err
		}
		for _, de := range des {
			vde, err := FromEntry(path.Join(dirPath, de.Name()), de)
			if err != nil {
				return nil, err
			}
			vdir.Entries = append(vdir.Entries, vde)
		}
		sort.Slice(vdir.Entries, func(i, j int) bool {
			return vdir.Entries[i].Name() < vdir.Entries[j].Name()
		})
	}
	return vdir, nil
}

func fromFile(path string, file fs.File, stat fs.FileInfo) (entry *File, err error) {
	// Read the data fully
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return &File{
		Path:    path,
		Data:    data,
		ModTime: stat.ModTime(),
		Mode:    stat.Mode(),
	}, nil
}
