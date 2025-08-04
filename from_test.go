package virt

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func TestFrom(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644))
	fsys := OS(dir)
	file, err := fsys.Open("a.txt")
	is.NoErr(err)
	defer file.Close()
	vfile, err := From(fsys, "a.txt")
	is.NoErr(err)
	fi, err := file.Stat()
	is.NoErr(err)
	is.Equal(vfile.Path, "a.txt")
	is.Equal(vfile.Data, []byte("a"))
	is.Equal(vfile.ModTime, fi.ModTime())
	is.Equal(vfile.Mode, fi.Mode())
	is.Equal(vfile.Entry().Size, fi.Size())
}

func TestFromFile(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644))
	fsys := os.DirFS(dir)
	file, err := fsys.Open("a.txt")
	is.NoErr(err)
	defer file.Close()
	vfile, err := FromFile("a.txt", file)
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
	vfile, err := FromFile(".", file)
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

func TestFromSymlink(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "to.txt"), []byte("to"), 0600))
	is.NoErr(os.Symlink("to.txt", filepath.Join(dir, "from.txt")))

	fsys := OS(dir)
	file, err := From(fsys, "from.txt")
	is.NoErr(err)
	is.Equal(file.Path, "from.txt")
	is.Equal(file.Data, []byte("to.txt"))
	is.Equal(file.ModTime.IsZero(), false)
	is.Equal(file.Mode, fs.FileMode(0755|fs.ModeSymlink))
	is.Equal(file.Size(), int64(len("to.txt")))
}

func TestFromDirEntrySymlink(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "to.txt"), []byte("to"), 0600))
	is.NoErr(os.Symlink("to.txt", filepath.Join(dir, "from.txt")))

	fsys := os.DirFS(dir)
	des, err := fs.ReadDir(fsys, ".")
	is.NoErr(err)
	is.Equal(len(des), 2)

	is.Equal(des[0].Name(), "from.txt")
	info, err := des[0].Info()
	is.NoErr(err)
	is.Equal(info.Mode(), fs.FileMode(0755|fs.ModeSymlink))

	vde, err := FromEntry("from.txt", des[0])
	is.NoErr(err)
	is.Equal(vde.Path, "from.txt")
	is.Equal(vde.ModTime.IsZero(), false)
	is.Equal(vde.Mode, fs.FileMode(0755|fs.ModeSymlink))
	is.Equal(vde.Size, int64(len("to.txt")))
}
