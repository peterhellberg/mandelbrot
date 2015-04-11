package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten"

	"github.com/peterhellberg/mandelbrot"
)

var (
	scale   = flag.Int("s", 3, "Scaling factor")
	width   = flag.Int("w", 256, "Width of the image")
	height  = flag.Int("h", 192, "Height of the image")
	inside  = flag.String("i", "000000", "Inside color")
	outside = flag.String("o", "ffffff", "Outside color")

	count int
	m     *image.RGBA
)

func iterations() int {
	return time.Now().Second()%15 + 5
}

func update(screen *ebiten.Image) error {
	count++

	m = mandelbrot.New(*width, *height, iterations(),
		mandelbrot.Colors(hex(*inside), hex(*outside))).Image()

	mandelbrotImage, err := ebiten.NewImageFromImage(m, ebiten.FilterLinear)
	if err == nil {
		op := &ebiten.DrawImageOptions{}

		op.GeoM.Translate(-float64(*width)/2, -float64(*height)/2)
		op.GeoM.Rotate(float64(count%360) * 2 * math.Pi / 360)
		op.GeoM.Translate(float64(*width)/2, float64(*height)/2)

		if err := screen.DrawImage(mandelbrotImage, op); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	flag.Parse()

	m = mandelbrot.New(*width, *height, iterations(),
		mandelbrot.Colors(hex(*inside), hex(*outside))).Image()

	if err := ebiten.Run(update, *width, *height, *scale, "Plasma GUI"); err != nil {
		log.Fatal(err)
	}
}

func hex(scol string) color.Color {
	var r, g, b uint8

	n, err := fmt.Sscanf(scol, "%02x%02x%02x", &r, &g, &b)

	if err != nil || n != 3 {
		return color.Black
	}

	return color.RGBA{r, g, b, 0xFF}
}
