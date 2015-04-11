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

package graphics

import (
	"errors"
	"github.com/hajimehoshi/ebiten/internal/graphics/internal/opengl"
	"image"
	"image/draw"
)

func adjustImageForTexture(img image.Image) *image.RGBA {
	width, height := img.Bounds().Size().X, img.Bounds().Size().Y
	adjustedImageBounds := image.Rectangle{
		image.ZP,
		image.Point{
			NextPowerOf2Int(width),
			NextPowerOf2Int(height),
		},
	}
	if adjustedImage, ok := img.(*image.RGBA); ok && img.Bounds() == adjustedImageBounds {
		return adjustedImage
	}

	adjustedImage := image.NewRGBA(adjustedImageBounds)
	dstBounds := image.Rectangle{
		image.ZP,
		img.Bounds().Size(),
	}
	draw.Draw(adjustedImage, dstBounds, img, image.ZP, draw.Src)
	return adjustedImage
}

type Texture struct {
	native opengl.Texture
	width  int
	height int
}

func (t *Texture) Size() (width, height int) {
	return t.width, t.height
}

func NewTexture(c *opengl.Context, width, height int, filter opengl.Filter) (*Texture, error) {
	w := NextPowerOf2Int(width)
	h := NextPowerOf2Int(height)
	if w < 4 {
		return nil, errors.New("width must be equal or more than 4.")
	}
	if h < 4 {
		return nil, errors.New("height must be equal or more than 4.")
	}
	native, err := c.NewTexture(w, h, nil, filter)
	if err != nil {
		return nil, err
	}
	return &Texture{native, width, height}, nil
}

func NewTextureFromImage(c *opengl.Context, img image.Image, filter opengl.Filter) (*Texture, error) {
	origSize := img.Bounds().Size()
	if origSize.X < 4 {
		return nil, errors.New("width must be equal or more than 4.")
	}
	if origSize.Y < 4 {
		return nil, errors.New("height must be equal or more than 4.")
	}
	adjustedImage := adjustImageForTexture(img)
	size := adjustedImage.Bounds().Size()
	native, err := c.NewTexture(size.X, size.Y, adjustedImage.Pix, filter)
	if err != nil {
		return nil, err
	}
	return &Texture{native, origSize.X, origSize.Y}, nil
}

func (t *Texture) Dispose(c *opengl.Context) {
	c.DeleteTexture(t.native)
}

func (t *Texture) ReplacePixels(c *opengl.Context, p []uint8) error {
	c.BindTexture(t.native)
	c.TexSubImage2D(p, t.width, t.height)
	return nil
}
