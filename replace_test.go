package main_test

import (
	"strconv"
	"strings"
	"testing"

	main "github.com/shu-go/eksemel"
	"github.com/shu-go/gotwant"
)

type replacetestdata struct {
	input string
	xpath string
	value string
	ennet string

	indent int

	err         error
	out, errout string
}

func testreplace(t *testing.T, data []replacetestdata) {
	t.Helper()

	for i, d := range data {
		in, out, errout := prepare(d.input)
		err := main.Replace(
			main.NewFakeCloseReader(in),
			out,
			errout,
			d.xpath,
			d.value,
			d.ennet,
			main.OutputConfig{Indent: strings.Repeat(" ", d.indent), EmptyElement: true},
		)

		seq := strconv.Itoa(i+1) + "/" + strconv.Itoa(len(data))

		gotwant.TestError(t, err, d.err, gotwant.Desc(seq))
		gotwant.Test(t, readAll(out), d.out, gotwant.Desc(seq))
		gotwant.Test(t, readAll(errout), d.errout, gotwant.Desc(seq))
	}
}

func TestReplace(t *testing.T) {
	testreplace(t, []replacetestdata{
		{
			input: ``,
			xpath: `/root`,
			ennet: `root>hoge`,
			out:   ``,
		},
		{
			input: xmlpi + `<root><fuga/></root>`,
			xpath: `/root`,
			ennet: `root>hoge`,
			out:   xmlpi + `<root><hoge/></root>`,
		},
		{
			input: xmlpi + `<root><hoge/></root>`,
			xpath: `/root`,
			value: "loot",
			out:   xmlpi + `<loot><hoge/></loot>`,
		},
		{
			input: xmlpi + `<root><hoge/><hoge/><fuga><hoge/></fuga></root>`,
			xpath: `/root/hoge`,
			value: "moge",
			out:   xmlpi + `<root><moge/><moge/><fuga><hoge/></fuga></root>`,
		},
		{
			input: xmlpi + `<root><hoge/><hoge/><fuga><hoge/></fuga></root>`,
			xpath: `//hoge`,
			value: "moge",
			out:   xmlpi + `<root><moge/><moge/><fuga><moge/></fuga></root>`,
		},
		{
			input: xmlpi + `<root><!--comment--></root>`,
			xpath: `/root/comment()`,
			value: "komento",
			out:   xmlpi + `<root><!--komento--></root>`,
		},
	})
}
