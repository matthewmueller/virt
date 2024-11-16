package virt_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/virt"
)

func TestWriteFSEmpty(t *testing.T) {
	is := is.New(t)
	from := virt.Tree{}
	to := virt.Tree{}
	is.NoErr(virt.WriteFS(from, to))
	is.Equal(len(from), 0)
	is.Equal(len(to), 0)
}

func TestWriteFSEmptySub(t *testing.T) {
	is := is.New(t)
	from := virt.Tree{}
	to := virt.Tree{}
	err := virt.WriteFS(from, to, "sub", "dir")
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(len(from), 0)
	is.Equal(len(to), 0)
}

func TestWriteFSFromTo(t *testing.T) {
	is := is.New(t)
	from := virt.Tree{
		"sub/dir/a.txt": &virt.File{Data: []byte("a")},
		"sub/dir/b.txt": &virt.File{Data: []byte("b")},
		"sub/c.txt":     &virt.File{Data: []byte("c")},
		"d.txt":         &virt.File{Data: []byte("d")},
	}
	to := virt.Tree{
		"e.txt": &virt.File{Data: []byte("e")},
	}
	err := virt.WriteFS(from, to)
	is.NoErr(err)
	is.Equal(len(from), 4)
	is.Equal(len(to), 7)
	code, err := fs.ReadFile(to, "sub/dir/a.txt")
	is.NoErr(err)
	is.Equal(string(code), "a")
	code, err = fs.ReadFile(to, "sub/dir/b.txt")
	is.NoErr(err)
	is.Equal(string(code), "b")
	code, err = fs.ReadFile(to, "sub/c.txt")
	is.NoErr(err)
	is.Equal(string(code), "c")
	code, err = fs.ReadFile(to, "d.txt")
	is.NoErr(err)
	is.Equal(string(code), "d")
	code, err = fs.ReadFile(to, "e.txt")
	is.NoErr(err)
	is.Equal(string(code), "e")
}

func TestWriteFSFromToSub(t *testing.T) {
	is := is.New(t)
	from := virt.Tree{
		"sub/dir/a.txt": &virt.File{Data: []byte("a")},
		"sub/dir/b.txt": &virt.File{Data: []byte("b")},
		"sub/c.txt":     &virt.File{Data: []byte("c")},
		"d.txt":         &virt.File{Data: []byte("d")},
	}
	to := virt.Tree{
		"sub/e.txt": &virt.File{Data: []byte("e")},
	}
	err := virt.WriteFS(from, to, "sub")
	is.NoErr(err)
	is.Equal(len(from), 4)
	is.Equal(len(to), 5)
	code, err := fs.ReadFile(to, "sub/dir/a.txt")
	is.NoErr(err)
	is.Equal(string(code), "a")
	code, err = fs.ReadFile(to, "sub/dir/b.txt")
	is.NoErr(err)
	is.Equal(string(code), "b")
	code, err = fs.ReadFile(to, "sub/c.txt")
	is.NoErr(err)
	is.Equal(string(code), "c")
	code, err = fs.ReadFile(to, "sub/e.txt")
	is.NoErr(err)
	is.Equal(string(code), "e")
}

func TestWrite(t *testing.T) {
	is := is.New(t)
	from := virt.Tree{
		"sub/dir/a.txt": &virt.File{Data: []byte("a")},
		"sub/dir/b.txt": &virt.File{Data: []byte("b")},
		"sub/c.txt":     &virt.File{Data: []byte("c")},
		"d.txt":         &virt.File{Data: []byte("d")},
	}
	toDir := t.TempDir()
	is.NoErr(os.WriteFile(filepath.Join(toDir, "e.txt"), []byte("e"), 0644))
	err := virt.Write(from, toDir)
	is.NoErr(err)
	des, err := os.ReadDir(toDir)
	is.NoErr(err)
	is.Equal(len(des), 3)
	is.Equal(des[0].Name(), "d.txt")
	is.Equal(des[1].Name(), "e.txt")
	is.Equal(des[2].Name(), "sub")
}

func TestWriteExecutable(t *testing.T) {
	is := is.New(t)
	from := virt.Tree{
		"exe": &virt.File{Data: []byte("exe"), Mode: 0755},
	}
	toDir := t.TempDir()
	err := virt.Write(from, toDir)
	is.NoErr(err)
	info, err := fs.Stat(virt.OS(toDir), "exe")
	is.NoErr(err)
	is.Equal(info.Mode(), fs.FileMode(0755))
}
