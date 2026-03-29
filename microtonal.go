package main

import (
	"fmt"
	"math"
)

type Tuning struct {
	Divisions int
	RightStep int
	UpStep    int
	Palette   map[int]uint8
}

func (tuning *Tuning) RowColToNote(row, col int) int {
	return (row-1)*tuning.UpStep + (col-1)*tuning.RightStep
}

func (tuning *Tuning) NoteToColor(note int) uint8 {
	note %= tuning.Divisions
	if note < 0 {
		note += tuning.Divisions
	}
	return tuning.Palette[note%tuning.Divisions]
}

func (tuning *Tuning) NoteToFreqRatio(note int) float64 {
	octaves := float64(note) / float64(tuning.Divisions)
	return math.Pow(2.0, octaves)
}

func (tuning *Tuning) OnEvent(ev Event, pad *Launchpad) {

	baseNote := tuning.RowColToNote(ev.Row, ev.Col)
	color := tuning.NoteToColor(baseNote)
	if ev.Down {
		color = Red
	}

	pad.ForEachPhysicalKey(func(row, col int) {
		note := tuning.RowColToNote(row, col)
		if note == baseNote {
			pad.DrawOneIndexed(row, col, color)
		}
	})

	freq := pad.BaseFreq * tuning.NoteToFreqRatio(baseNote)
	if ev.Down {
		pad.Synth.PlayNote(freq, ev.Velocity)
	} else {
		pad.Synth.StopNote(freq)
	}
}

func (pad *Launchpad) RedrawAllNotes() {
	pad.ForEachPhysicalKey(func(row, col int) {
		note := pad.Tuning.RowColToNote(row, col)
		pad.DrawOneIndexed(row, col, pad.Tuning.NoteToColor(note))
	})
}

func (pad *Launchpad) SetupTuning(divisions int) {
	if pad.Tuning == nil || divisions != pad.Tuning.Divisions {
		tuning := MakeTuning(divisions)
		pad.Tuning = tuning
		pad.OnEvent = tuning.OnEvent

		fmt.Println("Switched to", divisions, "EDO")
	}

	pad.RedrawAllNotes()
}
