package main

import (
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func main() {
	defer midi.CloseDriver()

	pad := SetupLaunchpad()
	pad.DrawSprites()
	pad.Shutdown()
}
