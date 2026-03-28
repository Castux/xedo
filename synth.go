package main

import (
	"math"
	"slices"
	"sync"

	"github.com/gordonklaus/portaudio"
)

type Voice struct {
	Freq       float64
	Volume     float64
	Ticks      int
	KeyOffTime int

	Attack float64
	Decay  float64

	Dead bool
}

func (voice *Voice) GenerateSample(sampleRate float64) (float32, float32) {
	voice.Ticks++

	t := float64(voice.Ticks) / sampleRate
	phase := 2 * math.Pi * voice.Freq * t

	volume := voice.Volume
	if t <= voice.Attack {
		volume *= t / voice.Attack
	}

	keyOff := float64(voice.KeyOffTime) / sampleRate
	if t >= keyOff {
		volume *= max(0.0, 1.0-(t-keyOff)/voice.Decay)
	}
	if t >= keyOff+voice.Decay {
		voice.Dead = true
	}

	sample := float32(math.Sin(phase) * volume)
	return sample, sample
}

type Synth struct {
	Stream     *portaudio.Stream
	SampleRate float64

	Mutex        sync.Mutex
	NotesPlaying []*Voice
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

	return &synth
}

func (synth *Synth) PlayNote(freq float64, volume float64) {
	synth.Mutex.Lock()
	defer synth.Mutex.Unlock()

	voice := &Voice{
		Freq:       freq,
		Volume:     volume,
		Ticks:      0,
		KeyOffTime: math.MaxInt,

		Attack: 0.1,
		Decay:  0.2,
	}

	synth.NotesPlaying = append(synth.NotesPlaying, voice)
}

func (synth *Synth) StopNote(freq float64) {
	synth.Mutex.Lock()
	defer synth.Mutex.Unlock()

	for _, voice := range synth.NotesPlaying {
		if voice.Freq == freq {
			voice.KeyOffTime = voice.Ticks
		}
	}
}

func (synth *Synth) GenerateAudio(out [][]float32) {
	synth.Mutex.Lock()
	defer synth.Mutex.Unlock()

	numSamples := len(out[0])

	for i := range numSamples {
		out[0][i] = 0.0
		out[1][i] = 0.0
	}

	for _, voice := range synth.NotesPlaying {
		if voice.Dead {
			continue
		}

		for i := range numSamples {
			left, right := voice.GenerateSample(synth.SampleRate)
			out[0][i] += left
			out[1][i] += right
		}
	}

	synth.NotesPlaying = slices.DeleteFunc(synth.NotesPlaying, func(voice *Voice) bool {
		return voice.Dead
	})
}

func (synth *Synth) Shutdown() {
	synth.Stream.Stop()
	synth.Stream.Close()
	portaudio.Terminate()
}
