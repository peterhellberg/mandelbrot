package main

import (
	"flag"
	"fmt"
	"image/color"
	"image/png"
	"os"
	"os/exec"

	"github.com/peterhellberg/mandelbrot"
)

var (
	width      = flag.Int("w", 640, "Width of the image")
	height     = flag.Int("h", 480, "Height of the image")
	iterations = flag.Int("n", 30, "Number of iterations to run")
	inside     = flag.String("i", "000000", "Inside color")
	outside    = flag.String("o", "ffffff", "Outside color")
	filename   = flag.String("f", "mandelbrot.png", "Filename of the image")

	show = flag.Bool("show", false, "Show the generated image")
)

func main() {
	flag.Parse()

	m := mandelbrot.New(*width, *height, *iterations,
		mandelbrot.Colors(hex(*inside), hex(*outside))).Image()

	if file, err := os.Create(*filename); err == nil {
		defer file.Close()

		if err := png.Encode(file, m); err == nil {
			open(*filename)
		}
	}
}

func open(fn string) {
	if *show {
		exec.Command("open", fn).Run()
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
