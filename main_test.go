package main_test

import (
	"bytes"
	"strconv"
	"strings"
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

	indent int

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
			main.OutputConfig{Indent: strings.Repeat(" ", d.indent), EmptyElement: true},
		)

		seq := strconv.Itoa(i+1) + "/" + strconv.Itoa(len(data))

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
			{
				input: xmlpi + `<root><hoge/></root>`,
				xpath: `/root/hoge`,
				name:  `a`,
				out:   xmlpi + `<root><hoge><a/></hoge></root>`,
			},
			{ /*with text*/
				input: xmlpi + `<root><hoge/></root>`,
				xpath: `/root/hoge`,
				name:  `a`,
				value: `a text`,
				out:   xmlpi + `<root><hoge><a>a text</a></hoge></root>`,
			},
			{ /*with text*/
				input:   xmlpi + `<root><hoge/></root>`,
				xpath:   `/root/hoge`,
				name:    `a`,
				value:   `a text`,
				sibling: true,
				out:     xmlpi + `<root><hoge/><a>a text</a></root>`,
			},
			{
				input:   xmlpi + `<root><hoge/></root>`,
				xpath:   `/root/hoge`,
				name:    `a`,
				sibling: true,
				out:     xmlpi + `<root><hoge/><a/></root>`,
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
			{ /*text*/
				input:   xmlpi + `<root><hoge/></root>`,
				xpath:   `/root/hoge`,
				name:    `#text`,
				value:   `textvalue`,
				sibling: true,
				out:     xmlpi + `<root><hoge/>textvalue</root>`,
			},
			{ /*comment*/
				input: xmlpi + `<root><hoge/></root>`,
				xpath: `/root/hoge`,
				name:  `#comment`,
				value: `comment text`,
				out:   xmlpi + `<root><hoge><!--comment text--></hoge></root>`,
			},
			{ /*comment*/
				input:   xmlpi + `<root><hoge/></root>`,
				xpath:   `/root/hoge`,
				name:    `#comment`,
				value:   `comment text`,
				sibling: true,
				out:     xmlpi + `<root><hoge/><!--comment text--></root>`,
			},
			{ /*cdata-section*/
				input: xmlpi + `<root><hoge/></root>`,
				xpath: `/root/hoge`,
				name:  `#cdata-section`,
				value: `cdata text`,
				out:   xmlpi + `<root><hoge><![CDATA[cdata text]]></hoge></root>`,
			},
			{ /*cdata-section*/
				input:   xmlpi + `<root><hoge/></root>`,
				xpath:   `/root/hoge`,
				name:    `#cdata-section`,
				value:   `cdata text`,
				sibling: true,
				out:     xmlpi + `<root><hoge/><![CDATA[cdata text]]></root>`,
			},
			{ /*multiple*/
				input: xmlpi + `<root><hoge/><hoge/><hoge><z/></hoge></root>`,
				xpath: `/root/hoge`,
				name:  `a`,
				out:   xmlpi + `<root><hoge><a/></hoge><hoge><a/></hoge><hoge><z/><a/></hoge></root>`,
			},
			{ /*multiple*/
				input:  xmlpi + `<root><hoge/><hoge/><hoge><z/></hoge></root>`,
				xpath:  `/root/hoge`,
				name:   `a`,
				indent: 1,
				out: xmlpi + `
<root>
 <hoge>
  <a/>
 </hoge>
 <hoge>
  <a/>
 </hoge>
 <hoge>
  <z/>
  <a/>
 </hoge>
</root>
`,
			},
			{ /*multiple*/
				input:   xmlpi + `<root><hoge/><hoge/><hoge><z/></hoge></root>`,
				xpath:   `/root/hoge`,
				name:    `a`,
				sibling: true,
				out:     xmlpi + `<root><hoge/><a/><hoge/><a/><hoge><z/></hoge><a/></root>`,
				//out: xmlpi + `<root><hoge/><hoge/><hoge><z/></hoge><a/><a/><a/></root>`,
			},
			{ /*multiple*/
				input: xmlpi + `<root><hoge/><hoge/><hoge><z/></hoge></root>`,
				xpath: `/root/hoge`,
				name:  `@a`,
				value: `v`,
				out:   xmlpi + `<root><hoge a="v"/><hoge a="v"/><hoge a="v"><z/></hoge></root>`,
			},
			{ /*multiple*/
				input: xmlpi + `<root><hoge/><hoge/><hoge><z/></hoge></root>`,
				xpath: `/root/hoge`,
				name:  `@a`,
				value: `v`,
				out:   xmlpi + `<root><hoge a="v"/><hoge a="v"/><hoge a="v"><z/></hoge></root>`,
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
