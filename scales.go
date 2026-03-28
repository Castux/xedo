package main

const (
	Green  uint8 = 26
	Blue   uint8 = 37
	Orange uint8 = 9
	Red    uint8 = 5
)

var Major = ScaleInfo{
	Divisions: 12,
	RightStep: 2,
	UpStep:    1,
	Palette: map[int]uint8{
		0:  Green,
		2:  Blue,
		4:  Blue,
		5:  Blue,
		7:  Blue,
		9:  Blue,
		11: Blue,
	},
}

var Minor = ScaleInfo{
	Divisions: 12,
	RightStep: 2,
	UpStep:    1,
	Palette: map[int]uint8{
		0:  Green,
		2:  Blue,
		3:  Blue,
		5:  Blue,
		7:  Blue,
		8:  Blue,
		10: Blue,
	},
}
