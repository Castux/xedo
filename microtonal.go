package main

import (
	"fmt"
	"math"

	"gitlab.com/gomidi/midi/v2"
)

type ScaleInfo struct {
	Divisions int
	RightStep int
	UpStep    int
	Palette   map[int]uint8
}

func (setup *ScaleInfo) KeyToNote(key uint8) int {
	row, col := KeyToRowCol(key)
	return (row-1)*setup.UpStep + (col-1)*setup.RightStep
}

func (scale *ScaleInfo) NoteToColor(note int) uint8 {
	return scale.Palette[note%scale.Divisions]
}

func (scale *ScaleInfo) NoteToFreq(note int) float64 {
	octaves := float64(note) / float64(scale.Divisions)
	return 440.0 * math.Pow(2.0, octaves)
}

func (scale *ScaleInfo) OnEvent(msg midi.Message, ts int32, pad *Launchpad) {
	var ch, key, vel uint8
	var down bool

	switch {
	case msg.GetNoteStart(&ch, &key, &vel):
		down = true
	case msg.GetNoteEnd(&ch, &key):
		down = false
	default:
		return
	}

	baseNote := scale.KeyToNote(key)
	color := scale.NoteToColor(baseNote)
	if down {
		color = Red
	}

	for row := 1; row <= 8; row++ {
		for col := 1; col <= 8; col++ {
			key2 := KeyFromRowCol(row, col)
			note := scale.KeyToNote(key2)
			if note == baseNote {
				pad.DrawOneIndexed(row, col, color)
			}
		}
	}

	freq := scale.NoteToFreq(baseNote)
	if down {
		fmt.Printf("Key: %d, note: %d, freq: %f\n", key, baseNote, freq)
		pad.Synth.PlayNote(freq)
	} else {
		fmt.Printf("Off %f\n", freq)
		pad.Synth.StopNote(freq)
	}
}

func (pad *Launchpad) SetupScale(scale *ScaleInfo) {
	pad.OnEvent = scale.OnEvent

	for row := 8; row >= 1; row-- {
		for col := 1; col <= 8; col++ {
			key := KeyFromRowCol(row, col)
			note := scale.KeyToNote(uint8(key))
			pad.DrawOneIndexed(row, col, scale.NoteToColor(note))
		}
	}
}
