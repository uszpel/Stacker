package main

import (
	"bytes"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

const (
	sampleRate = 44100
)

const (
	soundStateUnmuted = iota
	soundStateMuted
	soundStateMuteRequested
	soundStateUnmuteRequested
)

const (
	soundRotate = iota
	soundDrop
	soundLineComplete
)

type SoundGenerator struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	sounds       [][]byte
	state        int
}

func (s *SoundGenerator) Init() {
	s.audioContext = audio.NewContext(sampleRate)
	s.sounds = make([][]byte, 3)
	s.sounds[soundRotate] = mustLoadSound("assets/sounds/rotate.ogg")
	s.sounds[soundDrop] = mustLoadSound("assets/sounds/drop.ogg")
	s.sounds[soundLineComplete] = mustLoadSound("assets/sounds/line_complete.ogg")
	s.state = soundStateUnmuted
}

func (s *SoundGenerator) RequestStateChange() {
	switch s.state {
	case soundStateUnmuted:
		s.state = soundStateMuteRequested
	case soundStateMuted:
		s.state = soundStateUnmuteRequested
	}
}

func (s *SoundGenerator) UpdateState() {
	switch s.state {
	case soundStateMuteRequested:
		s.state = soundStateMuted
	case soundStateUnmuteRequested:
		s.state = soundStateUnmuted
	}
}

func (s *SoundGenerator) PlaySound(sound int) {
	if s.audioPlayer != nil {
		s.audioPlayer.Close()
	}
	if s.state == soundStateMuted || sound < 0 || sound >= len(s.sounds) {
		return
	}
	s.audioPlayer, _ = s.audioContext.NewPlayer(bytes.NewReader(s.sounds[sound]))
	s.audioPlayer.SetVolume(0.1)
	s.audioPlayer.Play()
}

func mustLoadSound(name string) []byte {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	stream, err := vorbis.DecodeWithoutResampling(f)
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(stream)
	if err != nil {
		panic(err)
	}
	return data
}
