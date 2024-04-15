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
	println(dump(parent))
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
	println(dump(parent))
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
			if styling {
				if !isOnelineText(n) {
					b.WriteByte('\n')
				}
			}
		}
		return
	case xmlquery.CharDataNode:
		b.WriteString("<![CDATA[")
		b.WriteString(n.Data)
		b.WriteString("]]>")
		if styling {
			b.WriteByte('\n')
		}
		return
	case xmlquery.CommentNode:
		b.WriteString("<!--")
		b.WriteString(n.Data)
		b.WriteString("-->")
		if styling {
			b.WriteByte('\n')
		}
		return
	case xmlquery.NotationNode:
		fmt.Fprintf(b, "<!%s>", n.Data)
		if styling {
			b.WriteByte('\n')
		}
		return
	case xmlquery.DeclarationNode:
		b.WriteString("<?" + n.Data)
	default:
		if n.Prefix == "" {
			b.WriteString("<" + n.Data)
		} else {
			fmt.Fprintf(b, "<%s:%s", n.Prefix, n.Data)
		}
	}

	for _, attr := range n.Attr {
		if attr.Name.Space != "" {
			fmt.Fprintf(b, ` %s:%s=`, attr.Name.Space, attr.Name.Local)
		} else {
			fmt.Fprintf(b, ` %s=`, attr.Name.Local)
		}
		b.WriteByte('"')
		b.WriteString(html.EscapeString(attr.Value))
		b.WriteByte('"')
	}
	if n.Type == xmlquery.DeclarationNode {
		b.WriteString("?>")
	} else {
		if n.FirstChild != nil || !config.EmptyElement {
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
		} else {
			b.WriteString("/>")
			if styling {
				b.WriteByte('\n')
			}
			return
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		outputXML(b, child, level+1, config)
	}
	if n.Type != xmlquery.DeclarationNode {
		if styling && !isOnelineText(n.FirstChild) {
			b.WriteString(strings.Repeat(config.Indent, level))
		}
		if n.Prefix == "" {
			fmt.Fprintf(b, "</%s>", n.Data)
		} else {
			fmt.Fprintf(b, "</%s:%s>", n.Prefix, n.Data)
		}
	}

	if styling {
		b.WriteByte('\n')
	}
}

func isOnelineText(n *xmlquery.Node) bool {
	return n == nil ||
		n.Type == xmlquery.TextNode &&
			strings.IndexByte(strings.TrimSpace(n.Data), '\n') == -1 &&
			n.NextSibling == nil &&
			n.PrevSibling == nil
}
