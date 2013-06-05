// Copyright 2012 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package gui

import (
	"excavation/engine"
)

type Image struct {
	name    string
	Overlay *engine.Overlay
}

func MakeImage(name, imagePath string, dimensions *engine.ScreenArea) *Image {
	newImage := new(Image)
	newImage.Overlay = engine.NewOverlay(imagePath, engine.NewColor(255, 255, 255, 255), dimensions)
	return newImage
}

func (i *Image) MouseArea() *engine.ScreenArea {
	return i.Overlay.Dimensions
}

func (i *Image) Update() {
	i.Overlay.Place()
}

func (i *Image) Name() string {
	return i.name
}

func (i *Image) Hover()           { return }
func (i *Image) Click(mouse int)  { return }
func (i *Image) Scroll(delta int) { return }
func (i *Image) Unload()          { return }
