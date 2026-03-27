package main

import (
	"fmt"
	"image/color"
//	"math/rand/v2"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)


type Launchpad struct {
	InPort  drivers.In
	OutPort drivers.Out
	Send  func(msg midi.Message) error
	Stop func()

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

	pad.Stop, err = midi.ListenTo(pad.InPort, pad.OnMidiEvent, midi.UseSysEx())
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

func (pad *Launchpad) DrawOne(row, col int, color color.Color) {

	if row < 1 || row > 8 || col < 1 || col > 8 {
		return
	}

	bytes := []byte{0x00, 0x20, 0x29, 0x02, 0x0C, 0x03}
	r,g,b,_ := color.RGBA()
	bytes = append(bytes,
		3, // RGB mode
		byte(row * 10 + col), // key
		byte(r >> 9),
		byte(g >> 9),
		byte(b >> 9),
	)

	pad.Send(midi.SysEx(bytes))
}

func (pad *Launchpad) DrawAll(colors []color.Color) {

	bytes := []byte{0x00, 0x20, 0x29, 0x02, 0x0C, 0x03}
	index := 0
	for row := 8 ; row >= 1; row-- {
		for col := 1 ; col <= 8; col++ {
			if index >= len(colors) {
				break
			}

			r,g,b,_ := colors[index].RGBA()
			index++

			bytes = append(bytes,
				3, // RGB mode
				byte(row * 10 + col), // key
				byte(r >> 9),
				byte(g >> 9),
				byte(b >> 9),
			)
		}
	}

	pad.Send(midi.SysEx(bytes))
}

func KeyToRowCol(key uint8) (int,int) {
	return int(key / 10), int(key % 10)
}

func (pad *Launchpad) OnMidiEvent(msg midi.Message, timestamps int32) {
	var ch, key, vel, controller, value uint8
	var absolute uint16

	switch {
	case msg.GetNoteStart(&ch, &key, &vel):
		fmt.Printf("starting note %s (%d), on channel %v with velocity %v\n", midi.Note(key), key, ch, vel)

		row,col := KeyToRowCol(key)
		pad.DrawOne(row, col, color.RGBA{0xff, 0, 0, 0})
		pad.DrawOne(row + 1, col, color.RGBA{0xff, 0, 0, 0})
		pad.DrawOne(row - 1, col, color.RGBA{0xff, 0, 0, 0})
		pad.DrawOne(row, col + 1, color.RGBA{0xff, 0, 0, 0})
		pad.DrawOne(row, col - 1, color.RGBA{0xff, 0, 0, 0})


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

func main() {
	defer midi.CloseDriver()

	pad := SetupLaunchpad()
	sprites := LoadSprites()

	n := 0
	for !pad.Exit {
		srow, scol := (n%100) / 10, (n%100) % 10
		colors := []color.Color{}
		for row := range 8 {
			for col := range 8 {
				colors = append(colors, sprites.At(
					scol * 12 + col,
					srow * 12 + row,
				))
			}
		}
		pad.DrawAll(colors)

		time.Sleep(1 * time.Second)
		n++
	}

	pad.Shutdown()
}
