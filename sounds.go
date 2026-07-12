package main

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

const (
	sampleRate = 48000
)

type SoundGenerator struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	rotateSound  *vorbis.Stream
}

func (s *SoundGenerator) Init() {
	s.audioContext = audio.NewContext(sampleRate)
	s.rotateSound = mustLoadSound("assets/sounds/rotate.ogg")
}

func (s *SoundGenerator) PlaySound() {
	if s.audioPlayer != nil {
		s.audioPlayer.Close()
		s.rotateSound.Seek(0, 0)
	}
	s.audioPlayer, _ = s.audioContext.NewPlayer(s.rotateSound)
	s.audioPlayer.Play()
}

func mustLoadSound(name string) *vorbis.Stream {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	stream, err := vorbis.DecodeF32(f)
	if err != nil {
		panic(err)
	}

	return stream
}
