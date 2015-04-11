// Copyright 2014 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ebiten_test

import (
	. "github.com/hajimehoshi/ebiten"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"math"
	"testing"
)

var ebitenImageBin = ""

func openImage(path string) (image.Image, error) {
	file, err := readFile(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func openEbitenImage(path string) (*Image, image.Image, error) {
	img, err := openImage(path)
	if err != nil {
		return nil, nil, err
	}

	eimg, err := NewImageFromImage(img, FilterNearest)
	if err != nil {
		return nil, nil, err
	}
	return eimg, img, nil
}

func diff(x, y uint8) uint8 {
	if x <= y {
		return y - x
	}
	return x - y
}

func TestImagePixels(t *testing.T) {
	img0, img, err := openEbitenImage("testdata/ebiten.png")
	if err != nil {
		t.Fatal(err)
		return
	}

	if got := img0.Bounds().Size(); got != img.Bounds().Size() {
		t.Errorf("img size: got %d; want %d", got, img.Bounds().Size())
	}

	for j := 0; j < img0.Bounds().Size().Y; j++ {
		for i := 0; i < img0.Bounds().Size().X; i++ {
			got := img0.At(i, j)
			want := color.RGBAModel.Convert(img.At(i, j))
			if got != want {
				t.Errorf("img0 At(%d, %d): got %#v; want %#v", i, j, got, want)
			}
		}
	}
}

func TestImageComposition(t *testing.T) {
	img2Color := color.NRGBA{0x24, 0x3f, 0x6a, 0x88}
	img3Color := color.NRGBA{0x85, 0xa3, 0x08, 0xd3}

	img1, _, err := openEbitenImage("testdata/ebiten.png")
	if err != nil {
		t.Fatal(err)
		return
	}

	w, h := img1.Bounds().Size().X, img1.Bounds().Size().Y

	img2, err := NewImage(w, h, FilterNearest)
	if err != nil {
		t.Fatal(err)
		return
	}

	img3, err := NewImage(w, h, FilterNearest)
	if err != nil {
		t.Fatal(err)
		return
	}

	img2.Fill(img2Color)
	img3.Fill(img3Color)
	img_12_3, err := NewImage(w, h, FilterNearest)
	if err != nil {
		t.Fatal(err)
		return
	}
	img2.DrawImage(img1, nil)
	img3.DrawImage(img2, nil)
	img_12_3.DrawImage(img3, nil)

	img2.Fill(img2Color)
	img3.Fill(img3Color)
	img_1_23, err := NewImage(w, h, FilterNearest)
	if err != nil {
		t.Fatal(err)
		return
	}
	img3.DrawImage(img2, nil)
	img3.DrawImage(img1, nil)
	img_1_23.DrawImage(img3, nil)

	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			c1 := img_12_3.At(i, j).(color.RGBA)
			c2 := img_1_23.At(i, j).(color.RGBA)
			if 1 < diff(c1.R, c2.R) || 1 < diff(c1.G, c2.G) || 1 < diff(c1.B, c2.B) || 1 < diff(c1.A, c2.A) {
				t.Errorf("img_12_3.At(%d, %d) = %#v; img_1_23.At(%[1]d, %[2]d) = %#[4]v", i, j, c1, c2)
			}
			if c1.A == 0 {
				t.Fatalf("img_12_3.At(%d, %d).A = 0; nothing is rendered?", i, j)
			}
			if c2.A == 0 {
				t.Fatalf("img_1_23.At(%d, %d).A = 0; nothing is rendered?", i, j)
			}
		}
	}
}

func TestImageSelf(t *testing.T) {
	img, _, err := openEbitenImage("testdata/ebiten.png")
	if err != nil {
		t.Fatal(err)
		return
	}
	if err := img.DrawImage(img, nil); err == nil {
		t.Fatalf("img.DrawImage(img, nil) doesn't return error; an error should be returned")
	}
}

func TestImageDotByDotInversion(t *testing.T) {
	img0, _, err := openEbitenImage("testdata/ebiten.png")
	if err != nil {
		t.Fatal(err)
		return
	}
	w, h := img0.Size()
	img1, err := NewImage(w, h, FilterNearest)
	if err != nil {
		t.Fatal(err)
		return
	}
	op := &DrawImageOptions{}
	op.GeoM.Rotate(2 * math.Pi / 2)
	op.GeoM.Translate(float64(w), float64(h))
	img1.DrawImage(img0, op)

	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			c0 := img0.At(i, j).(color.RGBA)
			c1 := img1.At(w-i-1, h-j-1).(color.RGBA)
			if c0 != c1 {
				t.Errorf("img0.At(%[1]d, %[2]d) should equal to img1.At(%[1]d, %[2]d) but not: %[3]v vs %[4]v", i, j, c0, c1)
			}
		}
	}
}

func TestReplacePixels(t *testing.T) {
	origImg, err := openImage("testdata/ebiten.png")
	if err != nil {
		t.Fatal(err)
		return
	}
	// Convert to RGBA
	img := image.NewRGBA(origImg.Bounds())
	draw.Draw(img, img.Bounds(), origImg, image.ZP, draw.Src)

	size := img.Bounds().Size()
	img0, err := NewImage(size.X, size.Y, FilterNearest)
	if err != nil {
		t.Fatal(err)
		return
	}

	img0.ReplacePixels(img.Pix)
	for j := 0; j < img0.Bounds().Size().Y; j++ {
		for i := 0; i < img0.Bounds().Size().X; i++ {
			got := img0.At(i, j)
			want := img.At(i, j)
			if got != want {
				t.Errorf("img0 At(%d, %d): got %#v; want %#v", i, j, got, want)
			}
		}
	}

	p := make([]uint8, 4*size.X*size.Y)
	for i, _ := range p {
		p[i] = 0x80
	}
	img0.ReplacePixels(p)
	for j := 0; j < img0.Bounds().Size().Y; j++ {
		for i := 0; i < img0.Bounds().Size().X; i++ {
			got := img0.At(i, j)
			want := color.RGBA{p[4*i], p[4*i+1], p[4*i+2], p[4*i+3]}
			if got != want {
				t.Errorf("img0 At(%d, %d): got %#v; want %#v", i, j, got, want)
			}
		}
	}
}

// TODO: Add more tests (e.g. DrawImage with color matrix)
