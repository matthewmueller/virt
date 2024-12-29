package virt_test

import (
	"io"
	"io/fs"
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/virt"
)

func TestTo(t *testing.T) {
	is := is.New(t)
	file := virt.To(&virt.File{
		Path: "a.txt",
		Data: []byte("aaa"),
		Mode: 0644,
	})
	defer file.Close()
	data, err := io.ReadAll(file)
	is.NoErr(err)
	is.Equal(string(data), "aaa")
	stat, err := file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "a.txt")
	is.Equal(stat.Size(), int64(3))
	is.Equal(stat.Mode(), fs.FileMode(0644))
}
