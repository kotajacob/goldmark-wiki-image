// Package wiki is an extension for goldmark.
// https://github.com/yuin/goldmark
package wiki

import (
	"bytes"

	"git.sr.ht/~kota/goldmark-wiki-image/ast"
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type wiki struct{}

// Wiki is a goldmark.Extender implementation.
var Wiki = &wiki{}

// New returns a new extension. Useless, but included for compatibility.
func New() goldmark.Extender {
	return &wiki{}
}

// Extend implements goldmark.Extender.
func (w *wiki) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewParser(), 199),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewHTMLRenderer(), 199),
	))
}

type wikiParser struct{}

// NewParser returns a new parser.InlineParser that can parse wiki link syntax.
func NewParser() parser.InlineParser {
	p := &wikiParser{}
	return p
}

// Trigger returns characters that trigger this parser.
func (p *wikiParser) Trigger() []byte {
	return []byte{'!'}
}

var (
	parseOpen  = []byte("![[")
	parseClose = []byte("]]")
)

// Parse a wiki image in the form: ![[click here]]
// "click here" will be both the link destination and label.
func (p *wikiParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	line, seg := block.PeekLine()
	if !bytes.HasPrefix(line, parseOpen) {
		return nil
	}

	stop := bytes.Index(line, parseClose)
	if stop < 0 {
		return nil // Link must close on the same line.
	}
	seg = text.NewSegment(seg.Start+3, seg.Start+stop)

	n := ast.NewWikiImage(block.Value(seg))
	if len(n.Destination) == 0 || seg.Len() == 0 {
		return nil // Ensure destination and label are not empty.
	}

	n.AppendChild(n, gast.NewTextSegment(seg))
	block.Advance(stop + 2)
	return n
}

type wikiHTMLRenderer struct{}

// NewHTMLRenderer returns a new HTMLRenderer.
func NewHTMLRenderer() renderer.NodeRenderer {
	r := &wikiHTMLRenderer{}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *wikiHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindWikiImage, r.renderWiki)
}

func (r *wikiHTMLRenderer) renderWiki(
	w util.BufWriter,
	source []byte,
	node gast.Node,
	entering bool,
) (gast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<img src=\"")
	} else {
		_, _ = w.WriteString("\" loading=\"lazy\" width=\"100%\">")
	}
	return gast.WalkContinue, nil
}
