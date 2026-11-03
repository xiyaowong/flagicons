package flagicons

import (
	"io/fs"
	"path"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestAllSVGsParse parses every SVG file embedded in both shapes: all 542
// flags must parse and declare a size.
func TestAllSVGsParse(t *testing.T) {
	for _, shape := range []string{"1x1", "4x3"} {
		files, err := fs.ReadDir(flags, shape)
		if err != nil {
			t.Fatalf("reading %s: %v", shape, err)
		}
		if len(files) == 0 {
			t.Fatalf("%s has no files", shape)
		}
		for _, f := range files {
			name := f.Name()
			if !strings.HasSuffix(name, ".svg") {
				t.Errorf("%s/%s: not an SVG", shape, name)
				continue
			}
			data, err := fs.ReadFile(flags, path.Join(shape, name))
			if err != nil {
				t.Errorf("reading %s/%s: %v", shape, name, err)
				continue
			}
			svg, err := ui.ParseSVG(data)
			if err != nil {
				t.Errorf("parsing %s/%s: %v", shape, name, err)
				continue
			}
			w, h := svg.Size()
			if w <= 0 || h <= 0 {
				t.Errorf("%s/%s: size %gx%g", shape, name, w, h)
			}
		}
	}
}

func TestCountries(t *testing.T) {
	countries, err := Countries()
	if err != nil {
		t.Fatalf("Countries: %v", err)
	}
	if len(countries) != 271 {
		t.Errorf("got %d countries, want 271", len(countries))
	}
	codes := map[string]bool{}
	for _, c := range countries {
		if c.Code == "" || c.Name == "" {
			t.Errorf("entry with empty code or name: %+v", c)
		}
		codes[c.Code] = true
	}
	for _, want := range []string{"jp", "gr", "de", "gb-eng", "eu", "xx"} {
		if !codes[want] {
			t.Errorf("code %q missing", want)
		}
	}
}

func TestCodes(t *testing.T) {
	codes, err := Codes()
	if err != nil {
		t.Fatalf("Codes: %v", err)
	}
	if len(codes) != 271 {
		t.Errorf("got %d codes, want 271", len(codes))
	}
}

// TestEveryCodeHasFlags parses the flag of every country.json code in both
// shapes.
func TestEveryCodeHasFlags(t *testing.T) {
	countries, err := Countries()
	if err != nil {
		t.Fatalf("Countries: %v", err)
	}
	for _, c := range countries {
		svg, err := Flags(c.Code)
		if err != nil {
			t.Errorf("Flags(%q): %v", c.Code, err)
			continue
		}
		if w, h := svg.Size(); w <= 0 || h <= 0 {
			t.Errorf("Flags(%q): size %gx%g", c.Code, w, h)
		}
		if _, err := Square(c.Code); err != nil {
			t.Errorf("Square(%q): %v", c.Code, err)
		}
	}
}

func TestFlagsUnknownCode(t *testing.T) {
	if _, err := Flags("nope"); err == nil {
		t.Error("Flags(nope) should fail")
	}
	if _, err := Square("nope"); err == nil {
		t.Error("Square(nope) should fail")
	}
}

func TestMustPanicsOnUnknownCode(t *testing.T) {
	for name, fn := range map[string]func(){"MustFlag": func() { MustFlag("nope") }, "MustSquare": func() { MustSquare("nope") }} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s(nope) should panic", name)
				}
			}()
			fn()
		}()
	}
}

// TestSquareFallsBack checks that a square flag differs from the 4x3 one
// when both exist, and that the fallback parses too.
func TestSquareFallsBack(t *testing.T) {
	sq, err := Square("gb-eng")
	if err != nil {
		t.Fatalf("Square(gb-eng): %v", err)
	}
	wide, err := Flags("gb-eng")
	if err != nil {
		t.Fatalf("Flags(gb-eng): %v", err)
	}
	sw, sh := sq.Size()
	ww, wh := wide.Size()
	if sw == ww && sh == wh {
		t.Errorf("Square(gb-eng) size %gx%g equals the 4x3 one", sw, sh)
	}

	fallback, err := Square("xx")
	if err != nil {
		t.Fatalf("Square(xx): %v", err)
	}
	// xx has a real 1x1 file; the code exists in both shapes, so both
	// return their own versions.
	fw, fh := fallback.Size()
	if fw != 512 || fh != 512 {
		t.Errorf("Square(xx) size %gx%g, want the 1x1 512x512", fw, fh)
	}

	// The 4x3 one is its own size.
	w2, h2 := wide.Size()
	if w2 != 640 || h2 != 480 {
		t.Errorf("Flags(xx) size %gx%g, want 640x480", w2, h2)
	}
}
