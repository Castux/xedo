package main

import (
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	"image/color"
	"math/rand/v2"
	"time"
)

type Launchpad struct {
	InPort  drivers.In
	OutPort drivers.Out
	Send    func(msg midi.Message) error
	OnEvent func(ev Event, pad *Launchpad)
	Stop    func()

	Synth *Synth

	Exit bool
}

type Event struct {
	Row, Col int
	Down bool
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
		default:
			return
		}
		row, col := pad.KeyToRowCol(key)

		ev := Event{
			Row: row,
			Col: col,
			Down: down,
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

	return &pad
}

func (pad *Launchpad) KeyToRowCol(key uint8) (int, int) {
	return int(key / 10), int(key % 10)
}

func (pad *Launchpad) KeyFromRowCol(row, col int) uint8 {
	return uint8(10*row + col)
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
	key := row*10 + col
	pad.Send(midi.NoteOn(0, uint8(key), color))
}

func (pad *Launchpad) DrawOne(row, col int, color color.Color) {

	if row < 1 || row > 8 || col < 1 || col > 8 {
		return
	}

	bytes := []byte{0x00, 0x20, 0x29, 0x02, 0x0C, 0x03}
	r, g, b, _ := color.RGBA()
	bytes = append(bytes,
		3,                // RGB mode
		byte(row*10+col), // key
		byte(r>>9),       // Go from 0xFFFF to 0x7F
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
