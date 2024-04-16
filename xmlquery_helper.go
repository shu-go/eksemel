package main

import (
	"bufio"
	"fmt"
	"html"
	"io"
	"strings"

	"github.com/antchfx/xmlquery"
)

type OutputConfig struct {
	EmptyElement bool
	Indent       string
}

func OutputXML(out io.Writer, n *xmlquery.Node, config OutputConfig) {
	b := bufio.NewWriter(out)
	level := 0

	if n.Type == xmlquery.DocumentNode {
		curr := n.FirstChild
		for curr != nil {
			outputXML(b, curr, level, config)
			curr = curr.NextSibling
		}
	} else {
		outputXML(b, n, level, config)
	}
	b.Flush()
}

func AddNode(n, newnode *xmlquery.Node, sibling bool) {
	if !sibling {
		xmlquery.AddChild(n, newnode)
		return
	}

	parent := n.Parent
	if n == parent.FirstChild {
		if nnext := n.NextSibling; nnext != nil {
			newnode.NextSibling = nnext
			nnext.PrevSibling = newnode
		} else {
			parent.LastChild = newnode
		}
	} else if n == parent.LastChild {
		parent.LastChild = newnode
	} else {
		nnext := n.NextSibling
		newnode.NextSibling = nnext
		nnext.PrevSibling = newnode
	}

	newnode.PrevSibling = n
	n.NextSibling = newnode

	newnode.Parent = parent
}

func outputXML(b *bufio.Writer, n *xmlquery.Node, level int, config OutputConfig) {
	if n.Type == xmlquery.TextNode && strings.TrimSpace(n.Data) == "" {
		return
	}

	styling := config.Indent != ""
	if styling && !isOnelineText(n) {
		b.WriteString(strings.Repeat(config.Indent, level))
	}

	switch n.Type {
	case xmlquery.TextNode:
		text := strings.TrimSpace(n.Data)
		if text != "" {
			b.WriteString(html.EscapeString(text))
			if !isOnelineText(n) {
				writeStylingNewLine(b, styling)
			}
		}
		return
	case xmlquery.CharDataNode:
		b.WriteString("<![CDATA[")
		b.WriteString(n.Data)
		b.WriteString("]]>")
		writeStylingNewLine(b, styling)
		return
	case xmlquery.CommentNode:
		b.WriteString("<!--")
		b.WriteString(n.Data)
		b.WriteString("-->")
		writeStylingNewLine(b, styling)
		return
	case xmlquery.NotationNode:
		fmt.Fprintf(b, "<!%s>", n.Data)
		writeStylingNewLine(b, styling)
		return
	case xmlquery.DeclarationNode:
		b.WriteString("<?" + n.Data)
	default:
		b.WriteByte('<')
		writeName(b, n.Prefix, n.Data)
	}

	for _, attr := range n.Attr {
		b.WriteByte(' ')
		writeName(b, attr.Name.Space, attr.Name.Local)
		b.WriteByte('=')
		b.WriteByte('"')
		b.WriteString(html.EscapeString(attr.Value))
		b.WriteByte('"')
	}

	if n.Type == xmlquery.DeclarationNode {
		b.WriteString("?>")
		writeStylingNewLine(b, styling)
		return
	}

	if n.FirstChild == nil && config.EmptyElement {
		b.WriteString("/>")
		writeStylingNewLine(b, styling)
		return
	}

	b.WriteString(">")

	if styling {
		newline := false
		curr := n.FirstChild
		for curr != nil {
			if !isOnelineText(curr) {
				newline = true
				break
			}
			curr = curr.NextSibling
		}
		if newline {
			b.WriteByte('\n')
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		outputXML(b, child, level+1, config)
	}

	if n.Type != xmlquery.DeclarationNode {
		if styling && !isOnelineText(n.FirstChild) {
			b.WriteString(strings.Repeat(config.Indent, level))
		}
		b.WriteString("</")
		writeName(b, n.Prefix, n.Data)
		b.WriteByte('>')
	}

	writeStylingNewLine(b, styling)
}

func isOnelineText(n *xmlquery.Node) bool {
	return n == nil ||
		n.Type == xmlquery.TextNode &&
			strings.IndexByte(strings.TrimSpace(n.Data), '\n') == -1 &&
			n.NextSibling == nil &&
			n.PrevSibling == nil
}

func writeName(b *bufio.Writer, space, name string) (int, error) {
	if space == "" {
		return b.WriteString(name)
	}
	b.WriteString(space)
	b.WriteByte(':')
	return b.WriteString(name)
}

func writeStylingNewLine(b *bufio.Writer, styling bool) error {
	if !styling {
		return nil
	}
	return b.WriteByte('\n')
}
