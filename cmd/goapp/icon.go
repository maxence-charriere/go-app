package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
)

type iconInfo struct {
	Name   string
	Width  float64
	Height float64
	Scale  float64
}

func generateIcons(original string, resizes []iconInfo) error {
	f, err := os.Open(original)
	if err != nil {
		return err
	}
	defer f.Close()

	var img image.Image
	if img, err = png.Decode(f); err != nil {
		return err
	}

	for _, r := range resizes {
		w := int(r.Width * r.Scale)
		h := int(r.Height * r.Scale)

		rimg := imaging.Fit(img, w, h, imaging.Lanczos)

		if w != h {
			background := imaging.New(w, h, color.NRGBA{0, 0, 0, 0})
			rimg = imaging.PasteCenter(background, rimg)
		}

		var icon *os.File
		if icon, err = os.Create(r.Name); err != nil {
			return err
		}
		defer icon.Close()

		if err = png.Encode(icon, rimg); err != nil {
			return err
		}
	}
	return nil
}
