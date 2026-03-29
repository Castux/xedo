package main

import (
	"flag"
)

func main() {
	demo := flag.Bool("demo", false, "start in demo mode")
	startDivs := flag.Int("tet", 12, "divisions of octave at startup")
	flag.Parse()

	synth := SetupSynth()
	defer synth.Shutdown()

	pad := SetupLaunchpad(synth, *demo, *startDivs)
	defer pad.Shutdown()

	for !pad.Exit {
	}
}
