package main

import (
	"math"
)

type ScaleInfo struct {
	Name      string
	Divisions int
	RightStep int
	UpStep    int
	Palette   map[int]uint8
}

func (setup *ScaleInfo) RowColToNote(row, col int) int {
	return (row-1)*setup.UpStep + (col-1)*setup.RightStep
}

func (scale *ScaleInfo) NoteToColor(note int) uint8 {
	note %= scale.Divisions
	if note < 0 {
		note += scale.Divisions
	}
	return scale.Palette[note%scale.Divisions]
}

func (scale *ScaleInfo) NoteToFreq(note int) float64 {
	octaves := float64(note) / float64(scale.Divisions)
	return 220.0 * math.Pow(2.0, octaves)
}

func (scale *ScaleInfo) OnEvent(ev Event, pad *Launchpad) {

	baseNote := scale.RowColToNote(ev.Row, ev.Col)
	color := scale.NoteToColor(baseNote)
	if ev.Down {
		color = Red
	}

	pad.ForEachPhysicalKey(func(row, col int) {
		note := scale.RowColToNote(row, col)
		if note == baseNote {
			pad.DrawOneIndexed(row, col, color)
		}
	})

	freq := scale.NoteToFreq(baseNote)
	if ev.Down {
		pad.Synth.PlayNote(freq, ev.Velocity)
	} else {
		pad.Synth.StopNote(freq)
	}
}

func (pad *Launchpad) SetupScale() {
	scale := Scales[pad.ScaleIndex]
	pad.OnEvent = scale.OnEvent

	pad.ForEachPhysicalKey(func(row, col int) {
		note := scale.RowColToNote(row, col)
		pad.DrawOneIndexed(row, col, scale.NoteToColor(note))
	})
}
