package main

import (
	"image"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
)

type iconInfo struct {
	Name   string
	Width  int
	Height int
	Scale  int
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
		rimg := imaging.Resize(img, r.Width*r.Scale, r.Height*r.Scale, imaging.Lanczos)

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
