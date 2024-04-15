package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/andrew-d/go-termutil"
	"github.com/antchfx/xmlquery"
	"github.com/shu-go/gli/v2"

	"github.com/shu-go/ennet"
)

type globalCmd struct {
	Replace replaceCmd
	Delete  deleteCmd
	Add     addCmd
}

type common struct {
	Indent       int  `cli:"indent=NUMBER" default:"4"`
	EmptyElement bool `cli:"empty" default:"true"`
}

type replaceCmd struct {
	_ struct{} `help:"--xpath //* --value newvalue hoge.xml"`

	XPath string `cli:"xpath" required:"true"`
	Value string `cli:"value"`
	Ennet string `cli:"ennet"`

	common
}

func (c replaceCmd) Before() error {
	if c.Value == "" && c.Ennet == "" {
		return errors.New("either --value or --ennet is required")
	}

	return nil
}

func Replace(input io.ReadCloser, output, errOutput io.Writer, xpath, value, abbrev string, config OutputConfig) error {
	doc, err := xmlquery.Parse(input)
	if err != nil {
		return err
	}
	input.Close()

	passthrough := false
	nodes, err := xmlquery.QueryAll(doc, xpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "xpath: %v\n", err)
		passthrough = true
	}

	if abbrev != "" {
		s, err := ennet.Expand(abbrev)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ennet: %v\n", err)
			passthrough = true
		}

		for _, n := range nodes {
			if passthrough {
				break
			}

			if n.Type != xmlquery.ElementNode {
				fmt.Fprintf(os.Stderr, "xpath: %v is not an element node\n", n.Data)
				break
			}

			b := bytes.NewBufferString(s)
			ewdoc, err := xmlquery.Parse(b)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ennet: %v\n", err)
				break
			}

			parent := n.Parent
			nn := ewdoc.FirstChild.NextSibling

			if parent.FirstChild == n {
				parent.FirstChild = nn
				nn.NextSibling = n.NextSibling
				if n.NextSibling != nil {
					n.NextSibling.PrevSibling = nn
				} else {
					// nn is the last
					parent.LastChild = nn
				}
			} else if parent.LastChild == n {
				parent.LastChild = nn
				nn.PrevSibling = n.PrevSibling
				if n.PrevSibling != nil {
					n.PrevSibling.NextSibling = nn
				}
			} else {
				nn.NextSibling = n.NextSibling
				nn.PrevSibling = n.PrevSibling
				n.PrevSibling.NextSibling = nn
				n.NextSibling.PrevSibling = nn
			}
		}

	} else {
		for _, n := range nodes {
			if n.Type == xmlquery.AttributeNode {
				n.Parent.SetAttr(n.Data, value)
			} else {
				n.Data = value
			}
		}
	}

	OutputXML(output, doc, config)

	return nil
}

func (c replaceCmd) Run(args []string) error {
	var input io.ReadCloser

	if !termutil.Isatty(os.Stdin.Fd()) {
		input = NewFakeCloseReader(os.Stdin)
	} else {
		if len(args) == 0 {
			return errors.New("input required")
		}

		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		input = f
	}

	return Replace(input, os.Stdout, os.Stderr, c.XPath, c.Value, c.Ennet, OutputConfig{Indent: strings.Repeat(" ", c.Indent), EmptyElement: c.EmptyElement})
}

type deleteCmd struct {
	_ struct{} `help:"--xpath //* hoge.xml"`

	XPath string `cli:"xpath" required:"true"`

	common
}

func Delete(input io.ReadCloser, output, errOutput io.Writer, xpath string, config OutputConfig) error {
	doc, err := xmlquery.Parse(input)
	if err != nil {
		return err
	}
	input.Close()

	nodes, err := xmlquery.QueryAll(doc, xpath)
	if err != nil {
		fmt.Fprintf(errOutput, "xpath: %v\n", err)
	} else {
		for _, n := range nodes {
			if n.Type == xmlquery.AttributeNode {
				n.Parent.RemoveAttr(n.Data)
			} else {
				xmlquery.RemoveFromTree(n)
			}
		}
	}

	OutputXML(output, doc, config)

	return nil
}

func (c deleteCmd) Run(args []string) error {
	var input io.ReadCloser

	if !termutil.Isatty(os.Stdin.Fd()) {
		input = NewFakeCloseReader(os.Stdin)
	} else {
		if len(args) == 0 {
			return errors.New("input required")
		}

		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		input = f
	}

	return Delete(input, os.Stdout, os.Stderr, c.XPath, OutputConfig{Indent: strings.Repeat(" ", c.Indent), EmptyElement: c.EmptyElement})
}

