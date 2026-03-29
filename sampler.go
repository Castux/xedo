package main

import (
	"cmp"
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"

	"github.com/jonchammer/audio-io/wave"
)

func LoadWav(path string) ([]float64, []float64) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := wave.NewReader(file)
	header, err := r.Header()
	if err != nil {
		panic(err)
	}

	sampleType, err := header.SampleType()
	if err != nil || sampleType != wave.SampleTypeInt16 {
		panic("Bad sample type")
	}

	if header.ChannelCount() != 2 {
		panic("Expected stereo samples")
	}

	if header.FrameRate() != 44100 {
		panic("Expected 44100 sample rate")
	}

	data := make([]int16, header.SampleCount())
	samplesRead, err := r.ReadInt16(data)
	if err != nil {
		panic(err)
	}
	if samplesRead != int(header.SampleCount()) {
		panic(fmt.Sprintf("Did not read all samples for %s", path))
	}

	left := make([]float64, samplesRead/2)
	right := make([]float64, samplesRead/2)
	scale := float64(math.MaxInt16)
	for i := range samplesRead / 2 {
		left[i] = float64(data[2*i+0]) / scale
		right[i] = float64(data[2*i+1]) / scale
	}

	return left, right
}

var SampleNameRegex = regexp.MustCompile(`(C|D#|F#|A)(\d)v10\.wav`)

var NameToPitchSteps = map[string]int{
	"C":  -9,
	"D#": -6,
	"F#": -3,
	"A":  0,
}

func NoteNameToPitch(note string, octave int) float64 {
	steps := NameToPitchSteps[note] + (octave-4)*12
	freq := 440.0 * math.Pow(2, float64(steps)/12)

	return freq
}

type Sample struct {
	Note   string
	Octave int
	Freq   float64
	Left   []float64
	Right   []float64
}

const SampleDir = "piano"
func LoadPianoSamples() []*Sample {
	dir, err := os.ReadDir(SampleDir)
	if err != nil {
		panic(err)
	}

	samples := []*Sample{}

	for _, entry := range dir {
		submatches := SampleNameRegex.FindStringSubmatch(entry.Name())
		if len(submatches) != 0 {
			note := submatches[1]
			octave, _ := strconv.Atoi(submatches[2])
			freq := NoteNameToPitch(note, octave)

			left, right := LoadWav(SampleDir + "/" + entry.Name())

			samples = append(samples, &Sample{note, octave, freq, left, right})
		}
	}

	slices.SortFunc(samples, func(a, b *Sample) int {
		return cmp.Compare(a.Freq, b.Freq)
	})

	return samples
}

type Sampler struct {
	Samples []*Sample
}

func MakeSampler() *Sampler {
	samples := LoadPianoSamples()

	return &Sampler{
		Samples: samples,
	}
}
