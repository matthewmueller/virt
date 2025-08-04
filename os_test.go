package virt_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/matryer/is"
	"github.com/matthewmueller/virt"
)

func TestOSRead(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644))
	is.NoErr(os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0644))
	is.NoErr(os.MkdirAll(filepath.Join(dir, "c"), 0755))
	is.NoErr(os.WriteFile(filepath.Join(dir, "c/c.txt"), []byte("d"), 0644))
	// Try reading the directory
	fsys := virt.OS(dir)
	err := fstest.TestFS(fsys, "a.txt", "b.txt", "c/c.txt")
	is.NoErr(err)
}

func TestOSRemoveAllOutsideFail(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644))
	is.NoErr(os.MkdirAll(filepath.Join(dir, "b"), 0755))
	fsys := virt.OS(filepath.Join(dir, "b"))
	err := fsys.RemoveAll("../a.txt")
	is.True(errors.Is(err, fs.ErrInvalid))
}

func TestOSRemoveAll(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644))
	fsys := virt.OS(dir)
	err := fsys.RemoveAll("a.txt")
	is.NoErr(err)
}

func TestTruncateDoesntChangeStat(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644))
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("b"), 0755))
	info, err := os.Stat(filepath.Join(dir, "a.txt"))
	is.NoErr(err)
	// Note that the mode is still 0644
	is.Equal(info.Mode(), os.FileMode(0644))
}

func TestSymlinkStat(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "to.txt"), []byte("to"), 0600))
	is.NoErr(os.Symlink(filepath.Join(dir, "to.txt"), filepath.Join(dir, "from.txt")))
	fsys := virt.OS(dir)

	// Direct Stat
	info, err := fsys.Stat("from.txt")
	is.NoErr(err)
	is.Equal(info.Name(), "from.txt")
	// Follows the symlink
	is.Equal(info.Mode(), fs.FileMode(0600))
	is.Equal(info.Size(), int64(2))

	// Through FS
	info, err = fs.Stat(fsys, "from.txt")
	is.NoErr(err)
	is.Equal(info.Name(), "from.txt")
	// Follows the symlink
	is.Equal(info.Mode(), fs.FileMode(0600))
	is.Equal(info.Size(), int64(2))
}

func TestSymlinkLstat(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(dir, "to.txt"), []byte("to"), 0600))
	is.NoErr(os.Symlink(filepath.Join(dir, "to.txt"), filepath.Join(dir, "from.txt")))
	fsys := virt.OS(dir)
	info, err := fsys.Lstat("from.txt")
	is.NoErr(err)
	is.Equal(info.Name(), "from.txt")
	// Doesn't follow the symlink
	is.Equal(info.Mode(), fs.FileMode(0755|fs.ModeSymlink))
	is.Equal(info.Size(), int64(len(filepath.Join(dir, "to.txt"))))
}
