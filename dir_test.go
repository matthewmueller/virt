package virt_test

import (
	"io/fs"
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/virt"
)

func TestDuplicateDirEntriesNotHandled(t *testing.T) {
	is := is.New(t)
	dir := &virt.File{
		Mode: fs.ModeDir,
		Entries: []*virt.DirEntry{
			(&virt.File{Path: "a", Mode: fs.ModeDir}).Entry(),
			(&virt.File{Path: "a", Mode: fs.ModeDir}).Entry(),
		},
	}
	file := virt.To(dir)
	defer file.Close()
	readDir, ok := file.(fs.ReadDirFile)
	is.True(ok)
	des, err := readDir.ReadDir(-1)
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "a")
	is.Equal(des[1].Name(), "a")
}
