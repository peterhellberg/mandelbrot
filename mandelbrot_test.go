package mandelbrot

import (
	"image"
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
	}{
		{320, 256, 15},
		{800, 600, 32},
	} {
		i := New(tt.w, tt.h, tt.i).Image()

		pt := image.Point{tt.w, tt.h}

		if i.Bounds().Max != pt {
			t.Errorf(`unexpected size %v, want %v`, i.Bounds().Max, pt)
		}
	}
}
