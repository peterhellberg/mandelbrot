package mandelbrot

import (
	"image"
	"image/color"
	"testing"
)

func TestNew(t *testing.T) {
	for _, tt := range []struct {
		w, h, i      int
		mi, ref, imf float64
	}{
		{320, 256, 15, 1.2, 0.009404388714733543, 0.009411764705882352},
		{800, 600, 32, 1.05, 0.0037546933667083854, 0.0037562604340567614},
	} {
		m := New(tt.w, tt.h, tt.i)

		if m.Width != tt.w {
			t.Errorf(`m.Width = %v, want %v`, m.Width, tt.w)
		}

		if m.Height != tt.h {
			t.Errorf(`m.Height = %v, want %v`, m.Height, tt.h)
		}

		if m.MaxIterations != tt.i {
			t.Errorf(`m.MaxIterations = %v, want %v`, m.MaxIterations, tt.i)
		}

		if m.MaxIm != tt.mi {
			t.Errorf(`m.MaxIm = %v, want %v`, m.MaxIm, tt.mi)
		}

		if m.ReFactor != tt.ref {
			t.Errorf(`m.ReFactor = %v, want %v`, m.ReFactor, tt.ref)
		}

		if m.ImFactor != tt.imf {
			t.Errorf(`m.ImFactor = %v, want %v`, m.ImFactor, tt.imf)
		}
	}
}

func TestImage(t *testing.T) {
	for _, tt := range []struct {
		w, h, i int
		c       color.Color
	}{
		{320, 256, 15, color.RGBA{0x1e, 0x8b, 0xb8, 0xff}},
		{800, 600, 32, color.RGBA{0xa7, 0xb8, 0x40, 0xff}},
	} {
		m := New(tt.w, tt.h, tt.i)

		m.InsideColor = tt.c

		i := m.Image()

		pt := image.Point{tt.w, tt.h}

		if i.Bounds().Max != pt {
			t.Errorf(`unexpected size %v, want %v`, i.Bounds().Max, pt)
		}

		if got := i.At(tt.w/2, tt.h/2); got != tt.c {
			t.Errorf(`unexpected color %v, want %v`, got, tt.c)
		}
	}
}

func TestOption(t *testing.T) {
	w := 32
	h := 32
	i := 15

	cf1 := Colors(color.RGBA{0, 0, 0, 0xff}, color.RGBA{0xff, 0xff, 0xff, 0xff})
	blue := color.RGBA{0, 0, 0xff, 0xff}
	cf2 := Colors(blue, color.RGBA{0xff, 0, 0, 0xff})
	m := New(w, h, i, cf1, cf2)

	img := m.Image()

	pt := image.Point{w, h}

	if img.Bounds().Max != pt {
		t.Errorf(`unexpected size %v, want %v`, img.Bounds().Max, pt)
	}

	if got := img.At(w/2, h/2); got != blue {
		t.Errorf(`unexpected color %v, want %v`, got, blue)
	}
}
