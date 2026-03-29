package main

import (
	"flag"
)

func main() {
	demo := flag.Bool("demo", false, "start in demo mode")
	baseFreq := flag.Float64("freq", 440.0, "base frequency")
	startDivs := flag.Int("tet", 12, "divisions of octave at startup")
	flag.Parse()

	synth := SetupSynth()
	defer synth.Shutdown()

	pad := SetupLaunchpad(synth, *demo, *baseFreq, *startDivs)
	defer pad.Shutdown()

	for !pad.Exit {
	}
}
