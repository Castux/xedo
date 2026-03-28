package main

import (
	"math"
	"sync"

	"github.com/gordonklaus/portaudio"
)

type Voice struct {
	Freq       float64
	Volume     float64
	Ticks      int
	KeyOffTime int
}

type Synth struct {
	Stream     *portaudio.Stream
	SampleRate float64

	Mutex        sync.Mutex
	NotesPlaying map[float64]*Voice
}

func SetupSynth() *Synth {
	err := portaudio.Initialize()
	if err != nil {
		panic(err)
	}

	synth := Synth{
		SampleRate: 44100,
	}

	synth.Stream, err = portaudio.OpenDefaultStream(0, 2, synth.SampleRate, 0, synth.GenerateAudio)
	if err != nil {
		panic(err)
	}
	err = synth.Stream.Start()
	if err != nil {
		panic(err)
	}

	synth.NotesPlaying = make(map[float64]*Voice)

	return &synth
}

func (synth *Synth) PlayNote(freq float64) {
	synth.Mutex.Lock()
	defer synth.Mutex.Unlock()

	if synth.NotesPlaying[freq] != nil {
		return
	}

	synth.NotesPlaying[freq] = &Voice{}
}

func (synth *Synth) StopNote(freq float64) {
	synth.Mutex.Lock()
	defer synth.Mutex.Unlock()

	delete(synth.NotesPlaying, freq)
}

func (synth *Synth) GenerateAudio(out [][]float32) {
	synth.Mutex.Lock()
	defer synth.Mutex.Unlock()

	numSamples := len(out[0])

	for i := range numSamples {
		out[0][i] = 0.0
		out[1][i] = 0.0
	}

	for freq, voice := range synth.NotesPlaying {
		for i := range numSamples {
			t := float64(voice.Ticks) / synth.SampleRate
			phase := 2 * math.Pi * freq * t

			out[0][i] += float32(math.Sin(phase))
			out[1][i] += float32(math.Sin(phase))

			voice.Ticks++
		}
	}

	numVoices := len(synth.NotesPlaying)
	if numVoices > 0 {
		for i := range numSamples {
			out[0][i] /= float32(numVoices)
			out[1][i] /= float32(numVoices)
		}
	}
}

func (synth *Synth) Shutdown() {
	synth.Stream.Stop()
	synth.Stream.Close()
	portaudio.Terminate()
}
