package main

import (
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func main() {
	defer midi.CloseDriver()

	synth := SetupSynth()
	defer synth.Shutdown()

	pad := SetupLaunchpad(synth)
	defer pad.Shutdown()

	pad.SetupScale(&Major)

	for !pad.Exit {

	}
}
