package main

import (
	"flag"
	"image/png"
	"os"
	"os/exec"

	"github.com/peterhellberg/mandelbrot"
)

var (
	width      = flag.Int("w", 640, "Width of the image")
	height     = flag.Int("h", 480, "Height of the image")
	iterations = flag.Int("n", 30, "Number of iterations to run")
	filename   = flag.String("f", "mandelbrot.png", "Filename of the image")

	show = flag.Bool("show", false, "Show the generated image")
)

func main() {
	flag.Parse()

	m := mandelbrot.New(*width, *height, *iterations).Image()

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
