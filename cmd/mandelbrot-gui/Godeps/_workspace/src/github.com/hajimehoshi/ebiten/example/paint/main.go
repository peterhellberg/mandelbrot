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

package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image"
	"image/color"
	"log"
	"math"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

var (
	count       int
	brushImage  *ebiten.Image
	canvasImage *ebiten.Image
)

func update(screen *ebiten.Image) error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		count++
	}

	mx, my := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(mx), float64(my))
		op.ColorM.Scale(1.0, 0.25, 0.25, 1.0)
		theta := 2.0 * math.Pi * float64(count%60) / 60.0
		op.ColorM.Concat(ebiten.RotateHue(theta))
		if err := canvasImage.DrawImage(brushImage, op); err != nil {
			return err
		}
	}

	if err := screen.DrawImage(canvasImage, nil); err != nil {
		return err
	}

	if err := ebitenutil.DebugPrint(screen, fmt.Sprintf("(%d, %d)", mx, my)); err != nil {
		return err
	}
	return nil
}

func main() {
	var err error
	const a0, a1, a2 = 0x40, 0xc0, 0xff
	pixels := []uint8{
		a0, a1, a1, a0,
		a1, a2, a2, a1,
		a1, a2, a2, a1,
		a0, a1, a1, a0,
	}
	brushImage, err = ebiten.NewImageFromImage(&image.Alpha{
		Pix:    pixels,
		Stride: 4,
		Rect:   image.Rect(0, 0, 4, 4),
	}, ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}

	canvasImage, err = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	canvasImage.Fill(color.White)

	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Paint (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}
}
