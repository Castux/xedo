# XTET

*A microtonal sampler/synthesizer and isomorphic layout for the Launchpad X*

Visually represents and plays [equal temperament tunings](https://en.wikipedia.org/wiki/Equal_temperament) (aka. TET or EDO).

## Requirements

* The Go compiler
* A Novation Launchpad X

## Usage

```
go build .
./xtet [--demo] [--tet <N>]
```

* Demo mode cycles through colorful pixel art
* `--tet` sets the initial divisions (defaults to 12)

# View

For each TET tuning, the pad represents a base note (At 440Hz and other octaves) in purple.

The key layout is isomorphic: moving right always increases by a fixed number of steps, and moving up similarly. For instance, for the standard 12TET, we move up by one semitone and right by one full tone. The program always fits a full octave horizontally.

The layout tiles infinitely in all directions, and a single note appears multiple times. When pressing a key, all notes of the same pitch are lit up in red.

It then colors an equivalent to a standard major scale within this tuning, in blue.

For TET tunings that can be exactly divided into 5 large intervals and 2 small intervals, it will map to the major scale exactly. For instance, the 19TET octave is split into 5 “tones” of 3 steps, and 2 “semitones” of 2 steps.

For the others, it will find the closest approximation.

## Shortcuts

* The four arrow buttons shift the pad’s view to access lower or higher octaves.
* The blue buttons at the top of the pad (“Session” and “Note”) are used to increase or decrease the number of divisions of the octave.
* The pink button (“Custom”) cycles through instruments (piano, sine, square, saw, triangle).
* The red button shuts down the player
* The bottom right button (“> Record arm”) toggles the pedal (keeping sounds playing even after releasing the keys).

## Credits

* Character pixel art by [Johan Vinet](https://johanvinet.tumblr.com/post/127476776680/here-are-100-characters-8x8-pixels-using-the).
* Piano samples: [Salamander Grand Piano](https://github.com/sfzinstruments/SalamanderGrandPiano) v3 by Alexander Holm, Creative Commons Attribution 3.0 Unported License.

## License

XTET is Copyright Noé Falzon 2026, and published under the [MIT license](LICENSE.md)
