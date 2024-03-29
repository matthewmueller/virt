package virt_test

import (
	"io/fs"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/virt"
)

func TestExclude(t *testing.T) {
	is := is.New(t)
	tree := virt.Tree{
		"view/a.txt": &virt.File{Data: []byte("a")},
		"view/b.txt": &virt.File{Data: []byte("b")},
		"bud/bud.go": &virt.File{Data: []byte("bud")},
	}
	fsys := virt.Exclude(tree, func(path string) bool {
		return path == "bud" || strings.HasPrefix(path, "bud/")
	})
	des, err := fs.ReadDir(fsys, ".")
	is.Equal(err, nil)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "view")
}
