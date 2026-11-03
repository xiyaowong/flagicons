// Package flagicons shows the country flags of lipis/flag-icons in the
// native UI of a MyGo app:
//
//	import "github.com/xiyaowong/flagicons"
//
//	ui.Image(c, flagicons.MustFlag("jp"))
//	ui.Image(c, flagicons.MustFlag("jp"))
//	ui.Image(c, flagicons.MustSquare("gb-eng"))
package flagicons

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"

	"github.com/egoist/mygo/ui"
)

//go:embed flags
var embedded embed.FS

// flags is the embedded flag files: 1x1, 4x3 and country.json.
var flags, _ = fs.Sub(embedded, "flags")

// Country is one entry of the country.json of flag-icons.
type Country struct {
	Name      string `json:"name"`
	Code      string `json:"code"`
	Capital   string `json:"capital"`
	Continent string `json:"continent"`
	// ISO tells whether the code is an ISO 3166-1 alpha-2 code.
	ISO bool `json:"iso"`
}

// Flag is a parsed flag of a code, shown in its own colors with ui.Image.
// Parse it once, not in the view, as ui.ParseSVG asks.
type Flag = ui.SVG

// readFlag reads the flag file of a shape ("1x1" or "4x3") and a code.
func readFlag(shape, code string) ([]byte, error) {
	return fs.ReadFile(flags, shape+"/"+code+".svg")
}

// Flags parses the flag of a code, such as "jp", "gb-eng" or "eu". An
// unknown code is an error matching fs.ErrNotExist.
func Flags(code string) (*Flag, error) {
	data, err := readFlag("4x3", code)
	if err != nil {
		return nil, fmt.Errorf("flag %q: %w", code, err)
	}
	svg, err := ui.ParseSVG(data)
	if err != nil {
		return nil, fmt.Errorf("flag %q: %w", code, err)
	}
	return svg, nil
}

// MustFlag is Flags for flags that are part of the program: it panics when
// the code has no flag.
func MustFlag(code string) *Flag {
	svg, err := Flags(code)
	if err != nil {
		panic(err)
	}
	return svg
}

// Square returns the square 1x1 flag of a code, an error matching
// fs.ErrNotExist when it has neither shape. Some flags have no square
// version; they show as the 4x3 one.
func Square(code string) (*Flag, error) {
	data, err := readFlag("1x1", code)
	if errors.Is(err, fs.ErrNotExist) {
		return Flags(code)
	}
	if err != nil {
		return nil, fmt.Errorf("square flag %q: %w", code, err)
	}
	svg, err := ui.ParseSVG(data)
	if err != nil {
		return nil, fmt.Errorf("square flag %q: %w", code, err)
	}
	return svg, nil
}

// MustSquare is Square for flags that are part of the program: it panics
// when the code has no flag.
func MustSquare(code string) *Flag {
	svg, err := Square(code)
	if err != nil {
		panic(err)
	}
	return svg
}

// Codes returns the codes of all flags: the ISO 3166-1 alpha-2 codes plus
// the extra flags of flag-icons, e.g. "gb-eng" and "eu".
func Codes() ([]string, error) {
	countries, err := Countries()
	if err != nil {
		return nil, err
	}
	codes := make([]string, len(countries))
	for i, c := range countries {
		codes[i] = c.Code
	}
	return codes, nil
}

// Countries returns the country.json of flag-icons: name, code, capital,
// continent, and whether the code is an ISO 3166-1 alpha-2 code.
func Countries() ([]Country, error) {
	data, err := fs.ReadFile(flags, "country.json")
	if err != nil {
		return nil, fmt.Errorf("reading country.json: %w", err)
	}
	var list []Country
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("parsing country.json: %w", err)
	}
	return list, nil
}
