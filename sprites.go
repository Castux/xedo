package main

import (
	"image"
	_ "image/png"
	"os"
)

func LoadSprites() (image.Image, error) {
	f, err := os.Open("sprites.png")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}
