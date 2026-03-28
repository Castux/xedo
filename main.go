package main

import (
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func main() {
	defer midi.CloseDriver()

	pad := SetupLaunchpad()
	defer pad.Shutdown()

	pad.SetupScale(&Minor)

	for !pad.Exit {

	}
}
