package main_test

import (
	"bytes"
	"testing"

	main "github.com/shu-go/eksemel"
	"github.com/shu-go/gotwant"
)

func TestAdd(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}
		errout := &bytes.Buffer{}
		main.Add(
			/*input*/ main.NewFakeCloseReader(in),
			/*output*/ out,
			/*errOutput*/ errout,
			/*xpath*/ `/`,
			/*name*/ ``,
			/*value*/ ``,
			/*ennet*/ `root>hoge`,
			/*sibling*/ false,
			main.OutputConfig{Indent: "", EmptyElement: true},
		)

		gotwant.Test(t, out.String(), `<root><hoge/></root>`)
	})
}
