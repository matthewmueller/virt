package virt

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func TestFromFile(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644))
	fsys := os.DirFS(dir)
	file, err := fsys.Open("a.txt")
	is.NoErr(err)
	defer file.Close()
	vfile, err := From("a.txt", file)
	is.NoErr(err)
	fi, err := file.Stat()
	is.NoErr(err)
	is.Equal(vfile.Path, "a.txt")
	is.Equal(vfile.Data, []byte("a"))
	is.Equal(vfile.ModTime, fi.ModTime())
	is.Equal(vfile.Mode, fi.Mode())
	is.Equal(vfile.Entry().Size, fi.Size())
}

func TestFromDir(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644))
	is.NoErr(os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0644))
	fsys := os.DirFS(dir)
	file, err := fsys.Open(".")
	is.NoErr(err)
	defer file.Close()
	vfile, err := From(".", file)
	is.NoErr(err)
	fi, err := file.Stat()
	is.NoErr(err)
	des, err := fs.ReadDir(fsys, ".")
	is.NoErr(err)
	is.Equal(vfile.Path, ".")
	is.Equal(vfile.ModTime, fi.ModTime())
	is.Equal(vfile.Mode, fi.Mode())
	is.Equal(vfile.Entry().Size, int64(0))
	is.Equal(len(vfile.Entries), len(des))
	for i, de := range vfile.Entries {
		is.Equal(de.Path, des[i].Name())
		info, err := des[i].Info()
		is.NoErr(err)
		is.Equal(de.ModTime, info.ModTime())
		is.Equal(de.Mode, info.Mode())
		is.Equal(de.Size, info.Size())
	}
}

func TestFromDirEntry(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644))
	fsys := os.DirFS(dir)
	des, err := fs.ReadDir(fsys, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)
	de := des[0]
	vde, err := FromEntry("a.txt", de)
	is.NoErr(err)
	is.Equal(vde.Path, "a.txt")
	is.Equal(de.Name(), vde.Name())
	is.Equal(de.Type(), vde.Type())
	is.Equal(de.IsDir(), vde.IsDir())
	dei, err := de.Info()
	is.NoErr(err)
	is.Equal(dei.ModTime(), vde.ModTime)
	is.Equal(dei.Mode(), vde.Mode)
	is.Equal(dei.Size(), vde.Size)
	vdei, err := vde.Info()
	is.NoErr(err)
	is.Equal(dei.ModTime(), vdei.ModTime())
	is.Equal(dei.Mode(), vdei.Mode())
	is.Equal(dei.Size(), vdei.Size())
}
