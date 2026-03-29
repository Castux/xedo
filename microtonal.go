package main

import (
	"fmt"
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

func (scale *ScaleInfo) NoteToFreqRatio(note int) float64 {
	octaves := float64(note) / float64(scale.Divisions)
	return math.Pow(2.0, octaves)
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

	freq := pad.BaseFreq * scale.NoteToFreqRatio(baseNote)
	if ev.Down {
		pad.Synth.PlayNote(freq, ev.Velocity)
	} else {
		pad.Synth.StopNote(freq)
	}
}

func (pad *Launchpad) RedrawAllNotes() {
	pad.ForEachPhysicalKey(func(row, col int) {
		note := pad.Scale.RowColToNote(row, col)
		pad.DrawOneIndexed(row, col, pad.Scale.NoteToColor(note))
	})
}

func (pad *Launchpad) SetupScale(divisions int) {
	if pad.Scale == nil || divisions != pad.Scale.Divisions {
		scale := MakeScale(divisions)
		pad.Scale = scale
		pad.OnEvent = scale.OnEvent

		fmt.Println("Switched to", divisions, "TET")
	}

	pad.RedrawAllNotes()
}
