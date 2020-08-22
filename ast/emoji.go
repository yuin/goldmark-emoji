// Package ast defines AST nodes that represetns emoji extension's elements.
package ast

import (
	"github.com/yuin/goldmark-emoji/definition"
	gast "github.com/yuin/goldmark/ast"
)

// Emoji represents an inline emoji.
type Emoji struct {
	gast.BaseInline

	Value *definition.Emoji
}

// Dump implements Node.Dump.
func (n *Emoji) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

// KindEmoji is a NodeKind of the emoji node.
var KindEmoji = gast.NewNodeKind("Emoji")

// Kind implements Node.Kind.
func (n *Emoji) Kind() gast.NodeKind {
	return KindEmoji
}

// NewEmoji returns a new Emoji node.
func NewEmoji(value *definition.Emoji) *Emoji {
	return &Emoji{
		Value: value,
	}
}
