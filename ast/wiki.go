// Package ast defines a wiki image AST node to represent the wiki extension's
// image element.
package ast

import (
	gast "github.com/yuin/goldmark/ast"
)

type WikiImage struct {
	gast.BaseInline

	// Destination is a destination(URL) of this image.
	Destination []byte
}

// Dump implements Node.Dump.
func (n *WikiImage) Dump(source []byte, level int) {
	m := map[string]string{}
	m["Destination"] = string(n.Destination)
	gast.DumpHelper(n, source, level, m, nil)
}

// KindWikiImage is a NodeKind of the WikiImage node.
var KindWikiImage = gast.NewNodeKind("WikiImage")

// Kind implements Node.Kind.
func (n *WikiImage) Kind() gast.NodeKind {
	return KindWikiImage
}

// NewWikiImage returns a new WikiImage node.
func NewWikiImage(dest []byte) *WikiImage {
	c := &WikiImage{
		BaseInline:  gast.BaseInline{},
		Destination: dest,
	}
	return c
}
