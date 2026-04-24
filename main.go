package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	demo := flag.Bool("demo", false, "start in demo mode")
	baseFreq := flag.Float64("freq", 440.0, "base frequency")
	startDivs := flag.Int("edo", 12, "divisions of octave at startup")
	flag.Parse()

	synth, err := SetupSynth()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to set up synth: %v\n", err)
		os.Exit(1)
	}
	defer synth.Shutdown()

	pad, err := SetupLaunchpad(synth, *demo, *baseFreq, *startDivs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to set up launchpad: %v\n", err)
		os.Exit(1)
	}
	defer pad.Shutdown()

	for !pad.Exit {
	}
}
