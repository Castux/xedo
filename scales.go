package main

const (
	Green  uint8 = 26
	Blue   uint8 = 37
	Orange uint8 = 9
	Red    uint8 = 5
	Yellow uint8 = 12
	Pink   uint8 = 4
)

var Scales = []*ScaleInfo{

	&ScaleInfo{
		Name:      "major",
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
	},

	&ScaleInfo{
		Name:      "minor",
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
	},

	&ScaleInfo{
		Name:      "19tet",
		Divisions: 19,
		RightStep: 3,
		UpStep:    1,
		Palette: map[int]uint8{
			0:  Green,
			1:  Yellow,
			2:  Pink,
			3:  Blue,
			4:  Yellow,
			5:  Pink,
			6:  Blue,
			7:  Yellow,
			8:  Blue,
			9:  Yellow,
			10: Pink,
			11: Blue,
			12: Yellow,
			13: Pink,
			14: Blue,
			15: Yellow,
			16: Pink,
			17: Blue,
			18: Yellow,
		},
	},
}
