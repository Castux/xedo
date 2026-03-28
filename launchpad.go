package main

import (
	"fmt"
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
	OnEvent func(msg midi.Message, timestamp int32, pad *Launchpad)
	Stop    func()

	Exit bool
}

func SetupLaunchpad() *Launchpad {
	var pad Launchpad
	var err error

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
		if pad.OnEvent != nil {
			pad.OnEvent(msg, timestamp, &pad)
		}
	})
	if err != nil {
		panic(err)
	}

	programmerMode := []byte{0x00, 0x20, 0x29, 0x02, 0x0C, 0x00, 0x7F}
	pad.Send(midi.SysEx(programmerMode))

	return &pad
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
	key := row * 10 + col
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

func KeyToRowCol(key uint8) (int, int) {
	return int(key / 10), int(key % 10)
}

func KeyFromRowCol(row, col int) uint8 {
	return uint8(10 * row + col)
}

func PrintMidiEvent(msg midi.Message, timestamp int32, pad *Launchpad) {
	var ch, key, vel, controller, value uint8
	var absolute uint16

	switch {
	case msg.GetNoteStart(&ch, &key, &vel):
		fmt.Printf("starting note %s (%d), on channel %v with velocity %v\n", midi.Note(key), key, ch, vel)
	case msg.GetNoteEnd(&ch, &key):
		fmt.Printf("ending note %s on channel %v\n", midi.Note(key), ch)
	case msg.GetControlChange(&ch, &controller, &value):
		fmt.Printf("control change (%d, %d) on channel %v\n", controller, value, ch)

		if controller == 98 {
			pad.Exit = true
			break
		}

	case msg.GetPitchBend(&ch, nil, &absolute):
		fmt.Printf("pitch bend %d in channel %v\n", absolute, ch)
	default:
		// ignore
	}
}

func (pad *Launchpad) DrawSprites() {
	pad.OnEvent = PrintMidiEvent
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
