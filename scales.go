package main

const (
	Green    uint8 = 26
	Blue     uint8 = 37
	DarkBlue uint8 = 45
	Orange   uint8 = 9
	Red      uint8 = 5
	Yellow   uint8 = 12
	Pink     uint8 = 4
)

var Scales = []*ScaleInfo{

	&ScaleInfo{
		Name:      "10tet",
		Divisions: 10,
		RightStep: 2,
		UpStep:    1,
		Palette: map[int]uint8{
			0: DarkBlue,
			2: Blue,
			4: Blue,
			6: Blue,
			8: Blue,

			1: Pink,
			3: Pink,
			5: Pink,
			7: Pink,
			9: Pink,
		},
	},

	&ScaleInfo{
		Name:      "12tet",
		Divisions: 12,
		RightStep: 2,
		UpStep:    1,
		Palette: map[int]uint8{
			0:  DarkBlue,
			2:  Blue,
			4:  Blue,
			5:  Blue,
			7:  Blue,
			9:  Blue,
			11: Blue,
		},
	},

	&ScaleInfo{
		Name:      "14tet",
		Divisions: 14,
		RightStep: 2,
		UpStep:    1,
		Palette: map[int]uint8{
			0:  DarkBlue,
			2:  Blue,
			5:  Blue,
			6:  Blue,
			8:  Blue,
			11: Blue,
			13: Blue,
		},
	},

	&ScaleInfo{
		Name:      "15tet",
		Divisions: 15,
		RightStep: 2,
		UpStep:    1,
		Palette: map[int]uint8{
			0:  DarkBlue,
			2:  Blue,
			5:  Blue,
			6:  Blue,
			9:  Blue,
			11: Blue,
			14: Blue,
		},
	},

	&ScaleInfo{
		Name:      "17tet",
		Divisions: 17,
		RightStep: 3,
		UpStep:    1,
		Palette: map[int]uint8{
			0:  DarkBlue,
			3:  Blue,
			6:  Blue,
			7:  Blue,
			10: Blue,
			13: Blue,
			16: Blue,

			1:  Yellow,
			4:  Yellow,
			8:  Yellow,
			11: Yellow,
			14: Yellow,

			2:  Pink,
			5:  Pink,
			9:  Pink,
			12: Pink,
			15: Pink,
		},
	},

	&ScaleInfo{
		Name:      "19tet",
		Divisions: 19,
		RightStep: 3,
		UpStep:    1,
		Palette: map[int]uint8{
			0:  DarkBlue,
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

	&ScaleInfo{
		Name:      "24tet",
		Divisions: 24,
		RightStep: 4,
		UpStep:    1,
		Palette: map[int]uint8{
			0:  DarkBlue,
			4:  Blue,
			8:  Blue,
			10: Blue,
			14: Blue,
			18: Blue,
			22: Blue,

			2:  Pink,
			6:  Pink,
			12: Pink,
			16: Pink,
			20: Pink,
		},
	},
}
