package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
)

type iconInfo struct {
	Name    string
	Width   float64
	Height  float64
	Scale   float64
	Padding bool
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
		ow := int(r.Width * r.Scale)
		oh := int(r.Height * r.Scale)

		p := 0
		if r.Padding {
			p = int((42 * r.Height / 100) * r.Scale)
		}

		w := ow - p
		h := oh - p

		rimg := imaging.Fit(img, w, h, imaging.Lanczos)

		if p != 0 {
			background := imaging.New(ow, oh, color.NRGBA{0, 0, 0, 0})
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
