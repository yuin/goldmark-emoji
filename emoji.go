// package emoji is a extension for the goldmark(http://github.com/yuin/goldmark).
package emoji

import (
	"fmt"
	"io"
	"strings"

	"github.com/yuin/goldmark-emoji/v2/definition"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/text"
	"github.com/yuin/goldmark/v2/util"
)

// Emoji represents an inline emoji.
type Emoji struct {
	ast.BaseInline

	shortName string
	Value     *definition.Emoji
}

// Dump implements Node.Dump.
func (n *Emoji) Dump(_ []byte) *ast.NodeDump {
	return ast.NewNodeDump(n, map[string]any{
		"ShortName": n.shortName,
		"Value":     fmt.Sprintf("%#v", n.Value),
	})
}

// KindEmoji is a NodeKind of the emoji node.
var KindEmoji = ast.NewNodeKind("Emoji")

// Kind implements Node.Kind.
func (n *Emoji) Kind() ast.NodeKind {
	return KindEmoji
}

// NewEmoji returns a new Emoji node.
func NewEmoji(shortName string, value *definition.Emoji) *Emoji {
	n := &Emoji{
		shortName: shortName,
		Value:     value,
	}
	n.Init(n)
	return n
}

// RenderingMethod indicates how emojis are rendered.
type RenderingMethod int

// RendererFunc will be used for rendering emojis.
type RendererFunc func(w io.Writer, source []byte, n *Emoji, config *HTMLRendererConfig)

const (
	// Entity renders an emoji as an html entity.
	Entity RenderingMethod = iota

	// Unicode renders an emoji as unicode character.
	Unicode

	// Twemoji renders an emoji as an img tag with [twemoji](https://github.com/twitter/twemoji).
	Twemoji

	// Func renders an emoji using RendererFunc.
	Func
)

// HTMLRendererConfig struct holds options for the emoji renderer.
type HTMLRendererConfig struct {
	// Method indicates how emojis are rendered.
	Method RenderingMethod

	// XHTML indicates whether to render in XHTML style.
	XHTML bool

	// TwemojiTemplate is a printf template for twemoji. This value is valid only when Method is set to Twemoji.
	// `printf` arguments are:
	//
	//     1: name (e.g. "face with tears of joy")
	//     2: file name without an extension (e.g. 1f646-2642)
	//     3: '/' if XHTML, otherwise ''
	//
	TwemojiTemplate string

	// RendererFunc is a RendererFunc that renders emojis. This value is valid only when Method is set to Func.
	RendererFunc RendererFunc
}

// DefaultTwemojiTemplate is a default value for RendererConfig.TwemojiTemplate.
const DefaultTwemojiTemplate = `<img class="emoji" draggable="false" alt="%[1]s" src="https://cdn.jsdelivr.net/gh/twitter/twemoji@latest/assets/72x72/%[2]s.png"%[3]s>`

// A HTMLRendererOption interface sets options for the emoji renderer.
type HTMLRendererOption interface {
	setEmojiOption(*HTMLRendererConfig)
}

var _ HTMLRendererOption = &withRenderingMethod{}

type withRenderingMethod struct {
	value RenderingMethod
}

// setEmojiOption implements RendererOption#setEmojiOption
func (o *withRenderingMethod) setEmojiOption(c *HTMLRendererConfig) {
	c.Method = o.value
}

// WithRenderingMethod is a functional option that indicates how emojis are rendered.
func WithRenderingMethod(a RenderingMethod) HTMLRendererOption {
	return &withRenderingMethod{a}
}

var _ HTMLRendererOption = &withTwemojiTemplate{}

type withTwemojiTemplate struct {
	value string
}

// setEmojiOption implements RendererOption#setEmojiOption
func (o *withTwemojiTemplate) setEmojiOption(c *HTMLRendererConfig) {
	c.TwemojiTemplate = o.value
}

// WithTwemojiTemplate is a functional option that changes a twemoji img tag.
func WithTwemojiTemplate(s string) HTMLRendererOption {
	return &withTwemojiTemplate{s}
}

var _ HTMLRendererOption = &withRendererFunc{}

type withRendererFunc struct {
	value RendererFunc
}

// setEmojiOption implements RendererOption#setEmojiOption
func (o *withRendererFunc) setEmojiOption(c *HTMLRendererConfig) {
	c.RendererFunc = o.value
}

// WithRendererFunc is a functional option that changes a renderer func.
func WithRendererFunc(f RendererFunc) HTMLRendererOption {
	return &withRendererFunc{f}
}

// ParserConfig struct is a data structure that holds configuration of
// the Emoji extension.
type ParserConfig struct {
	Emojis definition.Emojis
}

