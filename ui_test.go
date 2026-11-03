package flagicons_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/xiyaowong/flagicons"
)

// gallery shows flags the way an app does: 4x3 and square, labeled.
type gallery struct {
	code string
}

func (g *gallery) view(c *ui.Context) {
	ui.Column(c).Fill().Center().Gap(8).Children(func() {
		// 4x3 flag, labeled for finding it in tests.
		ui.Image(c, flagicons.MustFlag("jp")).Size(96, 72).Label("flag-4x3")
		// Square flag.
		ui.Image(c, flagicons.MustSquare("gr")).Size(48, 48).Label("flag-1x1")
	})
}

// TestViewRendersFlags shows both shapes of a flag in a tested view and
// checks the labels found each image and every frame rendered pixels.
func TestViewRendersFlags(t *testing.T) {
	app := &gallery{code: "jp"}
	tt := ui.NewTester(app.view, 320, 240)

	for _, label := range []string{"flag-4x3", "flag-1x1"} {
		box, ok := tt.Find(label)
		if !ok {
			t.Errorf("%s not found in the view", label)
			continue
		}
		if box.W <= 0 || box.H <= 0 {
			t.Errorf("%s has size %v", label, box)
		}
	}

	if img := tt.Image(); img == nil || img.Bounds().Empty() {
		t.Error("no frame rendered")
	}
}

// TestViewFlagsPaint checks that a flag paints non-empty pixels: a view
// showing only one flag on white must differ from a white window.
func TestViewFlagsPaint(t *testing.T) {
	blank := ui.Render(func(c *ui.Context) {}, 96, 72, 1)
	painted := ui.Render(func(c *ui.Context) {
		ui.Image(c, flagicons.MustFlag("jp"))
	}, 96, 72, 1)

	if painted == nil || blank == nil {
		t.Fatal("Render returned no image")
	}
	if sameImage(blank, painted) {
		t.Error("the flag view rendered as a blank window")
	}
}

// TestViewSquarePaint does the same for a square flag.
func TestViewSquarePaint(t *testing.T) {
	painted := ui.Render(func(c *ui.Context) {
		ui.Image(c, flagicons.MustSquare("de"))
	}, 72, 72, 1)
	blank := ui.Render(func(c *ui.Context) {}, 72, 72, 1)

	if sameImage(blank, painted) {
		t.Error("the square flag view rendered as a blank window")
	}
}

// TestViewFlagPixels checks what the test's image shows: the flag of Japan
// is a red disc on white, so the center must be red and a corner white.
func TestViewFlagPixels(t *testing.T) {
	img := ui.Render(func(c *ui.Context) {
		ui.Image(c, flagicons.MustFlag("jp")).Size(60, 45)
	}, 60, 45, 1)

	center := img.At(30, 22)
	cr, cg, cb, _ := center.RGBA()
	if cr < 0x8000 || cg > 0x4000 || cb > 0x4000 {
		t.Errorf("center of the flag of Japan is %v, want a red", center)
	}
	corner := img.At(2, 2)
	wr, wg, wb, _ := corner.RGBA()
	if wr < 0x8000 || wg < 0x8000 || wb < 0x8000 {
		t.Errorf("corner of the flag of Japan is %v, want a white", corner)
	}
}

// TestViewDynamicCode shows a view reading its flag from state, as an app
// would, and checks the pixels follow the state.
func TestViewDynamicCode(t *testing.T) {
	app := &gallery{code: "jp"}

	// Japan has a red disc; France has vertical bands, no red disc in the
	// middle (the blue and red stand at the edges).
	img := ui.Render(func(c *ui.Context) {
		ui.Image(c, flagicons.MustFlag(app.code)).Size(60, 45)
	}, 60, 45, 1)
	jp := img.At(30, 22)

	app.code = "fr"
	img = ui.Render(func(c *ui.Context) {
		ui.Image(c, flagicons.MustFlag(app.code)).Size(60, 45)
	}, 60, 45, 1)
	fr := img.At(30, 22)

	if sameColor(jp, fr) {
		t.Errorf("changing the code changed nothing: both centers are %v", fr)
	}
}

func sameImage(a, b *image.RGBA) bool {
	if a.Bounds() != b.Bounds() {
		return false
	}
	for y := a.Bounds().Min.Y; y < a.Bounds().Max.Y; y++ {
		for x := a.Bounds().Min.X; x < a.Bounds().Max.X; x++ {
			if a.At(x, y) != b.At(x, y) {
				return false
			}
		}
	}
	return true
}

func sameColor(a, b color.Color) bool {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
}

// TestImageSize checks the sizes Image shows flags at by default, from
// the SVGs' viewBoxes.
func TestImageSize(t *testing.T) {
	w, h := ui.ImageSize(flagicons.MustFlag("jp"))
	if w != 640 || h != 480 {
		t.Errorf("ImageSize(jp) = %gx%g, want 640x480", w, h)
	}
	w, h = ui.ImageSize(flagicons.MustSquare("de"))
	if w != 512 || h != 512 {
		t.Errorf("ImageSize(de 1x1) = %gx%g, want 512x512", w, h)
	}
}

// TestFlagsInRow lays flags out in a row, as a language picker would, and
// checks each lands beside the last.
func TestFlagsInRow(t *testing.T) {
	var boxes []ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Row(c).Gap(4).Children(func() {
			for _, code := range []string{"jp", "gr", "de"} {
				ui.Image(c, flagicons.MustSquare(code)).Size(24, 24).Label("flag-" + code)
			}
		})
	}, 200, 40)
	for _, code := range []string{"jp", "gr", "de"} {
		box, ok := tt.Find("flag-" + code)
		if !ok {
			t.Fatalf("flag-%s not found", code)
		}
		boxes = append(boxes, box)
	}
	for i := 1; i < len(boxes); i++ {
		if boxes[i].X <= boxes[i-1].X {
			t.Errorf("flag %d not right of the one before: %v then %v", i, boxes[i-1], boxes[i])
		}
	}
}