type addCmd struct {
	_ struct{} `help:"--xpath //* --name newnode --value newvalue hoge.xml"`

	XPath string `cli:"xpath" help:"parent" required:"true"`
	Name  string `cli:"name" help:"nodename, @attrname, #text, #cdata-section, #comment"`
	Value string `cli:"value" help:"if --name is a nodename, --value is its text"`

	Ennet string `cli:"ennet" help:"emmet-like abbreviation"`

	Sibling bool `cli:"sibling" default:"false" help:"as a LAST sibling"`

	common
}

func (c addCmd) Before() error {
	if c.Name == "" && c.Ennet == "" {
		return errors.New("either --name or --ennet is required")
	}

	return nil
}

func Add(input io.ReadCloser, output, errOutput io.Writer, xpath, name, value, abbrev string, sibling bool, config OutputConfig) error {
	doc, err := xmlquery.Parse(input)
	if err != nil {
		return fmt.Errorf("input: %w", err)
	}
	input.Close()

	passthrough := false
	nodes, err := xmlquery.QueryAll(doc, xpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "xpath: %v\n", err)
		passthrough = true
	}

	if abbrev != "" {
		s, err := ennet.Expand(abbrev)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ennet: %v\n", err)
			passthrough = true
		}

		for _, n := range nodes {
			if passthrough {
				break
			}

			b := bytes.NewBufferString(s)
			ewdoc, err := xmlquery.Parse(b)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ennet: %v\n", err)
				break
			}

			if sibling {
				xmlquery.AddSibling(n, ewdoc.FirstChild.NextSibling)
			} else {
				xmlquery.AddChild(n, ewdoc.FirstChild.NextSibling)
			}
		}
	} else {
		for _, n := range nodes {
			if strings.HasPrefix(name, "@") {
				xmlquery.AddAttr(n, name[1:], value)
			} else if name == "#text" {
				AddNode(n, &xmlquery.Node{
					Type: xmlquery.TextNode,
					Data: value,
				}, sibling)
			} else if name == "#cdata-section" {
				nn := &xmlquery.Node{
					Type: xmlquery.CharDataNode,
					Data: value,
				}
				xmlquery.AddChild(nn, &xmlquery.Node{
					Type: xmlquery.TextNode,
					Data: value,
				})
				AddNode(n, nn, sibling)

			} else if name == "#comment" {
				nn := &xmlquery.Node{
					Type: xmlquery.CommentNode,
					Data: value,
				}
				//xmlquery.AddChild(nn, &xmlquery.Node{
				//	Type: xmlquery.TextNode,
				//	Data: value,
				//})
				AddNode(n, nn, sibling)

			} else {
				nn := &xmlquery.Node{
					Data: name,
				}
				if value != "" {
					xmlquery.AddChild(nn, &xmlquery.Node{
						Type: xmlquery.TextNode,
						Data: value,
					})
				}
				AddNode(n, nn, sibling)
			}
		}

	}

	OutputXML(output, doc, config)

	return nil

}
func (c addCmd) Run(args []string) error {
	var input io.ReadCloser

	if !termutil.Isatty(os.Stdin.Fd()) {
		input = NewFakeCloseReader(os.Stdin)
	} else {
		if len(args) == 0 {
			return errors.New("input required")
		}

		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		input = f
	}

	return Add(input, os.Stdout, os.Stderr, c.XPath, c.Name, c.Value, c.Ennet, c.Sibling, OutputConfig{Indent: strings.Repeat(" ", c.Indent), EmptyElement: c.EmptyElement})
}

// Version is app version
var Version string

func main() {
	app := gli.NewWith(&globalCmd{})
	app.Name = "eksemel"
	app.Desc = "An XML file manipulator"
	app.Version = Version
	app.Usage = ``
	app.Copyright = "(C) 2024 Shuhei Kubota"
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func dump(n *xmlquery.Node) string {
	return dumpInner(n, 0)
}

func dumpInner(n *xmlquery.Node, indent int) string {
	s := strings.Repeat(" ", indent*2)
	if n.Type == xmlquery.ElementNode {
		s += n.Data
	} else if n.Type == xmlquery.TextNode {
		s += `"` + n.Data + `"`
	}

	if len(n.Attr) > 0 {
		slices.SortFunc(n.Attr, func(a, b xmlquery.Attr) int {
			return strings.Compare(a.Name.Local, b.Name.Local)
		})
		for _, a := range n.Attr {
			s += " @" + a.Name.Local + "=" + a.Value
		}
	}

	child := n.FirstChild
	for child != nil {
		s += "\n" + dumpInner(child, indent+1)
		child = child.NextSibling
	}
	return s
}
