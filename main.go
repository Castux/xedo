package main

import ()

func main() {
	synth := SetupSynth()
	defer synth.Shutdown()

	pad := SetupLaunchpad(synth)
	defer pad.Shutdown()

	for !pad.Exit {
	}
}
