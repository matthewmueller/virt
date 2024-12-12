package virt_test

import (
	"io/fs"
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/virt"
)

func TestMount(t *testing.T) {
	is := is.New(t)
	fsys := virt.Map{
		"a.txt":   "a",
		"b/b.txt": "b",
	}
	mounted := virt.Mount("nested/src", fsys)
	actual, err := virt.Print(mounted)
	is.NoErr(err)
	const expect = `.
└── nested
    └── src
        ├── a.txt
        └── b
            └── b.txt
`
	is.Equal(actual, expect)

	des, err := fs.ReadDir(mounted, "nested/src")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "a.txt")
	is.Equal(des[0].IsDir(), false)
	is.Equal(des[1].Name(), "b")
	is.Equal(des[1].IsDir(), true)

	des, err = fs.ReadDir(mounted, "nested/src/b")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "b.txt")
	is.Equal(des[0].IsDir(), false)

	data, err := fs.ReadFile(mounted, "nested/src/a.txt")
	is.NoErr(err)
	is.Equal(string(data), "a")

	data, err = fs.ReadFile(mounted, "nested/src/b/b.txt")
	is.NoErr(err)
	is.Equal(string(data), "b")
}

func TestMergeMount(t *testing.T) {
	is := is.New(t)
	a := virt.Map{
		"a.txt":   "a",
		"b/b.txt": "b",
	}
	c := virt.Map{
		"c.txt": "c",
	}
	fsys := virt.Merge(a, virt.Mount("nested/src", c))
	const expect = `.
├── a.txt
├── b
│   └── b.txt
└── nested
    └── src
        └── c.txt
`

	actual, err := virt.Print(fsys)
	is.NoErr(err)
	is.Equal(actual, expect)

	des, err := fs.ReadDir(fsys, ".")
	is.NoErr(err)

	is.Equal(len(des), 3)
	is.Equal(des[0].Name(), "a.txt")
	is.Equal(des[0].IsDir(), false)
	is.Equal(des[1].Name(), "b")
	is.Equal(des[1].IsDir(), true)
	is.Equal(des[2].Name(), "nested")
	is.Equal(des[2].IsDir(), true)
}
