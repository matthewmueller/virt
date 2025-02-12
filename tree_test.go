package virt_test

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/matryer/is"
	"github.com/matthewmueller/virt"
)

func TestTreeFSTree(t *testing.T) {
	is := is.New(t)
	fsys := virt.Tree{
		"bud/view/index.svelte": &virt.File{
			Data: []byte(`<h1>index</h1>`),
		},
		"bud/view/about/index.svelte": &virt.File{
			Data: []byte(`<h1>about</h1>`),
		},
	}
	err := fstest.TestFS(fsys, "bud/view/index.svelte", "bud/view/about/index.svelte")
	is.NoErr(err)
}

func TestTreeReadDir(t *testing.T) {
	is := is.New(t)
	fsys := virt.Tree{
		"bud/view/index.svelte": &virt.File{
			Data: []byte(`<h1>index</h1>`),
		},
		"bud/view/about/index.svelte": &virt.File{
			Data: []byte(`<h1>about</h1>`),
		},
	}
	des, err := fs.ReadDir(fsys, "bud/view")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "about")
	is.Equal(des[1].Name(), "index.svelte")
}

func TestTreeTreeWriteRead(t *testing.T) {
	is := is.New(t)
	fsys := virt.Tree{}

	// Create a directory
	err := fsys.MkdirAll("a/b", 0755)
	is.NoErr(err)
	stat, err := fs.Stat(fsys, "a/b")
	is.NoErr(err)
	is.Equal(stat.Name(), "b")
	is.Equal(stat.IsDir(), true) // a/b should be a directory

	// Write a file
	err = fsys.WriteFile("a/b/c.txt", []byte("c"), 0644)
	is.NoErr(err)
	code, err := fs.ReadFile(fsys, "a/b/c.txt")
	is.NoErr(err)
	is.Equal(string(code), `c`)
}

func TestTreeTreeRoot(t *testing.T) {
	is := is.New(t)
	fsys := virt.Tree{}
	des, err := fs.ReadDir(fsys, ".")
	is.NoErr(err)
	is.Equal(len(des), 0)
}

func TestTreeTreeDirAndFileWithin(t *testing.T) {
	is := is.New(t)
	fsys := virt.Tree{}
	err := fsys.MkdirAll("controller/hello", 0755)
	is.NoErr(err)
	err = fsys.WriteFile("controller/hello/controller.go", []byte("package hello"), 0644)
	is.NoErr(err)
	paths := []string{}
	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		is.NoErr(err)
		paths = append(paths, path)
		return nil
	})
	is.NoErr(err)
	is.Equal(len(paths), 4)
	is.Equal(paths[0], ".")
	is.Equal(paths[1], "controller")
	is.Equal(paths[2], "controller/hello")
	is.Equal(paths[3], "controller/hello/controller.go")
}

func TestTreeTreeOpenParentInvalid(t *testing.T) {
	is := is.New(t)
	fsys := virt.Tree{
		"a/b.txt": &virt.File{Data: []byte("b")},
		"a":       &virt.File{Mode: fs.ModeDir},
	}
	file, err := fsys.Open("../a/b.txt")
	is.True(errors.Is(err, fs.ErrInvalid))
	is.Equal(file, nil)
}

func TestTreeStatWithPerm(t *testing.T) {
	is := is.New(t)
	fsys := virt.Tree{}
	fsys.MkdirAll("a/b", 0755)
	stat, err := fs.Stat(fsys, "a/b")
	is.NoErr(err)
	is.Equal(stat.Mode(), fs.FileMode(0755|fs.ModeDir))
}

func TestTreeRemoveDir(t *testing.T) {
	is := is.New(t)
	fsys := virt.Tree{}
	fsys.MkdirAll("a/b", 0755)
	stat, err := fs.Stat(fsys, "a/b")
	is.NoErr(err)
	is.Equal(stat.Mode(), fs.FileMode(0755|fs.ModeDir))
}

