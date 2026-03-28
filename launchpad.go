package main

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

type Launchpad struct {
	InPort  drivers.In
	OutPort drivers.Out
	Send    func(msg midi.Message) error
	OnEvent func(ev Event, pad *Launchpad)
	Stop    func()

	RowOffset int
	ColOffset int

	Synth      *Synth
	ScaleIndex int

	Exit bool
}

type Event struct {
	Row, Col int
	Down     bool
	Velocity float64
}

func SetupLaunchpad(synth *Synth) *Launchpad {
	var pad Launchpad
	var err error

	pad.Synth = synth

	pad.InPort, err = midi.FindInPort("Launchpad X LPX MIDI Out")
	if err != nil {
		panic(err)
	}

	pad.OutPort, err = midi.FindOutPort("Launchpad X LPX MIDI In")
	if err != nil {
		panic(err)
	}

	pad.Send, err = midi.SendTo(pad.OutPort)
	if err != nil {
		panic(err)
	}

	pad.Stop, err = midi.ListenTo(pad.InPort, func(msg midi.Message, timestamp int32) {

		var ch, key, vel uint8
		var down bool

		switch {
		case msg.GetNoteStart(&ch, &key, &vel):
			down = true
		case msg.GetNoteEnd(&ch, &key):
			down = false
		case msg.GetControlChange(&ch, &key, &vel):
			if vel > 0 {
				switch key {
				case 91:
					pad.RowOffset++
				case 92:
					pad.RowOffset--
				case 93:
					pad.ColOffset--
				case 94:
					pad.ColOffset++
				case 98:
					pad.Exit = true
				case 19:
					pad.Synth.Shape++
					pad.Synth.Shape %= NumShapes
					fmt.Println("Synth switched to", ShapeNames[pad.Synth.Shape])
				case 29:
					pad.ScaleIndex++
					pad.ScaleIndex %= len(Scales)
					fmt.Println("Scale switched to", Scales[pad.ScaleIndex].Name)
				}
				pad.SetupScale()

			}
			return
		default:
			return
		}
		row, col := pad.KeyToRowCol(int(key))

		ev := Event{
			Row:      row,
			Col:      col,
			Down:     down,
			Velocity: float64(vel) / float64(0xff),
		}

		if pad.OnEvent != nil {
			pad.OnEvent(ev, &pad)
		}
	})
	if err != nil {
		panic(err)
	}

	programmerMode := []byte{0x00, 0x20, 0x29, 0x02, 0x0C, 0x00, 0x7F}
	pad.Send(midi.SysEx(programmerMode))

	pad.ScaleIndex = 2
	pad.SetupScale()

	return &pad
}

func (pad *Launchpad) KeyToRowCol(key int) (int, int) {
	return key/10 + pad.RowOffset, key%10 + pad.ColOffset
}

func (pad *Launchpad) KeyFromRowCol(row, col int) int {
	row -= pad.RowOffset
	col -= pad.ColOffset
	return 10*row + col
}

func (pad *Launchpad) ForEachPhysicalKey(cb func(row, col int)) {
	for row := range 8 {
		for col := range 8 {
			cb(row+1+pad.RowOffset, col+1+pad.ColOffset)
		}
	}
}

func (pad *Launchpad) Shutdown() {

	for key := 11; key <= 99; key++ {
		msg := midi.NoteOff(0, uint8(key))
		pad.Send(msg)
	}

	custom4 := []byte{0x00, 0x20, 0x29, 0x02, 0x0C, 0x00, 0x07}
	pad.Send(midi.SysEx(custom4))

	pad.Stop()
}

func (pad *Launchpad) DrawOneIndexed(row, col int, color uint8) {
	key := pad.KeyFromRowCol(row, col)
	pad.Send(midi.NoteOn(0, uint8(key), color))
}

func (pad *Launchpad) DrawOne(row, col int, color color.Color) {
	key := pad.KeyFromRowCol(row, col)
	if key < 0 || key > 0xff {
		fmt.Println("Bad key", key)
	}

	bytes := []byte{0x00, 0x20, 0x29, 0x02, 0x0C, 0x03}
	r, g, b, _ := color.RGBA()
	bytes = append(bytes,
		3,          // RGB mode
		uint8(key), // key
		byte(r>>9), // Go from 0xFFFF to 0x7F
		byte(g>>9),
		byte(b>>9),
	)

	pad.Send(midi.SysEx(bytes))
}

func (pad *Launchpad) DrawAll(colors []color.Color) {

	bytes := []byte{0x00, 0x20, 0x29, 0x02, 0x0C, 0x03}
	index := 0
	for row := 8; row >= 1; row-- {
		for col := 1; col <= 8; col++ {
			if index >= len(colors) {
				break
			}

			r, g, b, _ := colors[index].RGBA()
			index++

			bytes = append(bytes,
				3,                // RGB mode
				byte(row*10+col), // key
				byte(r>>9),       // Go from 0xFFFF to 0x7F
				byte(g>>9),
				byte(b>>9),
			)
		}
	}

	pad.Send(midi.SysEx(bytes))
}

func (pad *Launchpad) DrawSprites() {
	sprites := LoadSprites()

	for !pad.Exit {

		n := rand.IntN(100)

		srow, scol := n/10, n%10
		colors := []color.Color{}
		for row := range 8 {
			for col := range 8 {
				colors = append(colors, sprites.At(
					scol*12+col,
					srow*12+row,
				))
			}
		}
		pad.DrawAll(colors)

		time.Sleep(3 * time.Second)
		n++
	}
}
