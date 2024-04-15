package main_test

import (
	"bytes"
	"strconv"
	"testing"

	main "github.com/shu-go/eksemel"
	"github.com/shu-go/gotwant"
)

const xmlpi = `<?xml version="1.0" encoding="UTF-8"?>`

func prepare(input string) (in, out, errout *bytes.Buffer) {
	in = bytes.NewBufferString(input)
	out = &bytes.Buffer{}
	errout = &bytes.Buffer{}

	return in, out, errout
}

type addtestdata struct {
	input   string
	xpath   string
	name    string
	value   string
	ennet   string
	sibling bool

	err         error
	out, errout string
}

func testadd(t *testing.T, data []addtestdata) {
	t.Helper()

	for i, d := range data {
		in, out, errout := prepare(d.input)
		err := main.Add(
			main.NewFakeCloseReader(in),
			out,
			errout,
			d.xpath,
			d.name,
			d.value,
			d.ennet,
			d.sibling,
			main.OutputConfig{Indent: "", EmptyElement: true},
		)

		seq := strconv.Itoa(i)

		gotwant.TestError(t, err, d.err, gotwant.Desc(seq))
		gotwant.Test(t, out.String(), d.out, gotwant.Desc(seq))
		gotwant.Test(t, errout.String(), d.errout, gotwant.Desc(seq))
	}
}

func TestAdd(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		testadd(t, []addtestdata{
			{
				input: ``,
				xpath: `/`,
				ennet: `root>hoge`,
				out:   `<root><hoge/></root>`,
			},
			{
				input: xmlpi + `<root><hoge/></root>`,
				xpath: `/root/hoge`,
				ennet: `a>b`,
				out:   xmlpi + `<root><hoge><a><b/></a></hoge></root>`,
			},
			{
				input:   xmlpi + `<root><hoge/></root>`,
				xpath:   `/root/hoge`,
				ennet:   `a>b`,
				sibling: true,
				out:     xmlpi + `<root><hoge/><a><b/></a></root>`,
			},
			{ /*ignored*/
				input:   xmlpi + `<root><hoge/></root>`,
				xpath:   `/root/hoge/fuga`,
				ennet:   `a>b`,
				sibling: true,
				out:     xmlpi + `<root><hoge/></root>`,
			},
			{ /*attr*/
				input: xmlpi + `<root><hoge/></root>`,
				xpath: `/root/hoge`,
				name:  `@attr1`,
				value: `attrvalue1`,
				out:   xmlpi + `<root><hoge attr1="attrvalue1"/></root>`,
			},
			{ /*attr escaped*/
				input: xmlpi + `<root><hoge/></root>`,
				xpath: `/root/hoge`,
				name:  `@attr1`,
				value: `attrvalue"1'`,
				out:   xmlpi + `<root><hoge attr1="attrvalue&#34;1&#39;"/></root>`,
			},
			{ /*text*/
				input: xmlpi + `<root><hoge/></root>`,
				xpath: `/root/hoge`,
				name:  `#text`,
				value: `textvalue`,
				out:   xmlpi + `<root><hoge>textvalue</hoge></root>`,
			},
		})
	})

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
