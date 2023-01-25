package wiki

import (
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/testutil"
)

func TestWiki(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			Wiki,
		),
	)
	count := 0

	count++
	testutil.DoTestCase(markdown, testutil.MarkdownTestCase{
		No:          count,
		Description: "default",
		Markdown: strings.TrimSpace(`
		![[hello]]
		`),
		Expected: strings.TrimSpace(`
		<p><img src="hello" loading="lazy" width="100%"></p>
		`),
	}, t)
}
