package main

import (
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver

	"fmt"
)

func main() {
	defer midi.CloseDriver()

	synth := SetupSynth()
	defer synth.Shutdown()

	pad := SetupLaunchpad(synth)
	defer pad.Shutdown()

	for !pad.Exit {
		var divisions int
		fmt.Print("Switch to: ")

		_, err := fmt.Scanln(&divisions)
		if err == nil {
			pad.SetupScale(divisions)
		} else {
			fmt.Println(err)
		}
	}
}
