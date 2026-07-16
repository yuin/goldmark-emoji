goldmark-emoji
=========================

[![GoDev][godev-image]][godev-url]

[godev-image]: https://pkg.go.dev/badge/github.com/yuin/goldmark-emoji/v2
[godev-url]: https://pkg.go.dev/github.com/yuin/goldmark-emoji/v2

goldmark-emoji is an extension for [goldmark](http://github.com/yuin/goldmark) v2
that parses `:joy:` style emojis.

## Compatiblity
github.com/yuin/goldmark-emoji/v2 is compatible with goldmark(v2).

use `github.com/yuin/goldmark-meta` for goldmark(v1).

## Installation

```sh
go get github.com/yuin/goldmark-emoji/v2
```

## Usage

```go
import (
    "bytes"
    "fmt"

    "github.com/yuin/goldmark-emoji/v2"
    "github.com/yuin/goldmark/v2/parser"
    "github.com/yuin/goldmark/v2/renderer/html"
)

func main() {
    p := parser.New(parser.WithExtensions(emoji.Parser))
    r := html.New(html.WithExtensions(emoji.HTMLRenderer))
    source := []byte(`Joy :joy:`)
    doc := p.Parse(source)
    var buf bytes.Buffer
    if err := r.Render(&buf, source, doc); err != nil {
        panic(err)
    }
    fmt.Print(buf.String())
}
```

With options:

```go
p := parser.New(parser.WithExtensions(
    emoji.NewParser(emoji.WithEmojis(myEmojis)),
))
r := html.New(html.WithExtensions(
    emoji.NewHTMLRenderer(emoji.WithRenderingMethod(emoji.Twemoji)),
))
```

See `emoji_test.go` for detailed usage.

### Options

goldmark/v2 separates parser and renderer extensions, so options are split accordingly.

**Parser options** (passed to `NewParser`):

| Option | Description |
| ------ | ----------- |
| `WithEmojis` | Emoji definitions. Defaults to the GitHub emoji set |

**Renderer options** (passed to `NewHTMLRenderer`):

| Option | Description |
| ------ | ----------- |
| `WithRenderingMethod` | `Entity` (default): HTML entities, `Unicode`: raw Unicode characters, `Twemoji`: `<img>` tag using [twemoji](https://github.com/twitter/twemoji), `Func`: custom Go function |
| `WithTwemojiTemplate` | printf template for the twemoji `<img>` tag |
| `WithRendererFunc` | Custom rendering function (used when method is `Func`) |

## License
MIT

## Author

Yusuke Inuzuka
