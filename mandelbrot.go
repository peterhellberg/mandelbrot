package mandelbrot

import (
	"image"
	"image/color"
)

// Mandelbrot set
type Mandelbrot struct {
	MinRe         float64
	MaxRe         float64
	MinIm         float64
	MaxIm         float64
	ReFactor      float64
	ImFactor      float64
	Width         int
	Height        int
	MaxIterations int
	InsideColor   color.Color
	OutsideColor  color.Color
}

// Function that changes a Mandelbrot struct
type ChangeFunc func(*Mandelbrot)

// New creates a new Mandelbrot set
func New(w, h, i int, options ...ChangeFunc) *Mandelbrot {
	minRe := -2.0
	maxRe := 1.0
	minIm := -1.2
	maxIm := minIm + (maxRe-minRe)*float64(h)/float64(w)
	reFactor := (maxRe - minRe) / float64(w-1)
	imFactor := (maxIm - minIm) / float64(h-1)

	m := &Mandelbrot{
		minRe, maxRe, minIm, maxIm, reFactor, imFactor, w, h, i,
		color.Black,
		color.White,
	}

	for _, option := range options {
		// Apply changes to the struct
		option(m)
	}

	return m
}

// Colors returns a function that changes the inside and outside colors of the set
func Colors(inside, outside color.Color) ChangeFunc {
	return func(m *Mandelbrot) {
		m.InsideColor = inside
		m.OutsideColor = outside
	}
}

// Image renders an image of the set
func (m *Mandelbrot) Image() *image.RGBA {
	i := image.NewRGBA(image.Rect(0, 0, m.Width, m.Height))

	for y := 0; y < m.Height; y++ {
		cIm := m.MaxIm - float64(y)*m.ImFactor

		for x := 0; x < m.Width; x++ {
			cRe := m.MinRe + float64(x)*m.ReFactor
			zRe := cRe
			zIm := cIm
			isInside := true

			for n := 0; n < m.MaxIterations; n++ {
				zRe2 := zRe * zRe
				zIm2 := zIm * zIm

				if zRe2+zIm2 > 4 {
					isInside = false
					break
				}

				zIm = 2*zRe*zIm + cIm
				zRe = zRe2 - zIm2 + cRe
			}

			if isInside {
				i.Set(x, y, m.InsideColor)
			} else {
				i.Set(x, y, m.OutsideColor)
			}
		}
	}

	return i
}
