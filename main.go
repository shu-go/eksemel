package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/go-xmlfmt/xmlfmt"
	"github.com/shu-go/gli/v2"
)

type globalCmd struct {
	Replace replaceCmd
	Delete  deleteCmd
	Add     addCmd
}

type replaceCmd struct {
	_ struct{} `help:"--xpath //* --value newvalue hoge.xml"`

	XPath string `cli:"xpath" required:"true"`
	Value string `cli:"value" required:"true"`
}

func (c replaceCmd) Run(args []string) error {
	if len(args) == 0 {
		return errors.New("input required")
	}

	f, err := os.Open(args[0])
	if err != nil {
		return err
	}
	doc, err := xmlquery.Parse(f)
	if err != nil {
		return err
	}
	f.Close()

	nodes, err := xmlquery.QueryAll(doc, c.XPath)
	if err != nil {
		return err
	}
	for _, n := range nodes {
		if n.Type == xmlquery.AttributeNode {
			n.Parent.SetAttr(n.Data, c.Value)
		} else {
			n.Data = c.Value
		}
	}

	fmt.Println(strings.TrimSpace(xmlfmt.FormatXML(doc.OutputXML(true), "", "  ", true)))

	return nil
}

type deleteCmd struct {
	_ struct{} `help:"--xpath //* hoge.xml"`

	XPath string `cli:"xpath" required:"true"`
}

func (c deleteCmd) Run(args []string) error {
	if len(args) == 0 {
		return errors.New("input required")
	}

	f, err := os.Open(args[0])
	if err != nil {
		return err
	}
	doc, err := xmlquery.Parse(f)
	if err != nil {
		return err
	}
	f.Close()

	nodes, err := xmlquery.QueryAll(doc, c.XPath)
	if err != nil {
		return err
	}
	for _, n := range nodes {
		if n.Type == xmlquery.AttributeNode {
			n.Parent.RemoveAttr(n.Data)
		} else {
			xmlquery.RemoveFromTree(n)
		}
	}

	fmt.Println(strings.TrimSpace(xmlfmt.FormatXML(doc.OutputXML(true), "", "  ", true)))

	return nil
}

type addCmd struct {
	_ struct{} `help:"--xpath //* --name newnode --value newvalue hoge.xml"`

	XPath   string `cli:"xpath" help:"parent" required:"true"`
	Name    string `cli:"name" required:"true" help:"nodename, @attrname, #text, #cdata-section, #comment"`
	Value   string `cli:"value" help:"if --name is a nodename, --value is its text"`
	Sibling bool   `cli:"sibling" default:"false" help:"as a LAST sibling"`
}

func (c addCmd) Run(args []string) error {
	if len(args) == 0 {
		return errors.New("input required")
	}

	f, err := os.Open(args[0])
	if err != nil {
		return err
	}
	doc, err := xmlquery.Parse(f)
	if err != nil {
		return err
	}
	f.Close()

	nodes, err := xmlquery.QueryAll(doc, c.XPath)
	if err != nil {
		return err
	}

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

	fmt.Println(strings.TrimSpace(xmlfmt.FormatXML(doc.OutputXML(true), "", "  ", true)))

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
