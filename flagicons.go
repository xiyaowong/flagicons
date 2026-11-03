// Package flagicons shows the country flags of lipis/flag-icons in the
// native UI of a MyGo app:
//
//	import "github.com/xiyaowong/flagicons"
//
//	ui.Image(c, flagicons.MustFlag("jp"))
//	ui.Image(c, flagicons.MustFlag("gb-eng").Square())
package flagicons

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"

	"github.com/egoist/mygo/ui"
)

//go:embed flags
var embedded embed.FS

// flags is the embedded flag files: 1x1, 4x3 and country.json.
var flags, flagsErr = fs.Sub(embedded, "flags")

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

// Flags parses the flag of a code, such as "jp", "gb-eng" or "eu".
func Flags(code string) (*Flag, error) {
	if flagsErr != nil {
		return nil, flagsErr
	}
	data, err := fs.ReadFile(flags, "4x3/"+code+".svg")
	if err != nil {
		return nil, err
	}
	return ui.ParseSVG(data)
}

// MustFlag is Flags for flags that are part of the program: it panics when
// the code has no flag.
func MustFlag(code string) *Flag {
	svg, err := Flags(code)
	if err != nil {
		panic(fmt.Sprintf("flagicons: %v", err))
	}
	return svg
}

// Square returns the square 1x1 flag of a code. Some flags have no square
// version; they show as the 4x3 one.
func Square(code string) (*Flag, error) {
	if flagsErr != nil {
		return nil, flagsErr
	}
	data, err := fs.ReadFile(flags, "1x1/"+code+".svg")
	if err == nil {
		return ui.ParseSVG(data)
	}
	return Flags(code)
}

// MustSquare is Square for flags that are part of the program: it panics
// when the code has no flag.
func MustSquare(code string) *Flag {
	svg, err := Square(code)
	if err != nil {
		panic(fmt.Sprintf("flagicons: %v", err))
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
	if flagsErr != nil {
		return nil, flagsErr
	}
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