func TestTreeReadWriteDelete(t *testing.T) {
	is := is.New(t)
	fsys := virt.Tree{
		"duo/view/index.svelte": &virt.File{
			Data: []byte(`<h1>index</h1>`),
		},
	}

	// Read duo/view/index.svelte
	code, err := fs.ReadFile(fsys, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)

	// stat duo/
	stat, err := fs.Stat(fsys, "duo")
	is.NoErr(err)
	is.Equal(stat.Name(), "duo")
	is.Equal(stat.IsDir(), true)
	is.Equal(stat.Mode(), fs.FileMode(fs.ModeDir))

	// mkdir duo/controller
	err = fsys.MkdirAll("duo/controller", 0755)
	is.NoErr(err)
	stat, err = fs.Stat(fsys, "duo/controller")
	is.NoErr(err)
	is.Equal(stat.Name(), "controller")
	is.Equal(stat.IsDir(), true)
	is.Equal(stat.Mode(), fs.FileMode(0755|fs.ModeDir))

	// write duo/controller/controller.go
	err = fsys.WriteFile("duo/controller/controller.go", []byte(`package controller`), 0644)
	is.NoErr(err)

	// read duo/controller/controller.go
	code, err = fs.ReadFile(fsys, "duo/controller/controller.go")
	is.NoErr(err)
	is.Equal(string(code), `package controller`)

	// remove duo/view
	err = fsys.RemoveAll("duo/view")
	is.NoErr(err)

	// Read duo/view/index.svelte
	code, err = fs.ReadFile(fsys, "duo/view/index.svelte")
	is.Equal(errors.Is(err, fs.ErrNotExist), true)
	is.Equal(code, nil)
}

func TestTreeMkdirWriteChild(t *testing.T) {
	is := is.New(t)
	fsys := virt.Tree{}
	err := fsys.MkdirAll("a/b", 0755)
	is.NoErr(err)
	err = fsys.WriteFile("a/b/c.txt", []byte("c"), 0644)
	is.NoErr(err)
	stat, err := fs.Stat(fsys, "a/b")
	is.NoErr(err)
	is.Equal(stat.Name(), "b")
	is.Equal(stat.IsDir(), true)
	is.Equal(stat.Mode(), fs.FileMode(0755|fs.ModeDir))
}

func TestTreeSub(t *testing.T) {
	is := is.New(t)
	fsys := &virt.Tree{
		"bud/view/index.svelte": &virt.File{
			Data: []byte(`<h1>index</h1>`),
		},
		"bud/view/about/index.svelte": &virt.File{
			Data: []byte(`<h1>about</h1>`),
		},
	}
	sub, err := virt.Sub(fsys, "bud/view")
	is.NoErr(err)

	err = fstest.TestFS(sub, "index.svelte", "about/index.svelte")
	is.NoErr(err)

	des, err := fs.ReadDir(sub, ".")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "about")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[1].Name(), "index.svelte")
	is.Equal(des[1].IsDir(), false)

	// Read bud/view/index.svelte
	code, err := fs.ReadFile(sub, "index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)

	// Read about/index.svelte
	code, err = fs.ReadFile(sub, "about/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>about</h1>`)

	// Remove index.svelte
	err = sub.RemoveAll("index.svelte")
	is.NoErr(err)
	code, err = fs.ReadFile(sub, "index.svelte")
	is.Equal(errors.Is(err, fs.ErrNotExist), true)
	is.Equal(code, nil)

	err = sub.RemoveAll("../index.svelte")
	is.Equal(errors.Is(err, fs.ErrInvalid), true)
}

func TestTreeWalkBinary(t *testing.T) {
	is := is.New(t)
	fsys := &virt.Tree{
		"chunky": &virt.File{
			Data: []byte(`some executable`),
			Mode: 0755,
		},
	}
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		is.NoErr(err)
		if path == "." {
			return nil
		}
		is.Equal(d.Name(), "chunky")
		is.Equal(d.IsDir(), false)
		info, err := d.Info()
		is.NoErr(err)
		is.Equal(info.Mode(), fs.FileMode(0755))
		return nil
	})
	is.NoErr(err)
}
