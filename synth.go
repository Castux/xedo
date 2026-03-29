package main

import (
	"math"
	"slices"
	"sync"

	"github.com/gordonklaus/portaudio"
)

type Voice interface {
	GenerateSample(sampleRate float64) (float32, float32)
	Frequency() float64
	KeyOff()
	IsDead() bool
}

const (
	Sine = iota
	Square
	Saw
	Triangle
	Piano
	NumShapes
)

var ShapeNames = []string{"sine", "square", "saw", "triangle", "piano"}

type SynthVoice struct {
	Freq       float64
	Volume     float64
	Ticks      int
	KeyOffTime int

	Shape  int
	Attack float64
	Decay  float64

	Dead bool
}

const SynthGain = 0.5
const QuietSynthGain = 1.25

func (voice *SynthVoice) GenerateSample(sampleRate float64) (float32, float32) {
	voice.Ticks++

	t := float64(voice.Ticks) / sampleRate
	period := 1.0 / voice.Freq

	sample := 0.0

	switch voice.Shape {
	case Sine:
		sample = math.Sin(2 * math.Pi * voice.Freq * t)
	case Square:
		sample = 1.0
		if math.Mod(t, period) >= period/2.0 {
			sample = -1.0
		}
	case Saw:
		sample = 2.0*math.Mod(t, period)/period - 1.0
	case Triangle:
		sample = 2.0 * math.Abs(t/period-math.Floor(t/period+0.5)) * QuietSynthGain
	}

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

	sample *= volume * SynthGain
	return float32(sample), float32(sample)
}

func (voice *SynthVoice) IsDead() bool {
	return voice.Dead
}

func (voice *SynthVoice) Frequency() float64 {
	return voice.Freq
}

func (voice *SynthVoice) KeyOff() {
	voice.KeyOffTime = voice.Ticks
}

type Synth struct {
	Stream     *portaudio.Stream
	SampleRate float64
	Shape      int
	Piano      *Sampler

	Mutex  sync.Mutex
	Voices []Voice
	Pedal  bool
}

func SetupSynth() *Synth {
	err := portaudio.Initialize()
	if err != nil {
		panic(err)
	}

	piano := MakeSampler()

	synth := Synth{
		SampleRate: 44100,
		Shape:      Piano,
		Piano:      piano,
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

	var voice Voice

	if synth.Shape == Piano {
		voice = synth.Piano.PlayNote(freq, volume)
	} else {
		voice = &SynthVoice{
			Freq:       freq,
			Volume:     volume,
			Ticks:      0,
			KeyOffTime: math.MaxInt,

			Shape:  synth.Shape,
			Attack: 0.05,
			Decay:  0.2,
		}
	}

	synth.Voices = append(synth.Voices, voice)
}

func (synth *Synth) StopNote(freq float64) {
	if synth.Pedal {
		return
	}

	synth.Mutex.Lock()
	defer synth.Mutex.Unlock()

	for _, voice := range synth.Voices {
		if voice.Frequency() == freq {
			voice.KeyOff()
		}
	}
}

func (synth *Synth) TogglePedal() {
	synth.Pedal = !synth.Pedal

	if !synth.Pedal {
		synth.Mutex.Lock()
		defer synth.Mutex.Unlock()

		for _, voice := range synth.Voices {
			voice.KeyOff()
		}
	}
}

const BaseGain = 0.75

func (synth *Synth) GenerateAudio(out [][]float32) {
	synth.Mutex.Lock()
	defer synth.Mutex.Unlock()

	numSamples := len(out[0])

	for i := range numSamples {
		out[0][i] = 0.0
		out[1][i] = 0.0
	}

	for _, voice := range synth.Voices {
		if voice.IsDead() {
			continue
		}

		for i := range numSamples {
			left, right := voice.GenerateSample(synth.SampleRate)
			out[0][i] += left
			out[1][i] += right
		}
	}

	for i := range numSamples {
		// The poor man's compressor
		out[0][i] *= BaseGain
		out[1][i] *= BaseGain

		if out[0][i] >= BaseGain {
			out[0][i] = BaseGain + (out[0][i]-BaseGain)/3.0
		}
		if out[1][i] >= BaseGain {
			out[1][i] = BaseGain + (out[1][i]-BaseGain)/3.0
		}
	}

	synth.Voices = slices.DeleteFunc(synth.Voices, func(voice Voice) bool {
		return voice.IsDead()
	})
}

func (synth *Synth) Shutdown() {
	synth.Stream.Stop()
	synth.Stream.Close()
	portaudio.Terminate()
}
