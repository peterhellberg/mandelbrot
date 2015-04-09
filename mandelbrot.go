package mandelbrot

import (
	"image"
	"image/color"
)

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

func New(w, h, i int) *Mandelbrot {
	minRe := -2.0
	maxRe := 1.0
	minIm := -1.2
	maxIm := minIm + (maxRe-minRe)*float64(h)/float64(w)
	reFactor := (maxRe - minRe) / float64(w-1)
	imFactor := (maxIm - minIm) / float64(h-1)

	return &Mandelbrot{
		minRe, maxRe, minIm, maxIm, reFactor, imFactor, w, h, i,
		color.RGBA{0xF9, 0x2A, 0x82, 0xFF},
		color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
	}
}

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
