package main

import (
	"image"
	_ "image/png"
	"os"
)

func LoadSprites() image.Image {
	f, _ := os.Open("sprites.png")
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return img
}
