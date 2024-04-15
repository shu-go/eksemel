package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
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
	Value string `cli:"value" required:"true"`

	common
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

	doc, err := xmlquery.Parse(input)
	if err != nil {
		return err
	}
	input.Close()

	nodes, err := xmlquery.QueryAll(doc, c.XPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "xpath: %v\n", err)
	} else {
		for _, n := range nodes {
			if n.Type == xmlquery.AttributeNode {
				n.Parent.SetAttr(n.Data, c.Value)
			} else {
				n.Data = c.Value
			}
		}
	}

	OutputXML(os.Stdout, doc, OutputConfig{Indent: strings.Repeat(" ", c.Indent), EmptyElement: c.EmptyElement})

	return nil
}

type deleteCmd struct {
	_ struct{} `help:"--xpath //* hoge.xml"`

	XPath string `cli:"xpath" required:"true"`

	common
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

	doc, err := xmlquery.Parse(input)
	if err != nil {
		return err
	}
	input.Close()

	nodes, err := xmlquery.QueryAll(doc, c.XPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "xpath: %v\n", err)
	} else {
		for _, n := range nodes {
			if n.Type == xmlquery.AttributeNode {
				n.Parent.RemoveAttr(n.Data)
			} else {
				xmlquery.RemoveFromTree(n)
			}
		}
	}

	OutputXML(os.Stdout, doc, OutputConfig{Indent: strings.Repeat(" ", c.Indent), EmptyElement: c.EmptyElement})

	return nil
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

	doc, err := xmlquery.Parse(input)
	if err != nil {
		return fmt.Errorf("input: %w", err)
	}
	input.Close()

	passthrough := false
	nodes, err := xmlquery.QueryAll(doc, c.XPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "xpath: %v\n", err)
		passthrough = true
	}

	if c.Ennet != "" {
		s, err := ennet.Expand(c.Ennet)
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

			if c.Sibling {
				xmlquery.AddSibling(n, ewdoc.FirstChild.NextSibling)
			} else {
				xmlquery.AddChild(n, ewdoc.FirstChild.NextSibling)
			}
		}
	} else {
		addfunc := xmlquery.AddChild
		if c.Sibling {
			addfunc = xmlquery.AddSibling
		}

		for _, n := range nodes {
			if strings.HasPrefix(c.Name, "@") {
				xmlquery.AddAttr(n, c.Name[1:], c.Value)
			} else if c.Name == "#text" {
				addfunc(n, &xmlquery.Node{
					Type: xmlquery.TextNode,
					Data: c.Value,
				})
			} else if c.Name == "#cdata-section" {
				nn := &xmlquery.Node{
					Type: xmlquery.CharDataNode,
					Data: c.Value,
				}
				addfunc(nn, &xmlquery.Node{
					Type: xmlquery.TextNode,
					Data: c.Value,
				})
				xmlquery.AddChild(n, nn)

			} else if c.Name == "#comment" {
				nn := &xmlquery.Node{
					Type: xmlquery.CommentNode,
					Data: c.Value,
				}
				//xmlquery.AddChild(nn, &xmlquery.Node{
				//	Type: xmlquery.TextNode,
				//	Data: c.Value,
				//})
				addfunc(n, nn)

			} else {
				nn := &xmlquery.Node{
					Data: c.Name,
				}
				if c.Value != "" {
					xmlquery.AddChild(nn, &xmlquery.Node{
						Type: xmlquery.TextNode,
						Data: c.Value,
					})
				}
				addfunc(n, nn)
			}
		}

	}

	OutputXML(os.Stdout, doc, OutputConfig{Indent: strings.Repeat(" ", c.Indent), EmptyElement: c.EmptyElement})

	return nil
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
