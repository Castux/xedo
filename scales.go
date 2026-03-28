package main

import (
	"fmt"
	"math"
)

const (
	Green    uint8 = 26
	Blue     uint8 = 37
	DarkBlue uint8 = 45
	Orange   uint8 = 9
	Red      uint8 = 5
	Yellow   uint8 = 12
	Pink     uint8 = 4
)

func SolveExactMapping(divisions int) (int, int) {

	for big := 1; big <= divisions/5; big++ {
		for small := 1; small < big; small++ {
			if 5*big+2*small == divisions {
				return big, small
			}
		}
	}

	return -1, -1
}

var MajorSemitones = []int{0, 2, 4, 5, 7, 9, 11}

func MakeScale(divisions int) *ScaleInfo {

	right := 1
	palette := make(map[int]uint8)

	big, small := SolveExactMapping(divisions)
	if big > 0 {
		pitch := 0
		palette[pitch] = DarkBlue
		pitch += big
		palette[pitch] = Blue
		pitch += big
		palette[pitch] = Blue
		pitch += small
		palette[pitch] = Blue
		pitch += big
		palette[pitch] = Blue
		pitch += big
		palette[pitch] = Blue
		pitch += big
		palette[pitch] = Blue
		pitch += small
		palette[pitch] = Blue

		right = big
	} else {
		for i, semi := range MajorSemitones {
			pitch := float64(semi) / 12.0
			closest := math.Round(pitch * float64(divisions))

			if i == 1 {
				right = int(closest)
			}

			palette[int(closest)] = Blue
		}
		palette[0] = DarkBlue
	}

	return &ScaleInfo{
		Name:      fmt.Sprintf("%dtet", divisions),
		Divisions: divisions,
		RightStep: right,
		UpStep:    1,
		Palette:   palette,
	}
}
