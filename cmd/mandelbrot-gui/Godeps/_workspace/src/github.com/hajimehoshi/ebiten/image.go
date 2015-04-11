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

package ebiten

import (
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten/internal/graphics"
	"github.com/hajimehoshi/ebiten/internal/graphics/internal/opengl"
	"image"
	"image/color"
)

// Image represents an image.
// The pixel format is alpha-premultiplied.
// Image implements image.Image.
type Image struct {
	framebuffer *graphics.Framebuffer
	texture     *graphics.Texture
	pixels      []uint8
	width       int
	height      int
}

// Size returns the size of the image.
func (i *Image) Size() (width, height int) {
	if i.width == 0 {
		i.width, i.height = i.framebuffer.Size()
	}
	return i.width, i.height
}

// Clear resets the pixels of the image into 0.
func (i *Image) Clear() (err error) {
	return i.Fill(color.Transparent)
}

// Fill fills the image with a solid color.
func (i *Image) Fill(clr color.Color) (err error) {
	i.pixels = nil
	useGLContext(func(c *opengl.Context) {
		err = i.framebuffer.Fill(c, clr)
	})
	return
}

// DrawImage draws the given image on the receiver image.
//
// This method accepts the options.
// The parts of the given image at the parts of the destination.
// After determining parts to draw, this applies the geometry matrix and the color matrix.
//
// Here are the default values:
//     ImageParts: (0, 0) - (source width, source height) to (0, 0) - (source width, source height)
//                 (i.e. the whole source image)
//     GeoM:       Identity matrix
//     ColorM:     Identity matrix (that changes no colors)
//
// Be careful that this method is potentially slow.
// It would be better if you could call this method fewer times.
func (i *Image) DrawImage(image *Image, options *DrawImageOptions) (err error) {
	if i == image {
		return errors.New("Image.DrawImage: image should be different from the receiver")
	}
	i.pixels = nil
	if options == nil {
		options = &DrawImageOptions{}
	}
	parts := options.ImageParts
	if parts == nil {
		// Check options.Parts for backward-compatibility.
		dparts := options.Parts
		if dparts != nil {
			parts = imageParts(dparts)
		} else {
			w, h := image.Size()
			parts = &wholeImage{w, h}
		}
	}
	w, h := image.Size()
	quads := &textureQuads{parts: parts, width: w, height: h}
	useGLContext(func(c *opengl.Context) {
		err = i.framebuffer.DrawTexture(c, image.texture, quads, &options.GeoM, &options.ColorM)
	})
	return
}

// DrawLine draws a line.
func (i *Image) DrawLine(x0, y0, x1, y1 int, clr color.Color) error {
	return i.DrawLines(&line{x0, y0, x1, y1, clr})
}

// DrawLines draws lines.
func (i *Image) DrawLines(lines Lines) (err error) {
	i.pixels = nil
	useGLContext(func(c *opengl.Context) {
		err = i.framebuffer.DrawLines(c, lines)
	})
	return
}

// DrawRect draws a rectangle.
func (i *Image) DrawRect(x, y, width, height int, clr color.Color) error {
	return i.DrawLines(&rectsAsLines{&rect{x, y, width, height, clr}})
}

// DrawRect draws rectangles.
func (i *Image) DrawRects(rects Rects) error {
	return i.DrawLines(&rectsAsLines{rects})
}

// DrawFilledRect draws a filled rectangle.
func (i *Image) DrawFilledRect(x, y, width, height int, clr color.Color) error {
	return i.DrawFilledRects(&rect{x, y, width, height, clr})
}

// DrawFilledRects draws filled rectangles on the image.
func (i *Image) DrawFilledRects(rects Rects) (err error) {
	i.pixels = nil
	useGLContext(func(c *opengl.Context) {
		err = i.framebuffer.DrawFilledRects(c, rects)
	})
	return
}

// Bounds returns the bounds of the image.
func (i *Image) Bounds() image.Rectangle {
	w, h := i.Size()
	return image.Rect(0, 0, w, h)
}

// ColorModel returns the color model of the image.
func (i *Image) ColorModel() color.Model {
	return color.RGBAModel
}

// At returns the color of the image at (x, y).
//
// This method loads pixels from VRAM to system memory if necessary.
func (i *Image) At(x, y int) color.Color {
	if i.pixels == nil {
		useGLContext(func(c *opengl.Context) {
			var err error
			i.pixels, err = i.framebuffer.Pixels(c)
			if err != nil {
				panic(err)
			}
		})
	}
	w, _ := i.Size()
	w = graphics.NextPowerOf2Int(w)
	idx := 4*x + 4*y*w
	r, g, b, a := i.pixels[idx], i.pixels[idx+1], i.pixels[idx+2], i.pixels[idx+3]
	return color.RGBA{r, g, b, a}
}

func (i *Image) dispose() {
	useGLContext(func(c *opengl.Context) {
		if i.framebuffer != nil {
			i.framebuffer.Dispose(c)
		}
		if i.texture != nil {
			i.texture.Dispose(c)
		}
	})
	i.pixels = nil
}

// ReplacePixels replaces the pixels of the image with p.
//
// The given p must represent RGBA pre-multiplied alpha values. len(p) must equal to 4 * (image width) * (image height).
//
// This function may be slow (as for implementation, this calls glTexSubImage2D).
func (i *Image) ReplacePixels(p []uint8) error {
	// Don't set i.pixels here because i.pixels is used not every time.

	i.pixels = nil
	w, h := i.Size()
	l := 4 * w * h
	if len(p) != l {
		return errors.New(fmt.Sprintf("p's length must be %d", l))
	}
	var err error
	useGLContext(func(c *opengl.Context) {
		err = i.texture.ReplacePixels(c, p)
	})
	return err
}

// A DrawImageOptions represents options to render an image on an image.
type DrawImageOptions struct {
	ImageParts ImageParts
	GeoM       GeoM
	ColorM     ColorM

	// Deprecated (as of 1.1.0-alpha): Use ImageParts instead.
	Parts []ImagePart
}