// A ParserOption interface sets options for the emoji parser.
type ParserOption interface {
	SetEmojiOption(*ParserConfig)
}

var _ ParserOption = &withEmojis{}

type withEmojis struct {
	value definition.Emojis
}

func (o *withEmojis) SetEmojiOption(c *ParserConfig) {
	c.Emojis = o.value
}

// WithEmojis is a functional option that defines links names to unicode emojis.
func WithEmojis(value definition.Emojis) ParserOption {
	return &withEmojis{
		value: value,
	}
}

type emojiParser struct {
	ParserConfig
}

func newInlineParser(opts ...ParserOption) parser.InlineParser {
	p := &emojiParser{
		ParserConfig: ParserConfig{
			Emojis: definition.Github(),
		},
	}
	for _, o := range opts {
		o.SetEmojiOption(&p.ParserConfig)
	}
	return p
}

func (s *emojiParser) Trigger() []byte {
	return []byte{':'}
}

func (s *emojiParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if len(line) < 1 {
		return nil
	}
	i := 1
	for ; i < len(line); i++ {
		c := line[i]
		if !util.IsAlphaNumeric(c) && c != '_' && c != '-' && c != '+' {
			break
		}
	}
	if i >= len(line) || line[i] != ':' {
		return nil
	}
	block.Advance(i + 1)
	shortName := util.BytesToReadOnlyString(line[1:i])
	emoji, ok := s.Emojis.Get(shortName)
	if !ok {
		return nil
	}
	return NewEmoji(shortName, emoji)
}

type emojiHTMLRenderer struct {
	HTMLRendererConfig
}

// NewHTMLRenderer returns a new HTMLRenderer.
func NewHTMLRenderer(opts ...HTMLRendererOption) html.Extension {
	r := &emojiHTMLRenderer{
		HTMLRendererConfig: HTMLRendererConfig{
			Method:          Entity,
			TwemojiTemplate: DefaultTwemojiTemplate,
			RendererFunc:    nil,
		},
	}
	for _, opt := range opts {
		opt.setEmojiOption(&r.HTMLRendererConfig)
	}
	return r
}

// RendererOptions implements html.Extension.
func (r *emojiHTMLRenderer) RendererOptions(cfg *html.Config) []html.Option {
	r.XHTML = cfg.XHTML
	return []html.Option{
		html.WithNodeRenderers(map[ast.NodeKind]html.NodeRenderer{
			KindEmoji: html.NodeRendererFunc(r.renderEmoji),
		}),
	}
}

const slash = " /"
const empty = ""

func (r *emojiHTMLRenderer) renderEmoji(w io.Writer, source []byte, n ast.Node, entering bool, rc renderer.Context) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	bw := w.(util.BufWriter)
	node := n.(*Emoji)
	if !node.Value.IsUnicode() && r.Method != Func {
		_, _ = fmt.Fprintf(bw, `<span title="%s">:%s:</span>`, util.EscapeHTML(util.StringToReadOnlyBytes(node.Value.Name)), node.shortName)
		return ast.WalkContinue, nil
	}

	switch r.Method {
	case Entity:
		for _, ru := range node.Value.Unicode {
			if ru == 0x200D {
				_, _ = bw.WriteString("&zwj;")
				continue
			}
			_, _ = fmt.Fprintf(bw, "&#x%x;", ru)
		}
	case Unicode:
		_, _ = fmt.Fprintf(bw, "%s", string(node.Value.Unicode))
	case Twemoji:
		s := slash
		if !r.XHTML {
			s = empty
		}
		values := []string{}
		for _, ru := range node.Value.Unicode {
			values = append(values, fmt.Sprintf("%x", ru))
		}
		_, _ = fmt.Fprintf(bw, r.TwemojiTemplate, util.EscapeHTML(util.StringToReadOnlyBytes(node.Value.Name)), strings.Join(values, "-"), s)
	case Func:
		r.RendererFunc(w, source, node, &r.HTMLRendererConfig)
	}
	return ast.WalkContinue, nil
}

type parserExtension struct {
	parserOpts []ParserOption
}

// NewParser returns a new parser.Extension for emojis.
func NewParser(opts ...ParserOption) parser.Extension {
	return &parserExtension{parserOpts: opts}
}

// ParserOptions implements parser.Extension.
func (e *parserExtension) ParserOptions(_ *parser.Config) []parser.Option {
	return []parser.Option{
		parser.WithInlineParsers(
			util.Prioritized(newInlineParser(e.parserOpts...), 999),
		),
	}
}

// Parser is the default parser.Extension for emojis (using Github emoji set).
var Parser parser.Extension = NewParser()

// HTMLRenderer is the default html.Extension for emojis (Entity rendering method).
var HTMLRenderer html.Extension = NewHTMLRenderer()
