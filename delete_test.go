package main_test

import (
	"strconv"
	"strings"
	"testing"

	main "github.com/shu-go/eksemel"
	"github.com/shu-go/gotwant"
)

type deletetestdata struct {
	input string
	xpath string

	indent int

	err         error
	out, errout string
}

func testdelete(t *testing.T, data []deletetestdata) {
	t.Helper()

	for i, d := range data {
		in, out, errout := prepare(d.input)
		err := main.Delete(
			main.NewFakeCloseReader(in),
			out,
			errout,
			d.xpath,
			main.OutputConfig{Indent: strings.Repeat(" ", d.indent), EmptyElement: true},
		)

		seq := strconv.Itoa(i+1) + "/" + strconv.Itoa(len(data))

		gotwant.TestError(t, err, d.err, gotwant.Desc(seq))
		gotwant.Test(t, readAll(out), d.out, gotwant.Desc(seq))
		gotwant.Test(t, readAll(errout), d.errout, gotwant.Desc(seq))
	}
}

func TestDelete(t *testing.T) {
	testdelete(t, []deletetestdata{
		{
			input: ``,
			xpath: `/`,
			out:   ``,
		},
		{
			input: xmlpi + `<root><hoge/></root>`,
			xpath: `/root/hoge`,
			out:   xmlpi + `<root/>`,
		},
		{
			input: xmlpi + `<root><hoge/><hoge/><hoge/></root>`,
			xpath: `/root/hoge`,
			out:   xmlpi + `<root/>`,
		},
		{
			input: xmlpi + `<root><hoge/><hoge/><hoge/></root>`,
			xpath: `/root/hoge[1]`,
			out:   xmlpi + `<root><hoge/><hoge/></root>`,
		},
		{
			input: xmlpi + `<root><hoge no="1"/><hoge no="2"/><hoge no="3"/></root>`,
			xpath: `/root/hoge[1]`,
			out:   xmlpi + `<root><hoge no="2"/><hoge no="3"/></root>`,
		},
		{
			input: xmlpi + `<root><hoge no="1"/><hoge no="2"/><hoge no="3"/></root>`,
			xpath: `/root/hoge[@no="2"]`,
			out:   xmlpi + `<root><hoge no="1"/><hoge no="3"/></root>`,
		},
		{
			input: xmlpi + `<root><hoge no="1"/><hoge no="2"/><fuga><hoge no="3"/></fuga></root>`,
			xpath: `//hoge`,
			out:   xmlpi + `<root><fuga/></root>`,
		},
	})
}
