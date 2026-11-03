# flagicons

Country flags of [lipis/flag-icons](https://github.com/lipis/flag-icons) for
the [native UI](https://mygo.egoist.dev/docs/ui) of
[MyGo](https://mygo.egoist.dev) apps.

```sh
go get github.com/xiyaowong/flagicons
```

```go
import "github.com/xiyaowong/flagicons"

ui.Image(c, flagicons.MustFlag("jp"))
ui.Image(c, flagicons.MustSquare("gb-eng"))
```

`Flag` is a `*ui.SVG`: parse it once, not in the view, as
[`ui.ParseSVG`](https://pkg.go.dev/github.com/egoist/mygo/ui#ParseSVG) asks,
and show it with
[`ui.Image`](https://pkg.go.dev/github.com/egoist/mygo/ui#Image), in its own
colors.

## Functions

- `Flags(code) (*Flag, error)` — the 4x3 flag of a code, e.g. "jp",
  "gb-eng" or "eu".
- `MustFlag(code) *Flag` — `Flags`, panicking on a code without a flag.
- `Square(code) (*Flag, error)` — the square 1x1 flag; codes without a
  square version show as the 4x3 one.
- `MustSquare(code) *Flag` — `Square`, panicking on a code without a flag.
- `Codes() ([]string, error)` — all 271 codes.
- `Countries() ([]Country, error)` — name, code, capital, continent and
  ISO-ness of each.

## License

MIT, as [flag-icons](LICENSE.flags) is.
