package emoji_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	. "github.com/yuin/goldmark-emoji/v2"
	"github.com/yuin/goldmark-emoji/v2/definition"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/testutil"
	"github.com/yuin/goldmark/v2/util"
)

func newMarkdown(parserOpts []ParserOption, rendererOpts []HTMLRendererOption) testutil.MarkdownToStringFunc {
	p := parser.New(parser.WithExtensions(NewParser(parserOpts...)))
	r := html.New(html.WithExtensions(NewHTMLRenderer(rendererOpts...)))
	return testutil.NewMarkdownToStringFunc(p, r)
}

func newMarkdownWithHTMLOpts(parserOpts []ParserOption, rendererOpts []HTMLRendererOption, htmlOpts ...html.Option) testutil.MarkdownToStringFunc {
	p := parser.New(parser.WithExtensions(NewParser(parserOpts...)))
	allHTMLOpts := append([]html.Option{html.WithExtensions(NewHTMLRenderer(rendererOpts...))}, htmlOpts...)
	r := html.New(allHTMLOpts...)
	return testutil.NewMarkdownToStringFunc(p, r)
}

func TestOptions(t *testing.T) {
	markdown := newMarkdown(nil, nil)
	count := 0

	count++
	testutil.DoTestCase(markdown, testutil.MarkdownTestCase{
		No:          count,
		Description: "default",
		Markdown: strings.TrimSpace(`
		Lucky :ok_man:
		`),
		Expected: strings.TrimSpace(`
		<p>Lucky &#x1f646;&zwj;&#x2642;&#xfe0f;</p>
		`),
	}, t)

	markdown = newMarkdown(nil, []HTMLRendererOption{WithRenderingMethod(Twemoji)})

	count++
	testutil.DoTestCase(markdown, testutil.MarkdownTestCase{
		No:          count,
		Description: "twemoji(HTML5)",
		Markdown: strings.TrimSpace(`
		Lucky :joy:
		`),
		Expected: strings.TrimSpace(`
		<p>Lucky <img class="emoji" draggable="false" alt="face with tears of joy" src="https://cdn.jsdelivr.net/gh/twitter/twemoji@latest/assets/72x72/1f602.png"></p>
		`),
	}, t)

	markdown = newMarkdownWithHTMLOpts(
		nil, []HTMLRendererOption{WithRenderingMethod(Twemoji)},
		html.WithXHTML(),
	)

	count++
	testutil.DoTestCase(markdown, testutil.MarkdownTestCase{
		No:          count,
		Description: "twemoji(XHTML)",
		Markdown: strings.TrimSpace(`
		Lucky :joy:
		`),
		Expected: strings.TrimSpace(`
		<p>Lucky <img class="emoji" draggable="false" alt="face with tears of joy" src="https://cdn.jsdelivr.net/gh/twitter/twemoji@latest/assets/72x72/1f602.png" /></p>
		`),
	}, t)

	markdown = newMarkdown(nil, []HTMLRendererOption{
		WithRenderingMethod(Twemoji),
		WithTwemojiTemplate(`<img class="myclass" draggable="false" alt="%[1]s" src="https://cdn.jsdelivr.net/gh/twitter/twemoji@latest/assets/36x36/%[2]s.png"%[3]s>`),
	})

	count++
	testutil.DoTestCase(markdown, testutil.MarkdownTestCase{
		No:          count,
		Description: "twemoji with customized template",
		Markdown: strings.TrimSpace(`
		Lucky :joy:
		`),
		Expected: strings.TrimSpace(`
        <p>Lucky <img class="myclass" draggable="false" alt="face with tears of joy" src="https://cdn.jsdelivr.net/gh/twitter/twemoji@latest/assets/36x36/1f602.png"></p>
		`),
	}, t)

	markdown = newMarkdown([]ParserOption{
		WithEmojis(definition.NewEmojis(definition.NewEmoji(
			"Standing man",
			[]rune{0x1f9cd, 0x200d, 0x2642, 0xfe0f},
			"man_standing",
		))),
	}, nil)

	count++
	testutil.DoTestCase(markdown, testutil.MarkdownTestCase{
		No:          count,
		Description: "twemoji with customized emoji definitions",
		Markdown: strings.TrimSpace(`
		Lucky :joy: :man_standing:
		`),
		Expected: strings.TrimSpace(`
		<p>Lucky :joy: &#x1f9cd;&zwj;&#x2642;&#xfe0f;</p>
		`),
	}, t)

	markdown = newMarkdown([]ParserOption{
		WithEmojis(
			definition.Github(
				definition.WithEmojis(
					definition.NewEmoji(
						"Standing man",
						[]rune{0x1f9cd, 0x200d, 0x2642, 0xfe0f},
						"man_standing",
					),
				),
			),
		),
	}, nil)

	count++
	testutil.DoTestCase(markdown, testutil.MarkdownTestCase{
		No:          count,
		Description: "twemoji with github emojis that are customized",
		Markdown: strings.TrimSpace(`
		Lucky :joy: :man_standing:
		`),
		Expected: strings.TrimSpace(`
        <p>Lucky &#x1f602; &#x1f9cd;&zwj;&#x2642;&#xfe0f;</p>
		`),
	}, t)

	markdown = newMarkdown(
		[]ParserOption{
			WithEmojis(
				definition.NewEmojis(
					definition.NewEmoji("Fast parrot", nil, "fastparrot"),
				),
			),
		},
		[]HTMLRendererOption{
			WithRenderingMethod(Func),
			WithRendererFunc(func(w io.Writer, source []byte, n *Emoji, config *HTMLRendererConfig) {
				bw := w.(util.BufWriter)
				_, _ = fmt.Fprintf(bw, `<img class="emoji" alt="%s" src="https://cultofthepartyparrot.com/parrots/hd/%s.gif>`, n.Value.Name, n.Value.ShortNames[0])
			}),
		},
	)

	count++
	testutil.DoTestCase(markdown, testutil.MarkdownTestCase{
		No:          count,
		Description: "Using RendererFunc to render original emojis",
		Markdown: strings.TrimSpace(`
		:fastparrot:
		`),
		Expected: strings.TrimSpace(`
		<p><img class="emoji" alt="Fast parrot" src="https://cultofthepartyparrot.com/parrots/hd/fastparrot.gif></p>
		`),
	}, t)

	markdown = newMarkdown(nil, []HTMLRendererOption{WithRenderingMethod(Unicode)})

	count++
	testutil.DoTestCase(markdown, testutil.MarkdownTestCase{
		No:          count,
		Description: "unicode",
		Markdown: strings.TrimSpace(`
		Lucky :joy:
		`),
		Expected: strings.TrimSpace(`
		<p>Lucky 😂</p>
		`),
	}, t)

	markdown = newMarkdown(
		[]ParserOption{WithEmojis(definition.NewEmojis(definition.NewEmoji("Fast parrot", nil, "fastparrot")))},
		[]HTMLRendererOption{WithRenderingMethod(Twemoji)},
	)

	count++
	testutil.DoTestCase(markdown, testutil.MarkdownTestCase{
		No:          count,
		Description: "Non-unicode emoji in twemoji",
		Markdown: strings.TrimSpace(`
		:fastparrot:
		`),
		Expected: strings.TrimSpace(`
		<p><span title="Fast parrot">:fastparrot:</span></p>
		`),
	}, t)

	markdown = newMarkdown(
		[]ParserOption{WithEmojis(definition.NewEmojis(definition.NewEmoji("Fast parrot", nil, "fastparrot")))},
		[]HTMLRendererOption{WithRenderingMethod(Entity)},
	)

	count++
	testutil.DoTestCase(markdown, testutil.MarkdownTestCase{
		No:          count,
		Description: "Non-unicode emoji in entity",
		Markdown: strings.TrimSpace(`
		:fastparrot:
		`),
		Expected: strings.TrimSpace(`
		<p><span title="Fast parrot">:fastparrot:</span></p>
		`),
	}, t)

	markdown = newMarkdown(
		[]ParserOption{WithEmojis(definition.NewEmojis(definition.NewEmoji("Fast parrot", nil, "fastparrot")))},
		[]HTMLRendererOption{WithRenderingMethod(Unicode)},
	)

	count++
	testutil.DoTestCase(markdown, testutil.MarkdownTestCase{
		No:          count,
		Description: "Non-unicode emoji in unicode",
		Markdown: strings.TrimSpace(`
		:fastparrot:
		`),
		Expected: strings.TrimSpace(`
		<p><span title="Fast parrot">:fastparrot:</span></p>
		`),
	}, t)
}
